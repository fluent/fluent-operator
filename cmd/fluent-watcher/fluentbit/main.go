package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fluent/fluent-operator/v3/pkg/filenotify"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBinPath        = "/fluent-bit/bin/fluent-bit"
	defaultCfgPath        = "/fluent-bit/etc/fluent-bit.conf"
	defaultSecretYamlPath = "/fluent-bit/config/fluent-bit.yaml"
	defaultYamlCfgPath    = "/fluent-bit/etc/fluent-bit.yaml"
	defaultParserCfgPath  = "/fluent-bit/etc/parser.conf"
	defaultWatchDir       = "/fluent-bit/config"
	defaultPollInterval   = 1 * time.Second
	defaultReloadDelay    = 1 * time.Second
	defaultReloadMaxDelay = 10 * time.Second

	// dataLinkName is the symlink the kubelet replaces to publish a new version of a mounted
	// Secret, ConfigMap or projected volume.
	dataLinkName = "..data"
)

func main() {
	var configPath string
	var externalPluginPath string
	var binPath string
	var watchPath string
	var poll bool
	var pollInterval time.Duration
	var parserCfgPath string
	var reloadDelay time.Duration
	var reloadMaxDelay time.Duration
	flag.StringVar(&binPath, "b", defaultBinPath, "The fluent bit binary path.")
	flag.StringVar(&externalPluginPath, "e", "", "Path to external plugin (shared lib)")
	flag.StringVar(&watchPath, "watch-path", defaultWatchDir, "The path to watch.")
	flag.BoolVar(&poll, "poll", false, "Use poll watcher instead of ionotify.")
	flag.DurationVar(&pollInterval, "poll-interval", defaultPollInterval, "Poll interval if using poll watcher.")
	flag.StringVar(&parserCfgPath, "R", defaultParserCfgPath, "Specify a parser config file")
	flag.DurationVar(&reloadDelay, "reload-delay", defaultReloadDelay,
		"How long to wait before the config is reloaded: after the new version is published "+
			"if the watch path is a Kubernetes volume, after the last change otherwise.")
	flag.DurationVar(&reloadMaxDelay, "reload-max-delay", defaultReloadMaxDelay,
		"The longest a reload is postponed by further changes of a watch path that is not a "+
			"Kubernetes volume. 0 means no limit.")

	// Deprecated flags to be removed in one of the next releases.
	var exitOnFailure bool
	var flbTerminationTimeout time.Duration
	flag.BoolVar(&exitOnFailure, "exit-on-failure", false, "Deprecated: This has no effect anymore.")
	flag.DurationVar(&flbTerminationTimeout, "flb-timeout", 0, "Deprecated: This has no effect anymore.")

	ctx := context.Background()
	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	l := log.NewLogfmtLogger(os.Stdout)
	l = log.With(l, "time", log.TimestampFormat(time.Now, time.RFC3339))

	// check  the config file format
	_, err := os.Stat(defaultSecretYamlPath)
	if os.IsNotExist(err) {
		_ = level.Info(l).Log("msg", "No fluent-bit secret yaml found, using classic one.")
		flag.StringVar(&configPath, "c", defaultCfgPath, "The classic config file path.")
	} else {
		_ = level.Info(l).Log("msg", "fluent-bit secret yaml found, using yaml one.")
		flag.StringVar(&configPath, "c", defaultYamlCfgPath, "The yaml config file path.")
	}

	if exitOnFailure {
		_ = level.Warn(l).
			Log("--exit-on-failure is deprecated. The process will exit no matter what if " +
				"fluent-bit exits so this can safely be removed.")
	}
	if flbTerminationTimeout > 0 {
		_ = level.Warn(l).
			Log("--flb-timeout is deprecated. Consider setting the terminationGracePeriod " +
				"field on the `(Cluster)FluentBit` instance.")
	}

	flag.Parse()

	// Remember which version of a mounted volume is published before fluent-bit reads it, so a
	// version published at any point after that is noticed.
	dataLink := findDataLink(watchPath)
	var generation string
	if dataLink != "" {
		// Events of a single file do not show a switch of "..data", watch its directory.
		watchPath = filepath.Dir(dataLink)
		if generation, err = os.Readlink(dataLink); err != nil {
			_ = level.Warn(l).Log("msg", "failed to read the published config version", "error", err)
		}
	}

	// First, launch the fluent-bit process.
	args := []string{"--enable-hot-reload", "-c", configPath}
	if externalPluginPath != "" {
		args = append(args, "-e", externalPluginPath)
	}
	cmd := exec.Command(binPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		_ = level.Error(l).Log("msg", "failed to start fluent-bit", "error", err)
		os.Exit(1)
	}
	_ = level.Info(l).Log("msg", "fluent-bit watcher started")

	grp, grpCtx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		// Watch the process. If it exits, we want to crash immediately.
		defer cancel()
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("failed to run fluent-bit: %w", err)
		}
		return nil
	})
	grp.Go(func() error {
		// Watch the config as it's loaded into the pod and trigger a config reload.
		var watcher filenotify.FileWatcher
		if poll {
			watcher = filenotify.NewPollingWatcher(pollInterval)
		} else {
			var err error
			watcher, err = filenotify.NewEventWatcher()
			if err != nil {
				return fmt.Errorf("failed to open event watcher: %w", err)
			}
		}

		if err := watcher.Add(watchPath); err != nil {
			return fmt.Errorf("failed to watch path %q: %w", watchPath, err)
		}

		opts := watchOptions{
			dataLink:   dataLink,
			generation: generation,
			delay:      reloadDelay,
			maxDelay:   reloadMaxDelay,
			logger:     l,
		}
		if poll {
			ticker := time.NewTicker(pollInterval)
			defer ticker.Stop()
			opts.tick = ticker.C
		}

		return watchAndReload(signalCtx, grpCtx, watcher, opts, func() error {
			_ = level.Info(l).Log("msg", "Config file changed, reloading...")
			if err := cmd.Process.Signal(syscall.SIGHUP); err != nil {
				return fmt.Errorf("failed to reload config: %w", err)
			}
			return nil
		})
	})

	select {
	case <-signalCtx.Done():
	case <-grpCtx.Done():
	}

	// Always try to gracefully shut down fluent-bit. This will allow `cmd.Wait` above to finish
	// and thus allow `grp.Wait` below to return.
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		_ = level.Error(l).Log("msg", "Failed to send SIGTERM to fluent-bit", "error", err)
		// Do not exit on error here. The process might've died and that's okay.
	}

	if err := grp.Wait(); err != nil {
		_ = level.Error(l).Log("msg", "Failure during the run time of fluent-bit", "error", err)
		os.Exit(1)
	}
}

// findDataLink returns the path of the "..data" symlink if watchPath is a volume mounted by the
// kubelet or a file in one, and an empty string otherwise.
func findDataLink(watchPath string) string {
	for _, dir := range []string{watchPath, filepath.Dir(watchPath)} {
		link := filepath.Join(dir, dataLinkName)
		if fi, err := os.Lstat(link); err == nil && fi.Mode()&os.ModeSymlink != 0 {
			return link
		}
	}
	return ""
}

type watchOptions struct {
	// dataLink is the path of the "..data" symlink of a volume mounted by the kubelet. Empty if
	// the watch path is not such a volume.
	dataLink string
	// generation is the target of dataLink fluent-bit was started with.
	generation string
	// delay is how long to wait before a reload: after a new target of dataLink is seen, or
	// after the last change if dataLink is empty.
	delay time.Duration
	// maxDelay limits how long further changes can postpone a reload, counted from the first
	// change. Only used if dataLink is empty. 0 means no limit.
	maxDelay time.Duration
	// tick makes the watcher compare the target of dataLink without an event. Set in poll mode.
	tick   <-chan time.Time
	logger log.Logger
}

// watchAndReload calls reload when the config in the watched path has changed.
//
// The kubelet prepares a new version of a mounted Secret or ConfigMap next to the current one
// and publishes it by replacing the "..data" symlink. Files appearing in the volume are not a
// config change yet: a reload at that moment makes fluent-bit load the previous version, and
// fluent-bit drops the signals that arrive while that reload is running. So for such a volume
// the only trigger is a new target of "..data"; events and poll ticks are just a reason to
// look at it. The reload follows opts.delay after the switch.
//
// A path without "..data" has no such marker. There reload is called once the path has been
// quiet for opts.delay, and no later than opts.maxDelay after the first change, so that a file
// which is written all the time cannot hold the reload back forever.
func watchAndReload(
	signalCtx, grpCtx context.Context,
	watcher filenotify.FileWatcher,
	opts watchOptions,
	reload func() error,
) error {
	// settle fires opts.delay after the last accepted change, deadline opts.maxDelay after the
	// first one. Both are nil while there is nothing to reload.
	var settle, deadline <-chan time.Time

	generation := opts.generation
	checkGeneration := func() {
		target, err := os.Readlink(opts.dataLink)
		if err != nil {
			// Not a reason to reload: whatever is there is not a published version.
			_ = level.Warn(opts.logger).Log("msg", "failed to read the published config version", "error", err)
			return
		}
		if target == generation {
			return
		}
		generation = target
		settle = time.After(opts.delay)
	}
	if opts.dataLink != "" {
		// A version published before the watch was added gives no event.
		checkGeneration()
	}

	for {
		select {
		case <-signalCtx.Done():
			return nil
		case <-grpCtx.Done():
			return nil
		case event := <-watcher.Events():
			if opts.dataLink != "" {
				checkGeneration()
				continue
			}
			if !isValidEvent(event) {
				continue
			}
			settle = time.After(opts.delay)
			if deadline == nil && opts.maxDelay > 0 {
				deadline = time.After(opts.maxDelay)
			}
		case <-opts.tick:
			if opts.dataLink != "" {
				checkGeneration()
			}
		case <-settle:
			settle, deadline = nil, nil
			if err := reload(); err != nil {
				return err
			}
		case <-deadline:
			settle, deadline = nil, nil
			if err := reload(); err != nil {
				return err
			}
		case err := <-watcher.Errors():
			return fmt.Errorf("failed the watcher: %w", err)
		}
	}
}

// Inspired by https://github.com/jimmidyson/configmap-reload
func isValidEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Create || event.Op == fsnotify.Write
}

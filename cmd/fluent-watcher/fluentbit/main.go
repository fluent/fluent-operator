package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fluent/fluent-operator/v2/pkg/filenotify"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBinPath      = "/fluent-bit/bin/fluent-bit"
	defaultCfgPath      = "/fluent-bit/etc/fluent-bit.conf"
	defaultWatchDir     = "/fluent-bit/config"
	defaultPollInterval = 1 * time.Second
)

func main() {
	var configPath string
	var externalPluginPath string
	var binPath string
	var watchPath string
	var poll bool
	var pollInterval time.Duration
	flag.StringVar(&binPath, "b", defaultBinPath, "The fluent bit binary path.")
	flag.StringVar(&configPath, "c", defaultCfgPath, "The config file path.")
	flag.StringVar(&externalPluginPath, "e", "", "Path to external plugin (shared lib)")
	flag.StringVar(&watchPath, "watch-path", defaultWatchDir, "The path to watch.")
	flag.BoolVar(&poll, "poll", false, "Use poll watcher instead of ionotify.")
	flag.DurationVar(&pollInterval, "poll-interval", defaultPollInterval, "Poll interval if using poll watcher.")

	// Deprecated flags to be removed in one of the next releases.
	var exitOnFailure bool
	var flbTerminationTimeout time.Duration
	flag.BoolVar(&exitOnFailure, "exit-on-failure", false, "Deprecated: This has no effect anymore.")
	flag.DurationVar(&flbTerminationTimeout, "flb-timeout", 0, "Deprecated: This has no effect anymore.")

	flag.Parse()

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "time", log.TimestampFormat(time.Now, time.RFC3339))

	if exitOnFailure {
		level.Warn(logger).Log("--exit-on-failure is deprecated. The process will exit no matter what if fluent-bit exits so this can safely be removed.")
	}
	if flbTerminationTimeout > 0 {
		level.Warn(logger).Log("--flb-timeout is deprecated. Consider setting the terminationGracePeriod field on the `(Cluster)FluentBit` instance.")
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
		_ = level.Error(logger).Log("msg", "failed to start fluent-bit", "error", err)
		os.Exit(1)
	}
	_ = level.Info(logger).Log("msg", "fluent-bit started")

	grp, grpCtx := errgroup.WithContext(context.Background())
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

		for {
			select {
			case <-signalCtx.Done():
				return nil
			case <-grpCtx.Done():
				return nil
			case event := <-watcher.Events():
				if !isValidEvent(event) {
					continue
				}
				_ = level.Info(logger).Log("msg", "Config file changed, reloading...")
				if err := cmd.Process.Signal(syscall.SIGHUP); err != nil {
					return fmt.Errorf("failed to reload config: %w", err)
				}
			case err := <-watcher.Errors():
				return fmt.Errorf("failed the watcher: %w", err)
			}
		}
	})

	select {
	case <-signalCtx.Done():
	case <-grpCtx.Done():
	}

	// Always try to gracefully shut down fluent-bit. This will allow `cmd.Wait` above to finish
	// and thus allow `grp.Wait` below to return.
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		_ = level.Error(logger).Log("msg", "Failed to send SIGTERM to fluent-bit", "error", err)
		// Do not exit on error here. The process might've died and that's okay.
	}

	if err := grp.Wait(); err != nil {
		_ = level.Error(logger).Log("msg", "Failure during the run time of fluent-bit", "error", err)
		os.Exit(1)
	}
}

// Inspired by https://github.com/jimmidyson/configmap-reload
func isValidEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Create || event.Op == fsnotify.Write
}

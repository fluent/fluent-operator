package main

import (
	"context"
	"flag"
	"math"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fluent/fluent-operator/pkg/filenotify"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
)

const (
	defaultBinPath      = "/fluent-bit/bin/fluent-bit"
	defaultCfgPath      = "/fluent-bit/etc/fluent-bit.conf"
	defaultWatchDir     = "/fluent-bit/config"
	defaultPollInterval = 1 * time.Second

	MaxDelayTime = 5 * time.Minute
	ResetTime    = 10 * time.Minute
)

var (
	logger       log.Logger
	cmd          *exec.Cmd
	mutex        sync.Mutex
	restartTimes int32
	timerCtx     context.Context
	timerCancel  context.CancelFunc
)

var configPath string
var externalPluginPath string
var binPath string
var watchPath string
var poll bool
var exitOnFailure bool
var pollInterval time.Duration

func main() {

	flag.StringVar(&binPath, "b", defaultBinPath, "The fluent bit binary path.")
	flag.StringVar(&configPath, "c", defaultCfgPath, "The config file path.")
	flag.StringVar(&externalPluginPath, "e", "", "Path to external plugin (shared lib)")
	flag.BoolVar(&exitOnFailure, "exit-on-failure", false, "If fluentbit exits with failure, also exit the watcher.")
	flag.StringVar(&watchPath, "watch-path", defaultWatchDir, "The path to watch.")
	flag.BoolVar(&poll, "poll", false, "Use poll watcher instead of ionotify.")
	flag.DurationVar(&pollInterval, "poll-interval", defaultPollInterval, "Poll interval if using poll watcher.")

	flag.Parse()

	logger = log.NewLogfmtLogger(os.Stdout)

	timerCtx, timerCancel = context.WithCancel(context.Background())

	var g run.Group
	{
		// Termination handler.
		g.Add(run.SignalHandler(context.Background(), os.Interrupt, syscall.SIGTERM))
	}
	{
		// Watch the Fluent bit, if the Fluent bit not exists or stopped, restart it.
		cancel := make(chan struct{})
		g.Add(
			func() error {

				for {
					select {
					case <-cancel:
						return nil
					default:
					}

					// Start fluent bit if it does not existed.
					start()
					// Wait for the fluent bit exit.
					err := wait()
					if exitOnFailure && err != nil {
						_ = level.Error(logger).Log("msg", "Fluent bit exited with error; exiting watcher")
						return err
					}

					timerCtx, timerCancel = context.WithCancel(context.Background())

					// After the fluent bit exit, fluent bit watcher restarts it with an exponential
					// back-off delay (1s, 2s, 4s, ...), that is capped at five minutes.
					backoff()
				}
			},
			func(err error) {
				close(cancel)
				stop()
				resetTimer()
			},
		)
	}
	{
		// Watch the config file, if the config file changed, stop Fluent bit.
		watcher, err := newWatcher(poll, pollInterval)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return
		}

		// Start watcher.
		err = watcher.Add(watchPath)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return
		}

		cancel := make(chan struct{})
		g.Add(
			func() error {

				for {
					select {
					case <-cancel:
						return nil
					case event := <-watcher.Events():
						if !isValidEvent(event) {
							continue
						}

						_ = level.Info(logger).Log("msg", "Config file changed, stopping Fluent Bit")

						// After the config file changed, it should stop the fluent bit,
						// and resets the restart backoff timer.
						stop()
						resetTimer()
						_ = level.Info(logger).Log("msg", "Config file changed, stopped Fluent Bit")
					case <-watcher.Errors():
						_ = level.Error(logger).Log("msg", "Watcher stopped")
						return nil
					}
				}
			},
			func(err error) {
				_ = watcher.Close()
				close(cancel)
			},
		)
	}

	if err := g.Run(); err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	_ = level.Info(logger).Log("msg", "See you next time!")
}

func newWatcher(poll bool, interval time.Duration) (filenotify.FileWatcher, error) {
	var err error
	var watcher filenotify.FileWatcher

	if poll {
		watcher = filenotify.NewPollingWatcher(interval)
	} else {
		watcher, err = filenotify.New(interval)
	}

	if err != nil {
		return nil, err
	}

	return watcher, nil
}

// Inspired by https://github.com/jimmidyson/configmap-reload
func isValidEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Create || event.Op == fsnotify.Write
}

func start() {

	mutex.Lock()
	defer mutex.Unlock()

	if cmd != nil {
		return
	}

	if externalPluginPath != "" {
		cmd = exec.Command(binPath, "-c", configPath, "-e", externalPluginPath)
	} else {
		cmd = exec.Command(binPath, "-c", configPath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		_ = level.Error(logger).Log("msg", "start Fluent bit error", "error", err)
		cmd = nil
		return
	}

	_ = level.Info(logger).Log("msg", "Fluent bit started")
}

func wait() error {
	mutex.Lock()
	if cmd == nil {
		mutex.Unlock()
		return nil
	}
	mutex.Unlock()

	startTime := time.Now()
	err := cmd.Wait()
	_ = level.Error(logger).Log("msg", "Fluent bit exited", "error", err)
	// Once the fluent bit has executed for 10 minutes without any problems,
	// it should resets the restart backoff timer.
	if time.Since(startTime) >= ResetTime {
		atomic.StoreInt32(&restartTimes, 0)
	}

	mutex.Lock()
	cmd = nil
	mutex.Unlock()
	return err
}

func backoff() {

	delayTime := time.Duration(math.Pow(2, float64(atomic.LoadInt32(&restartTimes)))) * time.Second
	if delayTime >= MaxDelayTime {
		delayTime = MaxDelayTime
	}

	_ = level.Info(logger).Log("msg", "backoff", "delay", delayTime)

	startTime := time.Now()

	timer := time.NewTimer(delayTime)
	defer timer.Stop()

	select {
	case <-timerCtx.Done():
		_ = level.Info(logger).Log("msg", "context cancel", "actual", time.Since(startTime), "expected", delayTime)

		atomic.StoreInt32(&restartTimes, 0)

		return
	case <-timer.C:
		_ = level.Info(logger).Log("msg", "backoff timer done", "actual", time.Since(startTime), "expected", delayTime)

		atomic.AddInt32(&restartTimes, 1)

		return
	}
}

func stop() {

	mutex.Lock()
	defer mutex.Unlock()

	if cmd == nil || cmd.Process == nil {
		_ = level.Info(logger).Log("msg", "Fluent Bit not running. No process to stop.")
		return
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = level.Info(logger).Log("msg", "Kill Fluent Bit error", "error", err)
	} else {
		_ = level.Info(logger).Log("msg", "Killed Fluent Bit")
	}
}

func resetTimer() {
	timerCancel()
	atomic.StoreInt32(&restartTimes, 0)
}

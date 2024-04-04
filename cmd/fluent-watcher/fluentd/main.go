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

	"github.com/fluent/fluent-operator/v2/pkg/filenotify"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
)

const (
	defaultBinPath      = "/usr/local/bundle/bin/fluentd"
	defaultCfgPath      = "/fluentd/etc/fluent.conf"
	defaultWatchDir     = "/fluentd/etc"
	defaultPluginPath   = "/fluentd/plugins"
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
var binPath string
var pluginPath string
var watchPath string
var poll bool
var exitOnFailure bool
var pollInterval time.Duration

func main() {

	flag.StringVar(&binPath, "b", defaultBinPath, "The fluentd binary path.")
	flag.StringVar(&configPath, "c", defaultCfgPath, "The config file path.")
	flag.StringVar(&pluginPath, "p", defaultPluginPath, "The config file path.")
	flag.BoolVar(&exitOnFailure, "exit-on-failure", false, "If fluentd exits with failure, also exit the watcher.")
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
		// Watch the Fluentd, if the Fluentd not exists or stopped, restart it.
		cancel := make(chan struct{})
		g.Add(
			func() error {

				for {
					select {
					case <-cancel:
						return nil
					default:
					}

					// Start fluentd if it does not existed.
					start()
					// Wait for the fluentd exit.
					err := wait()
					if exitOnFailure && err != nil {
						_ = level.Error(logger).Log("msg", "Fluentd exited with error; exiting watcher")
						return err
					}

					timerCtx, timerCancel = context.WithCancel(context.Background())

					// After the fluentd exit, fluentd watcher restarts it with an exponential
					// back-off delay (1s, 2s, 4s, ...), that is capped at five minutes.
					backoff()
				}
			},
			func(err error) {
				close(cancel)
				reloadOrStop()
				resetTimer()
			},
		)
	}
	{
		// Watch the config file, if the config file changed, stop Fluentd.
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

						_ = level.Info(logger).Log("msg", "Config file changed, gracefully reloading configuration")

						// After the config file changed, it should gracefully reload the fluentd,
						// and resets the restart backoff timer.
						reloadOrStop()
						resetTimer()
						_ = level.Info(logger).Log("msg", "Config file changed, gracefully reloaded configuration")
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
	return event.Op == fsnotify.Rename
}

func start() {

	mutex.Lock()
	defer mutex.Unlock()

	if cmd != nil {
		return
	}

	cmd = exec.Command(binPath, "-c", configPath, "-p", pluginPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		_ = level.Error(logger).Log("msg", "start Fluentd error", "error", err)
		cmd = nil
		return
	}

	_ = level.Info(logger).Log("msg", "Fluentd started")
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
	_ = level.Error(logger).Log("msg", "Fluentd exited", "error", err)
	// Once the fluentd has executed for 10 minutes without any problems,
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

func reloadOrStop() {
	mutex.Lock()
	defer mutex.Unlock()

	if cmd == nil || cmd.Process == nil {
		_ = level.Info(logger).Log("msg", "Fluentd not running. No process to reload or stop.")
		return
	}

	// Reloads the configuration file by gracefully re-constructing the data pipeline.
	// https://docs.fluentd.org/deployment/signals#sigusr2
	err := cmd.Process.Signal(syscall.SIGHUP)
	if err == nil {
		_ = level.Info(logger).Log("msg", "Gracefully reloaded Fluentd config")
		return
	}

	_ = level.Info(logger).Log("msg", "Gracefully reload Fluentd config error", "error", err)

	err = cmd.Process.Signal(syscall.SIGTERM)
	if err == nil {
		_ = level.Info(logger).Log("msg", "Killed Fluentd")
		return

	}

	_ = level.Info(logger).Log("msg", "Kill Fluentd error", "error", err)
}

func resetTimer() {
	timerCancel()
	atomic.StoreInt32(&restartTimes, 0)
}

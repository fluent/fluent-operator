package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"math"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

const (
	binPath      = "/fluent-bit/bin/fluent-bit"
	cfgPath      = "/fluent-bit/etc/fluent-bit.conf"
	watchDir     = "/fluent-bit/config"
	MaxDelayTime = time.Minute * 5
	ResetTime    = time.Minute * 10
)

var (
	logger       log.Logger
	cmd          *exec.Cmd
	mutex        sync.Mutex
	restartTimes int
	timer        *time.Timer
)

func main() {
	logger = log.NewLogfmtLogger(os.Stdout)

	timer = time.NewTimer(0)

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
					wait()
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
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return
		}

		// Start watcher.
		err = watcher.Add(watchDir)
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
					case event := <-watcher.Events:
						if !isValidEvent(event) {
							continue
						}

						// After the config file changed, it should stop the fluent bit,
						// and resets the restart backoff timer.
						stop()
						resetTimer()
						_ = level.Info(logger).Log("msg", "Config file changed, stop Fluent Bit")
					case <-watcher.Errors:
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

// Inspired by https://github.com/jimmidyson/configmap-reload
func isValidEvent(event fsnotify.Event) bool {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return false
	}
	//if filepath.Base(event.Name) != "..data" {
	//	return false
	//}
	return true
}

func start() {

	mutex.Lock()
	defer mutex.Unlock()

	if cmd != nil {
		return
	}

	cmd = exec.Command(binPath, "-c", cfgPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		_ = level.Error(logger).Log("msg", "start Fluent bit error", "error", err)
		cmd = nil
		return
	}

	_ = level.Info(logger).Log("msg", "Fluent bit started")
}

func wait() {

	if cmd == nil {
		return
	}

	startTime := time.Now()
	_ = level.Error(logger).Log("msg", "Fluent bit exited", "error", cmd.Wait())
	// Once the fluent bit has executed for 10 minutes without any problems,
	// it should resets the restart backoff timer.
	if time.Now().Sub(startTime) >= ResetTime {
		restartTimes = 0
	}

	mutex.Lock()
	cmd = nil
	mutex.Unlock()
}

func backoff() {

	delayTime := time.Duration(math.Pow(2, float64(restartTimes))) * time.Second
	if delayTime >= MaxDelayTime {
		delayTime = MaxDelayTime
	}
	timer.Reset(delayTime)

	startTime := time.Now()
	<-timer.C
	_ = level.Info(logger).Log("msg", "delay", "actual", time.Now().Sub(startTime), "expected", delayTime)
	restartTimes = restartTimes + 1
}

func stop() {

	mutex.Lock()
	defer mutex.Unlock()

	if cmd == nil || cmd.Process == nil {
		return
	}

	if err := cmd.Process.Kill(); err != nil {
		_ = level.Info(logger).Log("msg", "Kill Fluent Bit error", "error", err)
	} else {
		_ = level.Info(logger).Log("msg", "Killed Fluent Bit")
	}
}

func resetTimer() {

	if timer != nil {
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(0)
	}
	restartTimes = 0
}

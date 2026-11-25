package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fluent/fluent-operator/v3/pkg/filenotify"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/log"
)

const (
	testReloadDelay    = 200 * time.Millisecond
	testReloadMaxDelay = 3 * testReloadDelay
	testPollInterval   = 20 * time.Millisecond
	testConfigName     = "fluent-bit.yaml"
)

// startWatch runs watchAndReload on dir the way main does. Each reload sends the content the
// config file has at that moment, read through the path fluent-bit reads it by.
func startWatch(t *testing.T, dir string, poll bool) <-chan string {
	t.Helper()

	opts := watchOptions{
		dataLink: findDataLink(dir),
		delay:    testReloadDelay,
		maxDelay: testReloadMaxDelay,
		logger:   log.NewNopLogger(),
	}
	if opts.dataLink != "" {
		generation, err := os.Readlink(opts.dataLink)
		if err != nil {
			t.Fatal(err)
		}
		opts.generation = generation
	}

	var watcher filenotify.FileWatcher
	if poll {
		watcher = filenotify.NewPollingWatcher(testPollInterval)
		ticker := time.NewTicker(testPollInterval)
		t.Cleanup(ticker.Stop)
		opts.tick = ticker.C
	} else {
		var err error
		watcher, err = filenotify.NewEventWatcher()
		if err != nil {
			t.Fatalf("failed to create the watcher: %v", err)
		}
	}
	if err := watcher.Add(dir); err != nil {
		t.Fatalf("failed to watch %q: %v", dir, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	reloads := make(chan string, 64)
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = watchAndReload(ctx, ctx, watcher, opts, func() error {
			content, err := os.ReadFile(filepath.Join(dir, testConfigName))
			if err != nil {
				return err
			}
			reloads <- string(content)
			return nil
		})
	}()
	t.Cleanup(func() {
		cancel()
		<-done
		_ = watcher.Close()
	})

	return reloads
}

// stageVolume does the first half of a kubelet update of a Secret volume: a new timestamped
// directory with the payload. The config visible through "..data" is still the old one.
func stageVolume(t *testing.T, dir, version, content string) {
	t.Helper()

	tsDir := filepath.Join(dir, ".."+version)
	if err := os.Mkdir(tsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tsDir, testConfigName), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// publishVolume does the second half: an atomic switch of the "..data" symlink to the staged
// directory, the user-visible symlink if it is missing, removal of the old directory.
func publishVolume(t *testing.T, dir, version string) {
	t.Helper()

	dataLink := filepath.Join(dir, dataLinkName)
	oldTsDir, _ := os.Readlink(dataLink)
	tmpLink := filepath.Join(dir, "..data_tmp")
	if err := os.Symlink(".."+version, tmpLink); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmpLink, dataLink); err != nil {
		t.Fatal(err)
	}

	visible := filepath.Join(dir, testConfigName)
	if _, err := os.Lstat(visible); os.IsNotExist(err) {
		if err := os.Symlink(filepath.Join(dataLinkName, testConfigName), visible); err != nil {
			t.Fatal(err)
		}
	}
	if oldTsDir != "" {
		if err := os.RemoveAll(filepath.Join(dir, oldTsDir)); err != nil {
			t.Fatal(err)
		}
	}
}

func writeVolume(t *testing.T, dir, version, content string) {
	t.Helper()
	stageVolume(t, dir, version, content)
	publishVolume(t, dir, version)
}

func expectReload(t *testing.T, reloads <-chan string, want string) {
	t.Helper()
	select {
	case got := <-reloads:
		if got != want {
			t.Fatalf("reloaded config %q, want %q", got, want)
		}
	case <-time.After(10 * testReloadDelay):
		t.Fatalf("no reload for config %q", want)
	}
}

func expectNoReload(t *testing.T, reloads <-chan string, wait time.Duration) {
	t.Helper()
	select {
	case got := <-reloads:
		t.Fatalf("unexpected reload with config %q", got)
	case <-time.After(wait):
	}
}

func TestWatchAndReloadReadsTheSwitchedConfig(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1")
	reloads := startWatch(t, dir, false)

	writeVolume(t, dir, "v2", "v2")

	expectReload(t, reloads, "v2")
	expectNoReload(t, reloads, 3*testReloadDelay)
}

// The new version stays staged for longer than both delays. Nothing may be reloaded until it
// is published, and then there is one reload, with the new content.
func TestWatchAndReloadWaitsForASlowPublish(t *testing.T) {
	for name, poll := range map[string]bool{"events": false, "poll": true} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			writeVolume(t, dir, "v1", "v1")
			reloads := startWatch(t, dir, poll)

			stageVolume(t, dir, "v2", "v2")
			expectNoReload(t, reloads, testReloadMaxDelay+testReloadDelay)

			publishVolume(t, dir, "v2")
			expectReload(t, reloads, "v2")
			expectNoReload(t, reloads, 3*testReloadDelay)
		})
	}
}

// A version published between the start of fluent-bit and the start of the watch gives no
// event, it is found by comparing with the version fluent-bit was started with.
func TestWatchAndReloadNoticesAVersionPublishedBeforeTheWatch(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1")

	opts := watchOptions{
		dataLink:   findDataLink(dir),
		generation: "..v1",
		delay:      testReloadDelay,
		logger:     log.NewNopLogger(),
	}
	writeVolume(t, dir, "v2", "v2")

	watcher, err := filenotify.NewEventWatcher()
	if err != nil {
		t.Fatal(err)
	}
	if err := watcher.Add(dir); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	reloads := make(chan string, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = watchAndReload(ctx, ctx, watcher, opts, func() error {
			reloads <- "reloaded"
			return nil
		})
	}()
	t.Cleanup(func() {
		cancel()
		<-done
		_ = watcher.Close()
	})

	expectReload(t, reloads, "reloaded")
}

// silentWatcher never reports anything, like a poller that sees no difference in what it stats.
type silentWatcher struct{}

func (silentWatcher) Events() <-chan fsnotify.Event { return nil }
func (silentWatcher) Errors() <-chan error          { return nil }
func (silentWatcher) Add(string) error              { return nil }
func (silentWatcher) Remove(string) error           { return nil }
func (silentWatcher) Close() error                  { return nil }

// In poll mode the target of "..data" is compared on every tick, a published version is found
// without any event from the watcher.
func TestWatchAndReloadComparesTheLinkTargetOnATick(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1")

	ticker := time.NewTicker(testPollInterval)
	t.Cleanup(ticker.Stop)
	opts := watchOptions{
		dataLink:   findDataLink(dir),
		generation: "..v1",
		delay:      testReloadDelay,
		tick:       ticker.C,
		logger:     log.NewNopLogger(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	reloads := make(chan string, 16)
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = watchAndReload(ctx, ctx, silentWatcher{}, opts, func() error {
			content, err := os.ReadFile(filepath.Join(dir, testConfigName))
			if err != nil {
				return err
			}
			reloads <- string(content)
			return nil
		})
	}()
	t.Cleanup(func() {
		cancel()
		<-done
	})

	stageVolume(t, dir, "v2", "v2")
	expectNoReload(t, reloads, 2*testReloadDelay)

	publishVolume(t, dir, "v2")
	expectReload(t, reloads, "v2")
	expectNoReload(t, reloads, 3*testReloadDelay)
}

func TestWatchAndReloadCoalescesUpdatesInARow(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1")
	reloads := startWatch(t, dir, false)

	writeVolume(t, dir, "v2", "v2")
	time.Sleep(testReloadDelay / 2)
	writeVolume(t, dir, "v3", "v3")

	expectReload(t, reloads, "v3")
	expectNoReload(t, reloads, 3*testReloadDelay)
}

func TestWatchAndReloadReloadsAgainAfterALaterUpdate(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1")
	reloads := startWatch(t, dir, false)

	writeVolume(t, dir, "v2", "v2")
	expectReload(t, reloads, "v2")

	writeVolume(t, dir, "v3", "v3")
	expectReload(t, reloads, "v3")
	expectNoReload(t, reloads, 3*testReloadDelay)
}

// A directory that is not a kubelet volume, with a file next to the config that is written
// more often than the reload delay. The config change still has to be reloaded.
func TestWatchAndReloadIsNotHeldBackBySiblingWrites(t *testing.T) {
	dir := t.TempDir()
	config := filepath.Join(dir, testConfigName)
	if err := os.WriteFile(config, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	reloads := startWatch(t, dir, false)

	stop := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		ticker := time.NewTicker(testReloadDelay / 4)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				_ = os.WriteFile(filepath.Join(dir, "sibling"), []byte("noise"), 0o644)
			}
		}
	}()
	t.Cleanup(func() {
		close(stop)
		<-stopped
	})

	if err := os.WriteFile(config, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-reloads:
		if got != "v2" {
			t.Fatalf("reloaded config %q, want %q", got, "v2")
		}
	case <-time.After(testReloadMaxDelay + 2*testReloadDelay):
		t.Fatal("no reload while a sibling file keeps changing")
	}
}

func TestFindDataLink(t *testing.T) {
	volume := t.TempDir()
	writeVolume(t, volume, "v1", "v1")
	plain := t.TempDir()

	for name, tc := range map[string]struct{ watchPath, want string }{
		"volume":           {volume, filepath.Join(volume, dataLinkName)},
		"file in a volume": {filepath.Join(volume, testConfigName), filepath.Join(volume, dataLinkName)},
		"plain directory":  {plain, ""},
	} {
		if got := findDataLink(tc.watchPath); got != tc.want {
			t.Errorf("%s: findDataLink(%q) = %q, want %q", name, tc.watchPath, got, tc.want)
		}
	}
}

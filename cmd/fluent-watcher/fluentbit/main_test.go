package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fluent/fluent-operator/v3/pkg/filenotify"
)

const testReloadDelay = 200 * time.Millisecond

// startWatch runs watchAndReload on dir. Each reload sends the content the config file has at
// that moment, read through the "..data" symlink like fluent-bit reads it.
func startWatch(t *testing.T, dir string) <-chan string {
	t.Helper()

	watcher, err := filenotify.NewEventWatcher()
	if err != nil {
		t.Fatalf("failed to create the watcher: %v", err)
	}
	if err := watcher.Add(dir); err != nil {
		t.Fatalf("failed to watch %q: %v", dir, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	reloads := make(chan string, 16)
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = watchAndReload(ctx, ctx, watcher, testReloadDelay, func() error {
			content, err := os.ReadFile(filepath.Join(dir, "fluent-bit.yaml"))
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

// writeVolume updates dir the way the kubelet updates a Secret volume: a new timestamped
// directory, the payload, then an atomic switch of the "..data" symlink, then removal of the
// old directory. pause is slept between creating the directory and switching the symlink.
func writeVolume(t *testing.T, dir, version, content string, pause time.Duration) {
	t.Helper()

	tsDir := filepath.Join(dir, ".."+version)
	if err := os.Mkdir(tsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tsDir, "fluent-bit.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(pause)

	dataLink := filepath.Join(dir, "..data")
	oldTsDir, _ := os.Readlink(dataLink)
	tmpLink := filepath.Join(dir, "..data_tmp")
	if err := os.Symlink(filepath.Base(tsDir), tmpLink); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmpLink, dataLink); err != nil {
		t.Fatal(err)
	}

	visible := filepath.Join(dir, "fluent-bit.yaml")
	if _, err := os.Lstat(visible); os.IsNotExist(err) {
		if err := os.Symlink(filepath.Join("..data", "fluent-bit.yaml"), visible); err != nil {
			t.Fatal(err)
		}
	}
	if oldTsDir != "" {
		if err := os.RemoveAll(filepath.Join(dir, oldTsDir)); err != nil {
			t.Fatal(err)
		}
	}
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

func expectNoReload(t *testing.T, reloads <-chan string) {
	t.Helper()
	select {
	case got := <-reloads:
		t.Fatalf("unexpected extra reload with config %q", got)
	case <-time.After(3 * testReloadDelay):
	}
}

func TestWatchAndReloadReadsTheSwitchedConfig(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1", 0)
	reloads := startWatch(t, dir)

	// The pause keeps the new directory visible before "..data" points to it, which is the
	// window in which a reload on the first event reads the previous config.
	writeVolume(t, dir, "v2", "v2", testReloadDelay/4)

	expectReload(t, reloads, "v2")
	expectNoReload(t, reloads)
}

func TestWatchAndReloadCoalescesUpdatesInARow(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1", 0)
	reloads := startWatch(t, dir)

	writeVolume(t, dir, "v2", "v2", 0)
	time.Sleep(testReloadDelay / 2)
	writeVolume(t, dir, "v3", "v3", 0)

	expectReload(t, reloads, "v3")
	expectNoReload(t, reloads)
}

func TestWatchAndReloadReloadsAgainAfterALaterUpdate(t *testing.T) {
	dir := t.TempDir()
	writeVolume(t, dir, "v1", "v1", 0)
	reloads := startWatch(t, dir)

	writeVolume(t, dir, "v2", "v2", 0)
	expectReload(t, reloads, "v2")

	writeVolume(t, dir, "v3", "v3", 0)
	expectReload(t, reloads, "v3")
	expectNoReload(t, reloads)
}

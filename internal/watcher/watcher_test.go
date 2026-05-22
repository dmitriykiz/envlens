package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envlens/internal/watcher"
)

const pollInterval = 20 * time.Millisecond
const waitTimeout = 500 * time.Millisecond

func waitForEvent(ch <-chan watcher.Event, op string) (watcher.Event, bool) {
	deadline := time.After(waitTimeout)
	for {
		select {
		case ev := <-ch:
			if ev.Op == op {
				return ev, true
			}
		case <-deadline:
			return watcher.Event{}, false
		}
	}
}

func TestWatcher_DetectsModified(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("KEY=value\n"), 0644); err != nil {
		t.Fatal(err)
	}

	w := watcher.New([]string{path}, pollInterval)
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ev, ok := waitForEvent(w.Events, "modified")
	if !ok {
		t.Fatal("expected modified event, got none")
	}
	if ev.Path != path {
		t.Errorf("expected path %s, got %s", path, ev.Path)
	}
}

func TestWatcher_DetectsCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.new")

	w := watcher.New([]string{path}, pollInterval)
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("X=1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, ok := waitForEvent(w.Events, "created")
	if !ok {
		t.Fatal("expected created event, got none")
	}
}

func TestWatcher_DetectsDeleted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("A=1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	w := watcher.New([]string{path}, pollInterval)
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	_, ok := waitForEvent(w.Events, "deleted")
	if !ok {
		t.Fatal("expected deleted event, got none")
	}
}

func TestWatcher_StopHaltsPolling(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("K=v\n"), 0644)

	w := watcher.New([]string{path}, pollInterval)
	w.Start()
	w.Stop()

	time.Sleep(60 * time.Millisecond)
	_ = os.WriteFile(path, []byte("K=changed\n"), 0644)
	time.Sleep(60 * time.Millisecond)

	select {
	case ev := <-w.Events:
		if ev.Op == "modified" {
			t.Error("received event after Stop()")
		}
	default:
		// expected: no events after stop
	}
}

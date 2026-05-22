package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Event represents a file system change detected on a .env file.
type Event struct {
	Path    string
	Op      string // "modified" | "created" | "deleted"
	OccurredAt time.Time
}

// Watcher polls a set of .env file paths for changes.
type Watcher struct {
	paths    []string
	interval time.Duration
	states   map[string]fileState
	Events   chan Event
	Errors   chan error
	stop     chan struct{}
}

type fileState struct {
	modTime time.Time
	exists  bool
}

// New creates a Watcher that polls the given paths every interval.
func New(paths []string, interval time.Duration) *Watcher {
	return &Watcher{
		paths:    paths,
		interval: interval,
		states:   make(map[string]fileState),
		Events:   make(chan Event, 32),
		Errors:   make(chan error, 8),
		stop:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	w.snapshot() // capture initial state
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the background polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) snapshot() {
	for _, p := range w.paths {
		info, err := os.Stat(p)
		if err != nil {
			w.states[p] = fileState{exists: false}
			continue
		}
		w.states[p] = fileState{modTime: info.ModTime(), exists: true}
	}
}

func (w *Watcher) poll() {
	for _, p := range w.paths {
		clean := filepath.Clean(p)
		info, err := os.Stat(clean)
		prev, known := w.states[p]

		if err != nil {
			if known && prev.exists {
				w.emit(p, "deleted")
				w.states[p] = fileState{exists: false}
			} else if !known {
				w.errors(fmt.Errorf("watcher: cannot stat %s: %w", p, err))
			}
			continue
		}

		if !known || !prev.exists {
			w.emit(p, "created")
		} else if info.ModTime().After(prev.modTime) {
			w.emit(p, "modified")
		}
		w.states[p] = fileState{modTime: info.ModTime(), exists: true}
	}
}

func (w *Watcher) emit(path, op string) {
	select {
	case w.Events <- Event{Path: path, Op: op, OccurredAt: time.Now()}:
	default:
	}
}

func (w *Watcher) errors(err error) {
	select {
	case w.Errors <- err:
	default:
	}
}

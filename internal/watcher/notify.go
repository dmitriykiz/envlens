package watcher

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Notifier formats and writes watcher Events to an output stream.
type Notifier struct {
	out    io.Writer
	prefix string
}

// NewNotifier creates a Notifier that writes formatted event lines to out.
// prefix is prepended to every line (e.g. "[envlens]").
func NewNotifier(out io.Writer, prefix string) *Notifier {
	return &Notifier{out: out, prefix: prefix}
}

// Drain reads from the Events channel until it is closed or the done channel
// is signalled, writing a formatted line for each event.
func (n *Notifier) Drain(events <-chan Event, done <-chan struct{}) {
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			n.write(ev)
		case <-done:
			return
		}
	}
}

// Format returns a human-readable string for a single Event.
func Format(ev Event) string {
	return fmt.Sprintf("%s %-8s %s",
		ev.OccurredAt.Format(time.RFC3339),
		strings.ToUpper(ev.Op),
		ev.Path,
	)
}

func (n *Notifier) write(ev Event) {
	line := Format(ev)
	if n.prefix != "" {
		line = n.prefix + " " + line
	}
	fmt.Fprintln(n.out, line)
}

// DrainErrors reads from the Errors channel until done, writing each error.
func (n *Notifier) DrainErrors(errors <-chan error, done <-chan struct{}) {
	for {
		select {
		case err, ok := <-errors:
			if !ok {
				return
			}
			if n.prefix != "" {
				fmt.Fprintf(n.out, "%s ERROR %v\n", n.prefix, err)
			} else {
				fmt.Fprintf(n.out, "ERROR %v\n", err)
			}
		case <-done:
			return
		}
	}
}

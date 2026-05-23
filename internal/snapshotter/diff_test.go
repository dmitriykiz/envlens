package snapshotter_test

import (
	"testing"
	"time"

	"github.com/envlens/internal/parser"
	"github.com/envlens/internal/snapshotter"
)

func makeSnap(label string, files map[string][]parser.Entry) *snapshotter.Snapshot {
	return &snapshotter.Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Files:     files,
	}
}

func TestDiff_NoDifferences(t *testing.T) {
	entries := []parser.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	base := makeSnap("base", map[string][]parser.Entry{"app.env": entries})
	curr := makeSnap("curr", map[string][]parser.Entry{"app.env": entries})
	delta := snapshotter.Diff(base, curr)
	if len(delta.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(delta.Changes))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	base := makeSnap("b", map[string][]parser.Entry{"a.env": {{Key: "X", Value: "1"}}})
	curr := makeSnap("c", map[string][]parser.Entry{"a.env": {{Key: "X", Value: "1"}, {Key: "Y", Value: "2"}}})
	delta := snapshotter.Diff(base, curr)
	if !hasChange(delta, "added", "Y") {
		t.Error("expected 'added' change for key Y")
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	base := makeSnap("b", map[string][]parser.Entry{"a.env": {{Key: "X", Value: "1"}, {Key: "Z", Value: "3"}}})
	curr := makeSnap("c", map[string][]parser.Entry{"a.env": {{Key: "X", Value: "1"}}})
	delta := snapshotter.Diff(base, curr)
	if !hasChange(delta, "removed", "Z") {
		t.Error("expected 'removed' change for key Z")
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	base := makeSnap("b", map[string][]parser.Entry{"a.env": {{Key: "PORT", Value: "8080"}}})
	curr := makeSnap("c", map[string][]parser.Entry{"a.env": {{Key: "PORT", Value: "9090"}}})
	delta := snapshotter.Diff(base, curr)
	if !hasChange(delta, "changed", "PORT") {
		t.Error("expected 'changed' change for key PORT")
	}
}

func TestDiff_Labels(t *testing.T) {
	base := makeSnap("before", map[string][]parser.Entry{})
	curr := makeSnap("after", map[string][]parser.Entry{})
	delta := snapshotter.Diff(base, curr)
	if delta.BaseLabel != "before" || delta.CurrentLabel != "after" {
		t.Errorf("labels wrong: %q / %q", delta.BaseLabel, delta.CurrentLabel)
	}
}

func hasChange(d *snapshotter.Delta, changeType, key string) bool {
	for _, c := range d.Changes {
		if c.Change == changeType && c.Key == key {
			return true
		}
	}
	return false
}

package snapshotter

import (
	"github.com/envlens/internal/parser"
)

// KeyDelta describes a key that changed between two snapshots.
type KeyDelta struct {
	File   string
	Key    string
	OldVal string
	NewVal string
	Change string // "added", "removed", "changed"
}

// Delta holds all differences between two snapshots.
type Delta struct {
	BaseLabel    string
	CurrentLabel string
	Changes      []KeyDelta
}

// Diff computes the difference between a base snapshot and a current snapshot.
func Diff(base, current *Snapshot) *Delta {
	delta := &Delta{
		BaseLabel:    base.Label,
		CurrentLabel: current.Label,
	}

	allFiles := union(keys(base.Files), keys(current.Files))
	for _, file := range allFiles {
		baseIdx := index(base.Files[file])
		currIdx := index(current.Files[file])

		for k, bv := range baseIdx {
			cv, exists := currIdx[k]
			if !exists {
				delta.Changes = append(delta.Changes, KeyDelta{File: file, Key: k, OldVal: bv, Change: "removed"})
			} else if cv != bv {
				delta.Changes = append(delta.Changes, KeyDelta{File: file, Key: k, OldVal: bv, NewVal: cv, Change: "changed"})
			}
		}
		for k, cv := range currIdx {
			if _, exists := baseIdx[k]; !exists {
				delta.Changes = append(delta.Changes, KeyDelta{File: file, Key: k, NewVal: cv, Change: "added"})
			}
		}
	}
	return delta
}

func index(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if !e.Comment {
			m[e.Key] = e.Value
		}
	}
	return m
}

func keys(m map[string][]parser.Entry) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func union(a, b []string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for _, v := range append(a, b...) {
		seen[v] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

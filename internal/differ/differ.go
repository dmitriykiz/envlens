// Package differ computes the diff between two sets of env entries,
// identifying keys that were added, removed, or changed between versions.
package differ

import "github.com/yourorg/envlens/internal/parser"

// DiffResult holds the outcome of comparing two env entry slices.
type DiffResult struct {
	Added   []parser.Entry // keys present in Next but not in Base
	Removed []parser.Entry // keys present in Base but not in Next
	Changed []Change       // keys present in both but with different values
}

// Change describes a key whose value differs between Base and Next.
type Change struct {
	Key      string
	OldValue string
	NewValue string
}

// Diff compares base entries against next entries and returns a DiffResult
// describing what was added, removed, or changed.
func Diff(base, next []parser.Entry) DiffResult {
	baseMap := indexEntries(base)
	nextMap := indexEntries(next)

	var result DiffResult

	// Detect removed and changed keys.
	for key, baseEntry := range baseMap {
		if nextEntry, ok := nextMap[key]; ok {
			if nextEntry.Value != baseEntry.Value {
				result.Changed = append(result.Changed, Change{
					Key:      key,
					OldValue: baseEntry.Value,
					NewValue: nextEntry.Value,
				})
			}
		} else {
			result.Removed = append(result.Removed, baseEntry)
		}
	}

	// Detect added keys.
	for key, nextEntry := range nextMap {
		if _, ok := baseMap[key]; !ok {
			result.Added = append(result.Added, nextEntry)
		}
	}

	return result
}

// indexEntries builds a map from key to Entry for fast lookup.
// When duplicate keys exist, the last occurrence wins.
func indexEntries(entries []parser.Entry) map[string]parser.Entry {
	m := make(map[string]parser.Entry, len(entries))
	for _, e := range entries {
		m[e.Key] = e
	}
	return m
}

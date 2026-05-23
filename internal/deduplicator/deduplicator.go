// Package deduplicator removes duplicate entries from a set of parsed env
// entries, keeping a configurable occurrence (first or last) and reporting
// every key that was deduplicated.
package deduplicator

import "github.com/yourorg/envlens/internal/parser"

// Strategy controls which entry survives when duplicates are found.
type Strategy int

const (
	// KeepFirst retains the first occurrence of a duplicate key.
	KeepFirst Strategy = iota
	// KeepLast retains the last occurrence of a duplicate key.
	KeepLast
)

// Result holds the deduplicated entries together with a report of every key
// that had duplicates removed.
type Result struct {
	Entries    []parser.Entry
	Duplicates map[string]int // key → number of extra copies removed
}

// Deduplicate scans entries for keys that appear more than once and returns a
// Result containing only unique keys. The Strategy argument decides whether the
// first or last value wins.
func Deduplicate(entries []parser.Entry, strategy Strategy) Result {
	seen := make(map[string]int)   // key → index in out slice
	out := make([]parser.Entry, 0, len(entries))
	removed := make(map[string]int)

	for _, e := range entries {
		if e.Comment {
			out = append(out, e)
			continue
		}

		idx, exists := seen[e.Key]
		if !exists {
			seen[e.Key] = len(out)
			out = append(out, e)
			continue
		}

		// Duplicate found.
		removed[e.Key]++
		if strategy == KeepLast {
			out[idx] = e
		}
	}

	return Result{
		Entries:    out,
		Duplicates: removed,
	}
}

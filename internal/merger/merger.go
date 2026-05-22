package merger

import (
	"fmt"
	"sort"

	"github.com/envlens/internal/parser"
)

// ConflictStrategy determines how key conflicts are resolved during a merge.
type ConflictStrategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst ConflictStrategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error when a conflict is detected.
	StrategyError
)

// Conflict records a key that appeared in more than one source file.
type Conflict struct {
	Key    string
	Files  []string
	Values []string
}

// Result holds the merged entries and any conflicts that were detected.
type Result struct {
	Entries   []parser.Entry
	Conflicts []Conflict
}

// Merge combines entries from multiple files into a single ordered slice.
// The strategy parameter controls what happens when the same key appears in
// more than one file.
func Merge(fileEntries map[string][]parser.Entry, strategy ConflictStrategy) (*Result, error) {
	seen := make(map[string]*Conflict)
	order := []string{}
	merged := make(map[string]parser.Entry)

	// Iterate files in deterministic order.
	files := make([]string, 0, len(fileEntries))
	for f := range fileEntries {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		for _, entry := range fileEntries[file] {
			if existing, exists := merged[entry.Key]; exists {
				if strategy == StrategyError {
					return nil, fmt.Errorf("conflict: key %q defined in multiple files", entry.Key)
				}
				c := seen[entry.Key]
				c.Files = append(c.Files, file)
				c.Values = append(c.Values, entry.Value)
				if strategy == StrategyLast {
					merged[entry.Key] = entry
				}
				_ = existing
			} else {
				merged[entry.Key] = entry
				order = append(order, entry.Key)
				seen[entry.Key] = &Conflict{
					Key:    entry.Key,
					Files:  []string{file},
					Values: []string{entry.Value},
				}
			}
		}
	}

	result := &Result{}
	for _, key := range order {
		result.Entries = append(result.Entries, merged[key])
		c := seen[key]
		if len(c.Files) > 1 {
			result.Conflicts = append(result.Conflicts, *c)
		}
	}
	return result, nil
}

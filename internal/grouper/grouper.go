// Package grouper organises parsed env entries into named groups
// based on key-prefix conventions (e.g. DB_, AWS_, APP_).
package grouper

import (
	"sort"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Group holds all entries that share a common prefix.
type Group struct {
	Prefix  string
	Entries []parser.Entry
}

// Options controls how grouping is performed.
type Options struct {
	// MinGroupSize discards groups with fewer entries than this value.
	MinGroupSize int
	// Separator is the character that delimits a prefix from the rest of the
	// key. Defaults to "_".
	Separator string
	// UngroupedLabel is the name given to entries that have no prefix.
	UngroupedLabel string
}

func defaultOptions(o Options) Options {
	if o.Separator == "" {
		o.Separator = "_"
	}
	if o.UngroupedLabel == "" {
		o.UngroupedLabel = "(ungrouped)"
	}
	if o.MinGroupSize < 1 {
		o.MinGroupSize = 1
	}
	return o
}

// ByPrefix groups entries by the first segment of their key when split on
// opts.Separator. Entries whose key contains no separator are placed in the
// ungrouped bucket.
func ByPrefix(entries []parser.Entry, opts Options) []Group {
	opts = defaultOptions(opts)

	buckets := map[string][]parser.Entry{}
	for _, e := range entries {
		if e.Comment {
			continue
		}
		prefix := prefixOf(e.Key, opts.Separator)
		if prefix == "" {
			prefix = opts.UngroupedLabel
		}
		buckets[prefix] = append(buckets[prefix], e)
	}

	groups := make([]Group, 0, len(buckets))
	for prefix, ents := range buckets {
		if len(ents) < opts.MinGroupSize {
			continue
		}
		groups = append(groups, Group{Prefix: prefix, Entries: ents})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})
	return groups
}

// Flatten is the inverse of ByPrefix: it returns all entries from all groups
// in prefix-sorted, then key-sorted order.
func Flatten(groups []Group) []parser.Entry {
	var out []parser.Entry
	for _, g := range groups {
		out = append(out, g.Entries...)
	}
	return out
}

func prefixOf(key, sep string) string {
	idx := strings.Index(key, sep)
	if idx <= 0 {
		return ""
	}
	return key[:idx]
}

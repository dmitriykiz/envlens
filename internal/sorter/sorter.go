// Package sorter provides utilities for sorting and grouping .env file entries
// by key name, value, or category (sensitive vs. non-sensitive).
package sorter

import (
	"sort"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Order defines the sort direction.
type Order int

const (
	Ascending  Order = iota
	Descending
)

// By defines the field to sort on.
type By int

const (
	ByKey By = iota
	ByValue
)

// Options controls sorting behaviour.
type Options struct {
	By    By
	Order Order
	// GroupSensitive moves sensitive keys to the top when true.
	GroupSensitive bool
}

var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

// Sort returns a new slice of entries sorted according to opts.
// The original slice is not modified.
func Sort(entries []parser.Entry, opts Options) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	sort.SliceStable(out, func(i, j int) bool {
		if opts.GroupSensitive {
			si := isSensitive(out[i].Key)
			sj := isSensitive(out[j].Key)
			if si != sj {
				return si // sensitive first
			}
		}

		var a, b string
		switch opts.By {
		case ByValue:
			a, b = out[i].Value, out[j].Value
		default:
			a, b = out[i].Key, out[j].Key
		}

		if opts.Order == Descending {
			return a > b
		}
		return a < b
	})

	return out
}

// GroupByPrefix splits entries into groups keyed by the first segment of the
// key name when split on sep (e.g. "_"). Keys without sep land in "".
func GroupByPrefix(entries []parser.Entry, sep string) map[string][]parser.Entry {
	groups := make(map[string][]parser.Entry)
	for _, e := range entries {
		prefix := ""
		if idx := strings.Index(e.Key, sep); idx > 0 {
			prefix = e.Key[:idx]
		}
		groups[prefix] = append(groups[prefix], e)
	}
	return groups
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

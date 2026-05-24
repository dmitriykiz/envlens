// Package normalizer standardises .env entries by applying a consistent
// set of key and value transformations such as uppercasing keys, trimming
// whitespace, and collapsing redundant quoting.
package normalizer

import (
	"strings"

	"github.com/your-org/envlens/internal/parser"
)

// Options controls which normalisation passes are applied.
type Options struct {
	// UppercaseKeys converts every key to UPPER_CASE.
	UppercaseKeys bool
	// TrimSpace removes leading/trailing whitespace from keys and values.
	TrimSpace bool
	// StripInlineComments removes trailing inline comments (e.g. VALUE # note).
	StripInlineComments bool
	// UnquoteValues removes surrounding single or double quotes from values.
	UnquoteValues bool
}

// DefaultOptions returns an Options struct with all passes enabled.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:       true,
		TrimSpace:           true,
		StripInlineComments: true,
		UnquoteValues:       true,
	}
}

// Result holds a single normalised entry alongside metadata about what changed.
type Result struct {
	Original   parser.Entry
	Normalized parser.Entry
	Changed    bool
	Changes    []string
}

// Normalize applies the configured passes to each entry and returns a slice of
// Result values. Comment and blank entries are passed through unchanged.
func Normalize(entries []parser.Entry, opts Options) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		if e.IsComment || e.IsBlank {
			results = append(results, Result{Original: e, Normalized: e})
			continue
		}
		results = append(results, normalizeEntry(e, opts))
	}
	return results
}

func normalizeEntry(e parser.Entry, opts Options) Result {
	n := e
	var changes []string

	if opts.TrimSpace {
		if k := strings.TrimSpace(n.Key); k != n.Key {
			n.Key = k
			changes = append(changes, "trimmed key whitespace")
		}
		if v := strings.TrimSpace(n.Value); v != n.Value {
			n.Value = v
			changes = append(changes, "trimmed value whitespace")
		}
	}

	if opts.UppercaseKeys {
		if up := strings.ToUpper(n.Key); up != n.Key {
			n.Key = up
			changes = append(changes, "uppercased key")
		}
	}

	if opts.StripInlineComments {
		if idx := strings.Index(n.Value, " #"); idx != -1 {
			n.Value = strings.TrimSpace(n.Value[:idx])
			changes = append(changes, "stripped inline comment")
		}
	}

	if opts.UnquoteValues {
		if v, ok := stripQuotes(n.Value); ok {
			n.Value = v
			changes = append(changes, "removed surrounding quotes")
		}
	}

	return Result{
		Original:   e,
		Normalized: n,
		Changed:    len(changes) > 0,
		Changes:    changes,
	}
}

func stripQuotes(s string) (string, bool) {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1], true
		}
	}
	return s, false
}

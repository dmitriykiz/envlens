package trimmer

import (
	"strings"

	"github.com/envlens/internal/parser"
)

// Options controls trimming behaviour.
type Options struct {
	// TrimKeys removes leading/trailing whitespace from key names.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// TrimQuotes strips surrounding single or double quotes from values.
	TrimQuotes bool
}

// DefaultOptions returns an Options with all trimming enabled.
func DefaultOptions() Options {
	return Options{
		TrimKeys:   true,
		TrimValues: true,
		TrimQuotes: true,
	}
}

// Result holds the trimmed entries and a count of changes made.
type Result struct {
	Entries  []parser.Entry
	Changed  int
}

// Trim applies the given Options to each entry and returns a Result.
// Comment and blank entries are passed through unchanged.
func Trim(entries []parser.Entry, opts Options) Result {
	out := make([]parser.Entry, 0, len(entries))
	changed := 0

	for _, e := range entries {
		if e.Comment || e.Key == "" {
			out = append(out, e)
			continue
		}

		origKey := e.Key
		origVal := e.Value

		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.TrimQuotes {
			e.Value = trimQuotes(e.Value)
		}

		if e.Key != origKey || e.Value != origVal {
			changed++
		}

		out = append(out, e)
	}

	return Result{Entries: out, Changed: changed}
}

// trimQuotes removes a matching pair of surrounding quotes (single or double).
func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

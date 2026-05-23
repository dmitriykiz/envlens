// Package filter provides utilities for selecting a subset of env entries
// based on key patterns, value patterns, or sensitivity classification.
package filter

import (
	"regexp"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Options controls which entries are retained by Filter.
type Options struct {
	// KeyPattern, if non-empty, keeps only entries whose key matches the regexp.
	KeyPattern string
	// ValuePattern, if non-empty, keeps only entries whose value matches the regexp.
	ValuePattern string
	// SensitiveOnly retains only entries whose key is considered sensitive.
	SensitiveOnly bool
	// ExcludeComments drops comment/blank entries when true.
	ExcludeComments bool
}

var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|secret|token|api_?key|private_?key|auth)`),
}

// Filter returns the subset of entries that satisfy all active options.
func Filter(entries []parser.Entry, opts Options) ([]parser.Entry, error) {
	var keyRe, valRe *regexp.Regexp
	var err error

	if opts.KeyPattern != "" {
		if keyRe, err = regexp.Compile(opts.KeyPattern); err != nil {
			return nil, err
		}
	}
	if opts.ValuePattern != "" {
		if valRe, err = regexp.Compile(opts.ValuePattern); err != nil {
			return nil, err
		}
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if opts.ExcludeComments && (e.IsComment || strings.TrimSpace(e.Raw) == "") {
			continue
		}
		if keyRe != nil && !keyRe.MatchString(e.Key) {
			continue
		}
		if valRe != nil && !valRe.MatchString(e.Value) {
			continue
		}
		if opts.SensitiveOnly && !isSensitive(e.Key) {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}

func isSensitive(key string) bool {
	for _, re := range sensitivePatterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

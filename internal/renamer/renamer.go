// Package renamer provides utilities for bulk-renaming keys across a set of
// parsed .env entries. It supports exact-match and pattern-based renaming,
// and returns a detailed record of every substitution made.
package renamer

import (
	"fmt"
	"regexp"

	"github.com/yourorg/envlens/internal/parser"
)

// Rule describes a single rename operation.
type Rule struct {
	// From is a literal key name or a regular expression (when IsRegexp is true).
	From string
	// To is the replacement key name. When IsRegexp is true it may contain
	// regexp expand tokens such as $1.
	To string
	// IsRegexp controls whether From is compiled as a regular expression.
	IsRegexp bool
}

// Change records one key that was renamed.
type Change struct {
	File    string
	OldKey  string
	NewKey  string
	Line    int
}

// Result is returned by Rename and bundles the rewritten entries together
// with the list of changes that were applied.
type Result struct {
	Entries []parser.Entry
	Changes []Change
}

// Rename applies rules to entries from a single file and returns a Result.
// Rules are evaluated in order; the first matching rule wins.
func Rename(file string, entries []parser.Entry, rules []Rule) (Result, error) {
	type compiled struct {
		rule Rule
		re   *regexp.Regexp
	}

	compiled_rules := make([]compiled, 0, len(rules))
	for _, r := range rules {
		var re *regexp.Regexp
		if r.IsRegexp {
			var err error
			re, err = regexp.Compile(r.From)
			if err != nil {
				return Result{}, fmt.Errorf("renamer: invalid regexp %q: %w", r.From, err)
			}
		}
		compiled_rules = append(compiled_rules, compiled{rule: r, re: re})
	}

	out := make([]parser.Entry, len(entries))
	var changes []Change

	for i, e := range entries {
		out[i] = e
		if e.IsComment || e.Key == "" {
			continue
		}
		for _, cr := range compiled_rules {
			var newKey string
			if cr.rule.IsRegexp {
				if !cr.re.MatchString(e.Key) {
					continue
				}
				newKey = cr.re.ReplaceAllString(e.Key, cr.rule.To)
			} else {
				if e.Key != cr.rule.From {
					continue
				}
				newKey = cr.rule.To
			}
			if newKey == e.Key {
				break
			}
			changes = append(changes, Change{File: file, OldKey: e.Key, NewKey: newKey, Line: e.Line})
			out[i].Key = newKey
			break
		}
	}

	return Result{Entries: out, Changes: changes}, nil
}

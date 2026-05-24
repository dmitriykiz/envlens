package sanitizer

import (
	"regexp"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Rule describes a single sanitization transformation.
type Rule struct {
	// KeyPattern, if non-empty, restricts the rule to matching keys.
	KeyPattern string
	// StripChars removes any of the listed characters from the value.
	StripChars string
	// ReplacePattern is a regexp applied to the value.
	ReplacePattern string
	// Replacement is the substitution string for ReplacePattern matches.
	Replacement string
	// TrimSpace trims leading/trailing whitespace from the value.
	TrimSpace bool
}

// Result records what was changed for a single entry.
type Result struct {
	File    string
	Line    int
	Key     string
	OldVal  string
	NewVal  string
	Changed bool
}

// Sanitize applies rules to every entry and returns the modified entries
// together with a per-entry result log.
func Sanitize(entries []parser.Entry, rules []Rule) ([]parser.Entry, []Result) {
	compiledKeys := compileKeyPatterns(rules)
	compiledVals := compileValuePatterns(rules)

	out := make([]parser.Entry, len(entries))
	results := make([]Result, 0, len(entries))

	for i, e := range entries {
		out[i] = e
		if e.Comment || e.Key == "" {
			continue
		}
		newVal := e.Value
		for j, rule := range rules {
			if compiledKeys[j] != nil && !compiledKeys[j].MatchString(e.Key) {
				continue
			}
			if rule.TrimSpace {
				newVal = strings.TrimSpace(newVal)
			}
			if rule.StripChars != "" {
				newVal = stripChars(newVal, rule.StripChars)
			}
			if compiledVals[j] != nil {
				newVal = compiledVals[j].ReplaceAllString(newVal, rule.Replacement)
			}
		}
		results = append(results, Result{
			File:    e.File,
			Line:    e.Line,
			Key:     e.Key,
			OldVal:  e.Value,
			NewVal:  newVal,
			Changed: newVal != e.Value,
		})
		out[i].Value = newVal
	}
	return out, results
}

func compileKeyPatterns(rules []Rule) []*regexp.Regexp {
	out := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		if r.KeyPattern != "" {
			out[i] = regexp.MustCompile(r.KeyPattern)
		}
	}
	return out
}

func compileValuePatterns(rules []Rule) []*regexp.Regexp {
	out := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		if r.ReplacePattern != "" {
			out[i] = regexp.MustCompile(r.ReplacePattern)
		}
	}
	return out
}

func stripChars(s, chars string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(chars, r) {
			return -1
		}
		return r
	}, s)
}

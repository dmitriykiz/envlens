// Package validator checks .env entries against a set of rules defined
// in a schema file (e.g. .env.schema), ensuring required keys are present
// and values conform to expected formats.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Rule describes the validation constraints for a single key.
type Rule struct {
	Key      string
	Required bool
	// Pattern is an optional regex the value must fully match.
	Pattern string
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Message string
}

// Validate checks a slice of parsed entries against the provided rules.
// It returns a (possibly empty) slice of Violations.
func Validate(entries []parser.Entry, rules []Rule) []Violation {
	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	var violations []Violation

	for _, rule := range rules {
		value, exists := index[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != "" {
			matched, err := regexp.MatchString("^"+rule.Pattern+"$", value)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !matched {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", truncate(value, 20), rule.Pattern),
				})
			}
		}
	}

	return violations
}

// truncate shortens s to at most n characters for display purposes.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return strings.TrimRight(s[:n], " ") + "..."
}

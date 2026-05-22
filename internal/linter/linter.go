package linter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envlens/internal/parser"
)

// Rule represents a linting rule applied to env entries.
type Rule struct {
	Name    string
	Message string
	Check   func(e parser.Entry) bool
}

// Violation describes a single rule violation found in a file.
type Violation struct {
	File    string
	Line    int
	Key     string
	Rule    string
	Message string
}

var defaultRules = []Rule{
	{
		Name:    "no-empty-value",
		Message: "key has an empty value",
		Check:   func(e parser.Entry) bool { return strings.TrimSpace(e.Value) == "" },
	},
	{
		Name:    "no-whitespace-in-key",
		Message: "key contains whitespace",
		Check:   func(e parser.Entry) bool { return strings.ContainsAny(e.Key, " \t") },
	},
	{
		Name:    "uppercase-key",
		Message: "key is not fully uppercase",
		Check:   func(e parser.Entry) bool { return e.Key != strings.ToUpper(e.Key) },
	},
	{
		Name:    "no-special-chars-in-key",
		Message: "key contains characters other than alphanumerics and underscores",
		Check:   func(e parser.Entry) bool { return !regexp.MustCompile(`^[A-Z0-9_]+$`).MatchString(e.Key) },
	},
}

// Lint runs all default linting rules against the entries parsed from the
// given file path and returns any violations found.
func Lint(filePath string, entries []parser.Entry) ([]Violation, error) {
	if filePath == "" {
		return nil, fmt.Errorf("linter: file path must not be empty")
	}
	var violations []Violation
	for _, entry := range entries {
		for _, rule := range defaultRules {
			if rule.Check(entry) {
				violations = append(violations, Violation{
					File:    filePath,
					Line:    entry.Line,
					Key:     entry.Key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return violations, nil
}

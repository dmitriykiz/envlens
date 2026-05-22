// Package linter provides rule-based linting for .env file entries.
//
// It applies a set of built-in rules to a slice of parsed entries and
// returns a list of Violation values describing any problems found.
//
// Built-in rules:
//
//   - no-empty-value: flags keys whose value is empty or blank.
//   - no-whitespace-in-key: flags keys that contain spaces or tabs.
//   - uppercase-key: flags keys that are not fully uppercase.
//   - no-special-chars-in-key: flags keys containing characters other
//     than uppercase letters, digits, and underscores.
//
// Example usage:
//
//	entries, err := parser.ParseFile(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	violations, err := linter.Lint(".env", entries)
//	for _, v := range violations {
//		fmt.Printf("%s:%d [%s] %s — %s\n", v.File, v.Line, v.Rule, v.Key, v.Message)
//	}
package linter

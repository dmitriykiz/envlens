// Package sanitizer provides value-level sanitization for .env entries.
//
// It applies a sequence of configurable [Rule] values to each entry's value,
// supporting operations such as:
//
//   - Trimming surrounding whitespace
//   - Stripping individual characters (e.g. control characters)
//   - Regexp-based search-and-replace
//
// Rules can be scoped to a subset of keys via a key-matching regexp, allowing
// fine-grained control over which entries are affected.
//
// Usage:
//
//	rules := []sanitizer.Rule{
//		{TrimSpace: true},
//		{KeyPattern: "URL", ReplacePattern: `\s+`, Replacement: ""},
//	}
//	cleaned, results := sanitizer.Sanitize(entries, rules)
//	sanitizer.WriteReport(os.Stdout, results)
//
// The original entries slice is never mutated; Sanitize returns a new slice.
package sanitizer

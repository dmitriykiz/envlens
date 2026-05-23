package auditor

import (
	"regexp"
	"strings"
)

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// isValidKeyFormat returns true when key matches UPPER_SNAKE_CASE conventions.
func isValidKeyFormat(key string) bool {
	return validKeyRe.MatchString(key)
}

// isSensitiveKey reports whether key contains any of the given patterns
// (case-insensitive substring match).
func isSensitiveKey(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// CountBySeverity returns a map from Severity to the number of findings at that level.
func CountBySeverity(findings []Finding) map[Severity]int {
	counts := make(map[Severity]int)
	for _, f := range findings {
		counts[f.Severity]++
	}
	return counts
}

// FilterBySeverity returns only the findings that match the requested severity.
func FilterBySeverity(findings []Finding, sev Severity) []Finding {
	var out []Finding
	for _, f := range findings {
		if f.Severity == sev {
			out = append(out, f)
		}
	}
	return out
}

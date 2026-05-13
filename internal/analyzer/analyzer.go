// Package analyzer provides functionality for analyzing parsed .env files
// to detect missing keys, duplicates, and potentially sensitive values.
package analyzer

import (
	"strings"

	"github.com/user/envlens/internal/parser"
)

// sensitivePatterns holds substrings that suggest a key may be sensitive.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL",
	"AUTH", "ACCESS_KEY", "SIGNING",
}

// Result holds the analysis outcome for a single .env file.
type Result struct {
	FilePath   string
	Duplicates []string
	Sensitive  []string
	Missing    []string // keys present in reference but absent here
}

// Analyze inspects a slice of entries from a single file and returns a Result.
// referenceKeys, if non-nil, is used to compute missing keys.
func Analyze(filePath string, entries []parser.Entry, referenceKeys []string) Result {
	result := Result{FilePath: filePath}

	seen := make(map[string]int)
	for _, e := range entries {
		seen[e.Key]++
	}

	for key, count := range seen {
		if count > 1 {
			result.Duplicates = append(result.Duplicates, key)
		}
		if isSensitive(key) {
			result.Sensitive = append(result.Sensitive, key)
		}
	}

	for _, ref := range referenceKeys {
		if seen[ref] == 0 {
			result.Missing = append(result.Missing, ref)
		}
	}

	return result
}

// isSensitive returns true if the key contains any known sensitive pattern.
func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

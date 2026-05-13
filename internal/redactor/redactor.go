package redactor

import (
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// MaskChar is the character used to mask sensitive values.
const MaskChar = "*"

// MaskedEntry represents a parsed env entry with potentially masked value.
type MaskedEntry struct {
	Key      string
	Value    string
	Masked   bool
	Line     int
}

// sensitivePatterns holds substrings that indicate a key is sensitive.
var sensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// Redact takes a slice of parsed entries and returns MaskedEntry slice,
// masking the values of any keys that match sensitive patterns.
func Redact(entries []parser.Entry) []MaskedEntry {
	result := make([]MaskedEntry, 0, len(entries))
	for _, e := range entries {
		me := MaskedEntry{
			Key:  e.Key,
			Line: e.Line,
		}
		if isSensitiveKey(e.Key) {
			me.Value = maskValue(e.Value)
			me.Masked = true
		} else {
			me.Value = e.Value
			me.Masked = false
		}
		result = append(result, me)
	}
	return result
}

// isSensitiveKey returns true if the key contains any sensitive pattern.
func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// maskValue replaces all characters in a value with MaskChar,
// preserving the length up to a maximum of 8 asterisks.
func maskValue(value string) string {
	if value == "" {
		return ""
	}
	length := len(value)
	if length > 8 {
		length = 8
	}
	return strings.Repeat(MaskChar, length)
}

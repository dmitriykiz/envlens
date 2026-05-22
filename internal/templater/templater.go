// Package templater generates .env.example files from existing .env entries
// by stripping sensitive values and replacing them with placeholder hints.
package templater

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// sensitivePatterns lists substrings that indicate a key holds a secret.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

// Options controls template generation behaviour.
type Options struct {
	// Placeholder is the value written for sensitive keys.
	// Defaults to "<your_value_here>" when empty.
	Placeholder string
	// KeepNonSensitive preserves the original value for non-sensitive keys.
	KeepNonSensitive bool
}

// Generate reads entries and returns the text of a .env.example file.
func Generate(entries []parser.Entry, opts Options) string {
	if opts.Placeholder == "" {
		opts.Placeholder = "<your_value_here>"
	}

	var sb strings.Builder
	for _, e := range entries {
		if e.Comment != "" {
			fmt.Fprintf(&sb, "%s\n", e.Comment)
			continue
		}
		value := opts.Placeholder
		if opts.KeepNonSensitive && !isSensitive(e.Key) {
			value = e.Value
		} else if isSensitive(e.Key) {
			value = opts.Placeholder
		} else if !opts.KeepNonSensitive {
			value = opts.Placeholder
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, value)
	}
	return sb.String()
}

// WriteFile writes the generated template to destPath.
func WriteFile(destPath string, entries []parser.Entry, opts Options) error {
	content := Generate(entries, opts)
	return os.WriteFile(destPath, []byte(content), 0o644)
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

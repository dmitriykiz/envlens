// Package interpolator provides shell-style variable interpolation
// for .env file values, supporting ${VAR} and $VAR syntax with optional
// default expressions such as ${VAR:-default}.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Result holds the outcome of interpolating a single entry.
type Result struct {
	Key      string
	Original string
	Resolved string
	Warnings []string
}

var refRe = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Interpolate expands variable references within entry values using the
// provided entries as the lookup table. Unknown references are left as-is
// and a warning is recorded. Entries are processed in order so earlier
// definitions are visible to later ones.
func Interpolate(entries []parser.Entry) []Result {
	env := make(map[string]string, len(entries))
	results := make([]Result, 0, len(entries))

	for _, e := range entries {
		resolved, warnings := expand(e.Value, env)
		env[e.Key] = resolved
		results = append(results, Result{
			Key:      e.Key,
			Original: e.Value,
			Resolved: resolved,
			Warnings: warnings,
		})
	}
	return results
}

func expand(value string, env map[string]string) (string, []string) {
	var warnings []string
	result := refRe.ReplaceAllStringFunc(value, func(match string) string {
		inner := strings.TrimPrefix(match, "$")
		inner = strings.TrimPrefix(inner, "{")
		inner = strings.TrimSuffix(inner, "}")

		// Handle ${VAR:-default} syntax.
		if idx := strings.Index(inner, ":-"); idx != -1 {
			varName := inner[:idx]
			defaultVal := inner[idx+2:]
			if v, ok := env[varName]; ok && v != "" {
				return v
			}
			return defaultVal
		}

		if v, ok := env[inner]; ok {
			return v
		}
		warnings = append(warnings, fmt.Sprintf("unresolved reference: %s", inner))
		return match
	})
	return result, warnings
}

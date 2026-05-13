// Package reporter formats and outputs the results of an envlens audit.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envlens/internal/analyzer"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the audit results and metadata for rendering.
type Report struct {
	FilePath string
	Result   analyzer.Result
}

// WriteText writes a human-readable audit report to w.
func WriteText(w io.Writer, reports []Report) error {
	for _, r := range reports {
		fmt.Fprintf(w, "=== %s ===\n", r.FilePath)

		if len(r.Result.Duplicates) == 0 &&
			len(r.Result.SensitiveKeys) == 0 &&
			len(r.Result.MissingKeys) == 0 {
			fmt.Fprintln(w, "  ✓ No issues found.")
			continue
		}

		if len(r.Result.Duplicates) > 0 {
			fmt.Fprintln(w, "  [DUPLICATE KEYS]")
			for _, k := range r.Result.Duplicates {
				fmt.Fprintf(w, "    - %s\n", k)
			}
		}

		if len(r.Result.SensitiveKeys) > 0 {
			fmt.Fprintln(w, "  [SENSITIVE KEYS]")
			for _, k := range r.Result.SensitiveKeys {
				fmt.Fprintf(w, "    - %s\n", k)
			}
		}

		if len(r.Result.MissingKeys) > 0 {
			fmt.Fprintln(w, "  [MISSING KEYS]")
			for _, k := range r.Result.MissingKeys {
				fmt.Fprintf(w, "    - %s\n", k)
			}
		}
	}
	return nil
}

// Summary returns a one-line summary string for a report.
func Summary(r Report) string {
	parts := []string{}
	if n := len(r.Result.Duplicates); n > 0 {
		parts = append(parts, fmt.Sprintf("%d duplicate(s)", n))
	}
	if n := len(r.Result.SensitiveKeys); n > 0 {
		parts = append(parts, fmt.Sprintf("%d sensitive", n))
	}
	if n := len(r.Result.MissingKeys); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing", n))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s: OK", r.FilePath)
	}
	return fmt.Sprintf("%s: %s", r.FilePath, strings.Join(parts, ", "))
}

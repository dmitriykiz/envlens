package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envlens/internal/analyzer"
	"github.com/yourorg/envlens/internal/redactor"
)

// WriteText writes a human-readable audit report to w.
// If redactedEntries is non-nil, sensitive values are shown masked.
func WriteText(w io.Writer, result analyzer.Result, masked []redactor.MaskedEntry) error {
	if len(result.Duplicates) == 0 && len(result.SensitiveKeys) == 0 && len(result.MissingKeys) == 0 {
		_, err := fmt.Fprintln(w, "✔ No issues found.")
		return err
	}

	if len(result.Duplicates) > 0 {
		fmt.Fprintln(w, "Duplicate keys:")
		for _, d := range result.Duplicates {
			fmt.Fprintf(w, "  - %s (lines: %s)\n", d.Key, joinInts(d.Lines))
		}
	}

	if len(result.SensitiveKeys) > 0 {
		fmt.Fprintln(w, "Sensitive keys detected:")
		for _, k := range result.SensitiveKeys {
			maskedVal := findMasked(masked, k)
			if maskedVal != "" {
				fmt.Fprintf(w, "  - %s = %s\n", k, maskedVal)
			} else {
				fmt.Fprintf(w, "  - %s\n", k)
			}
		}
	}

	if len(result.MissingKeys) > 0 {
		fmt.Fprintln(w, "Missing keys (present in reference but not found):")
		for _, k := range result.MissingKeys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	return nil
}

// Summary writes a one-line summary of the audit result to w.
func Summary(w io.Writer, result analyzer.Result) error {
	issue := len(result.Duplicates) + len(result.SensitiveKeys) + len(result.MissingKeys)
	if issue == 0 {
		_, err := fmt.Fprintln(w, "Summary: OK — 0 issues detected.")
		return err
	}
	_, err := fmt.Fprintf(w, "Summary: %d issue(s) detected — duplicates: %d, sensitive: %d, missing: %d\n",
		issue,
		len(result.Duplicates),
		len(result.SensitiveKeys),
		len(result.MissingKeys),
	)
	return err
}

// joinInts formats a slice of ints as a comma-separated string.
func joinInts(nums []int) string {
	parts := make([]string, len(nums))
	for i, n := range nums {
		parts[i] = fmt.Sprintf("%d", n)
	}
	return strings.Join(parts, ", ")
}

// findMasked looks up the masked value for a given key.
func findMasked(entries []redactor.MaskedEntry, key string) string {
	for _, me := range entries {
		if me.Key == key && me.Masked {
			return me.Value
		}
	}
	return ""
}

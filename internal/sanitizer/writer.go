package sanitizer

import (
	"fmt"
	"io"
)

// WriteReport writes a human-readable sanitization report to w.
func WriteReport(w io.Writer, results []Result) {
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	fmt.Fprintf(w, "Sanitizer report: %d entr%s examined, %d changed\n",
		len(results), pluralY(len(results)), changed)
	if changed == 0 {
		fmt.Fprintln(w, "  No changes required.")
		return
	}
	for _, r := range results {
		if !r.Changed {
			continue
		}
		fmt.Fprintf(w, "  [%s:%d] %s\n", r.File, r.Line, r.Key)
		fmt.Fprintf(w, "    before: %q\n", r.OldVal)
		fmt.Fprintf(w, "    after:  %q\n", r.NewVal)
	}
}

// Summary returns a one-line summary string.
func Summary(results []Result) string {
	total := len(results)
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	if changed == 0 {
		return fmt.Sprintf("sanitizer: %d entr%s checked, none modified", total, pluralY(total))
	}
	return fmt.Sprintf("sanitizer: %d entr%s checked, %d modified", total, pluralY(total), changed)
}

func pluralY(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}

package scorer

import (
	"fmt"
	"io"
	"strings"
)

// WriteReport writes a human-readable score report to w.
func WriteReport(w io.Writer, r Report) {
	fmt.Fprintf(w, "File : %s\n", r.File)
	fmt.Fprintf(w, "Score: %d/100  Grade: %s\n", r.Score, r.Grade)

	if len(r.Penalties) == 0 {
		fmt.Fprintln(w, "Status: no issues detected ✓")
		return
	}

	fmt.Fprintln(w, "Penalties:")
	for _, p := range r.Penalties {
		fmt.Fprintf(w, "  -%d  %s\n", p.Points, p.Reason)
	}
}

// WriteSummary writes a compact multi-file summary table to w.
func WriteSummary(w io.Writer, reports []Report) {
	if len(reports) == 0 {
		fmt.Fprintln(w, "No files scored.")
		return
	}

	// Determine column width for file names.
	maxLen := len("File")
	for _, r := range reports {
		if len(r.File) > maxLen {
			maxLen = len(r.File)
		}
	}

	header := fmt.Sprintf("%-*s  %5s  %5s", maxLen, "File", "Score", "Grade")
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, r := range reports {
		fmt.Fprintf(w, "%-*s  %5d  %5s\n", maxLen, r.File, r.Score, r.Grade)
	}

	// Aggregate stats.
	total := 0
	for _, r := range reports {
		total += r.Score
	}
	avg := total / len(reports)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))
	fmt.Fprintf(w, "%-*s  %5d  %5s\n", maxLen, "Average", avg, gradeFor(avg))
}

package rotator

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteReport writes a human-readable rotation report to w.
func WriteReport(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "rotator: no keys matched rotation rules")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tLINE\tKEY\tROTATED\tOLD\tNEW")
	for _, r := range results {
		rotated := "no"
		if r.Rotated {
			rotated = "yes"
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%s\n",
			r.File, r.Line, r.Key, rotated,
			truncate(r.OldValue, 24), truncate(r.NewValue, 24),
		)
	}
	tw.Flush()
}

// Summary returns a one-line summary of the rotation run.
func Summary(results []Result) string {
	total := len(results)
	rotated := 0
	for _, r := range results {
		if r.Rotated {
			rotated++
		}
	}
	return fmt.Sprintf("rotator: %d key(s) evaluated, %d rotated", total, rotated)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

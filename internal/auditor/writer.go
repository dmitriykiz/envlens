package auditor

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteReport writes a human-readable audit report to w.
func WriteReport(w io.Writer, r Report) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "envlens audit report\n")
	fmt.Fprintf(tw, "generated: %s\n", r.AuditedAt.Format(time.RFC3339))
	fmt.Fprintf(tw, "files scanned: %d\n\n", len(r.Files))

	if len(r.Findings) == 0 {
		fmt.Fprintln(tw, "no issues found — all checks passed")
		tw.Flush()
		return
	}

	fmt.Fprintf(tw, "%-10s\t%-30s\t%5s\t%s\n", "SEVERITY", "FILE", "LINE", "MESSAGE (KEY)")
	fmt.Fprintf(tw, "%-10s\t%-30s\t%5s\t%s\n", "--------", "----", "----", "-------------")
	for _, f := range r.Findings {
		fmt.Fprintf(tw, "%-10s\t%-30s\t%5d\t%s (%s)\n",
			f.Severity, truncatePath(f.File, 30), f.Line, f.Message, f.Key)
	}
	tw.Flush()

	counts := CountBySeverity(r.Findings)
	fmt.Fprintf(w, "\ntotals — critical: %d  warning: %d  info: %d\n",
		counts[SeverityCritical], counts[SeverityWarning], counts[SeverityInfo])
}

func truncatePath(p string, max int) string {
	if len(p) <= max {
		return p
	}
	return "..." + p[len(p)-(max-3):]
}

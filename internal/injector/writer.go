package injector

import (
	"fmt"
	"io"
	"strings"
)

// WriteReport writes a human-readable summary of injection results to w.
func WriteReport(w io.Writer, results []Result) {
	injected := 0
	skipped := 0
	failed := 0

	for _, r := range results {
		switch {
		case r.Err != nil:
			fmt.Fprintf(w, "  [ERROR]   %s — %v\n", r.Key, r.Err)
			failed++
		case r.Skipped:
			fmt.Fprintf(w, "  [SKIP]    %s (already set)\n", r.Key)
			skipped++
		default:
			fmt.Fprintf(w, "  [SET]     %s\n", r.Key)
			injected++
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintf(w, "Injected: %d  Skipped: %d  Errors: %d\n", injected, skipped, failed)
}

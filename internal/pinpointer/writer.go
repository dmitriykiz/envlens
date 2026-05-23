package pinpointer

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// WriteReport writes a human-readable location report to w.
func WriteReport(w io.Writer, locations []Location) error {
	if len(locations) == 0 {
		_, err := fmt.Fprintln(w, "pinpointer: no matching keys found.")
		return err
	}

	sorted := make([]Location, len(locations))
	copy(sorted, locations)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].File != sorted[j].File {
			return sorted[i].File < sorted[j].File
		}
		return sorted[i].Line < sorted[j].Line
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tLINE\tCOL\tKEY\tVALUE")
	fmt.Fprintln(tw, "----\t----\t---\t---\t-----")
	for _, loc := range sorted {
		val := loc.Value
		if len(val) > 40 {
			val = val[:37] + "..."
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%s\t%s\n",
			loc.File, loc.Line, loc.Column, loc.Key, val)
	}
	return tw.Flush()
}

// Summary returns a one-line summary string.
func Summary(locations []Location) string {
	if len(locations) == 0 {
		return "pinpointer: 0 matches found"
	}
	files := map[string]struct{}{}
	for _, l := range locations {
		files[l.File] = struct{}{}
	}
	return fmt.Sprintf("pinpointer: %d match(es) across %d file(s)",
		len(locations), len(files))
}

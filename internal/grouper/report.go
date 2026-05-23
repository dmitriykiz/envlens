package grouper

import (
	"fmt"
	"io"
	"strings"
)

// WriteGroupReport writes a human-readable summary of the provided groups to w.
// Each group is rendered as a labelled section listing its keys.
//
//	[DB] (2 keys)
//	  DB_HOST
//	  DB_PORT
//
//	[AWS] (3 keys)
//	  AWS_ACCESS_KEY_ID
//	  …
func WriteGroupReport(w io.Writer, groups []Group) error {
	if len(groups) == 0 {
		_, err := fmt.Fprintln(w, "no groups found")
		return err
	}

	for i, g := range groups {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		header := fmt.Sprintf("[%s] (%d %s)",
			g.Prefix,
			len(g.Entries),
			plural(len(g.Entries), "key", "keys"),
		)
		if _, err := fmt.Fprintln(w, header); err != nil {
			return err
		}
		for _, e := range g.Entries {
			line := "  " + e.Key
			if e.Value != "" {
				line += "=" + truncateValue(e.Value, 30)
			}
			if _, err := fmt.Fprintln(w, line); err != nil {
				return err
			}
		}
	}
	return nil
}

// GroupSummary returns a compact one-line-per-group summary string.
func GroupSummary(groups []Group) string {
	if len(groups) == 0 {
		return "(no groups)"
	}
	parts := make([]string, len(groups))
	for i, g := range groups {
		parts[i] = fmt.Sprintf("%s:%d", g.Prefix, len(g.Entries))
	}
	return strings.Join(parts, "  ")
}

func plural(n int, singular, pluralForm string) string {
	if n == 1 {
		return singular
	}
	return pluralForm
}

func truncateValue(v string, max int) string {
	if len(v) <= max {
		return v
	}
	return v[:max] + "…"
}

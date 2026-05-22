package merger

import (
	"fmt"
	"io"
	"strings"

	"github.com/envlens/internal/parser"
)

// WriteEnv serialises a slice of entries to w in standard .env file format.
// Keys with values containing spaces or special characters are quoted.
func WriteEnv(w io.Writer, entries []parser.Entry) error {
	for _, e := range entries {
		line, err := formatEntry(e)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// WriteConflictReport writes a human-readable summary of merge conflicts to w.
func WriteConflictReport(w io.Writer, conflicts []Conflict) error {
	if len(conflicts) == 0 {
		_, err := fmt.Fprintln(w, "No merge conflicts detected.")
		return err
	}
	_, err := fmt.Fprintf(w, "%d conflict(s) detected:\n", len(conflicts))
	if err != nil {
		return err
	}
	for _, c := range conflicts {
		_, err = fmt.Fprintf(w, "  KEY: %s\n", c.Key)
		if err != nil {
			return err
		}
		for i, f := range c.Files {
			_, err = fmt.Fprintf(w, "    %s => %q\n", f, c.Values[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// formatEntry returns a KEY=VALUE string, quoting the value when necessary.
func formatEntry(e parser.Entry) (string, error) {
	val := e.Value
	if needsQuoting(val) {
		escaped := strings.ReplaceAll(val, `"`, `\"`)
		val = fmt.Sprintf(`"%s"`, escaped)
	}
	return fmt.Sprintf("%s=%s", e.Key, val), nil
}

// needsQuoting reports whether a value must be wrapped in double-quotes.
func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	return strings.ContainsAny(v, " \t#\"'\\")
}

package cipher

import (
	"fmt"
	"io"
	"strings"

	"github.com/envlens/internal/parser"
)

// WriteEncryptedEnv writes entries to w in standard .env format.
// Encrypted values (prefixed with "enc:") are written as-is so the file
// can be safely committed to version control.
func WriteEncryptedEnv(w io.Writer, entries []parser.Entry) error {
	for _, e := range entries {
		if e.Comment {
			if _, err := fmt.Fprintln(w, e.Value); err != nil {
				return err
			}
			continue
		}
		if e.Key == "" {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
			continue
		}
		val := e.Value
		if needsQuoting(val) {
			val = `"` + val + `"`
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

// WriteReport prints a human-readable summary of which keys were encrypted.
func WriteReport(w io.Writer, original, encrypted []parser.Entry) {
	var changed []string
	for i, e := range encrypted {
		if i < len(original) && original[i].Value != e.Value {
			changed = append(changed, e.Key)
		}
	}
	if len(changed) == 0 {
		fmt.Fprintln(w, "cipher: no keys were encrypted")
		return
	}
	fmt.Fprintf(w, "cipher: encrypted %d key(s): %s\n", len(changed), strings.Join(changed, ", "))
}

func needsQuoting(v string) bool {
	return strings.ContainsAny(v, " \t#")
}

package pinner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// SavePin writes a Pin to a JSON file at the given path.
func SavePin(path string, pin Pin) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("pinner: create %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(pin)
}

// LoadPin reads a Pin from a JSON file at the given path.
func LoadPin(path string) (Pin, error) {
	f, err := os.Open(path)
	if err != nil {
		return Pin{}, fmt.Errorf("pinner: open %s: %w", path, err)
	}
	defer f.Close()

	var pin Pin
	if err := json.NewDecoder(f).Decode(&pin); err != nil {
		return Pin{}, fmt.Errorf("pinner: decode %s: %w", path, err)
	}
	return pin, nil
}

// WriteReport writes a human-readable drift report to w.
func WriteReport(w io.Writer, reports []DriftReport) {
	hasDrift := false
	for _, r := range reports {
		if len(r.Entries) > 0 {
			hasDrift = true
			break
		}
	}

	if !hasDrift {
		fmt.Fprintln(w, "pinner: no drift detected — all files match their baselines")
		return
	}

	for _, r := range reports {
		if len(r.Entries) == 0 {
			continue
		}
		fmt.Fprintf(w, "drift in %s:\n", r.File)
		for _, e := range r.Entries {
			switch e.Reason {
			case "added":
				fmt.Fprintf(w, "  + %s (added)\n", e.Key)
			case "removed":
				fmt.Fprintf(w, "  - %s (removed)\n", e.Key)
			case "changed":
				fmt.Fprintf(w, "  ~ %s (value changed)\n", e.Key)
			}
		}
	}
}

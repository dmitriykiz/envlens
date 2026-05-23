package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envlens/internal/parser"
)

// Snapshot represents a point-in-time capture of parsed .env entries
// from one or more files.
type Snapshot struct {
	CreatedAt time.Time                      `json:"created_at"`
	Label     string                         `json:"label"`
	Files     map[string][]parser.Entry      `json:"files"`
}

// Take creates a new Snapshot by parsing each of the given file paths.
// Files that cannot be parsed are skipped and an error is collected.
func Take(label string, paths []string) (*Snapshot, []error) {
	snap := &Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Files:     make(map[string][]parser.Entry, len(paths)),
	}
	var errs []error
	for _, p := range paths {
		entries, err := parser.ParseFile(p)
		if err != nil {
			errs = append(errs, fmt.Errorf("snapshotter: parse %q: %w", p, err))
			continue
		}
		snap.Files[filepath.Clean(p)] = entries
	}
	return snap, errs
}

// Save writes the snapshot as JSON to the given file path.
func Save(snap *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshotter: create %q: %w", dest, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshotter: encode: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot from a JSON file.
func Load(src string) (*Snapshot, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: open %q: %w", src, err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshotter: decode: %w", err)
	}
	return &snap, nil
}

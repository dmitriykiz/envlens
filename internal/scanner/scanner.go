// Package scanner provides functionality for discovering and scanning
// .env files across a monorepo directory tree.
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Result holds the parsed entries from a single .env file.
type Result struct {
	// Path is the absolute or relative path to the .env file.
	Path string
	// Entries contains the key-value pairs parsed from the file.
	Entries []parser.Entry
	// Err holds any error encountered while parsing this file.
	Err error
}

// Options configures the behaviour of ScanDir.
type Options struct {
	// Patterns is a list of filename globs to match (e.g. ".env", ".env.*").
	// Defaults to [".env"] when empty.
	Patterns []string
	// ExcludeDirs lists directory names that should be skipped entirely.
	ExcludeDirs []string
}

// ScanDir walks root recursively and parses every .env file that matches
// the configured patterns. It returns one Result per discovered file.
func ScanDir(root string, opts Options) ([]Result, error) {
	if len(opts.Patterns) == 0 {
		opts.Patterns = []string{".env"}
	}

	excluded := make(map[string]bool, len(opts.ExcludeDirs))
	for _, d := range opts.ExcludeDirs {
		excluded[d] = true
	}

	var results []Result

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}
		if d.IsDir() {
			if excluded[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !matchesAny(d.Name(), opts.Patterns) {
			return nil
		}

		entries, parseErr := parser.ParseFile(path)
		results = append(results, Result{
			Path:    path,
			Entries: entries,
			Err:     parseErr,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// matchesAny reports whether name matches at least one of the given glob patterns.
func matchesAny(name string, patterns []string) bool {
	for _, p := range patterns {
		matched, err := filepath.Match(p, name)
		if err == nil && matched {
			return true
		}
		// Also support simple prefix matching for patterns like ".env.*"
		if strings.HasSuffix(p, "*") && strings.HasPrefix(name, strings.TrimSuffix(p, "*")) {
			return true
		}
	}
	return false
}

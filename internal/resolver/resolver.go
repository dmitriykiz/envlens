package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/internal/parser"
)

// Resolution strategy controls how variable references are expanded.
type Strategy int

const (
	// StrategyStrict fails if a referenced variable is not defined.
	StrategyStrict Strategy = iota
	// StrategyPermissive leaves unresolved references as-is.
	StrategyPermissive
)

// Result holds the resolved entries and any warnings produced during resolution.
type Result struct {
	Entries  []parser.Entry
	Warnings []string
}

// Resolve expands variable references (e.g. ${VAR} or $VAR) within entry
// values, using the entries themselves as the lookup table and falling back to
// os.Getenv when a key is not found locally.
func Resolve(entries []parser.Entry, strategy Strategy) (Result, error) {
	index := buildIndex(entries)
	result := Result{}

	for _, e := range entries {
		if e.Comment {
			result.Entries = append(result.Entries, e)
			continue
		}
		expanded, warnings, err := expandValue(e.Value, index, strategy)
		if err != nil {
			return Result{}, fmt.Errorf("resolver: key %q: %w", e.Key, err)
		}
		result.Warnings = append(result.Warnings, warnings...)
		e.Value = expanded
		result.Entries = append(result.Entries, e)
	}
	return result, nil
}

// buildIndex creates a key→value map from a slice of entries.
func buildIndex(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if !e.Comment {
			m[e.Key] = e.Value
		}
	}
	return m
}

// expandValue replaces all ${VAR} and $VAR occurrences inside s.
func expandValue(s string, index map[string]string, strategy Strategy) (string, []string, error) {
	var warnings []string
	result := os.Expand(s, func(key string) string {
		key = strings.TrimSpace(key)
		if v, ok := index[key]; ok {
			return v
		}
		if v := os.Getenv(key); v != "" {
			return v
		}
		if strategy == StrategyStrict {
			// Signal unresolved by embedding a sentinel; checked below.
			return "\x00UNRESOLVED:" + key + "\x00"
		}
		warnings = append(warnings, fmt.Sprintf("unresolved variable: $%s", key))
		return "${" + key + "}"
	})
	if idx := strings.Index(result, "\x00UNRESOLVED:"); idx != -1 {
		raw := result[idx+len("\x00UNRESOLVED:"):]
		key := strings.SplitN(raw, "\x00", 2)[0]
		return "", nil, fmt.Errorf("undefined variable: $%s", key)
	}
	return result, warnings, nil
}

// Package patcher provides functionality for applying key-value patches
// to parsed .env entry sets, supporting add, update, and delete operations.
package patcher

import (
	"fmt"

	"github.com/yourorg/envlens/internal/parser"
)

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"    // add or update a key
	OpDelete Op = "delete" // remove a key
)

// Rule describes a single patch instruction.
type Rule struct {
	Op    Op
	Key   string
	Value string // only used for OpSet
}

// Result holds the outcome of applying a patch.
type Result struct {
	Entries []parser.Entry
	Added   []string
	Updated []string
	Deleted []string
	Skipped []string // delete ops where key was not found
}

// Patch applies the given rules to entries and returns a Result.
// Rules are applied in order; later rules override earlier ones for the
// same key when both are OpSet.
func Patch(entries []parser.Entry, rules []Rule) (Result, error) {
	for i, r := range rules {
		if r.Key == "" {
			return Result{}, fmt.Errorf("rule[%d]: key must not be empty", i)
		}
		if r.Op != OpSet && r.Op != OpDelete {
			return Result{}, fmt.Errorf("rule[%d]: unknown op %q", i, r.Op)
		}
	}

	index := buildIndex(entries)
	res := Result{}

	for _, r := range rules {
		switch r.Op {
		case OpSet:
			if pos, exists := index[r.Key]; exists {
				entries[pos].Value = r.Value
				res.Updated = append(res.Updated, r.Key)
			} else {
				entries = append(entries, parser.Entry{Key: r.Key, Value: r.Value})
				index[r.Key] = len(entries) - 1
				res.Added = append(res.Added, r.Key)
			}
		case OpDelete:
			if _, exists := index[r.Key]; exists {
				entries, index = removeKey(entries, r.Key)
				res.Deleted = append(res.Deleted, r.Key)
			} else {
				res.Skipped = append(res.Skipped, r.Key)
			}
		}
	}

	res.Entries = entries
	return res, nil
}

func buildIndex(entries []parser.Entry) map[string]int {
	idx := make(map[string]int, len(entries))
	for i, e := range entries {
		if !e.IsComment && e.Key != "" {
			idx[e.Key] = i
		}
	}
	return idx
}

func removeKey(entries []parser.Entry, key string) ([]parser.Entry, map[string]int) {
	out := make([]parser.Entry, 0, len(entries)-1)
	for _, e := range entries {
		if e.Key != key {
			out = append(out, e)
		}
	}
	return out, buildIndex(out)
}

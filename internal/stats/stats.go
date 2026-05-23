// Package stats provides aggregate statistics across a set of parsed .env files.
package stats

import (
	"sort"

	"github.com/example/envlens/internal/parser"
)

// Report holds aggregate statistics for a collection of .env files.
type Report struct {
	TotalFiles    int
	TotalKeys     int
	UniqueKeys    int
	CommentLines  int
	BlankLines    int
	SensitiveKeys int
	EmptyValues   int
	TopKeys       []KeyCount // top keys by occurrence across files
}

// KeyCount pairs a key with how many files it appears in.
type KeyCount struct {
	Key   string
	Count int
}

var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN",
	"API_KEY", "PRIVATE", "CREDENTIAL", "AUTH",
}

// Compute derives a Report from a map of filename → parsed entries.
func Compute(files map[string][]parser.Entry) Report {
	keyCounts := make(map[string]int)
	var r Report
	r.TotalFiles = len(files)

	for _, entries := range files {
		seen := make(map[string]bool)
		for _, e := range entries {
			if e.Comment {
				r.CommentLines++
				continue
			}
			if e.Key == "" {
				r.BlankLines++
				continue
			}
			r.TotalKeys++
			if e.Value == "" {
				r.EmptyValues++
			}
			if isSensitive(e.Key) {
				r.SensitiveKeys++
			}
			if !seen[e.Key] {
				seen[e.Key] = true
				keyCounts[e.Key]++
			}
		}
	}

	r.UniqueKeys = len(keyCounts)
	r.TopKeys = topN(keyCounts, 5)
	return r
}

func isSensitive(key string) bool {
	for _, p := range sensitivePatterns {
		if contains(key, p) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		len(s) > 0 && indexStr(s, sub) >= 0)
}

func indexStr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func topN(counts map[string]int, n int) []KeyCount {
	kcs := make([]KeyCount, 0, len(counts))
	for k, c := range counts {
		kcs = append(kcs, KeyCount{Key: k, Count: c})
	}
	sort.Slice(kcs, func(i, j int) bool {
		if kcs[i].Count != kcs[j].Count {
			return kcs[i].Count > kcs[j].Count
		}
		return kcs[i].Key < kcs[j].Key
	})
	if len(kcs) > n {
		return kcs[:n]
	}
	return kcs
}

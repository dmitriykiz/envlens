package stats_test

import (
	"testing"

	"github.com/example/envlens/internal/parser"
	"github.com/example/envlens/internal/stats"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompute_BasicCounts(t *testing.T) {
	files := map[string][]parser.Entry{
		"a.env": makeEntries("APP_NAME", "myapp", "PORT", "8080"),
		"b.env": makeEntries("APP_NAME", "other", "DEBUG", "true"),
	}
	r := stats.Compute(files)
	if r.TotalFiles != 2 {
		t.Errorf("TotalFiles: want 2, got %d", r.TotalFiles)
	}
	if r.TotalKeys != 4 {
		t.Errorf("TotalKeys: want 4, got %d", r.TotalKeys)
	}
	if r.UniqueKeys != 3 {
		t.Errorf("UniqueKeys: want 3, got %d", r.UniqueKeys)
	}
}

func TestCompute_SensitiveKeys(t *testing.T) {
	files := map[string][]parser.Entry{
		"a.env": makeEntries("API_KEY", "abc", "DB_PASSWORD", "secret", "HOST", "localhost"),
	}
	r := stats.Compute(files)
	if r.SensitiveKeys != 2 {
		t.Errorf("SensitiveKeys: want 2, got %d", r.SensitiveKeys)
	}
}

func TestCompute_EmptyValues(t *testing.T) {
	files := map[string][]parser.Entry{
		"a.env": makeEntries("KEY1", "", "KEY2", "val", "KEY3", ""),
	}
	r := stats.Compute(files)
	if r.EmptyValues != 2 {
		t.Errorf("EmptyValues: want 2, got %d", r.EmptyValues)
	}
}

func TestCompute_CommentsAndBlanks(t *testing.T) {
	files := map[string][]parser.Entry{
		"a.env": {
			{Key: "", Value: "", Comment: true},
			{Key: "", Value: "", Comment: false},
			{Key: "FOO", Value: "bar"},
		},
	}
	r := stats.Compute(files)
	if r.CommentLines != 1 {
		t.Errorf("CommentLines: want 1, got %d", r.CommentLines)
	}
	if r.BlankLines != 1 {
		t.Errorf("BlankLines: want 1, got %d", r.BlankLines)
	}
	if r.TotalKeys != 1 {
		t.Errorf("TotalKeys: want 1, got %d", r.TotalKeys)
	}
}

func TestCompute_TopKeys(t *testing.T) {
	files := map[string][]parser.Entry{
		"a.env": makeEntries("COMMON", "1", "RARE", "x"),
		"b.env": makeEntries("COMMON", "2", "ALSO", "y"),
		"c.env": makeEntries("COMMON", "3"),
	}
	r := stats.Compute(files)
	if len(r.TopKeys) == 0 {
		t.Fatal("expected TopKeys to be non-empty")
	}
	if r.TopKeys[0].Key != "COMMON" {
		t.Errorf("expected COMMON as top key, got %s", r.TopKeys[0].Key)
	}
	if r.TopKeys[0].Count != 3 {
		t.Errorf("expected count 3 for COMMON, got %d", r.TopKeys[0].Count)
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	r := stats.Compute(map[string][]parser.Entry{})
	if r.TotalFiles != 0 || r.TotalKeys != 0 || r.UniqueKeys != 0 {
		t.Errorf("expected zero report for empty input, got %+v", r)
	}
}

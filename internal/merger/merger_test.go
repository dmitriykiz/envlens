package merger

import (
	"testing"

	"github.com/envlens/internal/parser"
)

func makeFileEntries(kv ...string) []parser.Entry {
	if len(kv)%2 != 0 {
		panic("makeFileEntries requires key/value pairs")
	}
	var entries []parser.Entry
	for i := 0; i < len(kv); i += 2 {
		entries = append(entries, parser.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return entries
}

func TestMerge_NoConflicts(t *testing.T) {
	input := map[string][]parser.Entry{
		"a.env": makeFileEntries("HOST", "localhost", "PORT", "8080"),
		"b.env": makeFileEntries("DB", "postgres"),
	}
	res, err := Merge(input, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	input := map[string][]parser.Entry{
		"a.env": makeFileEntries("KEY", "first"),
		"b.env": makeFileEntries("KEY", "second"),
	}
	res, err := Merge(input, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "first" {
		t.Errorf("expected 'first', got %q", res.Entries[0].Value)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	input := map[string][]parser.Entry{
		"a.env": makeFileEntries("KEY", "first"),
		"b.env": makeFileEntries("KEY", "second"),
	}
	res, err := Merge(input, StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "second" {
		t.Errorf("expected 'second', got %q", res.Entries[0].Value)
	}
}

func TestMerge_StrategyError(t *testing.T) {
	input := map[string][]parser.Entry{
		"a.env": makeFileEntries("KEY", "v1"),
		"b.env": makeFileEntries("KEY", "v2"),
	}
	_, err := Merge(input, StrategyError)
	if err == nil {
		t.Fatal("expected error for conflicting key, got nil")
	}
}

func TestMerge_EmptyInput(t *testing.T) {
	res, err := Merge(map[string][]parser.Entry{}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(res.Entries))
	}
}

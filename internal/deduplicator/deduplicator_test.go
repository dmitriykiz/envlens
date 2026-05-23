package deduplicator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/deduplicator"
	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	if len(pairs)%2 != 0 {
		panic("makeEntries: pairs must be even")
	}
	out := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestDeduplicate_NoDuplicates(t *testing.T) {
	entries := makeEntries("FOO", "1", "BAR", "2")
	res := deduplicator.Deduplicate(entries, deduplicator.KeepFirst)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %v", res.Duplicates)
	}
}

func TestDeduplicate_KeepFirst(t *testing.T) {
	entries := makeEntries("FOO", "first", "BAR", "b", "FOO", "second")
	res := deduplicator.Deduplicate(entries, deduplicator.KeepFirst)

	if res.Duplicates["FOO"] != 1 {
		t.Fatalf("expected 1 duplicate for FOO, got %d", res.Duplicates["FOO"])
	}
	for _, e := range res.Entries {
		if e.Key == "FOO" && e.Value != "first" {
			t.Fatalf("KeepFirst: expected value 'first', got %q", e.Value)
		}
	}
}

func TestDeduplicate_KeepLast(t *testing.T) {
	entries := makeEntries("FOO", "first", "BAR", "b", "FOO", "second")
	res := deduplicator.Deduplicate(entries, deduplicator.KeepLast)

	if res.Duplicates["FOO"] != 1 {
		t.Fatalf("expected 1 duplicate for FOO, got %d", res.Duplicates["FOO"])
	}
	for _, e := range res.Entries {
		if e.Key == "FOO" && e.Value != "second" {
			t.Fatalf("KeepLast: expected value 'second', got %q", e.Value)
		}
	}
}

func TestDeduplicate_PreservesComments(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Value: "# section"},
		{Key: "A", Value: "1"},
		{Comment: true, Value: "# again"},
		{Key: "A", Value: "2"},
	}
	res := deduplicator.Deduplicate(entries, deduplicator.KeepFirst)

	commentCount := 0
	for _, e := range res.Entries {
		if e.Comment {
			commentCount++
		}
	}
	if commentCount != 2 {
		t.Fatalf("expected 2 comment entries preserved, got %d", commentCount)
	}
}

func TestDeduplicate_MultipleDuplicates(t *testing.T) {
	entries := makeEntries("X", "a", "X", "b", "X", "c")
	res := deduplicator.Deduplicate(entries, deduplicator.KeepFirst)

	if res.Duplicates["X"] != 2 {
		t.Fatalf("expected 2 extra copies removed, got %d", res.Duplicates["X"])
	}
	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 surviving entry, got %d", len(res.Entries))
	}
}

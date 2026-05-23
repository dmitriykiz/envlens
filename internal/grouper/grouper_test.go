package grouper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/grouper"
	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(kvs ...string) []parser.Entry {
	if len(kvs)%2 != 0 {
		panic("makeEntries: need even number of args")
	}
	out := make([]parser.Entry, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestByPrefix_BasicGrouping(t *testing.T) {
	entries := makeEntries(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"AWS_KEY", "abc",
		"AWS_SECRET", "xyz",
		"PORT", "8080",
	)
	groups := grouper.ByPrefix(entries, grouper.Options{})

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	// Groups are sorted by prefix: AWS, DB, (ungrouped)
	if groups[0].Prefix != "AWS" {
		t.Errorf("expected first prefix AWS, got %s", groups[0].Prefix)
	}
	if groups[1].Prefix != "DB" {
		t.Errorf("expected second prefix DB, got %s", groups[1].Prefix)
	}
	if groups[2].Prefix != "(ungrouped)" {
		t.Errorf("expected third prefix (ungrouped), got %s", groups[2].Prefix)
	}
}

func TestByPrefix_MinGroupSizeFilters(t *testing.T) {
	entries := makeEntries(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"AWS_KEY", "abc",
	)
	groups := grouper.ByPrefix(entries, grouper.Options{MinGroupSize: 2})

	if len(groups) != 1 {
		t.Fatalf("expected 1 group (DB only), got %d", len(groups))
	}
	if groups[0].Prefix != "DB" {
		t.Errorf("expected DB group, got %s", groups[0].Prefix)
	}
}

func TestByPrefix_CustomSeparator(t *testing.T) {
	entries := makeEntries(
		"DB.HOST", "localhost",
		"DB.PORT", "5432",
	)
	groups := grouper.ByPrefix(entries, grouper.Options{Separator: "."})

	if len(groups) != 1 || groups[0].Prefix != "DB" {
		t.Errorf("expected one DB group with dot separator, got %+v", groups)
	}
}

func TestByPrefix_SkipsCommentEntries(t *testing.T) {
	entries := []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "# a comment", Comment: true},
	}
	groups := grouper.ByPrefix(entries, grouper.Options{})
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
}

func TestFlatten_PreservesAllEntries(t *testing.T) {
	entries := makeEntries(
		"DB_HOST", "localhost",
		"AWS_KEY", "abc",
		"APP_NAME", "envlens",
	)
	groups := grouper.ByPrefix(entries, grouper.Options{})
	flat := grouper.Flatten(groups)

	if len(flat) != len(entries) {
		t.Errorf("expected %d entries after flatten, got %d", len(entries), len(flat))
	}
}

func TestByPrefix_EmptyInput(t *testing.T) {
	groups := grouper.ByPrefix(nil, grouper.Options{})
	if len(groups) != 0 {
		t.Errorf("expected empty groups for nil input")
	}
}

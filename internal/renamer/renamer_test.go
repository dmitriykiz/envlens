package renamer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/renamer"
)

func makeEntries(kvs ...string) []parser.Entry {
	if len(kvs)%2 != 0 {
		panic("makeEntries: need even number of arguments")
	}
	out := make([]parser.Entry, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1], Line: i/2 + 1})
	}
	return out
}

func TestRename_NoRules(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "DB_PORT", "5432")
	res, err := renamer.Rename("test.env", entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
	if res.Entries[0].Key != "DB_HOST" {
		t.Errorf("key should be unchanged")
	}
}

func TestRename_ExactMatch(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "APP_ENV", "prod")
	rules := []renamer.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}}
	res, err := renamer.Rename("test.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	if res.Changes[0].OldKey != "DB_HOST" || res.Changes[0].NewKey != "DATABASE_HOST" {
		t.Errorf("unexpected change: %+v", res.Changes[0])
	}
	if res.Entries[0].Key != "DATABASE_HOST" {
		t.Errorf("entry key not updated")
	}
	if res.Entries[1].Key != "APP_ENV" {
		t.Errorf("unrelated entry should be unchanged")
	}
}

func TestRename_RegexpMatch(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "DB_PORT", "5432", "APP_SECRET", "x")
	rules := []renamer.Rule{{From: `^DB_(.+)$`, To: "DATABASE_$1", IsRegexp: true}}
	res, err := renamer.Rename("test.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(res.Changes))
	}
	if res.Entries[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %s", res.Entries[0].Key)
	}
	if res.Entries[1].Key != "DATABASE_PORT" {
		t.Errorf("expected DATABASE_PORT, got %s", res.Entries[1].Key)
	}
	if res.Entries[2].Key != "APP_SECRET" {
		t.Errorf("APP_SECRET should be unchanged")
	}
}

func TestRename_InvalidRegexp(t *testing.T) {
	entries := makeEntries("KEY", "val")
	rules := []renamer.Rule{{From: `[invalid`, To: "NEW", IsRegexp: true}}
	_, err := renamer.Rename("test.env", entries, rules)
	if err == nil {
		t.Fatal("expected error for invalid regexp, got nil")
	}
}

func TestRename_SkipsComments(t *testing.T) {
	entries := []parser.Entry{
		{IsComment: true, Raw: "# comment", Line: 1},
		{Key: "FOO", Value: "bar", Line: 2},
	}
	rules := []renamer.Rule{{From: "FOO", To: "BAR"}}
	res, err := renamer.Rename("test.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].IsComment == false || res.Entries[0].Raw != "# comment" {
		t.Errorf("comment entry should be preserved verbatim")
	}
	if res.Entries[1].Key != "BAR" {
		t.Errorf("expected BAR, got %s", res.Entries[1].Key)
	}
}

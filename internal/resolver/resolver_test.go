package resolver

import (
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestResolve_NoReferences(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	res, err := Resolve(entries, StrategyStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "localhost" || res.Entries[1].Value != "5432" {
		t.Errorf("values should be unchanged")
	}
}

func TestResolve_ExpandsInternalReference(t *testing.T) {
	entries := makeEntries("BASE", "http://example.com", "URL", "${BASE}/api")
	res, err := Resolve(entries, StrategyStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "http://example.com/api"
	if res.Entries[1].Value != want {
		t.Errorf("got %q, want %q", res.Entries[1].Value, want)
	}
}

func TestResolve_StrictFailsOnMissing(t *testing.T) {
	entries := makeEntries("DSN", "postgres://${DB_USER}@localhost/db")
	_, err := Resolve(entries, StrategyStrict)
	if err == nil {
		t.Fatal("expected error for undefined variable")
	}
}

func TestResolve_PermissiveWarnsOnMissing(t *testing.T) {
	entries := makeEntries("DSN", "postgres://${DB_USER}@localhost/db")
	res, err := Resolve(entries, StrategyPermissive)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Warnings) == 0 {
		t.Error("expected at least one warning")
	}
	if res.Entries[0].Value != "postgres://${DB_USER}@localhost/db" {
		t.Errorf("permissive mode should leave unresolved references intact")
	}
}

func TestResolve_CommentPassthrough(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Raw: "# section header"},
		{Key: "KEY", Value: "val"},
	}
	res, err := Resolve(entries, StrategyStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Entries[0].Comment {
		t.Error("comment entry should pass through unchanged")
	}
}

func TestResolve_ChainedReferences(t *testing.T) {
	entries := makeEntries(
		"PROTO", "https",
		"HOST", "example.com",
		"BASE", "${PROTO}://${HOST}",
		"URL", "${BASE}/health",
	)
	res, err := Resolve(entries, StrategyStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://example.com/health"
	if res.Entries[3].Value != want {
		t.Errorf("got %q, want %q", res.Entries[3].Value, want)
	}
}

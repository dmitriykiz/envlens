package interpolator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/interpolator"
	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	if len(pairs)%2 != 0 {
		panic("makeEntries requires an even number of arguments")
	}
	out := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestInterpolate_NoReferences(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	results := interpolator.Interpolate(entries)
	if results[0].Resolved != "localhost" {
		t.Errorf("expected localhost, got %s", results[0].Resolved)
	}
	if results[1].Resolved != "5432" {
		t.Errorf("expected 5432, got %s", results[1].Resolved)
	}
}

func TestInterpolate_CurlyBraceRef(t *testing.T) {
	entries := makeEntries("HOST", "db.local", "DSN", "postgres://${HOST}/app")
	results := interpolator.Interpolate(entries)
	want := "postgres://db.local/app"
	if results[1].Resolved != want {
		t.Errorf("expected %q, got %q", want, results[1].Resolved)
	}
	if len(results[1].Warnings) != 0 {
		t.Errorf("unexpected warnings: %v", results[1].Warnings)
	}
}

func TestInterpolate_BareRef(t *testing.T) {
	entries := makeEntries("USER", "admin", "GREETING", "Hello $USER")
	results := interpolator.Interpolate(entries)
	if results[1].Resolved != "Hello admin" {
		t.Errorf("got %q", results[1].Resolved)
	}
}

func TestInterpolate_DefaultFallback(t *testing.T) {
	entries := makeEntries("DSN", "postgres://${HOST:-127.0.0.1}/db")
	results := interpolator.Interpolate(entries)
	want := "postgres://127.0.0.1/db"
	if results[0].Resolved != want {
		t.Errorf("expected %q, got %q", want, results[0].Resolved)
	}
}

func TestInterpolate_DefaultSkippedWhenSet(t *testing.T) {
	entries := makeEntries("HOST", "myhost", "DSN", "${HOST:-fallback}")
	results := interpolator.Interpolate(entries)
	if results[1].Resolved != "myhost" {
		t.Errorf("expected myhost, got %s", results[1].Resolved)
	}
}

func TestInterpolate_UnresolvedWarning(t *testing.T) {
	entries := makeEntries("URL", "http://${UNKNOWN_HOST}/path")
	results := interpolator.Interpolate(entries)
	if len(results[0].Warnings) == 0 {
		t.Error("expected a warning for unresolved reference")
	}
	if results[0].Resolved != "http://${UNKNOWN_HOST}/path" {
		t.Errorf("unresolved ref should remain verbatim, got %q", results[0].Resolved)
	}
}

func TestInterpolate_ChainedRefs(t *testing.T) {
	entries := makeEntries("A", "hello", "B", "${A}_world", "C", "${B}!")
	results := interpolator.Interpolate(entries)
	if results[2].Resolved != "hello_world!" {
		t.Errorf("chained expansion failed: got %q", results[2].Resolved)
	}
}

package filter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/filter"
	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "8080")
	got, err := filter.Filter(entries, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilter_KeyPattern(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens")
	got, err := filter.Filter(entries, filter.Options{KeyPattern: `^DB_`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 DB_ entries, got %d", len(got))
	}
}

func TestFilter_ValuePattern(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432", "TIMEOUT", "30s")
	got, err := filter.Filter(entries, filter.Options{ValuePattern: `^[0-9]+$`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Key != "PORT" {
		t.Fatalf("expected only PORT, got %+v", got)
	}
}

func TestFilter_SensitiveOnly(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "API_KEY", "s3cr3t", "DB_PASSWORD", "hunter2")
	got, err := filter.Filter(entries, filter.Options{SensitiveOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 sensitive entries, got %d", len(got))
	}
}

func TestFilter_ExcludeComments(t *testing.T) {
	entries := []parser.Entry{
		{Key: "", Raw: "# a comment", IsComment: true},
		{Key: "HOST", Value: "localhost"},
		{Key: "", Raw: ""},
	}
	got, err := filter.Filter(entries, filter.Options{ExcludeComments: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Key != "HOST" {
		t.Fatalf("expected only HOST entry, got %+v", got)
	}
}

func TestFilter_InvalidKeyPattern_ReturnsError(t *testing.T) {
	entries := makeEntries("HOST", "localhost")
	_, err := filter.Filter(entries, filter.Options{KeyPattern: `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regexp, got nil")
	}
}

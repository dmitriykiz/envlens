package normalizer

import (
	"testing"

	"github.com/your-org/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestNormalize_NoChanges(t *testing.T) {
	entries := makeEntries("APP_ENV", "production")
	results := Normalize(entries, DefaultOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Changed {
		t.Errorf("expected no change, got changes: %v", results[0].Changes)
	}
}

func TestNormalize_UppercasesKey(t *testing.T) {
	entries := makeEntries("app_env", "production")
	results := Normalize(entries, DefaultOptions())
	if results[0].Normalized.Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %q", results[0].Normalized.Key)
	}
	if !results[0].Changed {
		t.Error("expected Changed to be true")
	}
}

func TestNormalize_TrimWhitespace(t *testing.T) {
	entries := makeEntries("  DB_HOST  ", "  localhost  ")
	results := Normalize(entries, DefaultOptions())
	r := results[0].Normalized
	if r.Key != "DB_HOST" {
		t.Errorf("key not trimmed: %q", r.Key)
	}
	if r.Value != "localhost" {
		t.Errorf("value not trimmed: %q", r.Value)
	}
}

func TestNormalize_StripInlineComment(t *testing.T) {
	entries := makeEntries("LOG_LEVEL", "debug # set to info in prod")
	results := Normalize(entries, DefaultOptions())
	if results[0].Normalized.Value != "debug" {
		t.Errorf("expected 'debug', got %q", results[0].Normalized.Value)
	}
}

func TestNormalize_UnquoteDoubleQuotes(t *testing.T) {
	entries := makeEntries("SECRET", `"my-secret"`)
	results := Normalize(entries, DefaultOptions())
	if results[0].Normalized.Value != "my-secret" {
		t.Errorf("expected unquoted value, got %q", results[0].Normalized.Value)
	}
}

func TestNormalize_UnquoteSingleQuotes(t *testing.T) {
	entries := makeEntries("TOKEN", "'abc123'")
	results := Normalize(entries, DefaultOptions())
	if results[0].Normalized.Value != "abc123" {
		t.Errorf("expected unquoted value, got %q", results[0].Normalized.Value)
	}
}

func TestNormalize_PassthroughComments(t *testing.T) {
	entries := []parser.Entry{{Raw: "# a comment", IsComment: true}}
	results := Normalize(entries, DefaultOptions())
	if results[0].Changed {
		t.Error("comment entries should not be marked as changed")
	}
	if results[0].Normalized.Raw != "# a comment" {
		t.Errorf("comment raw value altered: %q", results[0].Normalized.Raw)
	}
}

func TestNormalize_DisabledPasses(t *testing.T) {
	opts := Options{UppercaseKeys: false, TrimSpace: false, StripInlineComments: false, UnquoteValues: false}
	entries := makeEntries("lower_key", `"quoted" # comment`)
	results := Normalize(entries, opts)
	if results[0].Changed {
		t.Errorf("expected no changes with all passes disabled, got: %v", results[0].Changes)
	}
}

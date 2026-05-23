package classifier_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/classifier"
	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestClassify_Credential(t *testing.T) {
	entries := makeEntries("API_SECRET", "abc123", "AUTH_TOKEN", "xyz")
	results := classifier.Classify(entries)
	for _, r := range results {
		if r.Category != classifier.CategoryCredential {
			t.Errorf("expected credential for %q, got %q", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_Database(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "DATABASE_URL", "postgres://...")
	results := classifier.Classify(entries)
	for _, r := range results {
		if r.Category != classifier.CategoryDatabase {
			t.Errorf("expected database for %q, got %q", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_Network(t *testing.T) {
	entries := makeEntries("SERVER_HOST", "0.0.0.0", "APP_PORT", "8080")
	results := classifier.Classify(entries)
	for _, r := range results {
		if r.Category != classifier.CategoryNetwork {
			t.Errorf("expected network for %q, got %q", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_Feature(t *testing.T) {
	entries := makeEntries("FEATURE_DARK_MODE", "true", "ENABLE_BETA", "false")
	results := classifier.Classify(entries)
	for _, r := range results {
		if r.Category != classifier.CategoryFeature {
			t.Errorf("expected feature for %q, got %q", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_Unknown(t *testing.T) {
	entries := makeEntries("APP_NAME", "myapp", "RETRY_COUNT", "3")
	results := classifier.Classify(entries)
	for _, r := range results {
		if r.Category != classifier.CategoryUnknown {
			t.Errorf("expected unknown for %q, got %q", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_SkipsCommentsAndBlanks(t *testing.T) {
	entries := []parser.Entry{
		{IsComment: true, Raw: "# comment"},
		{IsBlank: true},
		{Key: "APP_PORT", Value: "9000"},
	}
	results := classifier.Classify(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Entry.Key != "APP_PORT" {
		t.Errorf("unexpected key %q", results[0].Entry.Key)
	}
}

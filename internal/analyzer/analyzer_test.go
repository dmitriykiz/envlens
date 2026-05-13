package analyzer_test

import (
	"testing"

	"github.com/user/envlens/internal/analyzer"
	"github.com/user/envlens/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestAnalyze_NoDuplicatesNoSensitiveNoMissing(t *testing.T) {
	e := entries("APP_ENV", "production", "PORT", "8080")
	r := analyzer.Analyze(".env", e, nil)
	if len(r.Duplicates) != 0 || len(r.Sensitive) != 0 || len(r.Missing) != 0 {
		t.Errorf("expected clean result, got %+v", r)
	}
}

func TestAnalyze_DetectsDuplicates(t *testing.T) {
	e := entries("PORT", "8080", "PORT", "9090")
	r := analyzer.Analyze(".env", e, nil)
	if len(r.Duplicates) != 1 || r.Duplicates[0] != "PORT" {
		t.Errorf("expected duplicate PORT, got %v", r.Duplicates)
	}
}

func TestAnalyze_DetectsSensitiveKeys(t *testing.T) {
	e := entries("DB_PASSWORD", "s3cr3t", "APP_ENV", "dev")
	r := analyzer.Analyze(".env", e, nil)
	if len(r.Sensitive) != 1 || r.Sensitive[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD as sensitive, got %v", r.Sensitive)
	}
}

func TestAnalyze_DetectsMissingKeys(t *testing.T) {
	e := entries("APP_ENV", "dev")
	ref := []string{"APP_ENV", "PORT", "DB_URL"}
	r := analyzer.Analyze(".env", e, ref)
	if len(r.Missing) != 2 {
		t.Errorf("expected 2 missing keys, got %v", r.Missing)
	}
}

func TestAnalyze_SensitiveCaseInsensitive(t *testing.T) {
	e := entries("stripe_api_key", "sk_live_abc", "github_token", "ghp_xyz")
	r := analyzer.Analyze(".env", e, nil)
	if len(r.Sensitive) != 2 {
		t.Errorf("expected 2 sensitive keys, got %v", r.Sensitive)
	}
}

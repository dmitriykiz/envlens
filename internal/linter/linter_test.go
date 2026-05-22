package linter

import (
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEntries(kvs ...struct {
	key, value string
	line       int
}) []parser.Entry {
	out := make([]parser.Entry, len(kvs))
	for i, kv := range kvs {
		out[i] = parser.Entry{Key: kv.key, Value: kv.value, Line: kv.line}
	}
	return out
}

func TestLint_NoViolations(t *testing.T) {
	entries := makeEntries(
		struct{ key, value string; line int }{"APP_ENV", "production", 1},
		struct{ key, value string; line int }{"DB_HOST", "localhost", 2},
	)
	violations, err := Lint(".env", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	entries := makeEntries(struct{ key, value string; line int }{"APP_ENV", "", 1})
	violations, err := Lint(".env", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(violations, "no-empty-value") {
		t.Errorf("expected no-empty-value violation")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	entries := makeEntries(struct{ key, value string; line int }{"app_env", "dev", 1})
	violations, err := Lint(".env", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(violations, "uppercase-key") {
		t.Errorf("expected uppercase-key violation")
	}
}

func TestLint_KeyWithSpecialChars(t *testing.T) {
	entries := makeEntries(struct{ key, value string; line int }{"APP-ENV", "dev", 1})
	violations, err := Lint(".env", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(violations, "no-special-chars-in-key") {
		t.Errorf("expected no-special-chars-in-key violation")
	}
}

func TestLint_EmptyFilePath(t *testing.T) {
	_, err := Lint("", nil)
	if err == nil {
		t.Error("expected error for empty file path")
	}
}

func containsRule(violations []Violation, rule string) bool {
	for _, v := range violations {
		if v.Rule == rule {
			return true
		}
	}
	return false
}

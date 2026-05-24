package rotator

import (
	"strings"
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1], Line: i/2 + 1})
	}
	return out
}

func TestRotate_NoRules(t *testing.T) {
	entries := makeEntries("DB_PASS", "hunter2", "APP_ENV", "production")
	updated, results, err := Rotate("test.env", entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if updated[0].Value != "hunter2" {
		t.Errorf("value should be unchanged")
	}
}

func TestRotate_RedactStrategy(t *testing.T) {
	entries := makeEntries("API_SECRET", "abc123", "APP_ENV", "prod")
	rules := []Rule{
		{KeyPattern: `SECRET`, Strategy: StrategyRedact, Placeholder: "***"},
	}
	updated, results, err := Rotate("a.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Value != "***" {
		t.Errorf("expected redacted value, got %q", updated[0].Value)
	}
	if updated[1].Value != "prod" {
		t.Errorf("non-matching entry should be unchanged")
	}
	if len(results) != 1 || !results[0].Rotated {
		t.Errorf("expected one rotated result")
	}
}

func TestRotate_IncrementStrategy(t *testing.T) {
	entries := makeEntries("DB_PASS", "pass001")
	rules := []Rule{{KeyPattern: `DB_PASS`, Strategy: StrategyIncrement}}
	updated, _, err := Rotate("b.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Value != "pass002" {
		t.Errorf("expected incremented value, got %q", updated[0].Value)
	}
}

func TestRotate_IncrementStrategy_NoSuffix(t *testing.T) {
	entries := makeEntries("TOKEN", "mysecret")
	rules := []Rule{{KeyPattern: `TOKEN`, Strategy: StrategyIncrement}}
	updated, _, err := Rotate("c.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(updated[0].Value, "_2") {
		t.Errorf("expected _2 suffix, got %q", updated[0].Value)
	}
}

func TestRotate_CustomStrategy(t *testing.T) {
	entries := makeEntries("JWT_SECRET", "old")
	rules := []Rule{
		{
			KeyPattern: `JWT`,
			Strategy:   StrategyCustom,
			Generator:  func(_, _ string) (string, error) { return "brand-new", nil },
		},
	}
	updated, results, err := Rotate("d.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Value != "brand-new" {
		t.Errorf("expected custom value, got %q", updated[0].Value)
	}
	if results[0].OldValue != "old" {
		t.Errorf("OldValue should be preserved in result")
	}
}

func TestRotate_InvalidPattern(t *testing.T) {
	entries := makeEntries("KEY", "val")
	rules := []Rule{{KeyPattern: `[invalid`, Strategy: StrategyRedact}}
	_, _, err := Rotate("e.env", entries, rules)
	if err == nil {
		t.Error("expected error for invalid regexp")
	}
}

func TestRotate_SkipsComments(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Raw: "# this is a comment", Line: 1},
		{Key: "SECRET", Value: "val", Line: 2},
	}
	rules := []Rule{{KeyPattern: `SECRET`, Strategy: StrategyRedact}}
	updated, results, err := Rotate("f.env", entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Raw != "# this is a comment" {
		t.Error("comment entry should be untouched")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

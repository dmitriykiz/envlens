package validator_test

import (
	"testing"

	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/validator"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestValidate_NoViolations(t *testing.T) {
	entries := makeEntries("PORT", "8080", "HOST", "localhost")
	rules := []validator.Rule{
		{Key: "PORT", Required: true, Pattern: `[0-9]+`},
		{Key: "HOST", Required: true},
	}
	violations := validator.Validate(entries, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	entries := makeEntries("HOST", "localhost")
	rules := []validator.Rule{
		{Key: "PORT", Required: true},
	}
	violations := validator.Validate(entries, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %q", violations[0].Key)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	entries := makeEntries("PORT", "not-a-number")
	rules := []validator.Rule{
		{Key: "PORT", Required: true, Pattern: `[0-9]+`},
	}
	violations := validator.Validate(entries, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %q", violations[0].Key)
	}
}

func TestValidate_OptionalKeyAbsent(t *testing.T) {
	entries := makeEntries("HOST", "localhost")
	rules := []validator.Rule{
		{Key: "DEBUG", Required: false, Pattern: `true|false`},
	}
	violations := validator.Validate(entries, rules)
	if len(violations) != 0 {
		t.Fatalf("optional absent key should not produce violations, got %v", violations)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	entries := makeEntries("PORT", "abc")
	rules := []validator.Rule{
		{Key: "PORT", Required: true, Pattern: `[0-9]+`},
		{Key: "HOST", Required: true},
	}
	violations := validator.Validate(entries, rules)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(violations), violations)
	}
}

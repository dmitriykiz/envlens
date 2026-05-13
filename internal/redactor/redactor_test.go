package redactor_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/redactor"
)

func makeEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{
			Key:   pairs[i],
			Value: pairs[i+1],
			Line:  i/2 + 1,
		})
	}
	return entries
}

func TestRedact_NonSensitiveUnchanged(t *testing.T) {
	entries := makeEntries("APP_ENV", "production", "PORT", "8080")
	result := redactor.Redact(entries)
	for _, me := range result {
		if me.Masked {
			t.Errorf("key %q should not be masked", me.Key)
		}
		if me.Key == "APP_ENV" && me.Value != "production" {
			t.Errorf("expected value 'production', got %q", me.Value)
		}
	}
}

func TestRedact_SensitiveKeysMasked(t *testing.T) {
	entries := makeEntries(
		"DB_PASSWORD", "supersecret",
		"API_KEY", "abc123",
		"AUTH_TOKEN", "tok",
	)
	result := redactor.Redact(entries)
	for _, me := range result {
		if !me.Masked {
			t.Errorf("key %q should be masked", me.Key)
		}
		if me.Value == "" && me.Key != "" {
			// empty original values stay empty — skip length check
			continue
		}
		for _, ch := range me.Value {
			if string(ch) != redactor.MaskChar {
				t.Errorf("masked value for %q contains non-mask char: %q", me.Key, me.Value)
			}
		}
	}
}

func TestRedact_MaskLengthCappedAt8(t *testing.T) {
	entries := makeEntries("SECRET_KEY", "averylongsecretvalue")
	result := redactor.Redact(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if len(result[0].Value) != 8 {
		t.Errorf("expected masked length 8, got %d", len(result[0].Value))
	}
}

func TestRedact_EmptyValue(t *testing.T) {
	entries := makeEntries("DB_PASSWORD", "")
	result := redactor.Redact(entries)
	if result[0].Value != "" {
		t.Errorf("expected empty masked value, got %q", result[0].Value)
	}
	if !result[0].Masked {
		t.Error("expected Masked=true for sensitive key")
	}
}

func TestRedact_PreservesLineNumbers(t *testing.T) {
	entries := makeEntries("PORT", "3000", "API_KEY", "key123")
	result := redactor.Redact(entries)
	if result[0].Line != 1 || result[1].Line != 2 {
		t.Errorf("line numbers not preserved: got %d and %d", result[0].Line, result[1].Line)
	}
}

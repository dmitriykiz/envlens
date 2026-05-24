package sanitizer

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(kvs ...string) []parser.Entry {
	if len(kvs)%2 != 0 {
		panic("makeEntries: need even number of args")
	}
	out := make([]parser.Entry, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1], Line: i/2 + 1})
	}
	return out
}

func TestSanitize_NoRules(t *testing.T) {
	entries := makeEntries("FOO", "  bar  ", "BAZ", "qux")
	out, results := Sanitize(entries, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no changes, got change for %s", r.Key)
		}
	}
}

func TestSanitize_TrimSpace(t *testing.T) {
	entries := makeEntries("FOO", "  hello world  ")
	rules := []Rule{{TrimSpace: true}}
	out, results := Sanitize(entries, rules)
	if out[0].Value != "hello world" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true")
	}
}

func TestSanitize_StripChars(t *testing.T) {
	entries := makeEntries("KEY", "abc123!@#")
	rules := []Rule{{StripChars: "!@#"}}
	out, _ := Sanitize(entries, rules)
	if out[0].Value != "abc123" {
		t.Errorf("unexpected value %q", out[0].Value)
	}
}

func TestSanitize_RegexpReplace(t *testing.T) {
	entries := makeEntries("DB_URL", "postgres://user:pass word@host/db")
	rules := []Rule{{ReplacePattern: `\s+`, Replacement: ""}}
	out, _ := Sanitize(entries, rules)
	if strings.Contains(out[0].Value, " ") {
		t.Errorf("spaces not removed: %q", out[0].Value)
	}
}

func TestSanitize_KeyPatternScopes(t *testing.T) {
	entries := makeEntries("DB_HOST", "  host  ", "APP_NAME", "  app  ")
	rules := []Rule{{KeyPattern: `^DB_`, TrimSpace: true}}
	out, results := Sanitize(entries, rules)
	if out[0].Value != "host" {
		t.Errorf("DB_HOST should be trimmed, got %q", out[0].Value)
	}
	if out[1].Value != "  app  " {
		t.Errorf("APP_NAME should be unchanged, got %q", out[1].Value)
	}
	if results[0].Changed != true {
		t.Error("DB_HOST should be marked changed")
	}
	if results[1].Changed != false {
		t.Error("APP_NAME should not be marked changed")
	}
}

func TestSanitize_CommentEntriesSkipped(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Raw: "# a comment"},
		{Key: "FOO", Value: "  bar  "},
	}
	rules := []Rule{{TrimSpace: true}}
	out, results := Sanitize(entries, rules)
	if out[0].Raw != "# a comment" {
		t.Error("comment entry should pass through unchanged")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result (non-comment), got %d", len(results))
	}
}

func TestWriteReport_Summary(t *testing.T) {
	results := []Result{
		{Key: "A", OldVal: "x ", NewVal: "x", Changed: true, File: "a.env", Line: 1},
		{Key: "B", OldVal: "y", NewVal: "y", Changed: false, File: "a.env", Line: 2},
	}
	var sb strings.Builder
	WriteReport(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "1 changed") {
		t.Errorf("expected '1 changed' in output, got:\n%s", out)
	}
	s := Summary(results)
	if !strings.Contains(s, "1 modified") {
		t.Errorf("unexpected summary: %s", s)
	}
}

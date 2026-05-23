package trimmer

import (
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	if len(pairs)%2 != 0 {
		panic("makeEntries: need even number of args (key, value)")
	}
	out := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestTrim_NoChanges(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	r := Trim(entries, DefaultOptions())
	if r.Changed != 0 {
		t.Errorf("expected 0 changes, got %d", r.Changed)
	}
}

func TestTrim_KeyWhitespace(t *testing.T) {
	entries := makeEntries("  HOST  ", "localhost")
	r := Trim(entries, Options{TrimKeys: true})
	if r.Entries[0].Key != "HOST" {
		t.Errorf("expected trimmed key 'HOST', got %q", r.Entries[0].Key)
	}
	if r.Changed != 1 {
		t.Errorf("expected 1 change, got %d", r.Changed)
	}
}

func TestTrim_ValueWhitespace(t *testing.T) {
	entries := makeEntries("HOST", "  localhost  ")
	r := Trim(entries, Options{TrimValues: true})
	if r.Entries[0].Value != "localhost" {
		t.Errorf("expected trimmed value 'localhost', got %q", r.Entries[0].Value)
	}
}

func TestTrim_QuotedValues(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"hello world"`, "hello world"},
		{`'single'`, "single"},
		{`noquotes`, "noquotes"},
		{`"mismatch'`, `"mismatch'`},
	}
	for _, tc := range tests {
		entries := makeEntries("KEY", tc.input)
		r := Trim(entries, Options{TrimQuotes: true})
		if r.Entries[0].Value != tc.want {
			t.Errorf("input %q: expected %q, got %q", tc.input, tc.want, r.Entries[0].Value)
		}
	}
}

func TestTrim_CommentsPassThrough(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Raw: "# this is a comment"},
		{Key: "  KEY  ", Value: "  val  "},
	}
	r := Trim(entries, DefaultOptions())
	if r.Entries[0].Raw != "# this is a comment" {
		t.Error("comment entry should be unchanged")
	}
	if r.Entries[1].Key != "KEY" || r.Entries[1].Value != "val" {
		t.Errorf("unexpected entry: %+v", r.Entries[1])
	}
}

func TestTrim_ChangedCountAccurate(t *testing.T) {
	entries := makeEntries(
		"  A  ", "  1  ",
		"B", "2",
		"  C", "'quoted'",
	)
	r := Trim(entries, DefaultOptions())
	if r.Changed != 3 {
		t.Errorf("expected 3 changes, got %d", r.Changed)
	}
}

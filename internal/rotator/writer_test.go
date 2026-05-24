package rotator

import (
	"strings"
	"testing"
)

func TestWriteReport_NoResults(t *testing.T) {
	var sb strings.Builder
	WriteReport(&sb, nil)
	if !strings.Contains(sb.String(), "no keys matched") {
		t.Errorf("expected 'no keys matched' message, got: %q", sb.String())
	}
}

func TestWriteReport_WithResults(t *testing.T) {
	results := []Result{
		{File: "a.env", Line: 3, Key: "API_KEY", OldValue: "old", NewValue: "REDACTED", Rotated: true},
		{File: "a.env", Line: 5, Key: "APP_ENV", OldValue: "prod", NewValue: "prod", Rotated: false},
	}
	var sb strings.Builder
	WriteReport(&sb, results)
	out := sb.String()
	for _, want := range []string{"API_KEY", "REDACTED", "APP_ENV", "prod", "yes", "no"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestSummary_AllRotated(t *testing.T) {
	results := []Result{
		{Rotated: true},
		{Rotated: true},
	}
	s := Summary(results)
	if !strings.Contains(s, "2 rotated") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_NoneRotated(t *testing.T) {
	results := []Result{
		{Rotated: false},
	}
	s := Summary(results)
	if !strings.Contains(s, "0 rotated") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("hello", 10); got != "hello" {
		t.Errorf("short string should not be truncated: %q", got)
	}
	long := strings.Repeat("x", 30)
	got := truncate(long, 24)
	if !strings.HasSuffix(got, "...") {
		t.Errorf("long string should end with '...': %q", got)
	}
	if len(got) != 27 {
		t.Errorf("expected length 27, got %d", len(got))
	}
}

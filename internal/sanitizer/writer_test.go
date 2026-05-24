package sanitizer

import (
	"strings"
	"testing"
)

func TestWriteReport_NoChanges(t *testing.T) {
	results := []Result{
		{Key: "FOO", OldVal: "bar", NewVal: "bar", Changed: false},
	}
	var sb strings.Builder
	WriteReport(&sb, results)
	if !strings.Contains(sb.String(), "No changes required") {
		t.Errorf("expected no-changes message, got: %s", sb.String())
	}
}

func TestWriteReport_WithChanges(t *testing.T) {
	results := []Result{
		{Key: "SECRET", OldVal: "abc ", NewVal: "abc", Changed: true, File: "prod.env", Line: 5},
	}
	var sb strings.Builder
	WriteReport(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected key name in report, got: %s", out)
	}
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected filename in report, got: %s", out)
	}
}

func TestSummary_NoneModified(t *testing.T) {
	results := []Result{
		{Key: "A", Changed: false},
		{Key: "B", Changed: false},
	}
	s := Summary(results)
	if !strings.Contains(s, "none modified") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_SomeModified(t *testing.T) {
	results := []Result{
		{Key: "A", Changed: true},
		{Key: "B", Changed: false},
		{Key: "C", Changed: true},
	}
	s := Summary(results)
	if !strings.Contains(s, "2 modified") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestPluralY(t *testing.T) {
	if pluralY(1) != "y" {
		t.Error("expected 'y' for n=1")
	}
	if pluralY(2) != "ies" {
		t.Error("expected 'ies' for n=2")
	}
}

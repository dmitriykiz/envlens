package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/analyzer"
	"github.com/yourorg/envlens/internal/reporter"
)

func TestWriteText_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	reports := []reporter.Report{
		{
			FilePath: ".env",
			Result:   analyzer.Result{},
		},
	}
	if err := reporter.WriteText(&buf, reports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No issues found") {
		t.Errorf("expected clean message, got: %s", out)
	}
}

func TestWriteText_WithIssues(t *testing.T) {
	var buf bytes.Buffer
	reports := []reporter.Report{
		{
			FilePath: "services/api/.env",
			Result: analyzer.Result{
				Duplicates:    []string{"PORT"},
				SensitiveKeys: []string{"AWS_SECRET_KEY"},
				MissingKeys:   []string{"DATABASE_URL"},
			},
		},
	}
	if err := reporter.WriteText(&buf, reports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"PORT", "AWS_SECRET_KEY", "DATABASE_URL", "DUPLICATE", "SENSITIVE", "MISSING"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteText_MultipleReports(t *testing.T) {
	var buf bytes.Buffer
	reports := []reporter.Report{
		{
			FilePath: ".env",
			Result:   analyzer.Result{},
		},
		{
			FilePath: "services/api/.env",
			Result: analyzer.Result{
				Duplicates: []string{"PORT"},
			},
		},
	}
	if err := reporter.WriteText(&buf, reports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, ".env") {
		t.Errorf("expected first file path in output, got: %s", out)
	}
	if !strings.Contains(out, "services/api/.env") {
		t.Errorf("expected second file path in output, got: %s", out)
	}
}

func TestSummary_OK(t *testing.T) {
	r := reporter.Report{FilePath: ".env", Result: analyzer.Result{}}
	got := reporter.Summary(r)
	if !strings.Contains(got, "OK") {
		t.Errorf("expected OK summary, got: %s", got)
	}
}

func TestSummary_WithIssues(t *testing.T) {
	r := reporter.Report{
		FilePath: ".env",
		Result: analyzer.Result{
			Duplicates:  []string{"A", "B"},
			MissingKeys: []string{"C"},
		},
	}
	got := reporter.Summary(r)
	if !strings.Contains(got, "2 duplicate(s)") {
		t.Errorf("expected duplicate count in summary, got: %s", got)
	}
	if !strings.Contains(got, "1 missing") {
		t.Errorf("expected missing count in summary, got: %s", got)
	}
}

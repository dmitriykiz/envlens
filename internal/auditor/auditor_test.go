package auditor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envlens/internal/auditor"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestAudit_NoIssues(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=envlens\nLOG_LEVEL=debug\n")
	r, err := auditor.Audit([]string{path}, auditor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Findings) != 0 {
		t.Errorf("expected no findings, got %d", len(r.Findings))
	}
}

func TestAudit_DetectsEmptyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=\nLOG_LEVEL=debug\n")
	opts := auditor.DefaultOptions()
	r, err := auditor.Audit([]string{path}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	counts := auditor.CountBySeverity(r.Findings)
	if counts[auditor.SeverityWarning] < 1 {
		t.Errorf("expected at least one WARNING finding for empty value")
	}
}

func TestAudit_DetectsSensitivePlainText(t *testing.T) {
	path := writeTempEnv(t, "DB_PASSWORD=supersecret\n")
	r, err := auditor.Audit([]string{path}, auditor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	critical := auditor.FilterBySeverity(r.Findings, auditor.SeverityCritical)
	if len(critical) == 0 {
		t.Error("expected a CRITICAL finding for plain-text sensitive key")
	}
}

func TestAudit_DetectsKeyFormatViolation(t *testing.T) {
	path := writeTempEnv(t, "my-key=value\n")
	r, err := auditor.Audit([]string{path}, auditor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := auditor.FilterBySeverity(r.Findings, auditor.SeverityInfo)
	if len(info) == 0 {
		t.Error("expected an INFO finding for invalid key format")
	}
}

func TestAudit_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\n\nVALID_KEY=value\n")
	r, err := auditor.Audit([]string{path}, auditor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Findings) != 0 {
		t.Errorf("expected no findings for comment/blank lines, got %d", len(r.Findings))
	}
}

func TestWriteReport_ContainsSeverity(t *testing.T) {
	path := writeTempEnv(t, "DB_PASSWORD=exposed\nEMPTY=\n")
	r, err := auditor.Audit([]string{path}, auditor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var sb strings.Builder
	auditor.WriteReport(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "CRITICAL") {
		t.Error("report output should contain CRITICAL")
	}
	if !strings.Contains(out, "WARNING") {
		t.Error("report output should contain WARNING")
	}
}

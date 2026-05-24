package injector_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envlens/internal/injector"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestInject_SetsVariables(t *testing.T) {
	path := writeTempEnv(t, "INJECT_FOO=bar\nINJECT_BAZ=qux\n")
	t.Setenv("INJECT_FOO", "") // ensure clean state handled by t.Setenv cleanup

	results, err := injector.Inject(path, injector.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if got := os.Getenv("INJECT_BAZ"); got != "qux" {
		t.Errorf("INJECT_BAZ = %q, want %q", got, "qux")
	}
}

func TestInject_StrategySkip(t *testing.T) {
	t.Setenv("SKIP_KEY", "original")
	path := writeTempEnv(t, "SKIP_KEY=replaced\n")

	results, err := injector.Inject(path, injector.Options{Strategy: injector.StrategySkip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Error("expected result to be marked Skipped")
	}
	if got := os.Getenv("SKIP_KEY"); got != "original" {
		t.Errorf("SKIP_KEY = %q, want %q", got, "original")
	}
}

func TestInject_DryRun_DoesNotSetEnv(t *testing.T) {
	const key = "DRYRUN_KEY"
	os.Unsetenv(key)
	path := writeTempEnv(t, key+"=should_not_be_set\n")

	_, err := injector.Inject(path, injector.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := os.LookupEnv(key); ok {
		t.Errorf("DryRun should not have set %s", key)
	}
}

func TestInject_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nCOMMENT_TEST=yes\n")

	results, err := injector.Inject(path, injector.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestInject_FileNotFound(t *testing.T) {
	_, err := injector.Inject("/nonexistent/path/.env", injector.Options{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWriteReport_Output(t *testing.T) {
	results := []injector.Result{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux", Skipped: true},
	}
	var sb strings.Builder
	injector.WriteReport(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "[SET]") {
		t.Error("expected [SET] in report output")
	}
	if !strings.Contains(out, "[SKIP]") {
		t.Error("expected [SKIP] in report output")
	}
	if !strings.Contains(out, "Injected: 1") {
		t.Error("expected summary line in report output")
	}
}

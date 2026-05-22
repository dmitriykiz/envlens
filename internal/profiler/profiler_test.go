package profiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlens/internal/profiler"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestBuild_BasicCounts(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "APP_NAME=myapp\nDB_PASSWORD=secret\nEMPTY=\n")

	prof, err := profiler.Build([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prof.Files) != 1 {
		t.Fatalf("expected 1 file profile, got %d", len(prof.Files))
	}
	fp := prof.Files[0]
	if fp.TotalKeys != 3 {
		t.Errorf("TotalKeys: want 3, got %d", fp.TotalKeys)
	}
	if fp.SensitiveKeys != 1 {
		t.Errorf("SensitiveKeys: want 1, got %d", fp.SensitiveKeys)
	}
	if fp.EmptyValues != 1 {
		t.Errorf("EmptyValues: want 1, got %d", fp.EmptyValues)
	}
}

func TestBuild_UniqueKeysAcrossFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempEnv(t, dir, ".env.a", "KEY_A=1\nSHARED=x\n")
	p2 := writeTempEnv(t, dir, ".env.b", "KEY_B=2\nSHARED=y\n")

	prof, err := profiler.Build([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prof.UniqueKeys != 3 {
		t.Errorf("UniqueKeys: want 3, got %d", prof.UniqueKeys)
	}
	if prof.TotalKeys != 4 {
		t.Errorf("TotalKeys: want 4, got %d", prof.TotalKeys)
	}
}

func TestBuild_SensitivePatterns(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env",
		"API_KEY=abc\nAUTH_TOKEN=xyz\nPLAIN=val\n")

	prof, err := profiler.Build([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prof.SensitiveKeys != 2 {
		t.Errorf("SensitiveKeys: want 2, got %d", prof.SensitiveKeys)
	}
}

func TestBuild_EmptyFileList(t *testing.T) {
	prof, err := profiler.Build([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prof.TotalKeys != 0 || prof.UniqueKeys != 0 {
		t.Errorf("expected zero counts for empty input")
	}
}

func TestBuild_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "# comment\nFOO=bar\n")

	prof, err := profiler.Build([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prof.TotalKeys != 1 {
		t.Errorf("TotalKeys: want 1, got %d", prof.TotalKeys)
	}
}

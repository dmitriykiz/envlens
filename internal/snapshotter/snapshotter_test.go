package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envlens/internal/snapshotter"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestTake_BasicSnapshot(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "APP=hello\nDEBUG=true\n")

	snap, errs := snapshotter.Take("test", []string{p})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if snap.Label != "test" {
		t.Errorf("label = %q, want %q", snap.Label, "test")
	}
	if len(snap.Files) != 1 {
		t.Errorf("files count = %d, want 1", len(snap.Files))
	}
}

func TestTake_SkipsMissingFile(t *testing.T) {
	_, errs := snapshotter.Take("x", []string{"/no/such/file.env"})
	if len(errs) == 0 {
		t.Error("expected error for missing file, got none")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=val\n")

	snap, _ := snapshotter.Take("rt", []string{p})
	dest := filepath.Join(dir, "snap.json")
	if err := snapshotter.Save(snap, dest); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := snapshotter.Load(dest)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Label != snap.Label {
		t.Errorf("label mismatch: got %q want %q", loaded.Label, snap.Label)
	}
	if len(loaded.Files) != len(snap.Files) {
		t.Errorf("files count mismatch")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(bad, []byte("not json{"), 0o644)
	_, err := snapshotter.Load(bad)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

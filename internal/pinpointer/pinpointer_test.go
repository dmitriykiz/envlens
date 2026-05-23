package pinpointer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/pinpointer"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestPinpoint_FindsAllKeys(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=envlens\n")

	locs, err := pinpointer.Pinpoint([]string{f}, pinpointer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) != 3 {
		t.Fatalf("expected 3 locations, got %d", len(locs))
	}
}

func TestPinpoint_FiltersByKeyPattern(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=envlens\n")

	locs, err := pinpointer.Pinpoint([]string{f}, pinpointer.Options{KeyPattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) != 2 {
		t.Fatalf("expected 2 DB_ locations, got %d", len(locs))
	}
	for _, l := range locs {
		if !strings.HasPrefix(l.Key, "DB_") {
			t.Errorf("unexpected key %q", l.Key)
		}
	}
}

func TestPinpoint_CaseInsensitivePattern(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "SECRET_KEY=abc\nAPI_TOKEN=xyz\n")

	locs, err := pinpointer.Pinpoint([]string{f}, pinpointer.Options{KeyPattern: "secret", CaseSensitive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) != 1 || locs[0].Key != "SECRET_KEY" {
		t.Fatalf("expected SECRET_KEY match, got %+v", locs)
	}
}

func TestPinpoint_InvalidPattern(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "KEY=val\n")
	_, err := pinpointer.Pinpoint([]string{f}, pinpointer.Options{KeyPattern: "["})
	if err == nil {
		t.Fatal("expected error for invalid regexp, got nil")
	}
}

func TestPinpoint_LineNumbers(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "# comment\nFIRST=a\n\nSECOND=b\n")

	locs, err := pinpointer.Pinpoint([]string{f}, pinpointer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locs))
	}
	if locs[0].Line != 2 {
		t.Errorf("FIRST should be on line 2, got %d", locs[0].Line)
	}
	if locs[1].Line != 4 {
		t.Errorf("SECOND should be on line 4, got %d", locs[1].Line)
	}
}

package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp env file: %v", err)
	}
	return p
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" || entries[0].Value != "localhost" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	content := "# comment\n\nAPP_ENV=production\n"
	path := writeTempEnv(t, content)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "APP_ENV" {
		t.Errorf("unexpected key: %s", entries[0].Key)
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"\nTOKEN='bearer-token'\n`)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if len(e.Value) > 0 && (e.Value[0] == '"' || e.Value[0] == '\'') {
			t.Errorf("quotes not stripped for key %s: %q", e.Key, e.Value)
		}
	}
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "NODEQUALS\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for line without '=', got nil")
	}
}

func TestParseFile_LineNumbers(t *testing.T) {
	content := "# header\nFIRST=1\n\nSECOND=2\n"
	path := writeTempEnv(t, content)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Line != 2 {
		t.Errorf("expected FIRST on line 2, got %d", entries[0].Line)
	}
	if entries[1].Line != 4 {
		t.Errorf("expected SECOND on line 4, got %d", entries[1].Line)
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

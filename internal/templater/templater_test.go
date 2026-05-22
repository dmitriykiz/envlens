package templater_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/templater"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestGenerate_SensitiveValuesReplaced(t *testing.T) {
	entries := makeEntries("DB_PASSWORD", "s3cr3t", "APP_NAME", "myapp")
	out := templater.Generate(entries, templater.Options{KeepNonSensitive: false})

	if strings.Contains(out, "s3cr3t") {
		t.Error("sensitive value should have been replaced")
	}
	if !strings.Contains(out, "DB_PASSWORD=<your_value_here>") {
		t.Error("expected placeholder for DB_PASSWORD")
	}
}

func TestGenerate_KeepNonSensitive(t *testing.T) {
	entries := makeEntries("API_KEY", "abc123", "PORT", "8080")
	out := templater.Generate(entries, templater.Options{KeepNonSensitive: true})

	if !strings.Contains(out, "PORT=8080") {
		t.Error("non-sensitive value should be preserved when KeepNonSensitive=true")
	}
	if strings.Contains(out, "abc123") {
		t.Error("sensitive value must not appear in output")
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	entries := makeEntries("SECRET_TOKEN", "tok")
	out := templater.Generate(entries, templater.Options{Placeholder: "CHANGEME"})

	if !strings.Contains(out, "CHANGEME") {
		t.Error("custom placeholder should appear in output")
	}
}

func TestGenerate_CommentPassthrough(t *testing.T) {
	entries := []parser.Entry{
		{Comment: "# Database settings"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	out := templater.Generate(entries, templater.Options{KeepNonSensitive: true})

	if !strings.Contains(out, "# Database settings") {
		t.Error("comments should pass through unchanged")
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	entries := makeEntries("APP_ENV", "production", "DB_PASSWORD", "secret")
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env.example")

	if err := templater.WriteFile(dest, entries, templater.Options{}); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}
	if !strings.Contains(string(data), "APP_ENV=") {
		t.Error("expected APP_ENV key in written file")
	}
}

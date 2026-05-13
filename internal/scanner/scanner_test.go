package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/scanner"
)

// buildTree creates a temporary directory tree with the given files.
// files is a map of relative path -> file content.
func buildTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
	return root
}

func TestScanDir_FindsEnvFiles(t *testing.T) {
	root := buildTree(t, map[string]string{
		".env":          "FOO=bar\nBAZ=qux\n",
		"services/a/.env": "PORT=3000\n",
		"services/b/.env": "PORT=4000\nDEBUG=true\n",
	})

	results, err := scanner.ScanDir(root, scanner.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("parse error for %s: %v", r.Path, r.Err)
		}
	}
}

func TestScanDir_ExcludesDirs(t *testing.T) {
	root := buildTree(t, map[string]string{
		".env":              "FOO=bar\n",
		"node_modules/.env": "SECRET=x\n",
	})

	results, err := scanner.ScanDir(root, scanner.Options{
		ExcludeDirs: []string{"node_modules"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestScanDir_CustomPatterns(t *testing.T) {
	root := buildTree(t, map[string]string{
		".env":          "A=1\n",
		".env.local":    "B=2\n",
		".env.test":     "C=3\n",
		"other.txt":     "not an env file\n",
	})

	results, err := scanner.ScanDir(root, scanner.Options{
		Patterns: []string{".env", ".env.*"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestScanDir_EmptyDir(t *testing.T) {
	root := t.TempDir()
	results, err := scanner.ScanDir(root, scanner.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

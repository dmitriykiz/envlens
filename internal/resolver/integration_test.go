package resolver_test

import (
	"os"
	"testing"

	"github.com/envlens/internal/parser"
	"github.com/envlens/internal/resolver"
)

// TestResolve_FallsBackToOsEnv verifies that variables not present in the
// entry set are resolved from the host environment.
func TestResolve_FallsBackToOsEnv(t *testing.T) {
	const envKey = "ENVLENS_TEST_FALLBACK"
	t.Setenv(envKey, "from_os")

	entries := []parser.Entry{
		{Key: "VALUE", Value: "${" + envKey + "}"},
	}
	res, err := resolver.Resolve(entries, resolver.StrategyStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "from_os" {
		t.Errorf("got %q, want %q", res.Entries[0].Value, "from_os")
	}
}

// TestResolve_RoundTrip parses a real temp file and resolves its entries.
func TestResolve_RoundTrip(t *testing.T) {
	content := []byte("PROTO=https\nHOST=api.example.com\nBASE_URL=${PROTO}://${HOST}\n")
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	entries, err := parser.ParseFile(f.Name())
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	res, err := resolver.Resolve(entries, resolver.StrategyStrict)
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	var baseURL string
	for _, e := range res.Entries {
		if e.Key == "BASE_URL" {
			baseURL = e.Value
		}
	}
	want := "https://api.example.com"
	if baseURL != want {
		t.Errorf("BASE_URL: got %q, want %q", baseURL, want)
	}
}

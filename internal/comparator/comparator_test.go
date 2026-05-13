package comparator_test

import (
	"testing"

	"github.com/user/envlens/internal/comparator"
	"github.com/user/envlens/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompare_AllCommonNoMismatch(t *testing.T) {
	files := comparator.FileSet{
		"a.env": entries("HOST", "localhost", "PORT", "8080"),
		"b.env": entries("HOST", "localhost", "PORT", "8080"),
	}
	report := comparator.Compare(files)

	if len(report.KeysOnlyIn) != 0 {
		t.Errorf("expected no exclusive keys, got %v", report.KeysOnlyIn)
	}
	if len(report.ValueMismatches) != 0 {
		t.Errorf("expected no mismatches, got %v", report.ValueMismatches)
	}
	if len(report.CommonKeys) != 2 {
		t.Errorf("expected 2 common keys, got %d", len(report.CommonKeys))
	}
}

func TestCompare_DetectsValueMismatch(t *testing.T) {
	files := comparator.FileSet{
		"prod.env": entries("DB_HOST", "prod-db", "PORT", "5432"),
		"dev.env":  entries("DB_HOST", "localhost", "PORT", "5432"),
	}
	report := comparator.Compare(files)

	if _, ok := report.ValueMismatches["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be flagged as a value mismatch")
	}
	if _, ok := report.ValueMismatches["PORT"]; ok {
		t.Error("PORT values are equal; should not be flagged")
	}
}

func TestCompare_DetectsExclusiveKeys(t *testing.T) {
	files := comparator.FileSet{
		"a.env": entries("SHARED", "x", "ONLY_A", "1"),
		"b.env": entries("SHARED", "x", "ONLY_B", "2"),
	}
	report := comparator.Compare(files)

	if keys, ok := report.KeysOnlyIn["a.env"]; !ok || len(keys) != 1 || keys[0] != "ONLY_A" {
		t.Errorf("expected ONLY_A exclusive to a.env, got %v", report.KeysOnlyIn)
	}
	if keys, ok := report.KeysOnlyIn["b.env"]; !ok || len(keys) != 1 || keys[0] != "ONLY_B" {
		t.Errorf("expected ONLY_B exclusive to b.env, got %v", report.KeysOnlyIn)
	}
	if len(report.CommonKeys) != 1 || report.CommonKeys[0] != "SHARED" {
		t.Errorf("expected SHARED as the only common key, got %v", report.CommonKeys)
	}
}

func TestCompare_EmptyFileSet(t *testing.T) {
	report := comparator.Compare(comparator.FileSet{})

	if len(report.CommonKeys) != 0 {
		t.Errorf("expected no common keys for empty set")
	}
	if len(report.KeysOnlyIn) != 0 {
		t.Errorf("expected no exclusive keys for empty set")
	}
}

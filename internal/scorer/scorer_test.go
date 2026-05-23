package scorer

import (
	"testing"

	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestScore_PerfectFile(t *testing.T) {
	entries := makeEntries(
		"APP_NAME", "envlens",
		"LOG_LEVEL", "info",
		"PORT", "8080",
	)
	r := Score(".env", entries)
	if r.Score != 100 {
		t.Errorf("expected score 100, got %d", r.Score)
	}
	if r.Grade != GradeA {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
	if len(r.Penalties) != 0 {
		t.Errorf("expected no penalties, got %v", r.Penalties)
	}
}

func TestScore_EmptyValues(t *testing.T) {
	entries := makeEntries("APP_NAME", "", "PORT", "")
	r := Score(".env", entries)
	// 2 empty values × 5 = 10 points deducted
	if r.Score != 90 {
		t.Errorf("expected score 90, got %d", r.Score)
	}
	if r.Grade != GradeA {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
}

func TestScore_LowercaseKeys(t *testing.T) {
	entries := makeEntries("app_name", "x", "port", "8080")
	r := Score(".env", entries)
	// 2 lowercase × 4 = 8 deducted → 92
	if r.Score != 92 {
		t.Errorf("expected score 92, got %d", r.Score)
	}
}

func TestScore_DuplicateKeys(t *testing.T) {
	entries := []parser.Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "PORT", Value: "9090"},
	}
	r := Score(".env", entries)
	// 1 duplicate key × 10 = 10 deducted → 90
	if r.Score != 90 {
		t.Errorf("expected score 90, got %d", r.Score)
	}
}

func TestScore_SensitivePlaintext(t *testing.T) {
	entries := makeEntries("DB_PASSWORD", "hunter2", "API_KEY", "abc123")
	r := Score(".env", entries)
	// 2 sensitive × 8 = 16 deducted → 84
	if r.Score != 84 {
		t.Errorf("expected score 84, got %d", r.Score)
	}
	if r.Grade != GradeB {
		t.Errorf("expected grade B, got %s", r.Grade)
	}
}

func TestScore_ClampedToZero(t *testing.T) {
	var entries []parser.Entry
	for i := 0; i < 20; i++ {
		entries = append(entries, parser.Entry{Key: "db_secret", Value: "plaintext"})
	}
	r := Score(".env", entries)
	if r.Score != 0 {
		t.Errorf("expected score clamped to 0, got %d", r.Score)
	}
	if r.Grade != GradeF {
		t.Errorf("expected grade F, got %s", r.Grade)
	}
}

func TestScore_CommentsIgnored(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Raw: "# this is a comment"},
		{Key: "APP_ENV", Value: "production"},
	}
	r := Score(".env", entries)
	if r.Score != 100 {
		t.Errorf("expected score 100, got %d", r.Score)
	}
}

func TestScore_EmptyFile(t *testing.T) {
	r := Score(".env", []parser.Entry{})
	// An empty file has no violations, so it should score 100.
	if r.Score != 100 {
		t.Errorf("expected score 100 for empty file, got %d", r.Score)
	}
	if r.Grade != GradeA {
		t.Errorf("expected grade A for empty file, got %s", r.Grade)
	}
	if len(r.Penalties) != 0 {
		t.Errorf("expected no penalties for empty file, got %v", r.Penalties)
	}
}

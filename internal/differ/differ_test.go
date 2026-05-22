package differ_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/differ"
	"github.com/yourorg/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDiff_NoDifferences(t *testing.T) {
	base := makeEntries("APP_ENV", "production", "PORT", "8080")
	next := makeEntries("APP_ENV", "production", "PORT", "8080")

	result := differ.Diff(base, next)

	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Errorf("expected no diff, got added=%d removed=%d changed=%d",
			len(result.Added), len(result.Removed), len(result.Changed))
	}
}

func TestDiff_DetectsAddedKeys(t *testing.T) {
	base := makeEntries("APP_ENV", "production")
	next := makeEntries("APP_ENV", "production", "NEW_KEY", "value")

	result := differ.Diff(base, next)

	if len(result.Added) != 1 || result.Added[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in Added, got %+v", result.Added)
	}
}

func TestDiff_DetectsRemovedKeys(t *testing.T) {
	base := makeEntries("APP_ENV", "production", "OLD_KEY", "gone")
	next := makeEntries("APP_ENV", "production")

	result := differ.Diff(base, next)

	if len(result.Removed) != 1 || result.Removed[0].Key != "OLD_KEY" {
		t.Errorf("expected OLD_KEY in Removed, got %+v", result.Removed)
	}
}

func TestDiff_DetectsChangedValues(t *testing.T) {
	base := makeEntries("PORT", "8080")
	next := makeEntries("PORT", "9090")

	result := differ.Diff(base, next)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(result.Changed))
	}
	c := result.Changed[0]
	if c.Key != "PORT" || c.OldValue != "8080" || c.NewValue != "9090" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_EmptyBaseAndNext(t *testing.T) {
	result := differ.Diff(nil, nil)

	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Errorf("expected empty diff for nil inputs")
	}
}

func TestDiff_AllOperationsTogether(t *testing.T) {
	base := makeEntries("KEEP", "same", "CHANGE", "old", "REMOVE", "bye")
	next := makeEntries("KEEP", "same", "CHANGE", "new", "ADD", "hello")

	result := differ.Diff(base, next)

	if len(result.Added) != 1 || result.Added[0].Key != "ADD" {
		t.Errorf("unexpected Added: %+v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0].Key != "REMOVE" {
		t.Errorf("unexpected Removed: %+v", result.Removed)
	}
	if len(result.Changed) != 1 || result.Changed[0].Key != "CHANGE" {
		t.Errorf("unexpected Changed: %+v", result.Changed)
	}
}

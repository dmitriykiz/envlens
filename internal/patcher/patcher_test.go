package patcher_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/patcher"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestPatch_NoRules(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	res, err := patcher.Patch(entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Added)+len(res.Updated)+len(res.Deleted) != 0 {
		t.Error("expected no changes")
	}
}

func TestPatch_AddsNewKey(t *testing.T) {
	entries := makeEntries("HOST", "localhost")
	rules := []patcher.Rule{{Op: patcher.OpSet, Key: "PORT", Value: "8080"}}
	res, err := patcher.Patch(entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "PORT" {
		t.Errorf("expected PORT in Added, got %v", res.Added)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestPatch_UpdatesExistingKey(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	rules := []patcher.Rule{{Op: patcher.OpSet, Key: "PORT", Value: "9999"}}
	res, err := patcher.Patch(entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "PORT" {
		t.Errorf("expected PORT in Updated, got %v", res.Updated)
	}
	if res.Entries[1].Value != "9999" {
		t.Errorf("expected value 9999, got %s", res.Entries[1].Value)
	}
}

func TestPatch_DeletesExistingKey(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	rules := []patcher.Rule{{Op: patcher.OpDelete, Key: "PORT"}}
	res, err := patcher.Patch(entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Deleted) != 1 || res.Deleted[0] != "PORT" {
		t.Errorf("expected PORT in Deleted, got %v", res.Deleted)
	}
	if len(res.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(res.Entries))
	}
}

func TestPatch_DeleteMissingKeyIsSkipped(t *testing.T) {
	entries := makeEntries("HOST", "localhost")
	rules := []patcher.Rule{{Op: patcher.OpDelete, Key: "MISSING"}}
	res, err := patcher.Patch(entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in Skipped, got %v", res.Skipped)
	}
}

func TestPatch_InvalidOpReturnsError(t *testing.T) {
	entries := makeEntries("HOST", "localhost")
	rules := []patcher.Rule{{Op: "upsert", Key: "HOST", Value: "remotehost"}}
	_, err := patcher.Patch(entries, rules)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatch_EmptyKeyReturnsError(t *testing.T) {
	_, err := patcher.Patch(nil, []patcher.Rule{{Op: patcher.OpSet, Key: "", Value: "v"}})
	if err == nil {
		t.Error("expected error for empty key")
	}
}

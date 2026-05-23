package sorter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/sorter"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func keys(entries []parser.Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Key
	}
	return out
}

func TestSort_ByKeyAscending(t *testing.T) {
	in := makeEntries("ZEBRA", "z", "ALPHA", "a", "MANGO", "m")
	out := sorter.Sort(in, sorter.Options{By: sorter.ByKey, Order: sorter.Ascending})
	want := []string{"ALPHA", "MANGO", "ZEBRA"}
	got := keys(out)
	for i, k := range want {
		if got[i] != k {
			t.Errorf("position %d: want %s, got %s", i, k, got[i])
		}
	}
}

func TestSort_ByKeyDescending(t *testing.T) {
	in := makeEntries("ZEBRA", "z", "ALPHA", "a", "MANGO", "m")
	out := sorter.Sort(in, sorter.Options{By: sorter.ByKey, Order: sorter.Descending})
	want := []string{"ZEBRA", "MANGO", "ALPHA"}
	got := keys(out)
	for i, k := range want {
		if got[i] != k {
			t.Errorf("position %d: want %s, got %s", i, k, got[i])
		}
	}
}

func TestSort_ByValue(t *testing.T) {
	in := makeEntries("B", "banana", "A", "apple", "C", "cherry")
	out := sorter.Sort(in, sorter.Options{By: sorter.ByValue, Order: sorter.Ascending})
	if out[0].Value != "apple" || out[1].Value != "banana" || out[2].Value != "cherry" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestSort_GroupSensitiveFirst(t *testing.T) {
	in := makeEntries(
		"APP_NAME", "myapp",
		"DB_PASSWORD", "secret",
		"PORT", "8080",
		"API_KEY", "key123",
	)
	out := sorter.Sort(in, sorter.Options{By: sorter.ByKey, Order: sorter.Ascending, GroupSensitive: true})
	// First two should be sensitive
	for i := 0; i < 2; i++ {
		k := out[i].Key
		if k != "DB_PASSWORD" && k != "API_KEY" {
			t.Errorf("expected sensitive key at position %d, got %s", i, k)
		}
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	in := makeEntries("Z", "z", "A", "a")
	origFirst := in[0].Key
	sorter.Sort(in, sorter.Options{By: sorter.ByKey, Order: sorter.Ascending})
	if in[0].Key != origFirst {
		t.Error("Sort mutated the original slice")
	}
}

func TestGroupByPrefix(t *testing.T) {
	in := makeEntries(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_NAME", "envlens",
		"PORT", "8080",
	)
	groups := sorter.GroupByPrefix(in, "_")
	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 1 {
		t.Errorf("expected 1 APP entry, got %d", len(groups["APP"]))
	}
	if len(groups[""]) != 1 {
		t.Errorf("expected 1 ungrouped entry, got %d", len(groups[""]))
	}
}

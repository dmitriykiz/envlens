package pinner_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/envlens/internal/parser"
	"github.com/envlens/internal/pinner"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCreate_BuildsCorrectKeys(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "API_KEY", "secret")
	pin := pinner.Create("app.env", entries)

	if pin.File != "app.env" {
		t.Errorf("expected file app.env, got %s", pin.File)
	}
	if len(pin.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(pin.Keys))
	}
	if pin.Digest == "" {
		t.Error("expected non-empty digest")
	}
}

func TestCreate_SkipsComments(t *testing.T) {
	entries := []parser.Entry{
		{IsComment: true, Raw: "# comment"},
		{Key: "FOO", Value: "bar"},
	}
	pin := pinner.Create("x.env", entries)
	if len(pin.Keys) != 1 || pin.Keys[0] != "FOO" {
		t.Errorf("expected only FOO, got %v", pin.Keys)
	}
}

func TestDetect_NoDrift(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	pin := pinner.Create("db.env", entries)
	report := pinner.Detect(pin, entries)

	if len(report.Entries) != 0 {
		t.Errorf("expected no drift, got %v", report.Entries)
	}
}

func TestDetect_DetectsAdded(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	pin := pinner.Create("a.env", base)

	current := makeEntries("HOST", "localhost", "NEW_KEY", "value")
	report := pinner.Detect(pin, current)

	if len(report.Entries) != 1 || report.Entries[0].Key != "NEW_KEY" || report.Entries[0].Reason != "added" {
		t.Errorf("unexpected drift entries: %v", report.Entries)
	}
}

func TestDetect_DetectsRemoved(t *testing.T) {
	base := makeEntries("HOST", "localhost", "PORT", "5432")
	pin := pinner.Create("a.env", base)

	current := makeEntries("HOST", "localhost")
	report := pinner.Detect(pin, current)

	if len(report.Entries) != 1 || report.Entries[0].Reason != "removed" {
		t.Errorf("unexpected entries: %v", report.Entries)
	}
}

func TestDetect_DetectsChanged(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	pin := pinner.Create("a.env", base)

	current := makeEntries("HOST", "remotehost")
	report := pinner.Detect(pin, current)

	if len(report.Entries) != 1 || report.Entries[0].Reason != "changed" {
		t.Errorf("unexpected entries: %v", report.Entries)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	entries := makeEntries("KEY", "val")
	pin := pinner.Create("test.env", entries)

	tmp := filepath.Join(t.TempDir(), "pin.json")
	if err := pinner.SavePin(tmp, pin); err != nil {
		t.Fatalf("SavePin: %v", err)
	}

	loaded, err := pinner.LoadPin(tmp)
	if err != nil {
		t.Fatalf("LoadPin: %v", err)
	}
	if loaded.Digest != pin.Digest {
		t.Errorf("digest mismatch: want %s got %s", pin.Digest, loaded.Digest)
	}
}

func TestLoadPin_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not json"), 0644)
	_, err := pinner.LoadPin(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestWriteReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	pinner.WriteReport(&buf, []pinner.DriftReport{{File: "a.env", Entries: nil}})
	if !bytes.Contains(buf.Bytes(), []byte("no drift")) {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteReport_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	reports := []pinner.DriftReport{
		{
			File: "app.env",
			Entries: []pinner.DriftEntry{
				{Key: "SECRET", Reason: "changed"},
				{Key: "OLD_KEY", Reason: "removed"},
			},
		},
	}
	pinner.WriteReport(&buf, reports)
	out := buf.String()
	for _, want := range []string{"app.env", "SECRET", "changed", "OLD_KEY", "removed"} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
	_ = json.Valid // suppress unused import
}

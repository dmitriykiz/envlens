package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/exporter"
	"github.com/yourorg/envlens/internal/parser"
)

func makeFileEntries() map[string][]parser.Entry {
	return map[string][]parser.Entry{
		"service/.env": {
			{Key: "APP_NAME", Value: "envlens", Line: 1},
			{Key: "PORT", Value: "8080", Line: 2},
		},
		"worker/.env": {
			{Key: "QUEUE", Value: "default", Line: 1},
		},
	}
}

func TestExportJSON_ContainsAllEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.ExportJSON(&buf, makeFileEntries()); err != nil {
		t.Fatalf("ExportJSON returned error: %v", err)
	}

	var records []exporter.Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	keys := make(map[string]bool)
	for _, r := range records {
		keys[r.Key] = true
		if r.File == "" {
			t.Errorf("record %q has empty File field", r.Key)
		}
		if r.Line == 0 {
			t.Errorf("record %q has zero Line field", r.Key)
		}
	}

	for _, want := range []string{"APP_NAME", "PORT", "QUEUE"} {
		if !keys[want] {
			t.Errorf("expected key %q in JSON output", want)
		}
	}
}

func TestExportCSV_HeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.ExportCSV(&buf, makeFileEntries()); err != nil {
		t.Fatalf("ExportCSV returned error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least a header line")
	}

	header := lines[0]
	if !strings.HasPrefix(header, "file,key,value,line") {
		t.Errorf("unexpected CSV header: %q", header)
	}

	// header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (1 header + 3 rows), got %d", len(lines))
	}
}

func TestExportJSON_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.ExportJSON(&buf, map[string][]parser.Entry{}); err != nil {
		t.Fatalf("ExportJSON returned error on empty input: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty JSON output even for empty input")
	}
}

func TestExportCSV_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.ExportCSV(&buf, map[string][]parser.Entry{}); err != nil {
		t.Fatalf("ExportCSV returned error on empty input: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line for empty input, got %d lines", len(lines))
	}
}

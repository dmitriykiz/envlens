// Package exporter provides functionality to export parsed .env entries
// into various output formats such as JSON and CSV.
package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/envlens/internal/parser"
)

// Format represents the output format for exported data.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Record is a flat representation of a parsed env entry suitable for export.
type Record struct {
	File  string `json:"file"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Line  int    `json:"line"`
}

// ExportJSON writes all entries from the provided file-keyed map as a JSON
// array to w. Each entry is annotated with its source file path.
func ExportJSON(w io.Writer, fileEntries map[string][]parser.Entry) error {
	records := flatten(fileEntries)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(records); err != nil {
		return fmt.Errorf("exporter: json encode: %w", err)
	}
	return nil
}

// ExportCSV writes all entries from the provided file-keyed map as CSV rows
// to w. The header row contains: file, key, value, line.
func ExportCSV(w io.Writer, fileEntries map[string][]parser.Entry) error {
	records := flatten(fileEntries)
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"file", "key", "value", "line"}); err != nil {
		return fmt.Errorf("exporter: csv write header: %w", err)
	}

	for _, r := range records {
		row := []string{r.File, r.Key, r.Value, fmt.Sprintf("%d", r.Line)}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("exporter: csv write row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// flatten converts a map of file → entries into an ordered slice of Records.
func flatten(fileEntries map[string][]parser.Entry) []Record {
	var records []Record
	for file, entries := range fileEntries {
		for _, e := range entries {
			records = append(records, Record{
				File:  file,
				Key:   e.Key,
				Value: e.Value,
				Line:  e.Line,
			})
		}
	}
	return records
}

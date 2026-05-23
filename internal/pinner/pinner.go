// Package pinner provides functionality for pinning .env file snapshots
// to a baseline and detecting drift from that baseline over time.
package pinner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/envlens/internal/parser"
)

// Pin represents a pinned baseline for a single .env file.
type Pin struct {
	File    string            `json:"file"`
	Digest  string            `json:"digest"`
	Keys    []string          `json:"keys"`
	Hashes  map[string]string `json:"hashes"`
}

// DriftEntry describes a single key that has drifted from the baseline.
type DriftEntry struct {
	Key    string
	Reason string // "added", "removed", "changed"
}

// DriftReport holds all drift entries for a file.
type DriftReport struct {
	File    string
	Entries []DriftEntry
}

// Create builds a Pin from the parsed entries of a file.
func Create(file string, entries []parser.Entry) Pin {
	hashes := make(map[string]string, len(entries))
	keys := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.IsComment || e.Key == "" {
			continue
		}
		hashes[e.Key] = hashValue(e.Value)
		keys = append(keys, e.Key)
	}

	sort.Strings(keys)

	return Pin{
		File:   file,
		Digest: digestPin(hashes),
		Keys:   keys,
		Hashes: hashes,
	}
}

// Detect compares current entries against a baseline Pin and returns a DriftReport.
func Detect(pin Pin, entries []parser.Entry) DriftReport {
	current := make(map[string]string)
	for _, e := range entries {
		if e.IsComment || e.Key == "" {
			continue
		}
		current[e.Key] = hashValue(e.Value)
	}

	report := DriftReport{File: pin.File}

	for _, k := range pin.Keys {
		if h, ok := current[k]; !ok {
			report.Entries = append(report.Entries, DriftEntry{Key: k, Reason: "removed"})
		} else if h != pin.Hashes[k] {
			report.Entries = append(report.Entries, DriftEntry{Key: k, Reason: "changed"})
		}
	}

	pinned := make(map[string]struct{}, len(pin.Keys))
	for _, k := range pin.Keys {
		pinned[k] = struct{}{}
	}
	for k := range current {
		if _, ok := pinned[k]; !ok {
			report.Entries = append(report.Entries, DriftEntry{Key: k, Reason: "added"})
		}
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Key < report.Entries[j].Key
	})

	return report
}

func hashValue(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])
}

func digestPin(hashes map[string]string) string {
	keys := make([]string, 0, len(hashes))
	for k := range hashes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s;", k, hashes[k])
	}
	h := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(h[:])
}

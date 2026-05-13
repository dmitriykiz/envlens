// Package comparator provides functionality for comparing .env files
// across multiple environments to detect drift and inconsistencies.
package comparator

import (
	"sort"

	"github.com/user/envlens/internal/parser"
)

// FileSet maps a label (e.g. filename or environment name) to its parsed entries.
type FileSet map[string][]parser.Entry

// DriftReport holds the result of comparing two or more .env files.
type DriftReport struct {
	// KeysOnlyIn maps a label to keys that exist only in that file.
	KeysOnlyIn map[string][]string
	// CommonKeys lists keys present in every file.
	CommonKeys []string
	// ValueMismatches maps a key to the per-label values when files disagree.
	ValueMismatches map[string]map[string]string
}

// Compare analyses the provided FileSet and returns a DriftReport describing
// how the files differ from one another.
func Compare(files FileSet) DriftReport {
	report := DriftReport{
		KeysOnlyIn:      make(map[string][]string),
		ValueMismatches: make(map[string]map[string]string),
	}

	// Build per-label key→value maps.
	labelMaps := make(map[string]map[string]string, len(files))
	for label, entries := range files {
		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		labelMaps[label] = m
	}

	// Collect the union of all keys.
	keySet := make(map[string]struct{})
	for _, m := range labelMaps {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}

	labels := make([]string, 0, len(labelMaps))
	for l := range labelMaps {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	for key := range keySet {
		presentIn := []string{}
		values := make(map[string]string)

		for _, label := range labels {
			if v, ok := labelMaps[label][key]; ok {
				presentIn = append(presentIn, label)
				values[label] = v
			}
		}

		if len(presentIn) == len(labels) {
			// Key is in every file — check for value mismatches.
			report.CommonKeys = append(report.CommonKeys, key)
			if hasMismatch(values) {
				report.ValueMismatches[key] = values
			}
		} else {
			// Key is missing from at least one file.
			for _, label := range presentIn {
				report.KeysOnlyIn[label] = append(report.KeysOnlyIn[label], key)
			}
		}
	}

	sort.Strings(report.CommonKeys)
	return report
}

// hasMismatch returns true if not all values in the map are identical.
func hasMismatch(values map[string]string) bool {
	var ref string
	first := true
	for _, v := range values {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return true
		}
	}
	return false
}

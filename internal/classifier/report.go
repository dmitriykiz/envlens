package classifier

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteReport writes a human-readable classification report to w.
func WriteReport(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No entries to classify.")
		return
	}

	grouped := groupByCategory(results)
	categories := sortedCategories(grouped)

	for _, cat := range categories {
		items := grouped[cat]
		fmt.Fprintf(w, "[%s] (%d)\n", strings.ToUpper(string(cat)), len(items))
		for _, r := range items {
			fmt.Fprintf(w, "  %-30s  # %s\n", r.Entry.Key, r.Reason)
		}
		fmt.Fprintln(w)
	}
}

// Summary returns a map of category → count.
func Summary(results []Result) map[Category]int {
	out := make(map[Category]int)
	for _, r := range results {
		out[r.Category]++
	}
	return out
}

func groupByCategory(results []Result) map[Category][]Result {
	m := make(map[Category][]Result)
	for _, r := range results {
		m[r.Category] = append(m[r.Category], r)
	}
	return m
}

func sortedCategories(m map[Category][]Result) []Category {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	out := make([]Category, len(keys))
	for i, k := range keys {
		out[i] = Category(k)
	}
	return out
}

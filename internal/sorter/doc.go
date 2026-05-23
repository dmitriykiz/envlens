// Package sorter provides deterministic ordering and grouping of parsed .env
// entries for display, diffing, and export purposes.
//
// # Sorting
//
// Sort returns a new, sorted slice of [parser.Entry] values without modifying
// the caller's slice. Entries can be ordered by key or value, in ascending or
// descending direction:
//
//	out := sorter.Sort(entries, sorter.Options{
//		By:    sorter.ByKey,
//		Order: sorter.Ascending,
//	})
//
// # Sensitive grouping
//
// Setting Options.GroupSensitive to true causes keys that match common
// sensitive patterns (PASSWORD, TOKEN, API_KEY, etc.) to be floated to the
// top of the result, regardless of the primary sort field.
//
// # Prefix grouping
//
// GroupByPrefix partitions entries by the prefix that precedes a separator
// (typically "_"), making it easy to render environment files section by
// section:
//
//	groups := sorter.GroupByPrefix(entries, "_")
//	for prefix, group := range groups {
//		fmt.Println("[", prefix, "]")
//		for _, e := range group { fmt.Println(e.Key) }
//	}
package sorter

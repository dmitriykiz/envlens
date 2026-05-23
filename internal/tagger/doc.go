// Package tagger annotates parsed .env entries with semantic tags such as
// "sensitive", "deprecated", "internal", "optional", and "public".
//
// Tags are derived from two sources:
//
//  1. Built-in heuristics – when Options.AutoSensitive is true, keys that
//     match common patterns (e.g. SECRET, TOKEN, API_KEY) are automatically
//     labelled with TagSensitive.
//
//  2. User-defined rules – each Rule pairs a compiled regexp against one or
//     more Tags. All rules whose pattern matches the entry key are applied,
//     so a single key may carry multiple tags.
//
// Basic usage:
//
//	results := tagger.Tag(entries, tagger.Options{
//		AutoSensitive: true,
//		Rules: []tagger.Rule{
//			{
//				Pattern: regexp.MustCompile(`(?i)^LEGACY_`),
//				Tags:    []tagger.Tag{tagger.TagDeprecated},
//			},
//		},
//	})
//	for _, r := range results {
//		fmt.Println(r.Entry.Key, r.Tags)
//	}
package tagger

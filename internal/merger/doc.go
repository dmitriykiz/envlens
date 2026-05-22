// Package merger provides functionality for combining .env entries from
// multiple source files into a single unified set.
//
// # Overview
//
// In a monorepo it is common to layer environment configuration — for example
// a shared base .env supplemented by service-specific overrides.  The merger
// package handles this use-case by accepting a map of filename → []Entry and
// producing a merged Result that contains:
//
//   - Entries: the deduplicated, ordered list of key/value pairs.
//   - Conflicts: metadata about every key that appeared in more than one file.
//
// # Conflict strategies
//
// Three strategies are available via the ConflictStrategy type:
//
//   - StrategyFirst  – keep the value from the earliest (alphabetically) file.
//   - StrategyLast   – keep the value from the latest file.
//   - StrategyError  – abort and return an error on the first conflict.
//
// Example usage:
//
//	res, err := merger.Merge(fileEntries, merger.StrategyFirst)
//	if err != nil { ... }
//	for _, e := range res.Entries { fmt.Println(e.Key, e.Value) }
package merger

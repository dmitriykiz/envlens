// Package trimmer cleans up env entries by removing extraneous whitespace
// from keys and values, and optionally stripping surrounding quotes from
// values.
//
// # Usage
//
//	result := trimmer.Trim(entries, trimmer.DefaultOptions())
//	fmt.Printf("%d entries modified\n", result.Changed)
//
// # Options
//
// Options.TrimKeys   – strip leading/trailing space from key names.
// Options.TrimValues – strip leading/trailing space from values.
// Options.TrimQuotes – remove a matching pair of surrounding ' or " from values.
//
// All three flags are enabled by DefaultOptions.
//
// Comment and blank entries (Entry.Comment == true or Entry.Key == "") are
// passed through without modification and do not count towards the Changed
// total.
package trimmer

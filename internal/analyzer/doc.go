// Package analyzer inspects collections of parsed .env entries and produces
// structured reports highlighting:
//
//   - Duplicate keys: the same key appears more than once in a single file.
//   - Sensitive keys: keys whose names match common patterns for secrets,
//     passwords, tokens, or credentials that should never be committed.
//   - Missing keys: keys present in a reference set (e.g. .env.example) that
//     are absent from the file under inspection.
//
// Typical usage:
//
//	entries, err := parser.ParseFile(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	refEntries, _ := parser.ParseFile(".env.example")
//	refKeys := make([]string, len(refEntries))
//	for i, e := range refEntries {
//		refKeys[i] = e.Key
//	}
//	result := analyzer.Analyze(".env", entries, refKeys)
//	fmt.Println(result.Missing)
package analyzer

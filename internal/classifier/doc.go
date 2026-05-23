// Package classifier assigns semantic categories to .env file entries
// based on their key names.
//
// Entries are matched against a set of built-in regular-expression rules
// and assigned one of the following categories:
//
//   - credential  – secrets, tokens, passwords, API keys
//   - database    – connection strings, DB hosts, DSNs
//   - network     – hosts, ports, URLs, endpoints
//   - feature     – feature flags and toggles
//   - logging     – log level and output configuration
//   - unknown     – no pattern matched
//
// Usage:
//
//	entries, _ := parser.ParseFile(".env")
//	results := classifier.Classify(entries)
//	classifier.WriteReport(os.Stdout, results)
//
// The Summary function returns a map of category to entry count, suitable
// for dashboards or CI gate checks.
package classifier

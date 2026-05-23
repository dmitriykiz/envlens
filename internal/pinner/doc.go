// Package pinner implements baseline pinning for .env files.
//
// A "pin" is a cryptographic snapshot of the keys and hashed values present
// in an .env file at a specific point in time. Pins can be saved to disk as
// JSON files and later loaded to detect drift — keys that were added, removed,
// or whose values have changed since the baseline was recorded.
//
// Typical usage:
//
//	entries, _ := parser.ParseFile(".env")
//	pin := pinner.Create(".env", entries)
//	pinner.SavePin(".env.pin.json", pin)
//
//	// Later, check for drift:
//	loaded, _ := pinner.LoadPin(".env.pin.json")
//	current, _ := parser.ParseFile(".env")
//	report := pinner.Detect(loaded, current)
//	pinner.WriteReport(os.Stdout, []pinner.DriftReport{report})
//
// Value hashes use SHA-256 so that sensitive values are never stored in
// plain text within the pin file.
package pinner

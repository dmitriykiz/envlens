// Package profiler provides statistical profiling of .env files
// across a monorepo.
//
// It parses one or more .env files and aggregates metrics including:
//
//   - Total number of key-value entries per file and globally
//   - Count of sensitive keys (matching patterns like SECRET, TOKEN, etc.)
//   - Count of keys with empty values
//   - Number of globally unique keys across all files
//
// Usage:
//
//	paths := []string{".env", "services/api/.env"}
//	prof, err := profiler.Build(paths)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Total keys: %d\n", prof.TotalKeys)
//	fmt.Printf("Sensitive:  %d\n", prof.SensitiveKeys)
//	fmt.Printf("Unique:     %d\n", prof.UniqueKeys)
//
// The Profile and FileProfile types can be serialised with the
// exporter package for JSON or CSV output.
package profiler

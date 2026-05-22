// Package scanner walks a directory tree and discovers .env files matching
// configurable glob patterns, delegating the actual parsing to the parser
// package and returning structured results ready for analysis.
//
// Basic usage:
//
//	results, err := scanner.ScanDir("/path/to/monorepo", scanner.Options{
//		Patterns:    []string{".env", ".env.*"},
//		ExcludeDirs: []string{"node_modules", ".git", "vendor"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, r := range results {
//		if r.Err != nil {
//			log.Printf("skipping %s: %v", r.Path, r.Err)
//			continue
//		}
//		fmt.Printf("%s: %d keys\n", r.Path, len(r.Entries))
//	}
//
// # Options
//
// Patterns accepts standard glob syntax (e.g. ".env.*", "*.env"). If no
// patterns are provided, the scanner defaults to [".env"].
//
// ExcludeDirs is matched against each directory component of the path, so
// listing "node_modules" will skip it at any depth in the tree.
//
// The scanner does not perform any analysis; use the analyzer package to
// detect missing, duplicate, or sensitive keys across the collected results.
package scanner

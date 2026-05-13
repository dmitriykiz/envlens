// Package parser provides utilities for reading and parsing .env files
// used throughout the envlens audit tool.
//
// A .env file is expected to follow the common KEY=VALUE format:
//
//	# This is a comment
//	DB_HOST=localhost
//	DB_PORT=5432
//	SECRET="some quoted value"
//
// Rules applied during parsing:
//   - Blank lines are ignored.
//   - Lines beginning with '#' are treated as comments and ignored.
//   - Each non-comment, non-blank line must contain exactly one '=' character
//     separating the key from the value.
//   - Keys and values are trimmed of surrounding whitespace.
//   - Values optionally enclosed in single or double quotes have those quotes
//     stripped from the result.
//
// Usage:
//
//	entries, err := parser.ParseFile(".env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, e := range entries {
//	    fmt.Printf("%s = %s (line %d)\n", e.Key, e.Value, e.Line)
//	}
package parser

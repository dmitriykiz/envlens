// Package pinpointer locates the exact file, line, and column position of
// .env entries whose keys match a given pattern.
//
// # Overview
//
// In large monorepos a key may appear in many .env files. Pinpointer lets you
// answer "where exactly is DB_PASSWORD defined?" without manually grepping
// through every file.
//
// # Usage
//
//	locs, err := pinpointer.Pinpoint(files, pinpointer.Options{
//	    KeyPattern:    "SECRET|TOKEN|PASSWORD",
//	    CaseSensitive: false,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	pinpointer.WriteReport(os.Stdout, locs)
//
// # Output
//
// WriteReport renders a tab-aligned table with FILE, LINE, COL, KEY and a
// truncated VALUE column so sensitive data is never fully exposed in logs.
package pinpointer

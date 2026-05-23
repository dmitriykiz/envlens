// Package grouper partitions a flat list of env entries into named groups
// derived from key-prefix conventions.
//
// # Overview
//
// Many .env files follow a naming convention where related keys share a common
// prefix separated by an underscore (or another delimiter), for example:
//
//	DB_HOST=localhost
//	DB_PORT=5432
//	AWS_ACCESS_KEY_ID=…
//	AWS_SECRET_ACCESS_KEY=…
//
// The grouper package exploits this convention to produce a []Group value
// that associates each prefix with its matching entries.
//
// # Usage
//
//	groups := grouper.ByPrefix(entries, grouper.Options{
//	    MinGroupSize:   2,
//	    Separator:      "_",
//	    UngroupedLabel: "misc",
//	})
//	for _, g := range groups {
//	    fmt.Println(g.Prefix, len(g.Entries))
//	}
//
// Entries whose keys contain no separator are placed in the ungrouped bucket.
// Groups smaller than MinGroupSize are silently discarded.
package grouper

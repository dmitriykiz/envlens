// Package patcher applies declarative patch rules to a slice of parsed
// .env entries, enabling programmatic add, update, and delete operations
// without re-parsing the source file.
//
// # Overview
//
// A [Rule] pairs an [Op] (OpSet or OpDelete) with a target key and,
// for OpSet, the desired value. Rules are applied in order so that
// later rules for the same key take precedence.
//
// # Usage
//
//	entries, _ := parser.ParseFile(".env")
//
//	rules := []patcher.Rule{
//		{Op: patcher.OpSet,    Key: "LOG_LEVEL", Value: "debug"},
//		{Op: patcher.OpDelete, Key: "LEGACY_FLAG"},
//	}
//
//	result, err := patcher.Patch(entries, rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("added=%v updated=%v deleted=%v skipped=%v\n",
//		result.Added, result.Updated, result.Deleted, result.Skipped)
//
// The modified entries can then be written back with merger.WriteEnv or
// exported via the exporter package.
package patcher

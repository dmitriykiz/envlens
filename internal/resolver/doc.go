// Package resolver provides variable-interpolation for .env file entries.
//
// Many real-world .env files reference other variables defined in the same
// file or in the host environment, using the familiar shell syntax:
//
//	${VAR_NAME}   – braced form (recommended)
//	$VAR_NAME     – unbraced form
//
// # Usage
//
// Call [Resolve] with a slice of [parser.Entry] values and a [Strategy]:
//
//	result, err := resolver.Resolve(entries, resolver.StrategyStrict)
//
// The returned [Result] contains the expanded entries and any non-fatal
// warnings generated during resolution.
//
// # Strategies
//
//   - [StrategyStrict] – returns an error when a referenced variable cannot
//     be found in the entry set or in the host environment.
//   - [StrategyPermissive] – leaves unresolvable references as-is and
//     records a warning in [Result.Warnings].
//
// # Resolution order
//
// Resolution is performed in a single pass over the entries in the order
// they appear in the file. An entry's value is expanded using all previously
// defined entries plus the host environment, so later entries can reference
// earlier ones but not vice versa. The full key index is built before
// expansion begins, which means forward references are syntactically
// recognised but their values will be empty (or trigger an error under
// [StrategyStrict]).
package resolver

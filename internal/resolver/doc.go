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
// Resolution is performed in a single pass; forward references are supported
// because [os.Expand] is used as the expansion engine and the full index is
// built before expansion begins.
package resolver

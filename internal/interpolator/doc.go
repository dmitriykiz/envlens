// Package interpolator implements shell-style variable interpolation for
// .env file values.
//
// # Overview
//
// Many real-world .env files reference other variables defined earlier in
// the same file, e.g.:
//
//	BASE_URL=https://api.example.com
//	HEALTH_URL=${BASE_URL}/health
//
// Interpolate processes a slice of parser.Entry values in declaration order,
// building an internal lookup table as it goes so that each entry can
// reference any key defined before it.
//
// # Supported Syntax
//
//   - ${VAR}       — curly-brace reference
//   - $VAR         — bare reference (stops at first non-identifier character)
//   - ${VAR:-def}  — reference with default value when VAR is unset or empty
//
// # Warnings
//
// If a referenced variable cannot be resolved and no default is provided,
// the original placeholder is preserved verbatim and a warning string is
// added to the corresponding Result.Warnings slice.
package interpolator

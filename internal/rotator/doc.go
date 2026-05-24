// Package rotator provides value-rotation utilities for .env entries.
//
// A rotation rule pairs a key pattern (regular expression) with a Strategy
// that controls how the matched entry's value is replaced:
//
//   - StrategyRedact    — replaces the value with a configurable placeholder
//     (default: "REDACTED").  Useful for scrubbing secrets before committing
//     a sanitised snapshot.
//
//   - StrategyIncrement — appends or increments a numeric suffix on the
//     existing value (e.g. "v1" → "v2", "secret" → "secret_2").  Handy for
//     automated credential versioning in CI pipelines.
//
//   - StrategyCustom   — delegates value generation to a caller-supplied
//     Generator function, giving full control over the replacement logic.
//
// Usage:
//
//	rules := []rotator.Rule{
//	    {KeyPattern: `(?i)secret|token|key`, Strategy: rotator.StrategyRedact},
//	}
//	updated, results, err := rotator.Rotate(filename, entries, rules)
//	rotator.WriteReport(os.Stdout, results)
//
The package does not write files; callers are responsible for persisting the
returned updated entries via merger.WriteEnv or a similar writer.
package rotator

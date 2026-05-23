// Package injector loads key-value pairs from a .env file and injects them
// into the current process environment via os.Setenv.
//
// # Basic usage
//
//	results, err := injector.Inject(".env", injector.Options{
//		Strategy: injector.StrategySkip, // don't overwrite existing vars
//	})
//
// # Strategies
//
// StrategySkip (default) preserves any variable that is already present in the
// process environment. StrategyOverwrite replaces existing values with those
// found in the file.
//
// # Dry-run mode
//
// Setting Options.DryRun = true parses the file and returns Results without
// calling os.Setenv, which is useful for previewing what would change.
//
// # Reporting
//
// WriteReport renders a human-readable table of injection outcomes to any
// io.Writer, listing each key as SET, SKIP, or ERROR together with a summary
// line.
package injector

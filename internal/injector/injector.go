package injector

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/internal/parser"
)

// Strategy controls how conflicts between existing OS env vars and injected
// values are resolved.
type Strategy int

const (
	// StrategySkip leaves existing OS env vars untouched.
	StrategySkip Strategy = iota
	// StrategyOverwrite replaces existing OS env vars with values from the file.
	StrategyOverwrite
)

// Options configures the injection behaviour.
type Options struct {
	// Strategy determines what happens when a key already exists in the
	// process environment.
	Strategy Strategy
	// DryRun reports what would be injected without actually calling os.Setenv.
	DryRun bool
}

// Result holds the outcome of a single injection attempt.
type Result struct {
	Key      string
	Value    string
	Skipped  bool   // true when StrategySkip and key already existed
	Err      error
}

// Inject loads entries from envFile and sets them in the current process
// environment according to opts. It returns one Result per non-comment,
// non-blank entry found in the file.
func Inject(envFile string, opts Options) ([]Result, error) {
	entries, err := parser.ParseFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("injector: parse %q: %w", envFile, err)
	}

	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		if e.Comment || strings.TrimSpace(e.Key) == "" {
			continue
		}

		r := Result{Key: e.Key, Value: e.Value}

		if opts.Strategy == StrategySkip {
			if _, exists := os.LookupEnv(e.Key); exists {
				r.Skipped = true
				results = append(results, r)
				continue
			}
		}

		if !opts.DryRun {
			if err := os.Setenv(e.Key, e.Value); err != nil {
				r.Err = fmt.Errorf("os.Setenv: %w", err)
			}
		}

		results = append(results, r)
	}

	return results, nil
}

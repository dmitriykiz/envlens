package rotator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envlens/internal/parser"
)

// Strategy controls how a rotated value is generated.
type Strategy string

const (
	StrategyRedact  Strategy = "redact"  // replace with placeholder
	StrategyIncrement Strategy = "increment" // append/increment a numeric suffix
	StrategyCustom  Strategy = "custom"  // use caller-supplied generator
)

// Rule describes a single rotation rule.
type Rule struct {
	// KeyPattern is a regular expression matched against the entry key.
	KeyPattern string
	Strategy   Strategy
	// Placeholder is used when Strategy == StrategyRedact.
	Placeholder string
	// Generator is called when Strategy == StrategyCustom.
	Generator func(key, oldValue string) (string, error)
}

// Result holds the outcome for one rotated entry.
type Result struct {
	File     string
	Line     int
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
}

// Rotate applies rotation rules to the provided file entries and returns
// updated entries alongside a per-entry result log.
func Rotate(file string, entries []parser.Entry, rules []Rule) ([]parser.Entry, []Result, error) {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		re, err := regexp.Compile(r.KeyPattern)
		if err != nil {
			return nil, nil, fmt.Errorf("rotator: invalid pattern %q: %w", r.KeyPattern, err)
		}
		compiled[i] = re
	}

	updated := make([]parser.Entry, len(entries))
	var results []Result

	for idx, e := range entries {
		updated[idx] = e
		if e.Comment || e.Key == "" {
			continue
		}
		for i, re := range compiled {
			if !re.MatchString(e.Key) {
				continue
			}
			newVal, err := applyStrategy(rules[i], e.Key, e.Value)
			if err != nil {
				return nil, nil, fmt.Errorf("rotator: key %q: %w", e.Key, err)
			}
			results = append(results, Result{
				File:     file,
				Line:     e.Line,
				Key:      e.Key,
				OldValue: e.Value,
				NewValue: newVal,
				Rotated:  newVal != e.Value,
			})
			updated[idx].Value = newVal
			break
		}
	}
	return updated, results, nil
}

func applyStrategy(r Rule, key, old string) (string, error) {
	switch r.Strategy {
	case StrategyRedact:
		ph := r.Placeholder
		if ph == "" {
			ph = "REDACTED"
		}
		return ph, nil
	case StrategyIncrement:
		return increment(old), nil
	case StrategyCustom:
		if r.Generator == nil {
			return "", fmt.Errorf("custom strategy requires a Generator function")
		}
		return r.Generator(key, old)
	default:
		return old, nil
	}
}

var trailNum = regexp.MustCompile(`(\d+)$`)

func increment(v string) string {
	if m := trailNum.FindStringIndex(v); m != nil {
		num := v[m[0]:m[1]]
		n := 0
		fmt.Sscanf(num, "%d", &n)
		return v[:m[0]] + fmt.Sprintf("%0*d", len(num), n+1)
	}
	return strings.TrimRight(v, "_") + "_2"
}

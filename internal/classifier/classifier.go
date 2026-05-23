package classifier

import (
	"regexp"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Category represents a semantic classification for an env entry.
type Category string

const (
	CategoryCredential Category = "credential"
	CategoryDatabase   Category = "database"
	CategoryNetwork    Category = "network"
	CategoryFeature    Category = "feature"
	CategoryLogging    Category = "logging"
	CategoryUnknown    Category = "unknown"
)

// Result holds the classification result for a single entry.
type Result struct {
	Entry    parser.Entry
	Category Category
	Reason   string
}

var rules = []struct {
	pattern  *regexp.Regexp
	category Category
	reason   string
}{
	{regexp.MustCompile(`(?i)(secret|password|passwd|token|api_?key|private_?key|auth)`), CategoryCredential, "matches credential pattern"},
	{regexp.MustCompile(`(?i)(db_|database|postgres|mysql|mongo|redis|dsn|jdbc)`), CategoryDatabase, "matches database pattern"},
	{regexp.MustCompile(`(?i)(host|port|url|uri|endpoint|addr|address|domain)`), CategoryNetwork, "matches network pattern"},
	{regexp.MustCompile(`(?i)(feature_|flag_|enable_|disable_|toggle_)`), CategoryFeature, "matches feature flag pattern"},
	{regexp.MustCompile(`(?i)(log_|logging_|log$|loglevel|log_level)`), CategoryLogging, "matches logging pattern"},
}

// Classify assigns a Category to each non-comment, non-blank entry.
func Classify(entries []parser.Entry) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		if e.IsComment || e.IsBlank {
			continue
		}
		results = append(results, classifyEntry(e))
	}
	return results
}

func classifyEntry(e parser.Entry) Result {
	upper := strings.ToUpper(e.Key)
	for _, rule := range rules {
		if rule.pattern.MatchString(upper) {
			return Result{Entry: e, Category: rule.category, Reason: rule.reason}
		}
	}
	return Result{Entry: e, Category: CategoryUnknown, Reason: "no pattern matched"}
}

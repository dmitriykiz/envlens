package tagger

import (
	"regexp"
	"strings"

	"github.com/envlens/internal/parser"
)

// Tag represents a label attached to an env entry.
type Tag string

const (
	TagSensitive Tag = "sensitive"
	TagOptional  Tag = "optional"
	TagDeprecated Tag = "deprecated"
	TagInternal  Tag = "internal"
	TagPublic    Tag = "public"
)

// Result holds an entry alongside its assigned tags.
type Result struct {
	Entry parser.Entry
	Tags  []Tag
}

// Rule maps a key pattern (regexp) to a set of tags to apply.
type Rule struct {
	Pattern *regexp.Regexp
	Tags    []Tag
}

// Options configures the tagging behaviour.
type Options struct {
	// Rules are evaluated in order; all matching rules apply.
	Rules []Rule
	// AutoSensitive enables built-in sensitive key detection.
	AutoSensitive bool
}

var sensitivePattern = regexp.MustCompile(
	`(?i)(secret|password|passwd|token|api[_-]?key|private[_-]?key|auth|credential|cert|jwt)`,
)

// Tag applies tag rules to a slice of entries and returns annotated results.
func Tag(entries []parser.Entry, opts Options) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		if e.IsComment || e.IsBlank {
			results = append(results, Result{Entry: e})
			continue
		}
		tagSet := map[Tag]struct{}{}
		if opts.AutoSensitive && sensitivePattern.MatchString(e.Key) {
			tagSet[TagSensitive] = struct{}{}
		}
		for _, rule := range opts.Rules {
			if rule.Pattern != nil && rule.Pattern.MatchString(e.Key) {
				for _, t := range rule.Tags {
					tagSet[t] = struct{}{}
				}
			}
		}
		tags := make([]Tag, 0, len(tagSet))
		for t := range tagSet {
			tags = append(tags, t)
		}
		sortTags(tags)
		results = append(results, Result{Entry: e, Tags: tags})
	}
	return results
}

// HasTag reports whether a Result carries the given tag.
func HasTag(r Result, t Tag) bool {
	for _, tag := range r.Tags {
		if tag == t {
			return true
		}
	}
	return false
}

func sortTags(tags []Tag) {
	// simple insertion sort to keep output deterministic
	for i := 1; i < len(tags); i++ {
		for j := i; j > 0 && strings.Compare(string(tags[j-1]), string(tags[j])) > 0; j-- {
			tags[j-1], tags[j] = tags[j], tags[j-1]
		}
	}
}

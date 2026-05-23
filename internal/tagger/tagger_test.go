package tagger

import (
	"regexp"
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	if len(pairs)%2 != 0 {
		panic("makeEntries: pairs must be even")
	}
	out := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestTag_NoRulesNoAutoSensitive(t *testing.T) {
	entries := makeEntries("APP_NAME", "envlens", "DB_HOST", "localhost")
	results := Tag(entries, Options{})
	for _, r := range results {
		if len(r.Tags) != 0 {
			t.Errorf("expected no tags for %q, got %v", r.Entry.Key, r.Tags)
		}
	}
}

func TestTag_AutoSensitiveDetectsSecretKeys(t *testing.T) {
	entries := makeEntries(
		"APP_SECRET", "abc",
		"API_KEY", "xyz",
		"APP_NAME", "envlens",
	)
	results := Tag(entries, Options{AutoSensitive: true})
	if !HasTag(results[0], TagSensitive) {
		t.Errorf("expected APP_SECRET to be tagged sensitive")
	}
	if !HasTag(results[1], TagSensitive) {
		t.Errorf("expected API_KEY to be tagged sensitive")
	}
	if HasTag(results[2], TagSensitive) {
		t.Errorf("APP_NAME should not be tagged sensitive")
	}
}

func TestTag_CustomRuleApplied(t *testing.T) {
	entries := makeEntries("LEGACY_HOST", "old.host", "NEW_HOST", "new.host")
	rule := Rule{
		Pattern: regexp.MustCompile(`(?i)^LEGACY_`),
		Tags:    []Tag{TagDeprecated},
	}
	results := Tag(entries, Options{Rules: []Rule{rule}})
	if !HasTag(results[0], TagDeprecated) {
		t.Errorf("expected LEGACY_HOST to be tagged deprecated")
	}
	if HasTag(results[1], TagDeprecated) {
		t.Errorf("NEW_HOST should not be tagged deprecated")
	}
}

func TestTag_MultipleTagsFromMultipleRules(t *testing.T) {
	entries := makeEntries("INTERNAL_TOKEN", "tok")
	rules := []Rule{
		{Pattern: regexp.MustCompile(`(?i)TOKEN`), Tags: []Tag{TagSensitive}},
		{Pattern: regexp.MustCompile(`(?i)^INTERNAL_`), Tags: []Tag{TagInternal}},
	}
	results := Tag(entries, Options{Rules: rules})
	if !HasTag(results[0], TagSensitive) {
		t.Errorf("expected TagSensitive")
	}
	if !HasTag(results[0], TagInternal) {
		t.Errorf("expected TagInternal")
	}
}

func TestTag_SkipsCommentAndBlankEntries(t *testing.T) {
	entries := []parser.Entry{
		{IsComment: true, Raw: "# a comment"},
		{IsBlank: true},
		{Key: "DB_PASSWORD", Value: "secret"},
	}
	results := Tag(entries, Options{AutoSensitive: true})
	if len(results[0].Tags) != 0 {
		t.Errorf("comment entry should have no tags")
	}
	if len(results[1].Tags) != 0 {
		t.Errorf("blank entry should have no tags")
	}
	if !HasTag(results[2], TagSensitive) {
		t.Errorf("DB_PASSWORD should be sensitive")
	}
}

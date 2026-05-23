package scorer

import (
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Grade represents a letter-grade health score for a .env file.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// Report holds the numeric score and contributing penalty details.
type Report struct {
	File          string
	Score         int // 0–100
	Grade         Grade
	Penalties     []Penalty
}

// Penalty describes a single deduction applied to the score.
type Penalty struct {
	Reason string
	Points int
}

// Score evaluates the quality of a set of parsed entries for one file and
// returns a Report. The score starts at 100 and penalties are subtracted for
// each detected issue.
func Score(file string, entries []parser.Entry) Report {
	score := 100
	var penalties []Penalty

	penalise := func(reason string, pts int) {
		if pts > 0 {
			penalties = append(penalties, Penalty{Reason: reason, Points: pts})
			score -= pts
		}
	}

	// Count issues.
	empty, lower, special, sensitive := 0, 0, 0, 0
	seen := map[string]int{}

	for _, e := range entries {
		if e.Comment {
			continue
		}
		seen[e.Key]++
		if e.Value == "" {
			empty++
		}
		if e.Key != strings.ToUpper(e.Key) {
			lower++
		}
		if strings.ContainsAny(e.Key, " \t-.") {
			special++
		}
		if isSensitive(e.Key) && e.Value != "" {
			sensitive++
		}
	}

	dupes := 0
	for _, count := range seen {
		if count > 1 {
			dupes++
		}
	}

	penalise("empty values", empty*5)
	penalise("lowercase keys", lower*4)
	penalise("keys with special characters", special*6)
	penalise("duplicate keys", dupes*10)
	penalise("plaintext sensitive values", sensitive*8)

	if score < 0 {
		score = 0
	}

	return Report{
		File:      file,
		Score:     score,
		Grade:     gradeFor(score),
		Penalties: penalties,
	}
}

func gradeFor(score int) Grade {
	switch {
	case score >= 90:
		return GradeA
	case score >= 75:
		return GradeB
	case score >= 60:
		return GradeC
	case score >= 40:
		return GradeD
	default:
		return GradeF
	}
}

var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE_KEY", "AUTH", "CREDENTIAL", "DSN", "CERT",
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

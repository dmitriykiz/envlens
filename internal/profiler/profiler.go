// Package profiler analyses .env files and produces a usage profile
// summarising key statistics such as total keys, sensitive key ratio,
// empty value ratio, and per-file breakdowns.
package profiler

import (
	"path/filepath"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// FileProfile holds statistics for a single .env file.
type FileProfile struct {
	Path          string
	TotalKeys     int
	SensitiveKeys int
	EmptyValues   int
	UniqueKeys    int
}

// Profile is the aggregated result across all scanned files.
type Profile struct {
	Files         []FileProfile
	TotalKeys     int
	SensitiveKeys int
	EmptyValues   int
	UniqueKeys    int // across the whole set
}

var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "PRIVATE_KEY",
}

// Build reads each file path, parses it, and returns a Profile.
func Build(paths []string) (Profile, error) {
	globalKeys := make(map[string]struct{})
	var profile Profile

	for _, p := range paths {
		entries, err := parser.ParseFile(p)
		if err != nil {
			return Profile{}, err
		}

		fp := FileProfile{
			Path: filepath.Clean(p),
		}

		seen := make(map[string]struct{})
		for _, e := range entries {
			if e.Key == "" {
				continue
			}
			fp.TotalKeys++
			if isSensitive(e.Key) {
				fp.SensitiveKeys++
			}
			if e.Value == "" {
				fp.EmptyValues++
			}
			if _, dup := seen[e.Key]; !dup {
				fp.UniqueKeys++
				seen[e.Key] = struct{}{}
			}
			globalKeys[e.Key] = struct{}{}
		}

		profile.Files = append(profile.Files, fp)
		profile.TotalKeys += fp.TotalKeys
		profile.SensitiveKeys += fp.SensitiveKeys
		profile.EmptyValues += fp.EmptyValues
	}

	profile.UniqueKeys = len(globalKeys)
	return profile, nil
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pat := range sensitivePatterns {
		if strings.Contains(upper, pat) {
			return true
		}
	}
	return false
}

package pinpointer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourorg/envlens/internal/parser"
)

// Location describes where in a file a specific key was found.
type Location struct {
	File   string
	Key    string
	Line   int
	Column int
	Value  string
}

// Options controls Pinpoint behaviour.
type Options struct {
	// KeyPattern filters keys by regexp; empty means all keys.
	KeyPattern string
	// CaseSensitive controls whether KeyPattern matching is case-sensitive.
	CaseSensitive bool
}

// Pinpoint searches the given files for entries whose keys match the provided
// pattern and returns the exact file/line/column location of each match.
func Pinpoint(files []string, opts Options) ([]Location, error) {
	var re *regexp.Regexp
	if opts.KeyPattern != "" {
		pattern := opts.KeyPattern
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("pinpointer: invalid key pattern: %w", err)
		}
	}

	var locations []Location
	for _, file := range files {
		entries, err := parser.ParseFile(file)
		if err != nil {
			return nil, fmt.Errorf("pinpointer: parse %s: %w", file, err)
		}
		for _, e := range entries {
			if e.Key == "" {
				continue
			}
			if re != nil && !re.MatchString(e.Key) {
				continue
			}
			col := strings.Index(e.Raw, e.Key)
			if col < 0 {
				col = 0
			}
			locations = append(locations, Location{
				File:   file,
				Key:    e.Key,
				Line:   e.Line,
				Column: col + 1,
				Value:  e.Value,
			})
		}
	}
	return locations, nil
}

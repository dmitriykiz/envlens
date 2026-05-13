package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key      string
	Value    string
	Line     int
	FilePath string
}

// ParseFile reads a .env file and returns a slice of Entry.
// It skips blank lines and comments (lines starting with '#').
func ParseFile(filePath string) ([]Entry, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", filePath, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		line := strings.TrimSpace(raw)

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", filePath, lineNum, err)
		}

		entries = append(entries, Entry{
			Key:      key,
			Value:    value,
			Line:     lineNum,
			FilePath: filePath,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", filePath, err)
	}

	return entries, nil
}

// parseLine splits a "KEY=VALUE" line into its components.
func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("invalid line (no '=' found): %q", line)
	}

	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])

	// Strip optional surrounding quotes from value
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	if key == "" {
		return "", "", fmt.Errorf("empty key in line: %q", line)
	}

	return key, value, nil
}

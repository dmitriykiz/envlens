package auditor

import (
	"fmt"
	"time"

	"github.com/envlens/internal/parser"
)

// Severity represents the importance level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// Finding describes a single audit observation for a key-value entry.
type Finding struct {
	File     string
	Line     int
	Key      string
	Message  string
	Severity Severity
}

// Report is the result of auditing one or more .env files.
type Report struct {
	Files    []string
	Findings []Finding
	AuditedAt time.Time
}

// Options controls which audit checks are enabled.
type Options struct {
	CheckEmptyValues    bool
	CheckSensitivePlain bool
	CheckKeyFormat      bool
	SensitivePatterns   []string
}

// DefaultOptions returns a sensible default audit configuration.
func DefaultOptions() Options {
	return Options{
		CheckEmptyValues:    true,
		CheckSensitivePlain: true,
		CheckKeyFormat:      true,
		SensitivePatterns:   defaultSensitivePatterns,
	}
}

var defaultSensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE_KEY", "AUTH",
}

// Audit runs all enabled checks against the provided files and returns a Report.
func Audit(files []string, opts Options) (Report, error) {
	report := Report{
		Files:     files,
		AuditedAt: time.Now().UTC(),
	}

	for _, file := range files {
		entries, err := parser.ParseFile(file)
		if err != nil {
			return report, fmt.Errorf("auditor: parse %q: %w", file, err)
		}
		for _, entry := range entries {
			if entry.Comment {
				continue
			}
			if opts.CheckEmptyValues && entry.Value == "" {
				report.Findings = append(report.Findings, Finding{
					File:     file,
					Line:     entry.Line,
					Key:      entry.Key,
					Message:  "empty value",
					Severity: SeverityWarning,
				})
			}
			if opts.CheckSensitivePlain && isSensitiveKey(entry.Key, opts.SensitivePatterns) && entry.Value != "" {
				report.Findings = append(report.Findings, Finding{
					File:     file,
					Line:     entry.Line,
					Key:      entry.Key,
					Message:  "sensitive key has plain-text value",
					Severity: SeverityCritical,
				})
			}
			if opts.CheckKeyFormat && !isValidKeyFormat(entry.Key) {
				report.Findings = append(report.Findings, Finding{
					File:     file,
					Line:     entry.Line,
					Key:      entry.Key,
					Message:  "key contains invalid characters or is not UPPER_SNAKE_CASE",
					Severity: SeverityInfo,
				})
			}
		}
	}
	return report, nil
}

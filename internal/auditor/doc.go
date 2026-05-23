// Package auditor performs a comprehensive audit of one or more .env files,
// producing structured findings grouped by severity.
//
// # Overview
//
// The auditor runs a configurable set of checks:
//
//   - Empty values: keys that are declared but have no value assigned.
//   - Sensitive plain-text: keys whose names suggest secret material
//     (e.g. PASSWORD, TOKEN) but whose values are stored in clear text.
//   - Key format: keys that do not follow UPPER_SNAKE_CASE conventions.
//
// # Usage
//
//	opts := auditor.DefaultOptions()
//	report, err := auditor.Audit([]string{".env", "services/api/.env"}, opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//	auditor.WriteReport(os.Stdout, report)
//
// # Severity Levels
//
//   - CRITICAL – immediate attention required (e.g. exposed secrets).
//   - WARNING  – should be addressed before production (e.g. empty values).
//   - INFO     – style or convention issues.
package auditor

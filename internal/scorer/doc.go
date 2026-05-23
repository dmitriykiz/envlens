// Package scorer assigns a numeric health score (0–100) and a letter grade
// (A–F) to a parsed .env file.
//
// # Scoring model
//
// The score starts at 100. Penalties are subtracted for each detected
// quality issue:
//
//   - Empty value         –5 per occurrence
//   - Lowercase key       –4 per occurrence
//   - Key with spaces / special chars  –6 per occurrence
//   - Duplicate key       –10 per unique duplicated key
//   - Plaintext sensitive value  –8 per occurrence
//
// The final score is clamped to the range [0, 100].
//
// # Grades
//
//	90–100  A
//	75–89   B
//	60–74   C
//	40–59   D
//	 0–39   F
//
// # Usage
//
//	entries, _ := parser.ParseFile(".env")
//	report := scorer.Score(".env", entries)
//	scorer.WriteReport(os.Stdout, report)
package scorer

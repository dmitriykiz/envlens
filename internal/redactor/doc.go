// Package redactor provides utilities for masking sensitive values
// found in parsed .env file entries.
//
// The Redact function inspects each entry's key against a list of
// known sensitive patterns (e.g. PASSWORD, TOKEN, SECRET) and replaces
// the corresponding value with asterisks, preserving a maximum display
// length of 8 characters to avoid leaking value length information.
//
// Typical usage:
//
//	entries, err := parser.ParseFile(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	masked := redactor.Redact(entries)
//	for _, me := range masked {
//		fmt.Printf("%s=%s (masked=%v)\n", me.Key, me.Value, me.Masked)
//	}
//
// The original entries are never mutated; Redact returns a new slice
// of MaskedEntry values.
package redactor

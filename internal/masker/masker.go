// Package masker provides utilities for partially masking .env values
// in output, logs, and reports while preserving enough context for debugging.
package masker

import "strings"

// Style controls how a value is masked.
type Style int

const (
	// StyleFull replaces the entire value with asterisks (capped at 8).
	StyleFull Style = iota
	// StylePartial reveals the first and last character of the value.
	StylePartial
	// StylePrefix reveals the first four characters followed by asterisks.
	StylePrefix
)

// Options configures masking behaviour.
type Options struct {
	Style    Style
	MaskChar rune
	MaxMask  int // maximum number of mask characters; 0 means no limit
}

var defaultOptions = Options{
	Style:    StyleFull,
	MaskChar: '*',
	MaxMask:  8,
}

// Mask returns a masked representation of value using the provided options.
// If opts is nil the package defaults are used.
func Mask(value string, opts *Options) string {
	if opts == nil {
		opts = &defaultOptions
	}
	if value == "" {
		return ""
	}
	maskChar := string(opts.MaskChar)
	switch opts.Style {
	case StylePartial:
		return maskPartial(value, maskChar, opts.MaxMask)
	case StylePrefix:
		return maskPrefix(value, maskChar, opts.MaxMask)
	default:
		return maskFull(value, maskChar, opts.MaxMask)
	}
}

func maskFull(value, char string, max int) string {
	n := len(value)
	if max > 0 && n > max {
		n = max
	}
	return strings.Repeat(char, n)
}

func maskPartial(value, char string, max int) string {
	if len(value) <= 2 {
		return maskFull(value, char, max)
	}
	midLen := len(value) - 2
	if max > 0 && midLen > max {
		midLen = max
	}
	return string(value[0]) + strings.Repeat(char, midLen) + string(value[len(value)-1])
}

func maskPrefix(value, char string, max int) string {
	const reveal = 4
	if len(value) <= reveal {
		return value
	}
	suffix := len(value) - reveal
	if max > 0 && suffix > max {
		suffix = max
	}
	return value[:reveal] + strings.Repeat(char, suffix)
}

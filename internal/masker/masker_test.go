package masker_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/masker"
)

func TestMask_FullStyle_DefaultOptions(t *testing.T) {
	got := masker.Mask("supersecret", nil)
	if got != "********" {
		t.Errorf("expected 8 asterisks, got %q", got)
	}
}

func TestMask_FullStyle_ShortValue(t *testing.T) {
	got := masker.Mask("abc", nil)
	if got != "***" {
		t.Errorf("expected 3 asterisks, got %q", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	got := masker.Mask("", nil)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMask_PartialStyle(t *testing.T) {
	opts := &masker.Options{Style: masker.StylePartial, MaskChar: '*', MaxMask: 8}
	got := masker.Mask("hello", opts)
	// first='h', last='o', middle=3 stars
	if got != "h***o" {
		t.Errorf("expected h***o, got %q", got)
	}
}

func TestMask_PartialStyle_ShortValue(t *testing.T) {
	opts := &masker.Options{Style: masker.StylePartial, MaskChar: '*', MaxMask: 8}
	got := masker.Mask("ab", opts)
	if got != "**" {
		t.Errorf("expected **, got %q", got)
	}
}

func TestMask_PrefixStyle_RevealsFirst4(t *testing.T) {
	opts := &masker.Options{Style: masker.StylePrefix, MaskChar: '*', MaxMask: 8}
	got := masker.Mask("ABCDEFGHIJ", opts)
	if !strings.HasPrefix(got, "ABCD") {
		t.Errorf("expected prefix ABCD, got %q", got)
	}
	if !strings.Contains(got, "*") {
		t.Errorf("expected asterisks after prefix, got %q", got)
	}
}

func TestMask_PrefixStyle_ShortValue(t *testing.T) {
	opts := &masker.Options{Style: masker.StylePrefix, MaskChar: '*', MaxMask: 8}
	got := masker.Mask("ABC", opts)
	// shorter than reveal threshold — returned as-is
	if got != "ABC" {
		t.Errorf("expected ABC unchanged, got %q", got)
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	opts := &masker.Options{Style: masker.StyleFull, MaskChar: '#', MaxMask: 4}
	got := masker.Mask("password", opts)
	if got != "####" {
		t.Errorf("expected ####, got %q", got)
	}
}

func TestMask_NoMaxLimit(t *testing.T) {
	opts := &masker.Options{Style: masker.StyleFull, MaskChar: '*', MaxMask: 0}
	value := "verylongpassword"
	got := masker.Mask(value, opts)
	if len(got) != len(value) {
		t.Errorf("expected mask length %d, got %d", len(value), len(got))
	}
}

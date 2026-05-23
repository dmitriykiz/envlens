package cipher_test

import (
	"strings"
	"testing"

	"github.com/envlens/internal/cipher"
	"github.com/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

var testKey = []byte("0123456789abcdef") // 16 bytes

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	entries := makeEntries("DB_PASSWORD", "s3cr3t", "APP_NAME", "envlens")

	enc, err := cipher.Encrypt(entries, testKey)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if enc[0].Value == "s3cr3t" {
		t.Error("sensitive value should be encrypted")
	}
	if !strings.HasPrefix(enc[0].Value, "enc:") {
		t.Errorf("expected enc: prefix, got %q", enc[0].Value)
	}
	if enc[1].Value != "envlens" {
		t.Error("non-sensitive value should be unchanged")
	}

	dec, err := cipher.Decrypt(enc, testKey)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec[0].Value != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", dec[0].Value)
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	_, err := cipher.Encrypt(makeEntries("DB_PASSWORD", "x"), []byte("short"))
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	entries := []parser.Entry{{Key: "DB_PASSWORD", Value: "enc:invaliddataXXXX"}}
	_, err := cipher.Decrypt(entries, testKey)
	if err == nil {
		t.Fatal("expected decryption error for tampered ciphertext")
	}
}

func TestEncrypt_EmptyValueSkipped(t *testing.T) {
	entries := makeEntries("DB_PASSWORD", "")
	enc, err := cipher.Encrypt(entries, testKey)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc[0].Value != "" {
		t.Errorf("empty value should not be encrypted, got %q", enc[0].Value)
	}
}

func TestEncrypt_CommentEntryPassthrough(t *testing.T) {
	entries := []parser.Entry{{Key: "", Value: "# this is a comment", Comment: true}}
	enc, err := cipher.Encrypt(entries, testKey)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc[0].Value != "# this is a comment" {
		t.Errorf("comment entry should pass through unchanged")
	}
}

func TestDecrypt_NonEncryptedValuePassthrough(t *testing.T) {
	entries := makeEntries("APP_NAME", "envlens")
	dec, err := cipher.Decrypt(entries, testKey)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec[0].Value != "envlens" {
		t.Errorf("non-encrypted value should pass through, got %q", dec[0].Value)
	}
}

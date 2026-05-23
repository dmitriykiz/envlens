package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/envlens/internal/parser"
)

// ErrInvalidKey is returned when the AES key length is not 16, 24, or 32 bytes.
var ErrInvalidKey = errors.New("cipher: key must be 16, 24, or 32 bytes")

// ErrDecryptFailed is returned when decryption or authentication fails.
var ErrDecryptFailed = errors.New("cipher: decryption failed")

// Encrypt encrypts the value of each sensitive entry using AES-GCM.
// Non-sensitive entries are returned unchanged. Encrypted values are
// base64-encoded and prefixed with "enc:" to distinguish them.
func Encrypt(entries []parser.Entry, key []byte) ([]parser.Entry, error) {
	block, err := newBlock(key)
	if err != nil {
		return nil, err
	}
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if e.Comment || !isSensitiveKey(e.Key) || e.Value == "" {
			result[i] = e
			continue
		}
		enc, err := gcmEncrypt(block, []byte(e.Value))
		if err != nil {
			return nil, err
		}
		result[i] = parser.Entry{Key: e.Key, Value: "enc:" + enc, Comment: e.Comment}
	}
	return result, nil
}

// Decrypt decrypts entries previously encrypted by Encrypt.
// Values not prefixed with "enc:" are returned unchanged.
func Decrypt(entries []parser.Entry, key []byte) ([]parser.Entry, error) {
	block, err := newBlock(key)
	if err != nil {
		return nil, err
	}
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if e.Comment || !strings.HasPrefix(e.Value, "enc:") {
			result[i] = e
			continue
		}
		plain, err := gcmDecrypt(block, strings.TrimPrefix(e.Value, "enc:"))
		if err != nil {
			return nil, err
		}
		result[i] = parser.Entry{Key: e.Key, Value: plain, Comment: e.Comment}
	}
	return result, nil
}

func newBlock(key []byte) (cipher.Block, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKey
	}
	return aes.NewCipher(key)
}

func gcmEncrypt(block cipher.Block, plaintext []byte) (string, error) {
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := aead.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func gcmDecrypt(block cipher.Block, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrDecryptFailed
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := aead.NonceSize()
	if len(data) < ns {
		return "", ErrDecryptFailed
	}
	plain, err := aead.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", ErrDecryptFailed
	}
	return string(plain), nil
}

// Package cipher provides AES-GCM encryption and decryption for sensitive
// values found in .env files.
//
// # Overview
//
// Sensitive keys (those matching patterns such as PASSWORD, SECRET, TOKEN,
// API_KEY, etc.) can be encrypted in place so that .env files can be safely
// stored in version control or shared across teams without exposing secret
// material in plaintext.
//
// Encrypted values are base64-encoded AES-GCM ciphertexts prefixed with
// "enc:" so they are easily identifiable and round-trippable:
//
//	DB_PASSWORD=enc:<base64-ciphertext>
//
// # Usage
//
//	key := []byte(os.Getenv("ENVLENS_CIPHER_KEY")) // 16, 24, or 32 bytes
//
//	encrypted, err := cipher.Encrypt(entries, key)
//	decrypted, err := cipher.Decrypt(encrypted, key)
//
// Non-sensitive keys and comment lines are passed through unchanged.
// Empty values are never encrypted.
package cipher

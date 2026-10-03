// Package crypto provides the symmetric encryption used for provider
// credentials, webhook secrets and TOTP seeds at rest, plus hashing
// helpers for API keys and HMAC request signatures.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// KeySize is the AES-256 key length in bytes.
const KeySize = 32

// Version byte prefixed to every ciphertext so the format can evolve.
const versionV1 = 0x01

var (
	// ErrInvalidKey is returned when the master key has the wrong length.
	ErrInvalidKey = errors.New("crypto: master key must be 32 bytes (base64 or hex encoded)")
	// ErrCiphertext is returned for malformed or tampered ciphertext.
	ErrCiphertext = errors.New("crypto: invalid ciphertext")
)

// Cipher encrypts and decrypts with AES-256-GCM under one master key.
type Cipher struct {
	aead cipher.AEAD
}

// ParseKey decodes a master key given as base64 (std or URL, padded or
// not) or hex. The decoded key must be exactly 32 bytes.
func ParseKey(encoded string) ([]byte, error) {
	decoders := []func(string) ([]byte, error){
		hex.DecodeString,
		base64.StdEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.URLEncoding.DecodeString,
		base64.RawURLEncoding.DecodeString,
	}
	for _, d := range decoders {
		if b, err := d(encoded); err == nil && len(b) == KeySize {
			return b, nil
		}
	}
	return nil, ErrInvalidKey
}

// NewCipher builds a Cipher from a raw 32-byte key.
func NewCipher(key []byte) (*Cipher, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{aead: aead}, nil
}

// NewCipherFromString is NewCipher(ParseKey(encoded)).
func NewCipherFromString(encoded string) (*Cipher, error) {
	key, err := ParseKey(encoded)
	if err != nil {
		return nil, err
	}
	return NewCipher(key)
}

// Encrypt returns version || nonce || ciphertext+tag. aad binds the
// ciphertext to a context (e.g. the provider id) so blobs cannot be swapped
// between rows; pass nil when not needed.
func (c *Cipher) Encrypt(plaintext, aad []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("crypto: nonce: %w", err)
	}
	out := make([]byte, 0, 1+len(nonce)+len(plaintext)+c.aead.Overhead())
	out = append(out, versionV1)
	out = append(out, nonce...)
	return c.aead.Seal(out, nonce, plaintext, aad), nil
}

// Decrypt reverses Encrypt. The same aad must be supplied.
func (c *Cipher) Decrypt(blob, aad []byte) ([]byte, error) {
	ns := c.aead.NonceSize()
	if len(blob) < 1+ns+c.aead.Overhead() || blob[0] != versionV1 {
		return nil, ErrCiphertext
	}
	nonce, ct := blob[1:1+ns], blob[1+ns:]
	pt, err := c.aead.Open(nil, nonce, ct, aad)
	if err != nil {
		return nil, ErrCiphertext
	}
	return pt, nil
}

// GenerateKey returns a fresh random 32-byte key, base64 encoded, suitable
// for HABARCHY_MASTER_KEY.
func GenerateKey() (string, error) {
	key := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// RandomToken returns n random bytes encoded as URL-safe base64 without
// padding. Used for API key secrets and refresh tokens.
func RandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// SHA256 hashes s. API keys are stored as SHA256(plaintext).
func SHA256(s string) []byte {
	sum := sha256.Sum256([]byte(s))
	return sum[:]
}

// HMACSHA256Hex computes hex(HMAC-SHA256(secret, msg)). Used for request
// signatures (X-Signature) and webhook signatures (X-Habarchy-Signature).
func HMACSHA256Hex(secret, msg []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(msg)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyHMACSHA256Hex compares a hex signature in constant time.
func VerifyHMACSHA256Hex(secret, msg []byte, signatureHex string) bool {
	expected, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(msg)
	return hmac.Equal(mac.Sum(nil), expected)
}

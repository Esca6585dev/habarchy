package crypto

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	keyStr, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewCipherFromString(keyStr)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte(`{"token":"super-secret","chat_id":"123"}`)
	aad := []byte("provider:abc")

	blob, err := c.Encrypt(plain, aad)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(blob, plain) {
		t.Fatal("ciphertext must not contain plaintext")
	}
	got, err := c.Decrypt(blob, aad)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %q want %q", got, plain)
	}

	// Different nonce every time.
	blob2, _ := c.Encrypt(plain, aad)
	if bytes.Equal(blob, blob2) {
		t.Fatal("two encryptions of the same plaintext must differ")
	}

	// Wrong AAD, tampering and truncation must all fail.
	if _, err := c.Decrypt(blob, []byte("provider:other")); err == nil {
		t.Fatal("wrong aad must fail")
	}
	tampered := append([]byte(nil), blob...)
	tampered[len(tampered)-1] ^= 0xff
	if _, err := c.Decrypt(tampered, aad); err == nil {
		t.Fatal("tampered ciphertext must fail")
	}
	if _, err := c.Decrypt(blob[:5], aad); err == nil {
		t.Fatal("truncated ciphertext must fail")
	}

	// A different key must not decrypt.
	other, _ := NewCipherFromString(mustKey(t))
	if _, err := other.Decrypt(blob, aad); err == nil {
		t.Fatal("other key must fail")
	}
}

func TestParseKeyFormats(t *testing.T) {
	raw := bytes.Repeat([]byte{0xAB}, KeySize)
	for _, enc := range []string{
		hex.EncodeToString(raw),
		base64.StdEncoding.EncodeToString(raw),
		base64.RawStdEncoding.EncodeToString(raw),
		base64.RawURLEncoding.EncodeToString(raw),
	} {
		k, err := ParseKey(enc)
		if err != nil || !bytes.Equal(k, raw) {
			t.Errorf("ParseKey(%q) failed: %v", enc, err)
		}
	}
	if _, err := ParseKey("too-short"); err == nil {
		t.Error("short key must be rejected")
	}
	if _, err := NewCipher(make([]byte, 16)); err == nil {
		t.Error("16-byte key must be rejected (AES-256 only)")
	}
}

func TestHMAC(t *testing.T) {
	secret := []byte("whsec_test")
	msg := []byte(`1700000000.{"id":"m1"}`)
	sig := HMACSHA256Hex(secret, msg)
	if len(sig) != 64 {
		t.Fatalf("hex sig length %d", len(sig))
	}
	if !VerifyHMACSHA256Hex(secret, msg, sig) {
		t.Fatal("valid signature rejected")
	}
	if VerifyHMACSHA256Hex(secret, []byte("other"), sig) {
		t.Fatal("signature over other message accepted")
	}
	if VerifyHMACSHA256Hex(secret, msg, "zz") {
		t.Fatal("garbage signature accepted")
	}
}

func TestRandomToken(t *testing.T) {
	a, _ := RandomToken(32)
	b, _ := RandomToken(32)
	if a == b || len(a) != 43 {
		t.Fatalf("tokens: %q %q", a, b)
	}
	if len(SHA256(a)) != 32 {
		t.Fatal("sha256 length")
	}
}

func mustKey(t *testing.T) string {
	t.Helper()
	k, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return k
}

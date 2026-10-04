package webhooks

import (
	"testing"
	"time"
)

func TestSignVerify(t *testing.T) {
	secret := []byte("whsec_demo")
	body := []byte(`{"id":"1","type":"message.sent"}`)
	now := time.Unix(1_700_000_000, 0)
	h := Sign(secret, "1700000000", body)
	if !Verify(secret, h, body, 5*time.Minute, now) {
		t.Fatal("valid signature rejected")
	}
	if Verify(secret, h, []byte(`{"id":"2"}`), 5*time.Minute, now) {
		t.Fatal("tampered body accepted")
	}
	if Verify([]byte("other"), h, body, 5*time.Minute, now) {
		t.Fatal("wrong secret accepted")
	}
	if Verify(secret, h, body, 5*time.Minute, now.Add(10*time.Minute)) {
		t.Fatal("stale timestamp accepted")
	}
	if Verify(secret, "garbage", body, 0, now) {
		t.Fatal("garbage header accepted")
	}
}

func TestCheckURL(t *testing.T) {
	for _, ok := range []string{"https://tds.gov.tm/hooks", "http://example.com:8080/x?y=1"} {
		if err := checkURL(ok, false); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
	for _, bad := range []string{"ftp://x", "http://127.0.0.1/", "http://[::1]:80/", "http://0.0.0.0/", "http://10.1.2.3/", "http://localhost:3000/"} {
		if err := checkURL(bad, false); err == nil {
			t.Errorf("%s accepted", bad)
		}
	}
	if err := checkURL("http://127.0.0.1:9/", true); err != nil {
		t.Errorf("allowPrivate should accept loopback: %v", err)
	}
	if backoff(1) != 10*time.Second || backoff(99) != 3*time.Hour {
		t.Fatal("backoff schedule")
	}
}

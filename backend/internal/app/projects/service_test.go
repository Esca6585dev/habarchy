package projects

import (
	"errors"
	"testing"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Demo Project":        "demo-project",
		"  Türkmen  Şäher!! ": "turkmen-saher",
		"tds.gov.tm":          "tds-gov-tm",
		"ÄÖÜ ýňž":             "aou-ynz",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIPAllowed(t *testing.T) {
	list := []string{"10.0.0.0/8", "192.168.1.5"}
	if !IPAllowed(nil, "1.2.3.4") || !IPAllowed(list, "10.20.30.40") || !IPAllowed(list, "192.168.1.5") {
		t.Fatal("expected allowed")
	}
	if IPAllowed(list, "192.168.1.6") || IPAllowed(list, "garbage") {
		t.Fatal("expected denied")
	}
	if err := ValidateCIDRs([]string{"10.0.0.0/8", "::1"}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCIDRs([]string{"10.0.0.0/99"}); err == nil {
		t.Fatal("bad cidr accepted")
	}
}

func TestSignature(t *testing.T) {
	key := "hb_test_abc"
	now := time.Unix(1_700_000_000, 0)
	ts := "1700000000"
	body := []byte(`{"channel":"sms"}`)
	sig := Sign(key, ts, "post", "/api/v1/messages", body)

	in := SignatureInput{Signature: sig, Timestamp: ts, Method: "POST", Path: "/api/v1/messages", Body: body}
	if err := VerifySignature(key, in, now); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}
	if err := VerifySignature(key, in, now.Add(6*time.Minute)); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatal("stale timestamp accepted")
	}
	in.Body = []byte(`{"channel":"email"}`)
	if err := VerifySignature(key, in, now); err == nil {
		t.Fatal("tampered body accepted")
	}
	in.Body = body
	if err := VerifySignature("hb_test_other", in, now); err == nil {
		t.Fatal("other key accepted")
	}
	if err := VerifySignature(key, SignatureInput{}, now); err == nil {
		t.Fatal("missing signature accepted")
	}
	if err := VerifySignature(key, SignatureInput{Signature: sig, Timestamp: "yesterday"}, now); err == nil {
		t.Fatal("bad timestamp accepted")
	}
}

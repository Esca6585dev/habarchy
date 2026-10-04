package otp

import (
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestGenerate(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		c, err := generate(6)
		if err != nil || len(c) != 6 {
			t.Fatalf("%q %v", c, err)
		}
		for _, r := range c {
			if r < '0' || r > '9' {
				t.Fatalf("non digit in %q", c)
			}
		}
		seen[c] = true
	}
	if len(seen) < 45 {
		t.Fatalf("codes not random enough: %d unique", len(seen))
	}
}

func TestNormalizeAndDefaults(t *testing.T) {
	if to, err := normalize(domain.ChannelSMS, "65 12 34 56"); err != nil || to != "+99365123456" {
		t.Fatalf("%q %v", to, err)
	}
	if _, err := normalize(domain.ChannelSMS, "abc"); err == nil {
		t.Fatal("bad phone accepted")
	}
	if to, _ := normalize(domain.ChannelEmail, " A@B.TM "); to != "a@b.tm" {
		t.Fatal(to)
	}
	if defaultBody("", "ru") == defaultBody("", "tk") || defaultBody(domain.LocaleEN, "tk") == defaultBody(domain.LocaleTK, "tk") {
		t.Fatal("default bodies must differ per locale")
	}
	if isTemplateMissing(domain.ErrValidation.WithDetails(map[string]any{"template": "x"}), new(*domain.Error)) != true {
		t.Fatal("template missing not detected")
	}
	if isTemplateMissing(domain.ErrValidation.WithDetails(map[string]any{"body": "x"}), new(*domain.Error)) {
		t.Fatal("other validation error treated as template missing")
	}
}

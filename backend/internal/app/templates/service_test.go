package templates

import (
	"errors"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestInputValidate(t *testing.T) {
	in := Input{Key: " OTP ", Channel: domain.ChannelSMS, Body: "Kod: {{.code}}"}
	vars, err := in.validate(true)
	if err != nil || in.Key != "otp" || in.Locale != domain.LocaleTK || len(vars) != 1 {
		t.Fatalf("%v %+v %v", vars, in, err)
	}
	bad := Input{Key: "Bad Key!", Channel: "fax", Locale: "de", Subject: "x", Body: ""}
	_, err = bad.validate(true)
	var de *domain.Error
	if !errors.As(err, &de) || len(de.Details) != 4 {
		t.Fatalf("expected 4 validation details, got %v", err)
	}
	email := Input{Key: "welcome", Channel: domain.ChannelEmail, Body: "<p>hi</p>"}
	if _, err := email.validate(true); err == nil {
		t.Fatal("email without subject must fail")
	}
	syntax := Input{Key: "x", Channel: domain.ChannelPush, Body: "{{.oops"}
	if _, err := syntax.validate(true); !errors.Is(err, domain.ErrValidation) {
		t.Fatal("syntax error must be validation error")
	}
}

func TestPreviewRawReportsMissing(t *testing.T) {
	p, err := PreviewRaw(domain.ChannelSMS, "", "Salam {{.name}}, kod {{.code}}", map[string]any{"code": "1234"})
	if err != nil {
		t.Fatal(err)
	}
	if p.Body != "Salam {{name}}, kod 1234" || len(p.MissingVars) != 1 || p.MissingVars[0] != "name" {
		t.Fatalf("%+v", p)
	}
}

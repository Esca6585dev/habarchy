package providers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestValidateCredentials(t *testing.T) {
	ctx := context.Background()
	ok := map[domain.ProviderType]string{
		domain.ProviderHTTPSMS:     `{"url":"https://sms.example/send","body_template":"{\"to\":\"{{.To}}\"}"}`,
		domain.ProviderSMTP:        `{"host":"smtp.example","from_email":"no-reply@example.tm"}`,
		domain.ProviderFCM:         `{"service_account":{"type":"service_account","project_id":"demo"}}`,
		domain.ProviderTelegramBot: `{"bot_token":"123:abc"}`,
		domain.ProviderSMPP:        `{"host":"smsc.example","port":2775,"system_id":"user","password":"pw"}`,
	}
	for typ, creds := range ok {
		if err := ValidateCredentials(ctx, typ, json.RawMessage(creds)); err != nil {
			t.Errorf("%s: %v", typ, err)
		}
	}
	bad := map[domain.ProviderType]string{
		domain.ProviderHTTPSMS:     `{"method":"POST"}`,
		domain.ProviderSMTP:        `{"host":"x"}`,
		domain.ProviderFCM:         `{}`,
		domain.ProviderTelegramBot: `{"bot_token":"nocolon"}`,
		domain.ProviderSMPP:        `{"host":"x","system_id":"u","port":99999}`,
		"carrier_pigeon":           `{}`,
	}
	for typ, creds := range bad {
		if err := ValidateCredentials(ctx, typ, json.RawMessage(creds)); err == nil {
			t.Errorf("%s: expected error", typ)
		}
	}
	if err := ValidateCredentials(ctx, domain.ProviderSMTP, json.RawMessage(`not json`)); err == nil {
		t.Error("invalid json accepted")
	}
}

func TestMask(t *testing.T) {
	if mask("abc") != "••••" || mask("supersecrettoken") != "sup••••ken" {
		t.Fatal(mask("supersecrettoken"))
	}
}

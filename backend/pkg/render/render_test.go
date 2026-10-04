package render

import (
	"errors"
	"strings"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestValidateCollectsVars(t *testing.T) {
	vars, err := Validate("Salam {{.name}}", "Kod: {{.code}}. {{if .promo}}Promo: {{.promo}}{{end}} {{.name | upper}}")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(vars, ",") != "code,name,promo" {
		t.Fatalf("vars = %v", vars)
	}
	if _, err := Validate("", "{{.open"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("syntax error should be a validation error, got %v", err)
	}
}

func TestRender(t *testing.T) {
	res, err := Render(domain.ChannelSMS, "", "Siziň koduňyz: {{.code}}", map[string]any{"code": "4821"})
	if err != nil || res.Body != "Siziň koduňyz: 4821" {
		t.Fatalf("%+v %v", res, err)
	}
	_, err = Render(domain.ChannelSMS, "", "{{.code}} {{.name}}", map[string]any{"code": "1"})
	if !errors.Is(err, ErrMissingVars) || !strings.Contains(err.Error(), "name") {
		t.Fatalf("expected missing name, got %v", err)
	}
	// Email escapes data, SMS does not.
	html, err := Render(domain.ChannelEmail, "Hi {{.name}}", "<p>{{.name}}</p>", map[string]any{"name": "<b>x</b>"})
	if err != nil || html.Body != "<p>&lt;b&gt;x&lt;/b&gt;</p>" || html.Subject != "Hi <b>x</b>" {
		t.Fatalf("email: %+v %v", html, err)
	}
	sms, _ := Render(domain.ChannelSMS, "", "{{.name}}", map[string]any{"name": "<b>"})
	if sms.Body != "<b>" {
		t.Fatalf("sms should not escape: %q", sms.Body)
	}
	def, _ := Render(domain.ChannelPush, "", `{{default "Dost" .name}}`, map[string]any{"name": ""})
	if def.Body != "Dost" {
		t.Fatalf("default func: %q", def.Body)
	}
}

package contacts

import (
	"errors"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestNormalize(t *testing.T) {
	in := Input{Phone: " 865 12 34 56 ", Email: " Ali@Example.COM "}
	if err := in.normalize(); err != nil {
		t.Fatal(err)
	}
	if in.Phone != "+99365123456" || in.Email != "ali@example.com" || in.Locale != domain.LocaleTK || in.Tags == nil {
		t.Fatalf("%+v", in)
	}
	bad := Input{Phone: "12", Email: "nope", Locale: "xx"}
	err := bad.normalize()
	var de *domain.Error
	if !errors.As(err, &de) || len(de.Details) != 3 {
		t.Fatalf("expected 3 details, got %v", err)
	}
	if err := (&Input{}).normalize(); err == nil {
		t.Fatal("empty contact accepted")
	}
}

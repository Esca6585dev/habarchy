package delivery

import (
	"encoding/json"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
)

func TestEmailFrom(t *testing.T) {
	msg := &sqlcgen.Message{ToAddress: "a@b.tm", RenderedSubject: "Hi", RenderedBody: "<p>Salam</p>",
		Metadata: json.RawMessage(`{"reply_to":"r@b.tm","text":"Salam","attachments":[{"filename":"a.pdf","url":"https://x/a.pdf"}]}`)}
	em := emailFrom(msg)
	if em.HTML == "" || em.Text != "Salam" || em.ReplyTo != "r@b.tm" || len(em.Attachments) != 1 || em.Attachments[0].URL != "https://x/a.pdf" {
		t.Fatalf("%+v", em)
	}
	plain := emailFrom(&sqlcgen.Message{RenderedBody: "just text", Metadata: json.RawMessage(`{}`)})
	if plain.HTML != "" || plain.Text != "just text" {
		t.Fatalf("%+v", plain)
	}
	if m := metaStringMap(&sqlcgen.Message{Metadata: json.RawMessage(`{"data":{"a":1,"b":"x"}}`)}, "data"); m["a"] != "1" || m["b"] != "x" {
		t.Fatalf("%v", m)
	}
}

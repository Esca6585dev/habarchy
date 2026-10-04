package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
)

// TestGroupsBroadcast: groups created from the admin API with inline
// contacts, SMS broadcast via admin "send", WhatsApp broadcast through a
// fake Cloud API, and the public group endpoints with an API key.
func TestGroupsBroadcast(t *testing.T) {
	f := newFlowEnv(t)
	tok := bearer(f.adminTok)
	pbase := "/api/admin/projects/" + f.project.ID.String()

	// Fake WhatsApp Cloud API.
	var waMu sync.Mutex
	var waTo []string
	wa := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		waMu.Lock()
		waTo = append(waTo, fmt.Sprint(body["to"]))
		n := len(waTo)
		waMu.Unlock()
		fmt.Fprintf(w, `{"messages":[{"id":"wamid.%d"}]}`, n)
	}))
	t.Cleanup(wa.Close)
	f.expect(f.do("POST", pbase+"/providers", map[string]any{
		"name": "wa", "type": "whatsapp_cloud", "credentials": map[string]any{"access_token": "tok", "phone_number_id": "111", "endpoint": wa.URL},
	}, tok), 201)

	// Group with two inline contacts (one typed twice → deduplicated) and one existing contact.
	existing := decode[admin.ContactResponse](t, f.expect(f.do("POST", pbase+"/contacts", map[string]any{"name": "Aman", "phone": "+99365100001", "email": "aman@example.tm"}, tok), 201).Data)
	g := decode[admin.GroupResponse](t, f.expect(f.do("POST", pbase+"/groups", map[string]any{"name": "Işdeşler", "description": "office"}, tok), 201).Data)
	f.expect(f.do("POST", pbase+"/groups", map[string]any{"name": "Işdeşler"}, tok), 409)
	add := f.expect(f.do("POST", pbase+"/groups/"+g.ID.String()+"/members", map[string]any{
		"contact_ids": []string{existing.ID.String()},
		"contacts":    []map[string]any{{"name": "Maral", "phone": "+99365100002"}, {"name": "Maral again", "phone": "+993 65 10-00-02"}, {"name": "Only mail", "email": "mail@example.tm"}},
	}, tok), 200)
	var addRes struct {
		Added           int `json:"added"`
		CreatedContacts int `json:"created_contacts"`
	}
	_ = json.Unmarshal(add.Data, &addRes)
	if addRes.Added != 3 || addRes.CreatedContacts != 2 {
		t.Fatalf("add members: %+v", addRes)
	}
	var list []admin.GroupResponse
	_ = json.Unmarshal(f.expect(f.do("GET", pbase+"/groups", nil, tok), 200).Data, &list)
	if len(list) != 1 || list[0].MemberCount != 3 {
		t.Fatalf("groups %+v", list)
	}
	members := f.expect(f.do("GET", pbase+"/groups/"+g.ID.String()+"/members", nil, tok), 200)
	var mem []admin.ContactResponse
	_ = json.Unmarshal(members.Data, &mem)
	if len(mem) != 3 || mem[0].Name == "" {
		t.Fatalf("members %+v", mem)
	}

	// SMS broadcast: two members have phones, the mail-only one is rejected.
	r := f.expect(f.do("POST", pbase+"/messages/send", map[string]any{
		"channel": "sms", "body": "Salam {{.name}}! Ertir ýygnak 10:00.", "group_ids": []string{g.ID.String()},
	}, tok), 202)
	var meta struct {
		Accepted int             `json:"accepted"`
		Rejected json.RawMessage `json:"rejected"`
	}
	_ = json.Unmarshal(r.Meta, &meta)
	if meta.Accepted != 2 || string(meta.Rejected) == "[]" {
		t.Fatalf("sms send meta: %+v", meta)
	}
	var batch struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(r.Data, &batch)
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		f.gwCalls.mu.Lock()
		n := len(f.gwCalls.calls)
		f.gwCalls.mu.Unlock()
		if n >= 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	f.gwCalls.mu.Lock()
	texts := map[string]bool{}
	for _, call := range f.gwCalls.calls {
		texts[fmt.Sprint(call["text"])] = true
	}
	f.gwCalls.mu.Unlock()
	if !texts["Salam Aman! Ertir ýygnak 10:00."] || !texts["Salam Maral! Ertir ýygnak 10:00."] {
		t.Fatalf("gateway texts %v", texts)
	}

	// WhatsApp broadcast (falls back to the phone number when whatsapp is empty).
	r = f.expect(f.do("POST", pbase+"/messages/send", map[string]any{
		"channel": "whatsapp", "subject": "Habar", "body": "WhatsApp arkaly", "group_ids": []string{g.ID.String()},
	}, tok), 202)
	_ = json.Unmarshal(r.Meta, &meta)
	if meta.Accepted != 2 {
		t.Fatalf("wa accepted %d", meta.Accepted)
	}
	deadline = time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		waMu.Lock()
		n := len(waTo)
		waMu.Unlock()
		if n >= 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	waMu.Lock()
	if len(waTo) != 2 || (waTo[0] != "99365100001" && waTo[1] != "99365100001") {
		t.Fatalf("whatsapp recipients %v", waTo)
	}
	waMu.Unlock()
	// The message log shows the whatsapp messages as sent with the Graph id.
	log := f.expect(f.do("GET", pbase+"/messages?channel=whatsapp", nil, tok), 200)
	var msgs []map[string]any
	_ = json.Unmarshal(log.Data, &msgs)
	if len(msgs) != 2 {
		t.Fatalf("whatsapp log %d", len(msgs))
	}
	deadline = time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		_ = json.Unmarshal(f.expect(f.do("GET", pbase+"/messages?channel=whatsapp&status=sent", nil, tok), 200).Data, &msgs)
		if len(msgs) == 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if len(msgs) != 2 {
		t.Fatalf("whatsapp messages not sent: %v", msgs)
	}

	// Public API: groups with the API key, batch to a group.
	pg := f.expect(f.do("POST", "/api/v1/groups", map[string]any{"name": "Müşderiler"}, f.key()), 201)
	var pub struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(pg.Data, &pub)
	f.expect(f.do("POST", "/api/v1/groups/"+pub.ID+"/members", map[string]any{"external_ids": []string{"nope"}, "contacts": []map[string]any{{"phone": "+99365100003"}}}, f.key()), 200)
	rb := f.expect(f.do("POST", "/api/v1/messages/batch", map[string]any{"channel": "sms", "template": "otp", "data": map[string]any{"code": "7777", "minutes": 3}, "group_ids": []string{pub.ID}}, f.key()), 202)
	_ = json.Unmarshal(rb.Meta, &meta)
	if meta.Accepted != 1 {
		t.Fatalf("public batch accepted %d", meta.Accepted)
	}
	f.expect(f.do("DELETE", "/api/v1/groups/"+pub.ID+"/members/"+existing.ID.String(), nil, f.key()), 404)
	f.expect(f.do("DELETE", "/api/v1/groups/"+pub.ID, nil, f.key()), 204)
	f.expect(f.do("GET", "/api/v1/groups/"+pub.ID, nil, f.key()), 404)

	// Contact detail lists its groups; deleting the group keeps the contact.
	det := f.expect(f.do("GET", pbase+"/contacts/"+existing.ID.String(), nil, tok), 200)
	var detail struct {
		Groups []admin.GroupResponse `json:"groups"`
	}
	_ = json.Unmarshal(det.Data, &detail)
	if len(detail.Groups) != 1 {
		t.Fatalf("contact groups %+v", detail.Groups)
	}
	f.expect(f.do("DELETE", pbase+"/groups/"+g.ID.String(), nil, tok), 204)
	f.expect(f.do("GET", pbase+"/contacts/"+existing.ID.String(), nil, tok), 200)
}

package e2e

import (
	"encoding/json"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
)

// TestContactSoftDelete: delete keeps the row in the trash, it disappears
// from the active list and from its group (and from broadcasts), restore
// brings it back, purge removes it for good.
func TestContactSoftDelete(t *testing.T) {
	f := newFlowEnv(t)
	tok := bearer(f.adminTok)
	pbase := "/api/admin/projects/" + f.project.ID.String()

	g := decode[admin.GroupResponse](t, f.expect(f.do("POST", pbase+"/groups", map[string]any{"name": "Trash test"}, tok), 201).Data)
	a := decode[admin.ContactResponse](t, f.expect(f.do("POST", pbase+"/contacts", map[string]any{"name": "Aman", "phone": "+99365200001", "external_id": "u-200001"}, tok), 201).Data)
	b := decode[admin.ContactResponse](t, f.expect(f.do("POST", pbase+"/contacts", map[string]any{"name": "Maral", "phone": "+99365200002"}, tok), 201).Data)
	f.expect(f.do("POST", pbase+"/groups/"+g.ID.String()+"/members", map[string]any{"contact_ids": []string{a.ID.String(), b.ID.String()}}, tok), 200)

	count := func(path string) (rows []admin.ContactResponse, total int) {
		r := f.expect(f.do("GET", path, nil, tok), 200)
		_ = json.Unmarshal(r.Data, &rows)
		var meta struct {
			Total int `json:"total"`
		}
		_ = json.Unmarshal(r.Meta, &meta)
		return rows, meta.Total
	}

	// Soft delete Aman.
	f.expect(f.do("DELETE", pbase+"/contacts/"+a.ID.String(), nil, tok), 204)
	if rows, total := count(pbase + "/contacts"); total != 1 || len(rows) != 1 || rows[0].ID != b.ID {
		t.Fatalf("active list after delete: %d %+v", total, rows)
	}
	// It is in the trash, with deleted_at set.
	trash, ttotal := count(pbase + "/contacts/trash")
	if ttotal != 1 || len(trash) != 1 || trash[0].ID != a.ID || trash[0].DeletedAt == nil {
		t.Fatalf("trash: %d %+v", ttotal, trash)
	}
	// Gone from the group members and the member count.
	members, _ := count(pbase + "/groups/" + g.ID.String() + "/members")
	if len(members) != 1 || members[0].ID != b.ID {
		t.Fatalf("group members after delete: %+v", members)
	}
	gg := decode[admin.GroupResponse](t, f.expect(f.do("GET", pbase+"/groups/"+g.ID.String(), nil, tok), 200).Data)
	if gg.MemberCount != 1 {
		t.Fatalf("group member_count: %d", gg.MemberCount)
	}
	// Broadcasting to the group only reaches the surviving member.
	r := f.expect(f.do("POST", pbase+"/messages/send", map[string]any{"channel": "sms", "body": "hi", "group_ids": []string{g.ID.String()}}, tok), 202)
	var meta struct {
		Accepted int `json:"accepted"`
	}
	_ = json.Unmarshal(r.Meta, &meta)
	if meta.Accepted != 1 {
		t.Fatalf("broadcast reached %d, want 1", meta.Accepted)
	}
	// Editing a trashed contact is not found.
	if st := f.do("PUT", pbase+"/contacts/"+a.ID.String(), map[string]any{"name": "x", "phone": "+99365200001"}, tok).Status; st != 404 {
		t.Fatalf("edit trashed: %d", st)
	}

	// Restore Aman back into the group.
	restored := decode[admin.ContactResponse](t, f.expect(f.do("POST", pbase+"/contacts/"+a.ID.String()+"/restore", nil, tok), 200).Data)
	if restored.DeletedAt != nil {
		t.Fatalf("restored still deleted: %+v", restored)
	}
	if _, total := count(pbase + "/contacts"); total != 2 {
		t.Fatalf("active list after restore: %d", total)
	}
	if members, _ := count(pbase + "/groups/" + g.ID.String() + "/members"); len(members) != 2 {
		t.Fatalf("group members after restore: %+v", members)
	}

	// Re-importing the same external_id while the original is trashed creates a fresh active one.
	f.expect(f.do("DELETE", pbase+"/contacts/"+a.ID.String(), nil, tok), 204)
	fresh := decode[admin.ContactResponse](t, f.expect(f.do("POST", pbase+"/contacts", map[string]any{"name": "Aman 2", "phone": "+99365200003", "external_id": "u-200001"}, tok), 201).Data)
	if fresh.ID == a.ID {
		t.Fatalf("expected a new contact, got the trashed one")
	}

	// Purge the trashed Aman permanently; trash empties, restore now 404.
	f.expect(f.do("DELETE", pbase+"/contacts/"+a.ID.String()+"/purge", nil, tok), 204)
	if _, ttotal := count(pbase + "/contacts/trash"); ttotal != 0 {
		t.Fatalf("trash after purge: %d", ttotal)
	}
	if st := f.do("POST", pbase+"/contacts/"+a.ID.String()+"/restore", nil, tok).Status; st != 404 {
		t.Fatalf("restore after purge: %d", st)
	}

	// Purge-all empties the trash in one call.
	f.expect(f.do("DELETE", pbase+"/contacts/"+b.ID.String(), nil, tok), 204)
	f.expect(f.do("DELETE", pbase+"/contacts/"+fresh.ID.String(), nil, tok), 204)
	pa := f.expect(f.do("DELETE", pbase+"/contacts/trash", nil, tok), 200)
	var purged struct {
		Purged int `json:"purged"`
	}
	_ = json.Unmarshal(pa.Data, &purged)
	if purged.Purged != 2 {
		t.Fatalf("purge all: %+v", purged)
	}
}

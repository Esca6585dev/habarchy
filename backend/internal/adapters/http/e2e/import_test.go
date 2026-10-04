package e2e

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contactimport"
)

func (f *flowEnv) upload(t *testing.T, path, filename string, content []byte, fields map[string]string, headers map[string]string) resp {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	for k, v := range fields {
		_ = mw.WriteField(k, v)
	}
	fw, _ := mw.CreateFormFile("file", filename)
	_, _ = fw.Write(content)
	_ = mw.Close()
	req := httptest.NewRequest("POST", path, &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := f.app.Test(req, 15_000)
	if err != nil {
		t.Fatal(err)
	}
	var out resp
	_ = json.NewDecoder(res.Body).Decode(&out)
	out.Status = res.StatusCode
	return out
}

// TestContactImport: CSV upload with dry run, then a real import into a
// group; a second import of the same file only merges; the public API
// endpoint accepts a vCard.
func TestContactImport(t *testing.T) {
	f := newFlowEnv(t)
	tok := bearer(f.adminTok)
	pbase := "/api/admin/projects/" + f.project.ID.String()
	g := decode[admin.GroupResponse](t, f.expect(f.do("POST", pbase+"/groups", map[string]any{"name": "Import"}, tok), 201).Data)

	csv := []byte("Ady,Telefon,E-mail\nAman Amanow,+99365123456,aman@example.tm\nMaral,65123457,\nBroken,abc,\n")
	dry := f.upload(t, pbase+"/contacts/import", "list.csv", csv, map[string]string{"dry_run": "true", "group_id": g.ID.String()}, tok)
	if dry.Status != 200 {
		t.Fatalf("dry run: %d %v", dry.Status, dry.Error)
	}
	var res contactimport.Result
	_ = json.Unmarshal(dry.Data, &res)
	if !res.DryRun || res.Total != 2 || res.Created != 0 || len(res.Preview) != 2 {
		t.Fatalf("dry: %+v", res)
	}

	imported := f.upload(t, pbase+"/contacts/import", "list.csv", csv, map[string]string{"group_id": g.ID.String(), "tag": "excel-2026"}, tok)
	_ = json.Unmarshal(imported.Data, &res)
	if imported.Status != 200 || res.Created != 2 || res.GroupAdded != 2 {
		t.Fatalf("import: %d %+v %v", imported.Status, res, imported.Error)
	}
	again := f.upload(t, pbase+"/contacts/import", "list.csv", csv, nil, tok)
	_ = json.Unmarshal(again.Data, &res)
	if res.Created != 0 || res.Unchanged != 2 {
		t.Fatalf("re-import: %+v", res)
	}
	var members []admin.ContactResponse
	_ = json.Unmarshal(f.expect(f.do("GET", pbase+"/groups/"+g.ID.String()+"/members", nil, tok), 200).Data, &members)
	if len(members) != 2 || members[0].Name == "" || len(members[0].Tags) != 1 {
		t.Fatalf("members: %+v", members)
	}

	// Unsupported type.
	if r := f.upload(t, pbase+"/contacts/import", "list.pdf", []byte("x"), nil, tok); r.Status != 400 {
		t.Fatalf("pdf: %d", r.Status)
	}
	// Public API with the key: vCard.
	vcf := []byte("BEGIN:VCARD\nVERSION:3.0\nFN:Public Person\nTEL:+99365100009\nEND:VCARD\n")
	pub := f.upload(t, "/api/v1/contacts/import", "c.vcf", vcf, nil, f.key())
	_ = json.Unmarshal(pub.Data, &res)
	if pub.Status != 200 || res.Created != 1 {
		t.Fatalf("public import: %d %+v %v", pub.Status, res, pub.Error)
	}
	// Google not configured → clear message.
	r := f.do("GET", pbase+"/contacts/import/google/url", nil, tok)
	if r.Status != 400 {
		t.Fatalf("google url without config: %d", r.Status)
	}
}

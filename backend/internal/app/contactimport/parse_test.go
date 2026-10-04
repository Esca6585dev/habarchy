package contactimport_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contactimport"
)

func TestCSVWithHeaderAndFreeLines(t *testing.T) {
	csv := "Ady;Telefon;E-mail;Topar\nAman Amanow;+993 65 12-34-56;aman@example.tm;vip,işdeş\nMaral;65123457;;\n;;;\nBoş;;;\n"
	ins, _, err := contactimport.ParseText(strings.NewReader(csv))
	if err != nil || len(ins) != 2 {
		t.Fatalf("%v %+v", err, ins)
	}
	if ins[0].Name != "Aman Amanow" || ins[0].Phone != "+993651234-56"[:0]+"+99365123456" || ins[0].Email != "aman@example.tm" || len(ins[0].Tags) != 2 {
		t.Fatalf("row 1: %+v", ins[0])
	}
	if ins[1].Phone != "65123457" {
		t.Fatalf("row 2: %+v", ins[1])
	}
	free := "Aman Amanow, +99365123456\nmaral@example.tm Maral\nU0123ABCD Ops team\njust words\n"
	ins, _, err = contactimport.ParseText(strings.NewReader(free))
	if err != nil || len(ins) != 3 {
		t.Fatalf("free: %v %+v", err, ins)
	}
	if ins[0].Name != "Aman Amanow" || ins[0].Phone != "+99365123456" || ins[1].Email != "maral@example.tm" || ins[1].Name != "Maral" || ins[2].SlackID != "U0123ABCD" {
		t.Fatalf("free rows: %+v", ins)
	}
}

func TestXLSX(t *testing.T) {
	f := excelize.NewFile()
	_ = f.SetSheetRow("Sheet1", "A1", &[]any{"Name", "Phone", "Email"})
	_ = f.SetSheetRow("Sheet1", "A2", &[]any{"Aman", 99365123456, "aman@example.tm"})
	_ = f.SetSheetRow("Sheet1", "A3", &[]any{"Maral", "+99365123457", ""})
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	ins, _, err := contactimport.Parse("list.xlsx", &buf)
	if err != nil || len(ins) != 2 {
		t.Fatalf("%v %+v", err, ins)
	}
	if ins[0].Phone != "99365123456" || ins[0].Email != "aman@example.tm" || ins[1].Phone != "+99365123457" {
		t.Fatalf("%+v", ins)
	}
}

func TestDOCXTableAndParagraphs(t *testing.T) {
	doc := `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>
<w:p><w:r><w:t>Müşderiler</w:t></w:r></w:p>
<w:tbl><w:tr><w:tc><w:p><w:r><w:t>Ady</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>Telefon</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>Aman</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>+99365</w:t></w:r><w:r><w:t>123456</w:t></w:r></w:p></w:tc></w:tr></w:tbl>
<w:p><w:r><w:t>Maral, maral@example.tm</w:t></w:r></w:p>
</w:body></w:document>`
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("word/document.xml")
	_, _ = w.Write([]byte(doc))
	_ = zw.Close()
	ins, _, err := contactimport.Parse("clients.docx", &buf)
	if err != nil {
		t.Fatal(err)
	}
	// Header row detected ("Ady", "Telefon") → Aman via header mapping; the title paragraph has no address and is dropped; Maral via free line.
	if len(ins) != 2 || ins[0].Name != "Aman" || ins[0].Phone != "+99365123456" || ins[1].Email != "maral@example.tm" {
		t.Fatalf("%+v", ins)
	}
}

func TestVCard(t *testing.T) {
	vcf := "BEGIN:VCARD\r\nVERSION:3.0\r\nN:Amanow;Aman;;;\r\nFN:Aman Amanow\r\nitem1.TEL;type=CELL:+993 65 123456\r\nTEL;TYPE=WORK:+99312345678\r\nEMAIL;TYPE=INTERNET:Aman@Example.tm\r\nCATEGORIES:VIP,Işdeşler\r\nUID:abc-123\r\nEND:VCARD\r\n" +
		"BEGIN:VCARD\r\nVERSION:2.1\r\nN;ENCODING=QUOTED-PRINTABLE;CHARSET=UTF-8:;=4D=61=72=61=6C\r\nTEL;CELL:65123457\r\nEND:VCARD\r\n" +
		"BEGIN:VCARD\r\nVERSION:4.0\r\nFN:No address\r\nEND:VCARD\r\n"
	ins, _, err := contactimport.Parse("contacts.vcf", strings.NewReader(vcf))
	if err != nil || len(ins) != 2 {
		t.Fatalf("%v %+v", err, ins)
	}
	a := ins[0]
	if a.Name != "Aman Amanow" || a.Phone != "+99365123456" || a.WhatsApp != "+99312345678" || a.Email != "aman@example.tm" || a.ExternalID != "vcard:abc-123" || len(a.Tags) != 2 {
		t.Fatalf("card 1: %+v", a)
	}
	if ins[1].Name != "Maral" || ins[1].Phone != "65123457" {
		t.Fatalf("card 2: %+v", ins[1])
	}
	if _, _, err := contactimport.Parse("x.pdf", strings.NewReader("")); err == nil {
		t.Fatal("pdf accepted")
	}
}

func TestGoogleConnector(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			_ = r.ParseForm()
			if r.Form.Get("code") != "the-code" || r.Form.Get("client_secret") != "sec" {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
				return
			}
			_, _ = w.Write([]byte(`{"access_token":"tok"}`))
		case "/people":
			if r.Header.Get("Authorization") != "Bearer tok" {
				w.WriteHeader(401)
				return
			}
			if r.URL.Query().Get("pageToken") == "" {
				_, _ = w.Write([]byte(`{"connections":[{"resourceName":"people/c1","names":[{"displayName":"Aman Amanow"}],"phoneNumbers":[{"value":"+993 65 123456","canonicalForm":"+99365123456"}],"emailAddresses":[{"value":"Aman@example.tm"}]}],"nextPageToken":"p2"}`))
				return
			}
			_, _ = w.Write([]byte(`{"connections":[{"resourceName":"people/c2","names":[{"displayName":"No phone"}]},{"resourceName":"people/c3","phoneNumbers":[{"value":"65123457"}]}]}`))
		}
	}))
	defer srv.Close()
	g := contactimport.NewGoogleConnector("id", "sec", "https://h/cb", []byte("state-secret"))
	g.TokenURL, g.PeopleURL, g.HTTP = srv.URL+"/token", srv.URL+"/people", srv.Client()
	st := g.SignState(contactimport.State{ReturnTo: "/contacts"})
	if parsed, err := g.ParseState(st); err != nil || parsed.ReturnTo != "/contacts" {
		t.Fatalf("state: %v", err)
	}
	if _, err := g.ParseState(st + "x"); err == nil {
		t.Fatal("tampered state accepted")
	}
	if !strings.Contains(g.AuthorizationURL(st), "contacts.readonly") {
		t.Fatal("scope missing")
	}
	tok, err := g.Exchange(context.Background(), "the-code")
	if err != nil || tok != "tok" {
		t.Fatalf("exchange: %v", err)
	}
	ins, err := g.FetchContacts(context.Background(), tok)
	if err != nil || len(ins) != 2 || ins[0].Phone != "+99365123456" || ins[0].Email != "aman@example.tm" || ins[0].ExternalID != "google:c1" || ins[1].Phone != "65123457" {
		t.Fatalf("%v %+v", err, ins)
	}
}

func TestCardDAVICloudLikeFlow(t *testing.T) {
	var mu int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != "aman@icloud.com" || pass != "abcd-efgh-ijkl-mnop" {
			w.WriteHeader(401)
			return
		}
		mu++
		switch {
		case r.Method == "PROPFIND" && r.URL.Path == "/":
			w.WriteHeader(207)
			_, _ = w.Write([]byte(`<?xml version="1.0"?><d:multistatus xmlns:d="DAV:"><d:response><d:href>/</d:href><d:propstat><d:prop><d:current-user-principal><d:href>/123/principal/</d:href></d:current-user-principal></d:prop><d:status>HTTP/1.1 200 OK</d:status></d:propstat></d:response></d:multistatus>`))
		case r.Method == "PROPFIND" && r.URL.Path == "/123/principal/":
			// iCloud redirects the home set to another host; emulate with a redirect to the same server.
			w.Header().Set("Location", "/p01/123/principal/")
			w.WriteHeader(301)
		case r.Method == "PROPFIND" && r.URL.Path == "/p01/123/principal/":
			w.WriteHeader(207)
			_, _ = w.Write([]byte(`<d:multistatus xmlns:d="DAV:" xmlns:card="urn:ietf:params:xml:ns:carddav"><d:response><d:href>/p01/123/principal/</d:href><d:propstat><d:prop><card:addressbook-home-set><d:href>/p01/123/carddavhome/</d:href></card:addressbook-home-set></d:prop></d:propstat></d:response></d:multistatus>`))
		case r.Method == "PROPFIND" && r.URL.Path == "/p01/123/carddavhome/":
			w.WriteHeader(207)
			_, _ = w.Write([]byte(`<d:multistatus xmlns:d="DAV:" xmlns:card="urn:ietf:params:xml:ns:carddav"><d:response><d:href>/p01/123/carddavhome/</d:href><d:propstat><d:prop><d:resourcetype><d:collection/></d:resourcetype></d:prop></d:propstat></d:response><d:response><d:href>/p01/123/carddavhome/card/</d:href><d:propstat><d:prop><d:resourcetype><d:collection/><card:addressbook/></d:resourcetype><d:displayname>card</d:displayname></d:prop></d:propstat></d:response></d:multistatus>`))
		case r.Method == "REPORT" && r.URL.Path == "/p01/123/carddavhome/card/":
			body, _ := json.Marshal("BEGIN:VCARD\nVERSION:3.0\nFN:Aman Amanow\nTEL;type=CELL:+99365123456\nEMAIL:aman@example.tm\nEND:VCARD")
			var vcard string
			_ = json.Unmarshal(body, &vcard)
			w.WriteHeader(207)
			_, _ = w.Write([]byte(`<d:multistatus xmlns:d="DAV:" xmlns:card="urn:ietf:params:xml:ns:carddav"><d:response><d:href>/p01/123/carddavhome/card/1.vcf</d:href><d:propstat><d:prop><d:getetag>"1"</d:getetag><card:address-data>` + vcard + `</card:address-data></d:prop></d:propstat></d:response></d:multistatus>`))
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()
	c := contactimport.NewCardDAVClient(srv.Client())
	c.HTTP.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	ins, err := c.Fetch(context.Background(), srv.URL+"/", "aman@icloud.com", "abcd-efgh-ijkl-mnop")
	if err != nil || len(ins) != 1 || ins[0].Name != "Aman Amanow" || ins[0].Phone != "+99365123456" {
		t.Fatalf("%v %+v", err, ins)
	}
	if _, err := c.Fetch(context.Background(), srv.URL+"/", "aman@icloud.com", "wrong"); err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("wrong password: %v", err)
	}
	if _, err := c.Fetch(context.Background(), "http://example.com/", "u", "p"); err == nil {
		t.Fatal("plain http to remote host accepted")
	}
}

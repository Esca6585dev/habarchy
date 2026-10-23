package contactimport

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestParseSQLite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "src.sqlite")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE contacts(id INTEGER PRIMARY KEY, title TEXT, phone_number TEXT, is_sent INTEGER)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO contacts(title, phone_number, is_sent) VALUES ('Aman','+99365123456',0),('Merjen','65234567',1)`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	inputs, _, err := ParseSQLite(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("parse sqlite: %v", err)
	}
	if len(inputs) != 2 {
		t.Fatalf("want 2 contacts, got %d", len(inputs))
	}
	byPhone := map[string]string{}
	for _, in := range inputs {
		if in.Phone == "" {
			t.Errorf("row has empty phone: %+v", in)
		}
		byPhone[in.Phone] = in.Name
	}
	if byPhone["+99365123456"] != "Aman" {
		t.Errorf("expected title mapped to name; got %q", byPhone["+99365123456"])
	}
}

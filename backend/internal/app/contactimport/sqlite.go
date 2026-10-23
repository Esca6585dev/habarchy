package contactimport

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite" // pure-Go SQLite driver, registered as "sqlite"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
)

// ParseSQLite reads contacts from an uploaded SQLite database (.db / .sqlite).
// It picks the first user table that has a phone-like column and maps its
// recognized columns (phone, name, email, tags…) the same way CSV/XLSX do.
func ParseSQLite(r io.Reader) ([]contacts.Input, []RowError, error) {
	data, err := io.ReadAll(io.LimitReader(r, 50<<20))
	if err != nil {
		return nil, nil, err
	}
	f, err := os.CreateTemp("", "habarchy-import-*.sqlite")
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = os.Remove(f.Name()) }()
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	_ = f.Close()

	// Open read-only; immutable avoids creating -wal/-shm sidecars.
	db, err := sql.Open("sqlite", "file:"+f.Name()+"?mode=ro&immutable=1")
	if err != nil {
		return nil, nil, fmt.Errorf("sqlite: %w", err)
	}
	defer func() { _ = db.Close() }()

	tables, err := listTables(db)
	if err != nil {
		return nil, nil, fmt.Errorf("sqlite: %w", err)
	}
	for _, tbl := range tables {
		cols, err := tableColumns(db, tbl)
		if err != nil || len(cols) == 0 {
			continue
		}
		hasPhone := false
		for _, c := range cols {
			if headerMap[strings.ToLower(strings.TrimSpace(c))] == "phone" {
				hasPhone = true
				break
			}
		}
		if !hasPhone {
			continue
		}
		rows, err := readTable(db, tbl, cols)
		if err != nil {
			continue
		}
		if out := FromRows(rows); len(out) > 0 {
			return out, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("sqlite: no table with a recognizable phone column (phone/telefon/nomer) found")
}

func listTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err == nil {
			out = append(out, n)
		}
	}
	return out, rows.Err()
}

func tableColumns(db *sql.DB, tbl string) ([]string, error) {
	rows, err := db.Query(`PRAGMA table_info(` + quoteIdent(tbl) + `)`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

// readTable returns the rows of tbl as [][]string with a header row of the
// column names first, so FromRows can map them like a spreadsheet.
func readTable(db *sql.DB, tbl string, cols []string) ([][]string, error) {
	quoted := make([]string, len(cols))
	for i, c := range cols {
		quoted[i] = quoteIdent(c)
	}
	q := "SELECT " + strings.Join(quoted, ", ") + " FROM " + quoteIdent(tbl) + " LIMIT " + strconv.Itoa(MaxRows)
	rows, err := db.Query(q) //nolint:gosec // identifiers are quoted; limit is a constant
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := [][]string{append([]string(nil), cols...)}
	for rows.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		rec := make([]string, len(cols))
		for i := range vals {
			if vals[i].Valid {
				rec[i] = strings.TrimSpace(vals[i].String)
			}
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

func quoteIdent(id string) string {
	return `"` + strings.ReplaceAll(id, `"`, `""`) + `"`
}

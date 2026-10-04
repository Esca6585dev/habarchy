package contactimport

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
)

// ParseXLSX reads the first non-empty sheet of an Excel workbook.
func ParseXLSX(r io.Reader) ([]contacts.Input, []RowError, error) {
	data, err := io.ReadAll(io.LimitReader(r, 20<<20))
	if err != nil {
		return nil, nil, err
	}
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("xlsx: %w", err)
	}
	defer func() { _ = f.Close() }()
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet, excelize.Options{RawCellValue: true})
		if err != nil {
			return nil, nil, fmt.Errorf("xlsx: %w", err)
		}
		if len(rows) == 0 {
			continue
		}
		out := FromRows(rows)
		if len(out) > 0 {
			return out, nil, nil
		}
	}
	return nil, nil, errEmpty
}

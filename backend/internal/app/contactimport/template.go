package contactimport

import (
	"bytes"
	"encoding/csv"

	"github.com/xuri/excelize/v2"
)

// Column documents one import column for the "what the file should look like"
// guide shown before upload. The names match the header aliases the parser
// recognizes (see headerMap); only `phone` is required.
type Column struct {
	Key      string   `json:"key"`
	Required bool     `json:"required"`
	Aliases  []string `json:"aliases"`
	Example  string   `json:"example"`
	Note     string   `json:"note,omitempty"`
}

// TemplateColumns is the recommended column set, in order.
func TemplateColumns() []Column {
	return []Column{
		{Key: "phone", Required: true, Aliases: []string{"telefon", "nomer", "number", "msisdn"}, Example: "+99365123456", Note: "E.164 ýa-da 8 belgili türkmen nomeri (65123456). Diňe şu sütün hökmany."},
		{Key: "name", Required: false, Aliases: []string{"ady", "имя", "full name"}, Example: "Aman Myradow"},
		{Key: "email", Required: false, Aliases: []string{"e-mail", "poçta"}, Example: "aman@example.tm"},
		{Key: "tags", Required: false, Aliases: []string{"topar", "group", "группа"}, Example: "müşderi;vip", Note: "Birnäçe bellik üçin ; ýa , bilen bölüň."},
		{Key: "locale", Required: false, Aliases: []string{"dil", "lang"}, Example: "tk", Note: "tk, ru ýa en."},
		{Key: "external_id", Required: false, Aliases: []string{"id", "customer_id"}, Example: "cust-1001", Note: "Siziň ulgamyňyzdaky ID — gaýtadan import üçin."},
	}
}

// sampleRows is the header row plus two example rows (used for both the CSV
// and XLSX samples, so they stay identical).
func sampleRows() [][]string {
	return [][]string{
		{"phone", "name", "email", "tags", "locale"},
		{"+99365123456", "Aman Myradow", "aman@example.tm", "müşderi;vip", "tk"},
		{"65234567", "Merjen Nurowa", "", "müşderi", "tk"},
	}
}

// SampleCSV returns a downloadable example CSV (UTF-8 with BOM so Excel on
// Windows reads Turkmen letters correctly).
func SampleCSV() []byte {
	var buf bytes.Buffer
	buf.Write([]byte("\xef\xbb\xbf")) // BOM
	w := csv.NewWriter(&buf)
	_ = w.WriteAll(sampleRows())
	w.Flush()
	return buf.Bytes()
}

// SampleXLSX returns a downloadable example .xlsx with a header row and two
// example rows, matching what the importer expects.
func SampleXLSX() ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	const sheet = "Contacts"
	idx, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")

	header, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	for r, row := range sampleRows() {
		for cN, val := range row {
			cell, _ := excelize.CoordinatesToCellName(cN+1, r+1)
			_ = f.SetCellStr(sheet, cell, val)
			if r == 0 {
				_ = f.SetCellStyle(sheet, cell, cell, header)
			}
		}
	}
	_ = f.SetColWidth(sheet, "A", "E", 20)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

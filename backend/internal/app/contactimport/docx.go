package contactimport

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
)

// ParseDOCX extracts table rows and paragraphs from a Word document.
// Tables become cell rows; paragraphs are parsed as free lines.
func ParseDOCX(r io.Reader) ([]contacts.Input, []RowError, error) {
	data, err := io.ReadAll(io.LimitReader(r, 20<<20))
	if err != nil {
		return nil, nil, err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, nil, fmt.Errorf("docx: %w", err)
	}
	var doc io.ReadCloser
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			if doc, err = f.Open(); err != nil {
				return nil, nil, err
			}
			break
		}
	}
	if doc == nil {
		return nil, nil, fmt.Errorf("docx: word/document.xml missing")
	}
	defer func() { _ = doc.Close() }()
	rows, err := docxRows(doc)
	if err != nil {
		return nil, nil, fmt.Errorf("docx: %w", err)
	}
	out := FromRows(rows)
	if len(out) == 0 {
		return nil, nil, errEmpty
	}
	return out, nil, nil
}

// docxRows walks the XML: each <w:tr> becomes a row of cell texts, each
// paragraph outside a table becomes a single-cell row.
func docxRows(r io.Reader) ([][]string, error) {
	dec := xml.NewDecoder(r)
	var rows [][]string
	var row []string
	var cell, para strings.Builder
	inRow, inCell, inPara := false, false, false
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "tr":
				inRow, row = true, nil
			case "tc":
				inCell = true
				cell.Reset()
			case "p":
				inPara = true
				para.Reset()
			case "tab":
				if inCell {
					cell.WriteByte(' ')
				} else if inPara {
					para.WriteByte('\t')
				}
			}
		case xml.CharData:
			if inCell {
				cell.Write(t)
			} else if inPara {
				para.Write(t)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "tc":
				inCell = false
				row = append(row, strings.TrimSpace(cell.String()))
			case "tr":
				inRow = false
				if strings.Join(row, "") != "" {
					rows = append(rows, row)
				}
			case "p":
				if inPara && !inRow && !inCell {
					if s := strings.TrimSpace(para.String()); s != "" {
						rows = append(rows, []string{s})
					}
				}
				inPara = false
			}
		}
	}
	return rows, nil
}

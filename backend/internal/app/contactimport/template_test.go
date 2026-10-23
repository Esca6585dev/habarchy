package contactimport

import (
	"bytes"
	"testing"
)

func TestTemplateColumnsPhoneRequired(t *testing.T) {
	cols := TemplateColumns()
	if len(cols) == 0 || cols[0].Key != "phone" || !cols[0].Required {
		t.Fatalf("expected phone to be the first, required column; got %+v", cols)
	}
}

// The downloadable samples must parse back through the importer, so the guide
// we show can never drift from what the parser accepts.
func TestSampleCSVRoundTrips(t *testing.T) {
	inputs, _, err := Parse("habarchy-contacts-template.csv", bytes.NewReader(SampleCSV()))
	if err != nil {
		t.Fatalf("parse sample csv: %v", err)
	}
	if len(inputs) != 2 {
		t.Fatalf("expected 2 contacts from the sample, got %d", len(inputs))
	}
	for i, in := range inputs {
		if in.Phone == "" {
			t.Errorf("row %d: empty phone", i)
		}
	}
}

func TestSampleXLSXRoundTrips(t *testing.T) {
	b, err := SampleXLSX()
	if err != nil {
		t.Fatalf("build sample xlsx: %v", err)
	}
	inputs, _, err := Parse("habarchy-contacts-template.xlsx", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("parse sample xlsx: %v", err)
	}
	if len(inputs) != 2 {
		t.Fatalf("expected 2 contacts from the xlsx sample, got %d", len(inputs))
	}
}

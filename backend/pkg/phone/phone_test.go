package phone

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		in, region, want string
		wantErr          bool
	}{
		{"+99365123456", "", "+99365123456", false},
		{"99365123456", "", "+99365123456", false},
		{"865123456", "", "+99365123456", false},
		{"65 12 34 56", "TM", "+99365123456", false},
		{"+993 (65) 12-34-56", "", "+99365123456", false},
		{"0099365123456", "", "+99365123456", false},
		{"+7 916 123-45-67", "", "+79161234567", false},
		{"8 916 123 45 67", "RU", "+79161234567", false},
		{"+1 (415) 555-2671", "", "+14155552671", false},
		{"", "", "", true},
		{"hello", "", "", true},
		{"+993651", "", "", true},        // too short
		{"+9936512345678", "", "", true}, // too long
	}
	for _, c := range cases {
		got, err := Normalize(c.in, c.region)
		if c.wantErr {
			if err == nil {
				t.Errorf("Normalize(%q) expected error, got %q", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Normalize(%q) unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRegion(t *testing.T) {
	if r := Region("+99365123456"); r != "TM" {
		t.Errorf("Region = %q, want TM", r)
	}
	if r := Region("nope"); r != "" {
		t.Errorf("Region of garbage = %q, want empty", r)
	}
}

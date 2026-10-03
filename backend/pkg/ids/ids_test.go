package ids

import (
	"sort"
	"testing"
	"time"
)

func TestNewIsV7AndOrdered(t *testing.T) {
	var got []string
	for i := 0; i < 5; i++ {
		id := New()
		if id.Version() != 7 {
			t.Fatalf("version %d, want 7", id.Version())
		}
		got = append(got, id.String())
		time.Sleep(2 * time.Millisecond)
	}
	if !sort.StringsAreSorted(got) {
		t.Fatalf("v7 ids should sort by creation time: %v", got)
	}
}

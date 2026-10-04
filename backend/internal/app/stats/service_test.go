package stats

import (
	"testing"
	"time"
)

func TestLastDays(t *testing.T) {
	w := LastDays(7)
	if w.To.Sub(w.From) != 7*24*time.Hour {
		t.Fatalf("window %v", w)
	}
	if !w.To.After(time.Now()) {
		t.Fatal("window must include now")
	}
	if w.From.Hour() != 0 || w.From.Minute() != 0 {
		t.Fatal("window must start at midnight UTC")
	}
}

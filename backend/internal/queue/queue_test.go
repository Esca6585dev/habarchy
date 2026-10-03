package queue

import (
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func TestForChannel(t *testing.T) {
	for _, ch := range domain.AllChannels {
		q, task, err := ForChannel(ch)
		if err != nil || q == "" || task == "" {
			t.Fatalf("ForChannel(%s) = %q %q %v", ch, q, task, err)
		}
		if _, ok := Weights()[q]; !ok {
			t.Fatalf("queue %q has no weight", q)
		}
	}
	if _, _, err := ForChannel(domain.ChannelAuto); err == nil {
		t.Fatal("auto must not map to a queue")
	}
}

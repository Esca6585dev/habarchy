package events

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

func TestPubSub(t *testing.T) {
	url := os.Getenv("HABARCHY_TEST_REDIS_URL")
	if url == "" {
		t.Skip("HABARCHY_TEST_REDIS_URL not set")
	}
	opts, err := goredis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	r := goredis.NewClient(opts)
	defer r.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pid, other := uuid.New(), uuid.New()
	ch := Subscribe(ctx, r, []uuid.UUID{pid})
	time.Sleep(100 * time.Millisecond) // let the subscription register
	pub := NewPublisher(r)
	pub.Publish(ctx, Event{Type: "message.sent", ProjectID: other, Status: "sent"}) // not subscribed
	pub.Publish(ctx, Event{Type: "message.sent", ProjectID: pid, Status: "sent", To: "+99365"})
	select {
	case ev := <-ch:
		if ev.ProjectID != pid || ev.Status != "sent" || ev.At.IsZero() {
			t.Fatalf("%+v", ev)
		}
	case <-ctx.Done():
		t.Fatal("no event received")
	}
	var nilPub *Publisher
	nilPub.Publish(ctx, Event{}) // must not panic
}

// Package events publishes message status changes over Redis pub/sub so
// the admin dashboard can stream them live (SSE) from any API replica.
package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Channel prefix; one Redis channel per project.
const prefix = "habarchy:events:"

// Event is what dashboards receive.
type Event struct {
	Type      string          `json:"type"` // message.queued|processing|sent|delivered|failed|cancelled, webhook.delivered|failed
	ProjectID uuid.UUID       `json:"project_id"`
	MessageID *uuid.UUID      `json:"message_id,omitempty"`
	Channel   string          `json:"channel,omitempty"`
	Status    string          `json:"status,omitempty"`
	To        string          `json:"to,omitempty"`
	Provider  string          `json:"provider,omitempty"`
	ErrorCode string          `json:"error_code,omitempty"`
	At        time.Time       `json:"at"`
	Extra     json.RawMessage `json:"extra,omitempty"`
}

// Publisher sends events.
type Publisher struct {
	r *goredis.Client
}

// NewPublisher wraps a Redis client. A nil client disables publishing.
func NewPublisher(r *goredis.Client) *Publisher { return &Publisher{r: r} }

// Publish is fire-and-forget; failures are logged, never returned, because
// live updates must not affect message processing.
func (p *Publisher) Publish(ctx context.Context, ev Event) {
	if p == nil || p.r == nil {
		return
	}
	if ev.At.IsZero() {
		ev.At = time.Now().UTC()
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	if err := p.r.Publish(ctx, prefix+ev.ProjectID.String(), b).Err(); err != nil {
		log.Ctx(ctx).Debug().Err(err).Msg("publish event")
	}
}

// Subscribe streams events for the given projects until ctx is done.
// The returned channel is closed on exit.
func Subscribe(ctx context.Context, r *goredis.Client, projectIDs []uuid.UUID) <-chan Event {
	out := make(chan Event, 64)
	if r == nil || len(projectIDs) == 0 {
		close(out)
		return out
	}
	channels := make([]string, 0, len(projectIDs))
	for _, id := range projectIDs {
		channels = append(channels, prefix+id.String())
	}
	sub := r.Subscribe(ctx, channels...)
	go func() {
		defer close(out)
		defer func() { _ = sub.Close() }()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-ch:
				if !ok {
					return
				}
				var ev Event
				if json.Unmarshal([]byte(m.Payload), &ev) != nil {
					continue
				}
				select {
				case out <- ev:
				case <-ctx.Done():
					return
				default: // slow consumer: drop rather than block the subscriber
				}
			}
		}
	}()
	return out
}

// ---- chat ----

const chatUserPrefix = "habarchy:chat:u:"
const chatBroadcast = "habarchy:chat:all"

// ChatEvent is pushed to the chat SSE stream.
type ChatEvent struct {
	Type      string          `json:"type"` // message | message.deleted | channel
	ChannelID uuid.UUID       `json:"channel_id"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	At        time.Time       `json:"at"`
}

// PublishChatToUsers sends an event to each listed user's personal stream
// (private groups and direct messages).
func (p *Publisher) PublishChatToUsers(ctx context.Context, userIDs []uuid.UUID, ev ChatEvent) {
	p.publishChat(ctx, ev, func(b []byte) {
		for _, u := range userIDs {
			if p.r != nil {
				_ = p.r.Publish(ctx, chatUserPrefix+u.String(), b).Err()
			}
		}
	})
}

// PublishChatBroadcast sends an event to every connected user (public channels).
func (p *Publisher) PublishChatBroadcast(ctx context.Context, ev ChatEvent) {
	p.publishChat(ctx, ev, func(b []byte) {
		if p.r != nil {
			_ = p.r.Publish(ctx, chatBroadcast, b).Err()
		}
	})
}

func (p *Publisher) publishChat(_ context.Context, ev ChatEvent, send func([]byte)) {
	if p == nil || p.r == nil {
		return
	}
	if ev.At.IsZero() {
		ev.At = time.Now().UTC()
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	send(b)
}

// SubscribeChat streams chat events for one user (their personal channel
// plus the public broadcast) until ctx is done.
func SubscribeChat(ctx context.Context, r *goredis.Client, userID uuid.UUID) <-chan ChatEvent {
	out := make(chan ChatEvent, 64)
	if r == nil {
		close(out)
		return out
	}
	sub := r.Subscribe(ctx, chatUserPrefix+userID.String(), chatBroadcast)
	go func() {
		defer close(out)
		defer func() { _ = sub.Close() }()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-ch:
				if !ok {
					return
				}
				var ev ChatEvent
				if json.Unmarshal([]byte(m.Payload), &ev) != nil {
					continue
				}
				select {
				case out <- ev:
				case <-ctx.Done():
					return
				default:
				}
			}
		}
	}()
	return out
}

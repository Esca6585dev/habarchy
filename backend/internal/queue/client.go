package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// SendPayload is the body of send:* tasks. Only the id travels; the worker
// reloads the message so the task never holds stale content.
type SendPayload struct {
	MessageID uuid.UUID `json:"message_id"`
}

// WebhookPayload is the body of webhook:deliver tasks.
type WebhookPayload struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
}

// Client enqueues tasks; it implements ports.Queue.
type Client struct {
	c         *asynq.Client
	inspector *asynq.Inspector
	maxRetry  int
	whRetry   int
}

// NewClient wraps an asynq client.
func NewClient(opt asynq.RedisConnOpt, sendMaxRetry, webhookMaxRetry int) *Client {
	return &Client{c: asynq.NewClient(opt), inspector: asynq.NewInspector(opt), maxRetry: sendMaxRetry, whRetry: webhookMaxRetry}
}

// Close releases connections.
func (c *Client) Close() error {
	_ = c.inspector.Close()
	return c.c.Close()
}

// EnqueueSend schedules delivery of a message. The task id equals the
// message id so a queued message can be cancelled by deleting its task.
// High-priority messages get two extra retries; low-priority one fewer.
func (c *Client) EnqueueSend(ctx context.Context, messageID uuid.UUID, channel domain.Channel, priority domain.Priority, at *time.Time) error {
	queueName, taskType, err := ForChannel(channel)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(SendPayload{MessageID: messageID})
	maxRetry := c.maxRetry
	switch priority {
	case domain.PriorityHigh:
		maxRetry += 2
	case domain.PriorityLow:
		if maxRetry > 1 {
			maxRetry--
		}
	}
	opts := []asynq.Option{
		asynq.Queue(queueName),
		asynq.TaskID(SendTaskID(messageID)),
		asynq.MaxRetry(maxRetry),
		asynq.Timeout(2 * time.Minute),
		asynq.Retention(24 * time.Hour),
	}
	if at != nil && at.After(time.Now()) {
		opts = append(opts, asynq.ProcessAt(*at))
	}
	task := asynq.NewTask(taskType, body)
	_, err = c.c.EnqueueContext(ctx, task, opts...)
	if errors.Is(err, asynq.ErrTaskIDConflict) {
		// A resend of a message whose previous task is still retained
		// (completed / archived): drop the old one and enqueue again.
		_ = c.inspector.DeleteTask(queueName, SendTaskID(messageID))
		_, err = c.c.EnqueueContext(ctx, task, opts...)
	}
	if err != nil {
		return fmt.Errorf("queue: enqueue %s: %w", taskType, err)
	}
	return nil
}

// EnqueueWebhook schedules a webhook delivery attempt.
func (c *Client) EnqueueWebhook(ctx context.Context, deliveryID uuid.UUID) error {
	body, _ := json.Marshal(WebhookPayload{DeliveryID: deliveryID})
	_, err := c.c.EnqueueContext(ctx, asynq.NewTask(TaskWebhook, body),
		asynq.Queue(QueueWebhooks),
		asynq.MaxRetry(c.whRetry-1), // attempts = 1 + retries
		asynq.Timeout(30*time.Second),
		asynq.Retention(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("queue: enqueue webhook: %w", err)
	}
	return nil
}

// CancelSend removes a not-yet-started send task. It returns false when
// the task is already running or gone.
func (c *Client) CancelSend(_ context.Context, messageID uuid.UUID, channel domain.Channel) (bool, error) {
	queueName, _, err := ForChannel(channel)
	if err != nil {
		return false, err
	}
	err = c.inspector.DeleteTask(queueName, SendTaskID(messageID))
	switch {
	case err == nil:
		return true, nil
	case isNotFound(err) || isActive(err):
		return false, nil
	}
	return false, err
}

// SendTaskID derives the asynq task id of a message.
func SendTaskID(messageID uuid.UUID) string { return "send:" + messageID.String() }

// QueueStats summarises one queue for the health endpoint.
type QueueStats struct {
	Queue     string `json:"queue"`
	Pending   int    `json:"pending"`
	Active    int    `json:"active"`
	Scheduled int    `json:"scheduled"`
	Retry     int    `json:"retry"`
	Archived  int    `json:"archived"`
	Processed int    `json:"processed_today"`
	Failed    int    `json:"failed_today"`
}

// Stats returns per-queue depths.
func (c *Client) Stats() ([]QueueStats, error) {
	names, err := c.inspector.Queues()
	if err != nil {
		return nil, err
	}
	out := make([]QueueStats, 0, len(names))
	for _, n := range names {
		info, err := c.inspector.GetQueueInfo(n)
		if err != nil {
			return nil, err
		}
		out = append(out, QueueStats{Queue: n, Pending: info.Pending, Active: info.Active, Scheduled: info.Scheduled,
			Retry: info.Retry, Archived: info.Archived, Processed: info.Processed, Failed: info.Failed})
	}
	return out, nil
}

// SendRetryDelay implements the 10 s / 60 s / 300 s backoff from the spec
// for send tasks; later retries (high priority allows 5) keep 300 s.
func SendRetryDelay(n int, _ error, _ *asynq.Task) time.Duration {
	switch n {
	case 1:
		return 10 * time.Second
	case 2:
		return 60 * time.Second
	}
	return 300 * time.Second
}

// WebhookRetryDelay is exponential: 10 s, 30 s, 1 m, 5 m, 15 m, 1 h, 3 h.
func WebhookRetryDelay(n int, _ error, _ *asynq.Task) time.Duration {
	delays := []time.Duration{10 * time.Second, 30 * time.Second, time.Minute, 5 * time.Minute, 15 * time.Minute, time.Hour, 3 * time.Hour}
	if n < 1 {
		n = 1
	}
	if n > len(delays) {
		n = len(delays)
	}
	return delays[n-1]
}

// RetryDelay dispatches on task type so one asynq server can host both.
func RetryDelay(n int, err error, t *asynq.Task) time.Duration {
	if t.Type() == TaskWebhook {
		return WebhookRetryDelay(n, err, t)
	}
	return SendRetryDelay(n, err, t)
}

func isNotFound(err error) bool {
	return err != nil && (err.Error() == asynq.ErrTaskNotFound.Error() || err.Error() == asynq.ErrQueueNotFound.Error() || containsAny(err.Error(), "not found"))
}

func isActive(err error) bool {
	return err != nil && containsAny(err.Error(), "active", "running")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) && indexOf(s, sub) >= 0 {
			return true
		}
	}
	return false
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

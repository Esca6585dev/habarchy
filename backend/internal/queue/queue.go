// Package queue holds the asynq queue names, task types and shared helpers
// used by the API (producer), worker (consumer) and scheduler.
package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Queue names. Each channel has its own queue so a slow SMPP link cannot
// starve email, and webhooks never compete with deliveries.
const (
	QueueSMS      = "sms"
	QueueEmail    = "email"
	QueuePush     = "push"
	QueueTelegram = "telegram"
	QueueWebhooks = "webhooks"
)

// Task type names.
const (
	TaskSendSMS      = "send:sms"
	TaskSendEmail    = "send:email"
	TaskSendPush     = "send:push"
	TaskSendTelegram = "send:telegram"
	TaskWebhook      = "webhook:deliver"
)

// Weights returns the asynq queue priority map. Higher weight = more
// workers; asynq processes queues proportionally to these weights.
func Weights() map[string]int {
	return map[string]int{
		QueuePush:     6,
		QueueSMS:      6,
		QueueTelegram: 4,
		QueueEmail:    3,
		QueueWebhooks: 2,
	}
}

// ForChannel maps a channel to its queue and task type.
func ForChannel(ch domain.Channel) (queueName, taskType string, err error) {
	switch ch {
	case domain.ChannelSMS:
		return QueueSMS, TaskSendSMS, nil
	case domain.ChannelEmail:
		return QueueEmail, TaskSendEmail, nil
	case domain.ChannelPush:
		return QueuePush, TaskSendPush, nil
	case domain.ChannelTelegram:
		return QueueTelegram, TaskSendTelegram, nil
	}
	return "", "", fmt.Errorf("queue: no queue for channel %q", ch)
}

// PriorityOption maps a domain priority to an asynq option. asynq has no
// per-task priority inside a queue, so high-priority tasks get a larger
// retention and jump via a dedicated unique-id prefix later; for now the
// mapping is expressed as task deadline/timeout differences.
func PriorityOption(p domain.Priority) asynq.Option {
	switch p {
	case domain.PriorityHigh:
		return asynq.MaxRetry(5)
	case domain.PriorityLow:
		return asynq.MaxRetry(2)
	}
	return asynq.MaxRetry(3)
}

// RedisOpt adapts parsed go-redis options to asynq's connection options so
// both libraries share one HABARCHY_REDIS_URL.
func RedisOpt(o *goredis.Options) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Network:   o.Network,
		Addr:      o.Addr,
		Username:  o.Username,
		Password:  o.Password,
		DB:        o.DB,
		TLSConfig: o.TLSConfig,
	}
}

// AsynqLogger adapts zerolog to asynq.Logger.
type AsynqLogger struct{ Logger zerolog.Logger }

func (l AsynqLogger) Debug(args ...any) { l.Logger.Debug().Msg(fmt.Sprint(args...)) }
func (l AsynqLogger) Info(args ...any)  { l.Logger.Info().Msg(fmt.Sprint(args...)) }
func (l AsynqLogger) Warn(args ...any)  { l.Logger.Warn().Msg(fmt.Sprint(args...)) }
func (l AsynqLogger) Error(args ...any) { l.Logger.Error().Msg(fmt.Sprint(args...)) }
func (l AsynqLogger) Fatal(args ...any) { l.Logger.Fatal().Msg(fmt.Sprint(args...)) }

// Package webhooks delivers signed event notifications to project
// webhook URLs with retries.
package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
)

// Headers sent with every webhook.
const (
	HeaderSignature = "X-Habarchy-Signature"
	HeaderTimestamp = "X-Habarchy-Timestamp"
	HeaderEvent     = "X-Habarchy-Event"
	HeaderDelivery  = "X-Habarchy-Delivery"
	UserAgent       = "Habarchy-Webhooks/1.0"
)

// Service holds webhook use cases.
type Service struct {
	db       *postgres.DB
	cipher   ports.Cipher
	queue    ports.Queue
	client   *http.Client
	maxTries int
	now      func() time.Time
	// AllowPrivate permits webhook URLs on loopback / private addresses
	// (development and tests). Off in production: SSRF guard.
	AllowPrivate bool
}

// New creates the service. client may be nil.
func New(db *postgres.DB, cipher ports.Cipher, queue ports.Queue, maxAttempts int, client *http.Client) *Service {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse // never follow redirects with signed bodies
		}}
	}
	if maxAttempts <= 0 {
		maxAttempts = 8
	}
	return &Service{db: db, cipher: cipher, queue: queue, client: client, maxTries: maxAttempts, now: time.Now}
}

// Event is the JSON body delivered to the webhook URL.
type Event struct {
	ID        uuid.UUID           `json:"id"`
	Type      domain.WebhookEvent `json:"type"`
	CreatedAt time.Time           `json:"created_at"`
	ProjectID uuid.UUID           `json:"project_id"`
	Data      any                 `json:"data"`
}

// Emit records an event for a project and schedules delivery. Projects
// without a webhook URL are skipped silently.
func (s *Service) Emit(ctx context.Context, projectID uuid.UUID, event domain.WebhookEvent, messageID, batchID *uuid.UUID, data any) error {
	proj, err := s.db.Queries.GetProjectWebhook(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	if proj.WebhookUrl == "" {
		return nil
	}
	id := uuid.New()
	body, err := json.Marshal(Event{ID: id, Type: event, CreatedAt: s.now().UTC(), ProjectID: projectID, Data: data})
	if err != nil {
		return err
	}
	secret, err := s.secret(proj.WebhookSecretEnc, projectID)
	if err != nil {
		return err
	}
	ts := strconv.FormatInt(s.now().Unix(), 10)
	sig := Sign(secret, ts, body)
	row, err := s.db.Queries.CreateWebhookDelivery(ctx, sqlcgen.CreateWebhookDeliveryParams{
		ProjectID: projectID, MessageID: messageID, BatchID: batchID, Event: string(event), Url: proj.WebhookUrl, Payload: body, Signature: sig,
	})
	if err != nil {
		return err
	}
	return s.queue.EnqueueWebhook(ctx, row.ID)
}

// Sign computes the signature header value: "t=<ts>,v1=<hex hmac>" over
// "<ts>.<body>". With no secret configured the body is still signed with
// an empty key so the format stays stable.
func Sign(secret []byte, ts string, body []byte) string {
	return "t=" + ts + ",v1=" + crypto.HMACSHA256Hex(secret, []byte(ts+"."+string(body)))
}

// Verify checks a signature header as a receiving application would. It
// is exported for the SDKs' test suites and the docs.
func Verify(secret []byte, header string, body []byte, tolerance time.Duration, now time.Time) bool {
	var ts, v1 string
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			ts = kv[1]
		case "v1":
			v1 = kv[1]
		}
	}
	if ts == "" || v1 == "" {
		return false
	}
	sec, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return false
	}
	if d := now.Sub(time.Unix(sec, 0)); tolerance > 0 && (d > tolerance || d < -tolerance) {
		return false
	}
	return crypto.VerifyHMACSHA256Hex(secret, []byte(ts+"."+string(body)), v1)
}

// ErrRetry signals asynq to retry; the attempt was recorded already.
var ErrRetry = errors.New("webhook: delivery failed, will retry")

// Deliver performs one attempt for a stored delivery. It returns nil when
// delivered or permanently given up, ErrRetry when another attempt should
// follow.
func (s *Service) Deliver(ctx context.Context, deliveryID uuid.UUID) error {
	d, err := s.db.Queries.GetWebhookDeliveryByID(ctx, deliveryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	if d.DeliveredAt != nil {
		return nil
	}
	if err := checkURL(d.Url, s.AllowPrivate); err != nil {
		_ = s.db.Queries.RecordWebhookAttempt(ctx, sqlcgen.RecordWebhookAttemptParams{ID: d.ID, ResponseBody: err.Error()})
		return nil //nolint:nilerr // unsafe URL is permanent; recorded above, never retried
	}
	// The payload column is jsonb, so Postgres may re-serialise it. Sign the
	// exact bytes being sent, with a fresh timestamp for every attempt.
	proj, err := s.db.Queries.GetProjectWebhook(ctx, d.ProjectID)
	if err != nil {
		return err
	}
	secret, err := s.secret(proj.WebhookSecretEnc, d.ProjectID)
	if err != nil {
		return err
	}
	ts := strconv.FormatInt(s.now().Unix(), 10)
	sig := Sign(secret, ts, d.Payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.Url, bytes.NewReader(d.Payload))
	if err != nil {
		// Malformed URL: permanent, record and stop.
		_ = s.db.Queries.RecordWebhookAttempt(ctx, sqlcgen.RecordWebhookAttemptParams{ID: d.ID, ResponseBody: err.Error()})
		return nil //nolint:nilerr // permanent failure recorded above; retrying cannot help
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set(HeaderSignature, sig)
	req.Header.Set(HeaderTimestamp, ts)
	req.Header.Set(HeaderEvent, d.Event)
	req.Header.Set(HeaderDelivery, d.ID.String())

	var code *int
	var respBody string
	res, err := s.client.Do(req)
	if err != nil {
		respBody = err.Error()
	} else {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		_ = res.Body.Close()
		c := res.StatusCode
		code, respBody = &c, string(b)
	}
	success := code != nil && *code >= 200 && *code < 300
	attempts := d.Attempts + 1
	var delivered, next *time.Time
	now := s.now()
	if success {
		delivered = &now
	} else if int(attempts) < s.maxTries {
		n := now.Add(backoff(int(attempts)))
		next = &n
	}
	if err := s.db.Queries.RecordWebhookAttempt(ctx, sqlcgen.RecordWebhookAttemptParams{
		ID: d.ID, ResponseCode: code32(code), ResponseBody: truncate(respBody, 2000), NextRetryAt: next, DeliveredAt: delivered,
	}); err != nil {
		return err
	}
	if success {
		return nil
	}
	if next == nil {
		log.Ctx(ctx).Warn().Str("delivery", d.ID.String()).Str("url", d.Url).Msg("webhook given up after max attempts")
		return nil
	}
	return ErrRetry
}

// Resend resets a delivery so it is attempted again (admin action).
func (s *Service) Resend(ctx context.Context, projectID, deliveryID uuid.UUID) (*sqlcgen.WebhookDelivery, error) {
	row, err := s.db.Queries.ResetWebhookDelivery(ctx, sqlcgen.ResetWebhookDeliveryParams{ID: deliveryID, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("webhook delivery not found")
		}
		return nil, err
	}
	return &row, s.queue.EnqueueWebhook(ctx, row.ID)
}

func (s *Service) secret(enc []byte, projectID uuid.UUID) ([]byte, error) {
	if len(enc) == 0 {
		return []byte{}, nil
	}
	return s.cipher.Decrypt(enc, []byte("webhook:"+projectID.String()))
}

// backoff mirrors queue.WebhookRetryDelay for the next_retry_at column.
func backoff(attempt int) time.Duration {
	delays := []time.Duration{10 * time.Second, 30 * time.Second, time.Minute, 5 * time.Minute, 15 * time.Minute, time.Hour, 3 * time.Hour}
	if attempt < 1 {
		attempt = 1
	}
	if attempt > len(delays) {
		attempt = len(delays)
	}
	return delays[attempt-1]
}

// checkURL blocks obviously unsafe targets (SSRF): only http(s), and no
// loopback / link-local / private hosts unless allowPrivate is set.
func checkURL(raw string, allowPrivate bool) error {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return fmt.Errorf("webhook url must be http(s)")
	}
	host := raw[strings.Index(raw, "://")+3:]
	if i := strings.IndexAny(host, "/?#"); i >= 0 {
		host = host[:i]
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if allowPrivate {
		return nil
	}
	if strings.EqualFold(host, "localhost") {
		return fmt.Errorf("webhook url points to a local address")
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() || ip.IsPrivate() {
			return fmt.Errorf("webhook url points to a local address")
		}
	}
	return nil
}

func code32(c *int) *int32 {
	if c == nil {
		return nil
	}
	v := int32(*c) //nolint:gosec // HTTP status
	return &v
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

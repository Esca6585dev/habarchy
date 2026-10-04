// Package otp issues and verifies one-time codes. Codes are stored only as
// salted hashes in Redis, rate limited per address and per client IP, and
// never logged.
package otp

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	rds "github.com/Esca6585dev/habarchy/backend/internal/adapters/redis"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/phone"
)

// Limits are the OTP tunables from config.
type Limits struct {
	Length         int
	TTL            time.Duration
	MaxAttempts    int
	PerAddressHour int
	PerIPHour      int
}

// Service holds the OTP use cases.
type Service struct {
	redis    *rds.Client
	messages *messages.Service
	limits   Limits
}

// New creates the service.
func New(redis *rds.Client, msgs *messages.Service, limits Limits) *Service {
	if limits.Length < 4 || limits.Length > 10 {
		limits.Length = 6
	}
	if limits.TTL <= 0 {
		limits.TTL = 5 * time.Minute
	}
	if limits.MaxAttempts <= 0 {
		limits.MaxAttempts = 5
	}
	if limits.PerAddressHour <= 0 {
		limits.PerAddressHour = 5
	}
	if limits.PerIPHour <= 0 {
		limits.PerIPHour = 30
	}
	return &Service{redis: redis, messages: msgs, limits: limits}
}

// DefaultTemplateKey is used when the request names none.
const DefaultTemplateKey = "otp"

// SendInput is the payload of POST /otp/send.
type SendInput struct {
	Channel  domain.Channel
	To       string
	Length   int
	TTL      time.Duration
	Template string
	Locale   domain.Locale
	Data     map[string]any
	ClientIP string
}

// SendResult is returned to the caller; the code is not.
type SendResult struct {
	MessageID uuid.UUID      `json:"message_id"`
	To        string         `json:"to"`
	Channel   domain.Channel `json:"channel"`
	ExpiresIn int            `json:"expires_in"`
	Length    int            `json:"length"`
}

// Send generates a code, stores its hash and dispatches it through the
// normal message pipeline.
func (s *Service) Send(ctx context.Context, caller messages.Caller, in SendInput) (*SendResult, error) {
	if in.Channel == "" {
		in.Channel = domain.ChannelSMS
	}
	if in.Channel != domain.ChannelSMS && in.Channel != domain.ChannelEmail && in.Channel != domain.ChannelTelegram && in.Channel != domain.ChannelWhatsApp {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"channel": "sms, email or telegram"})
	}
	to, err := normalize(in.Channel, in.To)
	if err != nil {
		return nil, err
	}
	length := in.Length
	if length < 4 || length > 10 {
		length = s.limits.Length
	}
	ttl := in.TTL
	if ttl <= 0 || ttl > time.Hour {
		ttl = s.limits.TTL
	}

	// Rate limits: per address and per IP, hourly windows.
	scope := caller.Project.ID.String()
	if n, _, err := s.redis.Incr(ctx, "otprl:addr:"+scope+":"+to, 1, time.Hour); err != nil {
		return nil, err
	} else if n > int64(s.limits.PerAddressHour) {
		return nil, domain.ErrRateLimited.WithMessage("too many codes requested for this address; try later")
	}
	if in.ClientIP != "" {
		if n, _, err := s.redis.Incr(ctx, "otprl:ip:"+scope+":"+in.ClientIP, 1, time.Hour); err != nil {
			return nil, err
		} else if n > int64(s.limits.PerIPHour) {
			return nil, domain.ErrRateLimited.WithMessage("too many codes requested from this address; try later")
		}
	}

	code, err := generate(length)
	if err != nil {
		return nil, err
	}
	if err := s.redis.PutOTP(ctx, scope, to, rds.HashOTP(to, code), ttl); err != nil {
		return nil, err
	}

	data := make(map[string]any, len(in.Data)+2)
	for k, v := range in.Data {
		data[k] = v
	}
	data["code"] = code
	data["minutes"] = int(ttl.Minutes())
	tpl := in.Template
	if tpl == "" {
		tpl = DefaultTemplateKey
	}
	res, err := s.messages.Send(ctx, caller, messages.SendInput{
		Channel: in.Channel, To: messages.Recipient{Address: to}, Template: tpl, Data: data, Locale: in.Locale,
		Priority: domain.PriorityHigh, Metadata: map[string]any{"otp": true},
	})
	if err != nil {
		// Fall back to a built-in text when the project has no otp template.
		var de *domain.Error
		if isTemplateMissing(err, &de) {
			res, err = s.messages.Send(ctx, caller, messages.SendInput{
				Channel: in.Channel, To: messages.Recipient{Address: to}, Body: defaultBody(in.Locale, caller.Project.DefaultLocale),
				Subject: defaultSubject(in.Channel), Data: data, Priority: domain.PriorityHigh, Metadata: map[string]any{"otp": true},
			})
		}
		if err != nil {
			_ = s.redis.DelOTP(ctx, scope, to)
			return nil, err
		}
	}
	return &SendResult{MessageID: res.Message.ID, To: to, Channel: in.Channel, ExpiresIn: int(ttl.Seconds()), Length: length}, nil
}

// VerifyResult is returned by Verify.
type VerifyResult struct {
	Verified          bool `json:"verified"`
	AttemptsRemaining int  `json:"attempts_remaining"`
}

// Verify checks a code. Wrong codes count against MaxAttempts; the code is
// deleted on success or lockout.
func (s *Service) Verify(ctx context.Context, projectID uuid.UUID, rawTo, code string) (*VerifyResult, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"code": "required"})
	}
	to := strings.TrimSpace(rawTo)
	if p, err := phone.Normalize(to, ""); err == nil {
		to = p
	} else {
		to = strings.ToLower(to)
	}
	scope := projectID.String()
	rec, ok, err := s.redis.GetOTP(ctx, scope, to)
	if err != nil {
		return nil, err
	}
	if !ok || time.Now().After(rec.ExpiresAt) {
		return &VerifyResult{Verified: false, AttemptsRemaining: 0}, domain.ErrUnauthorized.WithMessage("code expired or not requested")
	}
	if rec.Attempts >= int64(s.limits.MaxAttempts) {
		_ = s.redis.DelOTP(ctx, scope, to)
		return &VerifyResult{Verified: false}, domain.ErrRateLimited.WithMessage("too many attempts; request a new code")
	}
	if subtle.ConstantTimeCompare([]byte(rec.Hash), []byte(rds.HashOTP(to, code))) == 1 {
		_ = s.redis.DelOTP(ctx, scope, to)
		return &VerifyResult{Verified: true}, nil
	}
	n, _ := s.redis.IncrOTPAttempts(ctx, scope, to)
	remaining := s.limits.MaxAttempts - int(n)
	if remaining <= 0 {
		_ = s.redis.DelOTP(ctx, scope, to)
		remaining = 0
	}
	return &VerifyResult{Verified: false, AttemptsRemaining: remaining}, domain.ErrUnauthorized.WithMessage("invalid code").WithDetails(map[string]any{"attempts_remaining": remaining})
}

func normalize(ch domain.Channel, raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	switch ch {
	case domain.ChannelSMS:
		p, err := phone.Normalize(raw, "")
		if err != nil {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid phone number")
		}
		return p, nil
	case domain.ChannelEmail:
		if !strings.Contains(raw, "@") {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid email address")
		}
		return strings.ToLower(raw), nil
	}
	if raw == "" {
		return "", domain.ErrInvalidRecipient
	}
	return raw, nil
}

func generate(length int) (string, error) {
	var b strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		fmt.Fprint(&b, n.Int64())
	}
	return b.String(), nil
}

func isTemplateMissing(err error, de **domain.Error) bool {
	if !errorsAs(err, de) || (*de).Code != domain.CodeValidationFailed {
		return false
	}
	_, ok := (*de).Details["template"]
	return ok
}

func defaultBody(locale domain.Locale, projectDefault string) string {
	if locale == "" {
		locale = domain.Locale(projectDefault)
	}
	switch locale {
	case domain.LocaleRU:
		return "Ваш код подтверждения: {{.code}}. Действителен {{.minutes}} мин."
	case domain.LocaleEN:
		return "Your verification code is {{.code}}. Valid for {{.minutes}} minutes."
	}
	return "Siziň tassyklama koduňyz: {{.code}}. {{.minutes}} minut güýjünde."
}

func defaultSubject(ch domain.Channel) string {
	if ch == domain.ChannelEmail {
		return "Tassyklama kody / Verification code"
	}
	return ""
}

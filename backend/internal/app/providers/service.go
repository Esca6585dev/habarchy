// Package providers manages delivery provider configuration (encrypted
// credentials) and turns stored rows into live adapter instances.
package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/email/smtp"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/push/fcm"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sandbox"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/slack"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/androidgw"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/httpgeneric"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/telegram"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/whatsapp"
	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
)

// Service holds provider use cases.
type Service struct {
	db     *postgres.DB
	cipher ports.Cipher
	http   *http.Client

	mu    sync.Mutex
	cache map[uuid.UUID]cached // built adapters keyed by provider id
	// SMPPFactory builds SMPP sessions; set by the worker (keeps the heavy
	// dependency out of the API binary). nil => smpp providers fail to build.
	SMPPFactory func(ctx context.Context, id uuid.UUID, creds json.RawMessage) (ports.SMSProvider, error)
	// GatewayFactory builds the Android gateway provider (needs the outbox
	// store); set by the worker. nil => android_sms providers fail to build.
	GatewayFactory func(id uuid.UUID, creds json.RawMessage) (ports.SMSProvider, error)
}

type cached struct {
	updatedAt time.Time
	adapter   any
}

// New creates the service.
func New(db *postgres.DB, cipher ports.Cipher) *Service {
	return &Service{db: db, cipher: cipher, http: &http.Client{Timeout: 20 * time.Second}, cache: map[uuid.UUID]cached{}}
}

// Input is the payload for Create and Update. Credentials is the plaintext
// JSON for the provider type; nil on Update keeps the stored one.
type Input struct {
	Name            string
	Type            domain.ProviderType
	Priority        *int
	IsActive        *bool
	Credentials     json.RawMessage
	RateLimitPerSec *int
}

// Create validates the credentials for the type, encrypts and stores them.
func (s *Service) Create(ctx context.Context, projectID uuid.UUID, in Input) (*sqlcgen.Provider, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"name": "required"})
	}
	if !in.Type.Valid() {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"type": "http_sms, smpp, android_sms, smtp, fcm, telegram_bot, whatsapp_cloud or slack"})
	}
	if in.Type == domain.ProviderAndroidSMS {
		var err error
		if in.Credentials, err = withGatewayKey(in.Credentials); err != nil {
			return nil, err
		}
	}
	if err := ValidateCredentials(ctx, in.Type, in.Credentials); err != nil {
		return nil, err
	}
	id := uuid.New()
	enc, err := s.cipher.Encrypt(in.Credentials, aad(id))
	if err != nil {
		return nil, err
	}
	priority, active, rate := 100, true, 0
	if in.Priority != nil {
		priority = *in.Priority
	}
	if in.IsActive != nil {
		active = *in.IsActive
	}
	if in.RateLimitPerSec != nil {
		rate = *in.RateLimitPerSec
	}
	var p sqlcgen.Provider
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		p, err = q.CreateProvider(ctx, sqlcgen.CreateProviderParams{
			ID: id, ProjectID: projectID, Name: strings.TrimSpace(in.Name), Channel: sqlcgen.Channel(in.Type.Channel()),
			Type: sqlcgen.ProviderType(in.Type), Priority: int32(priority), IsActive: active, CredentialsEnc: enc, RateLimitPerSec: int32(rate), //nolint:gosec // small ints
		})
		if err != nil {
			return err
		}
		if err := registerGatewayKey(ctx, q, &p, in.Credentials); err != nil {
			return err
		}
		audit.Record(ctx, q, &projectID, "provider.create", "provider", p.ID.String(), map[string]any{"name": p.Name, "type": p.Type})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// withGatewayKey fills in a random gateway_key when the admin left it out.
func withGatewayKey(creds json.RawMessage) (json.RawMessage, error) {
	m := map[string]any{}
	if len(creds) > 0 {
		if err := json.Unmarshal(creds, &m); err != nil {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"credentials": "must be a JSON object"})
		}
	}
	if k, _ := m["gateway_key"].(string); strings.TrimSpace(k) == "" {
		tok, err := crypto.RandomToken(24)
		if err != nil {
			return nil, err
		}
		m["gateway_key"] = "gw_" + tok
	}
	return json.Marshal(m)
}

// registerGatewayKey stores the hash of an android_sms provider's pairing
// key so the phone can be authenticated without decrypting every row.
func registerGatewayKey(ctx context.Context, q *sqlcgen.Queries, p *sqlcgen.Provider, creds json.RawMessage) error {
	if domain.ProviderType(p.Type) != domain.ProviderAndroidSMS {
		return nil
	}
	cfg, err := androidgw.ParseConfig(creds)
	if err != nil {
		return domain.ErrValidation.WithDetails(map[string]any{"credentials": err.Error()})
	}
	if _, err := q.UpsertGatewayDevice(ctx, sqlcgen.UpsertGatewayDeviceParams{ProviderID: p.ID, ProjectID: p.ProjectID, KeyHash: crypto.SHA256(cfg.GatewayKey)}); err != nil {
		if postgres.IsUniqueViolation(err) {
			return domain.ErrConflict.WithMessage("this gateway_key is already used by another provider")
		}
		return err
	}
	return nil
}

// Get returns one provider (credentials stay encrypted).
func (s *Service) Get(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Provider, error) {
	p, err := s.db.Queries.GetProvider(ctx, sqlcgen.GetProviderParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("provider not found")
		}
		return nil, err
	}
	return &p, nil
}

// GetByID returns a provider regardless of project (callbacks).
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*sqlcgen.Provider, error) {
	p, err := s.db.Queries.GetProviderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("provider not found")
		}
		return nil, err
	}
	return &p, nil
}

// List returns a project's providers.
func (s *Service) List(ctx context.Context, projectID uuid.UUID) ([]sqlcgen.Provider, error) {
	return s.db.Queries.ListProviders(ctx, projectID)
}

// Update changes settings and optionally replaces credentials.
func (s *Service) Update(ctx context.Context, projectID, id uuid.UUID, in Input) (*sqlcgen.Provider, error) {
	cur, err := s.Get(ctx, projectID, id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		cur.Name = strings.TrimSpace(in.Name)
	}
	if in.Priority != nil {
		cur.Priority = int32(*in.Priority) //nolint:gosec // small int
	}
	if in.IsActive != nil {
		cur.IsActive = *in.IsActive
	}
	if in.RateLimitPerSec != nil {
		cur.RateLimitPerSec = int32(*in.RateLimitPerSec) //nolint:gosec // small int
	}
	changes := map[string]any{"priority": cur.Priority, "is_active": cur.IsActive, "rate_limit_per_sec": cur.RateLimitPerSec}
	if in.Credentials != nil {
		if domain.ProviderType(cur.Type) == domain.ProviderAndroidSMS {
			if in.Credentials, err = withGatewayKey(in.Credentials); err != nil {
				return nil, err
			}
		}
		if err := ValidateCredentials(ctx, domain.ProviderType(cur.Type), in.Credentials); err != nil {
			return nil, err
		}
		if cur.CredentialsEnc, err = s.cipher.Encrypt(in.Credentials, aad(id)); err != nil {
			return nil, err
		}
		changes["credentials"] = "rotated"
	}
	var p sqlcgen.Provider
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		p, err = q.UpdateProvider(ctx, sqlcgen.UpdateProviderParams{
			ID: id, ProjectID: projectID, Name: cur.Name, Priority: cur.Priority, IsActive: cur.IsActive,
			CredentialsEnc: cur.CredentialsEnc, RateLimitPerSec: cur.RateLimitPerSec,
		})
		if err != nil {
			return err
		}
		if in.Credentials != nil {
			if err := registerGatewayKey(ctx, q, &p, in.Credentials); err != nil {
				return err
			}
		}
		audit.Record(ctx, q, &projectID, "provider.update", "provider", id.String(), changes)
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.invalidate(id)
	return &p, nil
}

// Delete removes a provider.
func (s *Service) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	err := s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		n, err := q.DeleteProvider(ctx, sqlcgen.DeleteProviderParams{ID: id, ProjectID: projectID})
		if err != nil {
			return err
		}
		if n == 0 {
			return domain.ErrNotFound.WithMessage("provider not found")
		}
		audit.Record(ctx, q, &projectID, "provider.delete", "provider", id.String(), nil)
		return nil
	})
	if err == nil {
		s.invalidate(id)
	}
	return err
}

// Decrypt returns the plaintext credentials. Only workers and test-send
// call it; it is never exposed over HTTP.
func (s *Service) Decrypt(p *sqlcgen.Provider) (json.RawMessage, error) {
	raw, err := s.cipher.Decrypt(p.CredentialsEnc, aad(p.ID))
	if err != nil {
		return nil, fmt.Errorf("provider %s: decrypt credentials: %w", p.ID, err)
	}
	return raw, nil
}

// Redacted returns the credentials with secret fields masked for display.
func (s *Service) Redacted(p *sqlcgen.Provider) (map[string]any, error) {
	raw, err := s.Decrypt(p)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	for k, v := range m {
		lk := strings.ToLower(k)
		switch {
		case strings.Contains(lk, "password"), strings.Contains(lk, "token"), strings.Contains(lk, "secret"), lk == "service_account":
			m[k] = mask(fmt.Sprint(v))
		case lk == "headers":
			if hdr, ok := v.(map[string]any); ok {
				for hk, hv := range hdr {
					if strings.Contains(strings.ToLower(hk), "auth") || strings.Contains(strings.ToLower(hk), "key") {
						hdr[hk] = mask(fmt.Sprint(hv))
					}
				}
			}
		}
	}
	return m, nil
}

func mask(s string) string {
	if len(s) <= 6 {
		return "••••"
	}
	return s[:3] + "••••" + s[len(s)-3:]
}

// ActiveForChannel lists active providers of a project for a channel in
// fallback order.
func (s *Service) ActiveForChannel(ctx context.Context, projectID uuid.UUID, ch domain.Channel) ([]sqlcgen.Provider, error) {
	return s.db.Queries.ListActiveProvidersForChannel(ctx, sqlcgen.ListActiveProvidersForChannelParams{ProjectID: projectID, Channel: sqlcgen.Channel(ch)})
}

// Build returns the adapter for a provider row, cached until the row's
// updated_at changes. The result is one of ports.SMSProvider,
// ports.EmailProvider, ports.PushProvider or ports.TelegramProvider.
func (s *Service) Build(ctx context.Context, p *sqlcgen.Provider) (any, error) {
	s.mu.Lock()
	if c, ok := s.cache[p.ID]; ok && c.updatedAt.Equal(p.UpdatedAt) {
		s.mu.Unlock()
		return c.adapter, nil
	}
	s.mu.Unlock()

	creds, err := s.Decrypt(p)
	if err != nil {
		return nil, err
	}
	adapter, err := s.build(ctx, p, creds)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[p.ID] = cached{updatedAt: p.UpdatedAt, adapter: adapter}
	s.mu.Unlock()
	return adapter, nil
}

func (s *Service) build(ctx context.Context, p *sqlcgen.Provider, creds json.RawMessage) (any, error) {
	switch domain.ProviderType(p.Type) {
	case domain.ProviderHTTPSMS:
		var cfg httpgeneric.Config
		if err := json.Unmarshal(creds, &cfg); err != nil {
			return nil, err
		}
		return httpgeneric.New(cfg, nil)
	case domain.ProviderSMTP:
		var cfg smtp.Config
		if err := json.Unmarshal(creds, &cfg); err != nil {
			return nil, err
		}
		return smtp.New(cfg)
	case domain.ProviderFCM:
		var cfg fcm.Config
		if err := json.Unmarshal(creds, &cfg); err != nil {
			return nil, err
		}
		return fcm.New(ctx, cfg, nil)
	case domain.ProviderTelegramBot:
		var cfg telegram.Config
		if err := json.Unmarshal(creds, &cfg); err != nil {
			return nil, err
		}
		return telegram.New(cfg, nil)
	case domain.ProviderSMPP:
		if s.SMPPFactory == nil {
			return nil, errors.New("smpp providers are only available in the worker")
		}
		return s.SMPPFactory(ctx, p.ID, creds)
	case domain.ProviderWhatsAppCloud:
		var cfg whatsapp.Config
		if err := json.Unmarshal(creds, &cfg); err != nil {
			return nil, err
		}
		return whatsapp.New(cfg, nil)
	case domain.ProviderSlack:
		var cfg slack.Config
		if err := json.Unmarshal(creds, &cfg); err != nil {
			return nil, err
		}
		return slack.New(cfg, nil)
	case domain.ProviderAndroidSMS:
		if s.GatewayFactory == nil {
			return nil, errors.New("android_sms providers are only available in the worker")
		}
		return s.GatewayFactory(p.ID, creds)
	}
	return nil, fmt.Errorf("unknown provider type %q", p.Type)
}

func (s *Service) invalidate(id uuid.UUID) {
	s.mu.Lock()
	delete(s.cache, id)
	s.mu.Unlock()
}

// Sandbox returns the no-op adapter for a channel (test API keys).
func Sandbox(ch domain.Channel) any {
	switch ch {
	case domain.ChannelEmail:
		return sandbox.Email{}
	case domain.ChannelPush:
		return sandbox.Push{}
	case domain.ChannelTelegram:
		return sandbox.Telegram{}
	case domain.ChannelWhatsApp, domain.ChannelSlack:
		return sandbox.Chat{}
	}
	return sandbox.Provider{}
}

// ValidateCredentials checks that creds is a well-formed config for the
// provider type by building the adapter without using it.
func ValidateCredentials(ctx context.Context, t domain.ProviderType, creds json.RawMessage) error {
	if len(creds) == 0 || !json.Valid(creds) {
		return domain.ErrValidation.WithDetails(map[string]any{"credentials": "must be a JSON object"})
	}
	var err error
	switch t {
	case domain.ProviderHTTPSMS:
		var cfg httpgeneric.Config
		if err = json.Unmarshal(creds, &cfg); err == nil {
			_, err = httpgeneric.New(cfg, http.DefaultClient)
		}
	case domain.ProviderSMTP:
		var cfg smtp.Config
		if err = json.Unmarshal(creds, &cfg); err == nil {
			_, err = smtp.New(cfg)
		}
	case domain.ProviderFCM:
		var cfg fcm.Config
		if err = json.Unmarshal(creds, &cfg); err == nil {
			_, err = fcm.New(ctx, cfg, http.DefaultClient) // no OAuth at validation time
		}
	case domain.ProviderTelegramBot:
		var cfg telegram.Config
		if err = json.Unmarshal(creds, &cfg); err == nil {
			_, err = telegram.New(cfg, http.DefaultClient)
		}
	case domain.ProviderSMPP:
		err = ValidateSMPP(creds)
	case domain.ProviderAndroidSMS:
		_, err = androidgw.ParseConfig(creds)
	case domain.ProviderWhatsAppCloud:
		var cfg whatsapp.Config
		if err = json.Unmarshal(creds, &cfg); err == nil {
			_, err = whatsapp.New(cfg, http.DefaultClient)
		}
	case domain.ProviderSlack:
		var cfg slack.Config
		if err = json.Unmarshal(creds, &cfg); err == nil {
			_, err = slack.New(cfg, http.DefaultClient)
		}
	default:
		err = fmt.Errorf("unknown provider type %q", t)
	}
	if err != nil {
		return domain.ErrValidation.WithDetails(map[string]any{"credentials": err.Error()})
	}
	return nil
}

// SMPPConfig is the credentials JSON of an smpp provider. The adapter lives
// in internal/adapters/providers/sms/smpp; the schema is here so the API
// can validate without importing the SMPP client.
type SMPPConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	SystemID       string `json:"system_id"`
	Password       string `json:"password"`
	SystemType     string `json:"system_type"`
	SourceAddr     string `json:"source_addr"`
	SourceTON      int    `json:"source_ton"`
	SourceNPI      int    `json:"source_npi"`
	DestTON        int    `json:"dest_ton"`
	DestNPI        int    `json:"dest_npi"`
	EnquireLinkSec int    `json:"enquire_link_sec"`
	UseTLS         bool   `json:"use_tls"`
	// RequestDLR asks the SMSC for delivery receipts (registered_delivery=1).
	RequestDLR bool `json:"request_dlr"`
}

// ValidateSMPP checks the SMPP schema.
func ValidateSMPP(creds json.RawMessage) error {
	var cfg SMPPConfig
	if err := json.Unmarshal(creds, &cfg); err != nil {
		return err
	}
	if cfg.Host == "" || cfg.SystemID == "" {
		return errors.New("smpp: host and system_id are required")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("smpp: port must be 1-65535")
	}
	return nil
}

func aad(id uuid.UUID) []byte { return []byte("provider:" + id.String()) }

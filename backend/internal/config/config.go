// Package config loads process configuration from environment variables.
// Every variable is prefixed with HABARCHY_ (e.g. HABARCHY_DATABASE_URL)
// and documented in backend/.env.example.
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Prefix is the environment variable prefix.
const Prefix = "HABARCHY"

// Config is the full configuration shared by api, worker and scheduler.
// Each binary only uses the parts it needs, but loading one struct keeps
// the documentation in a single place.
type Config struct {
	Env       string `envconfig:"ENV" default:"development"` // development | production | test
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogPretty bool   `envconfig:"LOG_PRETTY" default:"false"`

	// Sections are embedded so every variable keeps the flat HABARCHY_
	// prefix (envconfig would otherwise nest HABARCHY_DATABASE_...).
	HTTP
	Database
	Redis
	Auth
	Security
	Queue
	Limits
	Telemetry
}

// HTTP configures the API server.
type HTTP struct {
	Addr            string        `envconfig:"HTTP_ADDR" default:":8080"`
	PublicURL       string        `envconfig:"PUBLIC_URL" default:"http://localhost:8080"`
	ReadTimeout     time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"15s"`
	WriteTimeout    time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"30s"`
	IdleTimeout     time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"120s"`
	BodyLimitBytes  int           `envconfig:"HTTP_BODY_LIMIT_BYTES" default:"4194304"` // 4 MiB
	ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"20s"`
	CORSOrigins     []string      `envconfig:"CORS_ORIGINS" default:"http://localhost:3000"`
	TrustedProxies  []string      `envconfig:"TRUSTED_PROXIES" default:""`
}

// Database configures PostgreSQL.
type Database struct {
	URL             string        `envconfig:"DATABASE_URL" required:"true"`
	MaxConns        int32         `envconfig:"DB_MAX_CONNS" default:"20"`
	MinConns        int32         `envconfig:"DB_MIN_CONNS" default:"2"`
	MaxConnLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"1h"`
	AutoMigrate     bool          `envconfig:"DB_AUTO_MIGRATE" default:"false"`
}

// Redis configures the Redis connection shared by queues, rate limits,
// idempotency and OTP storage.
type Redis struct {
	URL string `envconfig:"REDIS_URL" default:"redis://localhost:6379/0"`
}

// Auth configures admin JWTs.
type Auth struct {
	JWTSecret       string        `envconfig:"JWT_SECRET" required:"true"`
	JWTIssuer       string        `envconfig:"JWT_ISSUER" default:"habarchy"`
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"720h"` // 30 days
}

// Security configures encryption at rest and request signing.
type Security struct {
	// MasterKey is a 32-byte key (base64 or hex) for AES-256-GCM.
	MasterKey          string        `envconfig:"MASTER_KEY" required:"true"`
	SignatureTolerance time.Duration `envconfig:"SIGNATURE_TOLERANCE" default:"5m"`
	IdempotencyTTL     time.Duration `envconfig:"IDEMPOTENCY_TTL" default:"24h"`
	// WebhookAllowPrivate lets webhook URLs target private / loopback hosts.
	WebhookAllowPrivate bool `envconfig:"WEBHOOK_ALLOW_PRIVATE" default:"false"`
}

// Queue configures asynq workers.
type Queue struct {
	Concurrency  int `envconfig:"WORKER_CONCURRENCY" default:"20"`
	MaxRetry     int `envconfig:"SEND_MAX_RETRY" default:"3"`
	WebhookRetry int `envconfig:"WEBHOOK_MAX_RETRY" default:"8"`
}

// Limits configures default rate limits and OTP behaviour.
type Limits struct {
	OTPLength          int           `envconfig:"OTP_LENGTH" default:"6"`
	OTPTTL             time.Duration `envconfig:"OTP_TTL" default:"5m"`
	OTPMaxAttempts     int           `envconfig:"OTP_MAX_ATTEMPTS" default:"5"`
	OTPPerAddressRate  int           `envconfig:"OTP_PER_ADDRESS_PER_HOUR" default:"5"`
	OTPPerIPRate       int           `envconfig:"OTP_PER_IP_PER_HOUR" default:"30"`
	BatchMaxRecipients int           `envconfig:"BATCH_MAX_RECIPIENTS" default:"1000"`
	// APIRatePerSec limits public API requests per API key; 0 disables.
	APIRatePerSec float64 `envconfig:"API_RATE_PER_SEC" default:"50"`
}

// Telemetry configures metrics and tracing.
type Telemetry struct {
	MetricsEnabled bool   `envconfig:"METRICS_ENABLED" default:"true"`
	TracingEnabled bool   `envconfig:"TRACING_ENABLED" default:"false"`
	OTLPEndpoint   string `envconfig:"OTLP_ENDPOINT" default:"localhost:4317"`
}

// Load reads configuration from the environment and validates it.
func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process(Prefix, &c); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &c, nil
}

// Validate checks cross-field invariants that envconfig cannot express.
func (c *Config) Validate() error {
	var errs []error
	switch c.Env {
	case "development", "production", "test":
	default:
		errs = append(errs, fmt.Errorf("ENV must be development, production or test, got %q", c.Env))
	}
	// envconfig treats a variable that is set but empty as present, so
	// required-ness is re-checked here.
	if c.Database.URL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if c.Security.MasterKey == "" {
		errs = append(errs, errors.New("MASTER_KEY is required"))
	}
	if len(c.Auth.JWTSecret) < 32 {
		errs = append(errs, errors.New("JWT_SECRET must be at least 32 characters"))
	}
	if c.Limits.OTPLength < 4 || c.Limits.OTPLength > 10 {
		errs = append(errs, errors.New("OTP_LENGTH must be between 4 and 10"))
	}
	if c.Queue.Concurrency < 1 {
		errs = append(errs, errors.New("WORKER_CONCURRENCY must be >= 1"))
	}
	if c.Database.MinConns > c.Database.MaxConns {
		errs = append(errs, errors.New("DB_MIN_CONNS must not exceed DB_MAX_CONNS"))
	}
	return errors.Join(errs...)
}

// IsProduction reports whether the process runs with production settings.
func (c *Config) IsProduction() bool { return c.Env == "production" }

// Usage prints the documented variables; used by `-help` flags.
func Usage() error {
	return envconfig.Usage(Prefix, &Config{})
}

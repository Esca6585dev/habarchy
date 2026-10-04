// Package auth implements admin authentication: password login, JWT
// access tokens, rotated refresh tokens and optional TOTP 2FA.
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/password"
)

// Service holds the dependencies of the auth use cases.
type Service struct {
	db     *postgres.DB
	cfg    config.Auth
	cipher ports.Cipher
	now    func() time.Time
	hash   func(string) (string, error)
}

// Option customises a Service.
type Option func(*Service)

// WithClock overrides time.Now (tests).
func WithClock(now func() time.Time) Option { return func(s *Service) { s.now = now } }

// WithPasswordHasher overrides the password hash function (tests use
// cheaper argon2 parameters).
func WithPasswordHasher(h func(string) (string, error)) Option {
	return func(s *Service) { s.hash = h }
}

// New creates the auth service.
func New(db *postgres.DB, cfg config.Auth, cipher ports.Cipher, opts ...Option) *Service {
	s := &Service{db: db, cfg: cfg, cipher: cipher, now: time.Now, hash: password.Hash}
	for _, o := range opts {
		o(s)
	}
	return s
}

// ErrTOTPRequired is returned by Login when the account has 2FA enabled and
// no code was supplied. The HTTP layer maps it to 401 with
// details.totp_required = true so the client can show the second step.
var ErrTOTPRequired = domain.ErrUnauthorized.WithMessage("totp code required").WithDetails(map[string]any{"totp_required": true})

var errBadCredentials = domain.ErrUnauthorized.WithMessage("invalid email or password")

// Session describes the client making the request (for refresh tokens).
type Session struct {
	UserAgent string
	IP        string
}

// Login verifies credentials (and the TOTP code when enabled) and issues a
// token pair.
func (s *Service) Login(ctx context.Context, email, pass, totpCode string, sess Session) (*Tokens, *sqlcgen.User, error) {
	user, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Burn comparable time so user enumeration via latency is harder.
			_, _ = password.Verify(pass, dummyHash)
			return nil, nil, errBadCredentials
		}
		return nil, nil, err
	}
	ok, err := password.Verify(pass, user.PasswordHash)
	if err != nil || !ok {
		return nil, nil, errBadCredentials
	}
	if !user.IsActive {
		return nil, nil, domain.ErrUnauthorized.WithMessage("account disabled")
	}
	if user.TotpEnabled {
		if totpCode == "" {
			return nil, nil, ErrTOTPRequired
		}
		if !s.verifyTOTP(user, totpCode) {
			return nil, nil, domain.ErrUnauthorized.WithMessage("invalid totp code")
		}
	}
	tokens, err := s.issue(ctx, &user, sess)
	if err != nil {
		return nil, nil, err
	}
	_ = s.db.Queries.TouchUserLogin(ctx, user.ID)
	return tokens, &user, nil
}

// Refresh rotates a refresh token. Presenting an already-rotated token is
// treated as theft: every session of the user is revoked.
func (s *Service) Refresh(ctx context.Context, refreshToken string, sess Session) (*Tokens, error) {
	row, err := s.db.Queries.GetRefreshTokenByHash(ctx, crypto.SHA256(refreshToken))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUnauthorized.WithMessage("invalid refresh token")
		}
		return nil, err
	}
	now := s.now()
	if row.RevokedAt != nil {
		_ = s.db.Queries.RevokeAllUserRefreshTokens(ctx, row.UserID)
		audit.Record(ctx, s.db.Queries, nil, "auth.refresh_reuse", "user", row.UserID.String(), nil)
		return nil, domain.ErrUnauthorized.WithMessage("refresh token reused; all sessions revoked")
	}
	if now.After(row.ExpiresAt) {
		return nil, domain.ErrUnauthorized.WithMessage("refresh token expired")
	}
	user, err := s.db.Queries.GetUserByID(ctx, row.UserID)
	if err != nil || !user.IsActive {
		return nil, domain.ErrUnauthorized.WithMessage("account disabled")
	}
	var tokens *Tokens
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		tokens, err = s.issueWith(ctx, q, &user, sess)
		if err != nil {
			return err
		}
		newRow, err := q.GetRefreshTokenByHash(ctx, crypto.SHA256(tokens.RefreshToken))
		if err != nil {
			return err
		}
		return q.RevokeRefreshToken(ctx, sqlcgen.RevokeRefreshTokenParams{ID: row.ID, ReplacedBy: &newRow.ID})
	})
	return tokens, err
}

// Logout revokes one refresh token. Unknown tokens are ignored.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	row, err := s.db.Queries.GetRefreshTokenByHash(ctx, crypto.SHA256(refreshToken))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	return s.db.Queries.RevokeRefreshToken(ctx, sqlcgen.RevokeRefreshTokenParams{ID: row.ID})
}

// LogoutAll revokes every refresh token of a user.
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.db.Queries.RevokeAllUserRefreshTokens(ctx, userID)
}

// CreateUser registers an admin user.
func (s *Service) CreateUser(ctx context.Context, email, pass, fullName string) (*sqlcgen.User, error) {
	hash, err := s.hash(pass)
	if err != nil {
		return nil, domain.ErrValidation.WithMessage(err.Error())
	}
	user, err := s.db.Queries.CreateUser(ctx, sqlcgen.CreateUserParams{Email: email, PasswordHash: hash, FullName: fullName})
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil, domain.ErrConflict.WithMessage("email already registered")
		}
		return nil, err
	}
	audit.Record(ctx, s.db.Queries, nil, "user.create", "user", user.ID.String(), map[string]any{"email": user.Email})
	return &user, nil
}

// ChangePassword verifies the current password, stores the new one and
// revokes all refresh tokens so other sessions must log in again.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, current, next string) error {
	user, err := s.db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if ok, err := password.Verify(current, user.PasswordHash); err != nil || !ok {
		return domain.ErrUnauthorized.WithMessage("current password is incorrect")
	}
	hash, err := s.hash(next)
	if err != nil {
		return domain.ErrValidation.WithMessage(err.Error())
	}
	if err := s.db.Queries.UpdateUserPassword(ctx, sqlcgen.UpdateUserPasswordParams{ID: userID, PasswordHash: hash}); err != nil {
		return err
	}
	audit.Record(ctx, s.db.Queries, nil, "user.password_change", "user", userID.String(), nil)
	return s.db.Queries.RevokeAllUserRefreshTokens(ctx, userID)
}

// TOTPSetup is returned by SetupTOTP; the secret is shown once.
type TOTPSetup struct {
	Secret string `json:"secret"`
	URL    string `json:"otpauth_url"`
}

// SetupTOTP generates a new secret and stores it encrypted but disabled.
// ConfirmTOTP with a valid code turns 2FA on.
func (s *Service) SetupTOTP(ctx context.Context, userID uuid.UUID) (*TOTPSetup, error) {
	user, err := s.db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.TotpEnabled {
		return nil, domain.ErrConflict.WithMessage("2fa is already enabled; disable it first")
	}
	key, err := totp.Generate(totp.GenerateOpts{Issuer: s.cfg.JWTIssuer, AccountName: user.Email, Algorithm: otp.AlgorithmSHA1})
	if err != nil {
		return nil, err
	}
	enc, err := s.cipher.Encrypt([]byte(key.Secret()), totpAAD(userID))
	if err != nil {
		return nil, err
	}
	if err := s.db.Queries.SetUserTOTP(ctx, sqlcgen.SetUserTOTPParams{ID: userID, TotpEnabled: false, TotpSecretEnc: enc}); err != nil {
		return nil, err
	}
	return &TOTPSetup{Secret: key.Secret(), URL: key.URL()}, nil
}

// ConfirmTOTP enables 2FA after the user proves they can generate codes.
func (s *Service) ConfirmTOTP(ctx context.Context, userID uuid.UUID, code string) error {
	user, err := s.db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.TotpSecretEnc == nil {
		return domain.ErrValidation.WithMessage("run totp setup first")
	}
	if !s.verifyTOTP(user, code) {
		return domain.ErrUnauthorized.WithMessage("invalid totp code")
	}
	if err := s.db.Queries.SetUserTOTP(ctx, sqlcgen.SetUserTOTPParams{ID: userID, TotpEnabled: true, TotpSecretEnc: user.TotpSecretEnc}); err != nil {
		return err
	}
	audit.Record(ctx, s.db.Queries, nil, "user.totp_enable", "user", userID.String(), nil)
	return nil
}

// DisableTOTP turns 2FA off after re-verifying the password.
func (s *Service) DisableTOTP(ctx context.Context, userID uuid.UUID, pass string) error {
	user, err := s.db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if ok, err := password.Verify(pass, user.PasswordHash); err != nil || !ok {
		return domain.ErrUnauthorized.WithMessage("password is incorrect")
	}
	if err := s.db.Queries.SetUserTOTP(ctx, sqlcgen.SetUserTOTPParams{ID: userID, TotpEnabled: false, TotpSecretEnc: nil}); err != nil {
		return err
	}
	audit.Record(ctx, s.db.Queries, nil, "user.totp_disable", "user", userID.String(), nil)
	return nil
}

func (s *Service) verifyTOTP(user sqlcgen.User, code string) bool {
	secret, err := s.cipher.Decrypt(user.TotpSecretEnc, totpAAD(user.ID))
	if err != nil {
		return false
	}
	ok, err := totp.ValidateCustom(code, string(secret), s.now(), totp.ValidateOpts{
		Period: 30, Skew: 1, Digits: otp.DigitsSix, Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && ok
}

func totpAAD(userID uuid.UUID) []byte { return []byte("totp:" + userID.String()) }

func (s *Service) issue(ctx context.Context, user *sqlcgen.User, sess Session) (*Tokens, error) {
	return s.issueWith(ctx, s.db.Queries, user, sess)
}

func (s *Service) issueWith(ctx context.Context, q *sqlcgen.Queries, user *sqlcgen.User, sess Session) (*Tokens, error) {
	now := s.now()
	access, err := s.signAccess(user.ID, user.Email, now)
	if err != nil {
		return nil, err
	}
	refresh, err := crypto.RandomToken(32)
	if err != nil {
		return nil, err
	}
	if _, err := q.CreateRefreshToken(ctx, sqlcgen.CreateRefreshTokenParams{
		UserID: user.ID, TokenHash: crypto.SHA256(refresh), UserAgent: sess.UserAgent, Ip: sess.IP,
		ExpiresAt: now.Add(s.cfg.RefreshTokenTTL),
	}); err != nil {
		return nil, err
	}
	return &Tokens{
		AccessToken: access, RefreshToken: refresh, TokenType: "Bearer",
		ExpiresIn: int(s.cfg.AccessTokenTTL.Seconds()),
	}, nil
}

// dummyHash is verified against when the user does not exist, so both
// branches cost roughly the same time.
const dummyHash = "$argon2id$v=19$m=65536,t=3,p=2$AAAAAAAAAAAAAAAAAAAAAA$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

// User loads a user by id.
func (s *Service) User(ctx context.Context, id uuid.UUID) (*sqlcgen.User, error) {
	user, err := s.db.Queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("user not found")
		}
		return nil, err
	}
	return &user, nil
}

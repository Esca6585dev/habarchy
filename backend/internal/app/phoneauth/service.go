// Package phoneauth logs end-user clients in with their phone number and a
// one-time code: Start sends an OTP through the operator SMPP short code
// (0796) scoped to a system "auth" project, and Verify checks the code,
// creates the user and a personal tenant on first login, and issues the
// normal JWT access + refresh tokens. Each phone owns its own tenant, so one
// client's numbers never merge with another's.
package phoneauth

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/app/otp"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/phone"
)

// Service holds the phone-login use cases.
type Service struct {
	db       *postgres.DB
	auth     *auth.Service
	otp      *otp.Service
	projects *projects.Service
	// authProject owns the SMS provider that login codes are sent through and
	// scopes the OTP store; its message log shows every login code sent.
	authProject sqlcgen.Project
}

// New wires the service. authProject is the tenant whose SMS provider (e.g.
// the 0796 SMPP) delivers login codes.
func New(db *postgres.DB, authSvc *auth.Service, otpSvc *otp.Service, projectSvc *projects.Service, authProject sqlcgen.Project) *Service {
	return &Service{db: db, auth: authSvc, otp: otpSvc, projects: projectSvc, authProject: authProject}
}

// StartResult is returned by Start.
type StartResult struct {
	Phone     string `json:"phone"`
	ExpiresIn int    `json:"expires_in"`
	Length    int    `json:"length"`
}

// Start normalizes the phone and sends a login OTP through the auth project's
// SMS provider. Rate limits (per phone and per IP) come from the OTP service.
func (s *Service) Start(ctx context.Context, rawPhone, clientIP string) (*StartResult, error) {
	p, err := phone.Normalize(rawPhone, "")
	if err != nil {
		return nil, domain.ErrInvalidRecipient.WithMessage("invalid phone number")
	}
	res, err := s.otp.Send(ctx, messages.Caller{Project: s.authProject}, otp.SendInput{
		Channel: domain.ChannelSMS, To: p, ClientIP: clientIP,
	})
	if err != nil {
		return nil, err
	}
	return &StartResult{Phone: p, ExpiresIn: res.ExpiresIn, Length: res.Length}, nil
}

// VerifyResult is returned by Verify.
type VerifyResult struct {
	Tokens  *auth.Tokens
	User    *sqlcgen.User
	Created bool
}

// Verify checks the code, finds or creates the user (with a personal tenant
// on first login) and issues tokens.
func (s *Service) Verify(ctx context.Context, rawPhone, code string, sess auth.Session) (*VerifyResult, error) {
	p, err := phone.Normalize(rawPhone, "")
	if err != nil {
		return nil, domain.ErrInvalidRecipient.WithMessage("invalid phone number")
	}
	vr, err := s.otp.Verify(ctx, s.authProject.ID, p, code)
	if err != nil {
		return nil, err
	}
	if !vr.Verified {
		return nil, domain.ErrUnauthorized.WithMessage("invalid code")
	}

	user, created, err := s.findOrCreate(ctx, p)
	if err != nil {
		return nil, err
	}
	tokens, u, err := s.auth.IssueTokens(ctx, user.ID, sess)
	if err != nil {
		return nil, err
	}
	return &VerifyResult{Tokens: tokens, User: u, Created: created}, nil
}

// findOrCreate returns the user for the phone, creating the user and their
// personal tenant on first login.
func (s *Service) findOrCreate(ctx context.Context, e164 string) (*sqlcgen.User, bool, error) {
	user, err := s.db.Queries.GetUserByPhone(ctx, e164)
	switch {
	case err == nil:
		return &user, false, nil
	case errors.Is(err, pgx.ErrNoRows):
		// fall through to creation
	default:
		return nil, false, err
	}
	created, err := s.db.Queries.CreatePhoneUser(ctx, sqlcgen.CreatePhoneUserParams{FullName: "", Phone: e164})
	if err != nil {
		if postgres.IsUniqueViolation(err) { // raced with another login
			if u, gerr := s.db.Queries.GetUserByPhone(ctx, e164); gerr == nil {
				return &u, false, nil
			}
		}
		return nil, false, err
	}
	// Personal tenant, named after the phone; ignore a slug conflict so a
	// retried first login does not fail.
	slug := "u-" + digitsOnly(e164)
	if _, perr := s.projects.Create(ctx, created.ID, projects.CreateInput{Name: e164, Slug: slug}); perr != nil {
		var de *domain.Error
		if !errors.As(perr, &de) || de.Code != domain.CodeConflict {
			return nil, false, perr
		}
	}
	return &created, true, nil
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

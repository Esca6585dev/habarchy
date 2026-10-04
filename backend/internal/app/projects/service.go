// Package projects implements tenant management: projects, memberships
// and API keys.
package projects

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
)

// Service holds the dependencies of the project use cases.
type Service struct {
	db     *postgres.DB
	cipher ports.Cipher
}

// New creates the service.
func New(db *postgres.DB, cipher ports.Cipher) *Service { return &Service{db: db, cipher: cipher} }

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// CreateInput is the payload for Create.
type CreateInput struct {
	Name          string
	Slug          string
	DailyQuota    int64
	MonthlyQuota  int64
	DefaultLocale domain.Locale
}

// Create makes a project and adds the creator as owner.
func (s *Service) Create(ctx context.Context, ownerID uuid.UUID, in CreateInput) (*sqlcgen.Project, error) {
	in.Slug = strings.ToLower(strings.TrimSpace(in.Slug))
	if in.Slug == "" {
		in.Slug = Slugify(in.Name)
	}
	if !slugRe.MatchString(in.Slug) {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"slug": "lowercase letters, digits and dashes only"})
	}
	if in.DefaultLocale == "" {
		in.DefaultLocale = domain.DefaultLocale
	}
	if !in.DefaultLocale.Valid() {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"default_locale": "must be tk, ru or en"})
	}
	var project sqlcgen.Project
	err := s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		project, err = q.CreateProject(ctx, sqlcgen.CreateProjectParams{
			Name: strings.TrimSpace(in.Name), Slug: in.Slug, DailyQuota: in.DailyQuota,
			MonthlyQuota: in.MonthlyQuota, DefaultLocale: string(in.DefaultLocale),
		})
		if err != nil {
			if postgres.IsUniqueViolation(err) {
				return domain.ErrConflict.WithMessage("slug already in use")
			}
			return err
		}
		if _, err := q.UpsertProjectMember(ctx, sqlcgen.UpsertProjectMemberParams{
			ProjectID: project.ID, UserID: ownerID, Role: sqlcgen.MemberRoleOwner,
		}); err != nil {
			return err
		}
		audit.Record(ctx, q, &project.ID, "project.create", "project", project.ID.String(), map[string]any{"name": project.Name, "slug": project.Slug})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// Membership is a project together with the caller's role in it.
type Membership struct {
	Project sqlcgen.Project
	Role    domain.Role
}

// ListForUser returns every project the user belongs to.
func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID) ([]Membership, error) {
	rows, err := s.db.Queries.ListProjectsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]Membership, 0, len(rows))
	for _, r := range rows {
		out = append(out, Membership{Project: r.Project, Role: domain.Role(r.Role)})
	}
	return out, nil
}

// Authorize loads the project and checks the user holds at least minRole.
// It is the single entry point for every project-scoped admin operation.
func (s *Service) Authorize(ctx context.Context, userID, projectID uuid.UUID, minRole domain.Role) (*Membership, error) {
	m, err := s.db.Queries.GetProjectMember(ctx, sqlcgen.GetProjectMemberParams{ProjectID: projectID, UserID: userID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Hide existence of projects the user is not part of.
			return nil, domain.ErrNotFound.WithMessage("project not found")
		}
		return nil, err
	}
	role := domain.Role(m.Role)
	if !role.AtLeast(minRole) {
		return nil, domain.ErrForbiddenScope.WithMessage(fmt.Sprintf("requires %s role", minRole))
	}
	p, err := s.db.Queries.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Membership{Project: p, Role: role}, nil
}

// UpdateInput is the payload for Update. Nil pointers leave a field as is.
type UpdateInput struct {
	Name             *string
	Status           *domain.ProjectStatus
	DailyQuota       *int64
	MonthlyQuota     *int64
	WebhookURL       *string
	WebhookSecret    *string // plaintext; stored encrypted
	DefaultLocale    *domain.Locale
	AllowedIPs       *[]string
	AutoChannelOrder *[]domain.Channel
}

// Update changes project settings (admin role).
func (s *Service) Update(ctx context.Context, userID, projectID uuid.UUID, in UpdateInput) (*sqlcgen.Project, error) {
	m, err := s.Authorize(ctx, userID, projectID, domain.RoleAdmin)
	if err != nil {
		return nil, err
	}
	p := m.Project
	changes := map[string]any{}
	if in.Name != nil {
		p.Name, changes["name"] = strings.TrimSpace(*in.Name), *in.Name
	}
	if in.Status != nil {
		if !in.Status.Valid() {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"status": "active, suspended or archived"})
		}
		p.Status, changes["status"] = sqlcgen.ProjectStatus(*in.Status), *in.Status
	}
	if in.DailyQuota != nil {
		p.DailyQuota, changes["daily_quota"] = *in.DailyQuota, *in.DailyQuota
	}
	if in.MonthlyQuota != nil {
		p.MonthlyQuota, changes["monthly_quota"] = *in.MonthlyQuota, *in.MonthlyQuota
	}
	if in.WebhookURL != nil {
		u := strings.TrimSpace(*in.WebhookURL)
		if u != "" && !strings.HasPrefix(u, "https://") && !strings.HasPrefix(u, "http://") {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"webhook_url": "must start with http:// or https://"})
		}
		p.WebhookUrl, changes["webhook_url"] = u, u
	}
	if in.DefaultLocale != nil {
		if !in.DefaultLocale.Valid() {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"default_locale": "must be tk, ru or en"})
		}
		p.DefaultLocale, changes["default_locale"] = string(*in.DefaultLocale), *in.DefaultLocale
	}
	if in.AllowedIPs != nil {
		if err := ValidateCIDRs(*in.AllowedIPs); err != nil {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"allowed_ips": err.Error()})
		}
		p.AllowedIps, changes["allowed_ips"] = *in.AllowedIPs, *in.AllowedIPs
	}
	if in.AutoChannelOrder != nil {
		order := make([]sqlcgen.Channel, 0, len(*in.AutoChannelOrder))
		for _, c := range *in.AutoChannelOrder {
			if !c.Valid() {
				return nil, domain.ErrValidation.WithDetails(map[string]any{"auto_channel_order": "unknown channel " + string(c)})
			}
			order = append(order, sqlcgen.Channel(c))
		}
		p.AutoChannelOrder, changes["auto_channel_order"] = order, *in.AutoChannelOrder
	}

	var out sqlcgen.Project
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		out, err = q.UpdateProject(ctx, sqlcgen.UpdateProjectParams{
			ID: p.ID, Name: p.Name, Status: p.Status, DailyQuota: p.DailyQuota, MonthlyQuota: p.MonthlyQuota,
			WebhookUrl: p.WebhookUrl, DefaultLocale: p.DefaultLocale, AllowedIps: p.AllowedIps, AutoChannelOrder: p.AutoChannelOrder,
		})
		if err != nil {
			return err
		}
		if in.WebhookSecret != nil {
			var enc []byte
			if *in.WebhookSecret != "" {
				if enc, err = s.cipher.Encrypt([]byte(*in.WebhookSecret), []byte("webhook:"+p.ID.String())); err != nil {
					return err
				}
			}
			if err := q.SetProjectWebhookSecret(ctx, sqlcgen.SetProjectWebhookSecretParams{ID: p.ID, WebhookSecretEnc: enc}); err != nil {
				return err
			}
			out.WebhookSecretEnc = enc
			changes["webhook_secret"] = "rotated"
		}
		audit.Record(ctx, q, &p.ID, "project.update", "project", p.ID.String(), changes)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a project and everything in it (owner role).
func (s *Service) Delete(ctx context.Context, userID, projectID uuid.UUID) error {
	if _, err := s.Authorize(ctx, userID, projectID, domain.RoleOwner); err != nil {
		return err
	}
	return s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		audit.Record(ctx, q, nil, "project.delete", "project", projectID.String(), nil)
		return q.DeleteProject(ctx, projectID)
	})
}

// Member is a project member with user details.
type Member struct {
	UserID   uuid.UUID   `json:"user_id"`
	Email    string      `json:"email"`
	FullName string      `json:"full_name"`
	Role     domain.Role `json:"role"`
	JoinedAt time.Time   `json:"joined_at"`
}

// ListMembers returns members (viewer role).
func (s *Service) ListMembers(ctx context.Context, userID, projectID uuid.UUID) ([]Member, error) {
	if _, err := s.Authorize(ctx, userID, projectID, domain.RoleViewer); err != nil {
		return nil, err
	}
	rows, err := s.db.Queries.ListProjectMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]Member, 0, len(rows))
	for _, r := range rows {
		out = append(out, Member{UserID: r.ProjectMember.UserID, Email: r.Email, FullName: r.FullName, Role: domain.Role(r.ProjectMember.Role), JoinedAt: r.ProjectMember.CreatedAt})
	}
	return out, nil
}

// SetMember adds a user by email or changes their role (admin role). Only
// owners may grant the owner role.
func (s *Service) SetMember(ctx context.Context, userID, projectID uuid.UUID, email string, role domain.Role) (*Member, error) {
	m, err := s.Authorize(ctx, userID, projectID, domain.RoleAdmin)
	if err != nil {
		return nil, err
	}
	if !role.Valid() {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"role": "owner, admin, developer or viewer"})
	}
	if role == domain.RoleOwner && m.Role != domain.RoleOwner {
		return nil, domain.ErrForbiddenScope.WithMessage("only owners can grant the owner role")
	}
	target, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("no user with that email")
		}
		return nil, err
	}
	if target.ID == userID && role != domain.RoleOwner {
		if n, _ := s.db.Queries.CountProjectOwners(ctx, projectID); n <= 1 && m.Role == domain.RoleOwner {
			return nil, domain.ErrConflict.WithMessage("a project needs at least one owner")
		}
	}
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		if _, err := q.UpsertProjectMember(ctx, sqlcgen.UpsertProjectMemberParams{ProjectID: projectID, UserID: target.ID, Role: sqlcgen.MemberRole(role)}); err != nil {
			return err
		}
		audit.Record(ctx, q, &projectID, "member.set", "user", target.ID.String(), map[string]any{"email": target.Email, "role": role})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Member{UserID: target.ID, Email: target.Email, FullName: target.FullName, Role: role}, nil
}

// RemoveMember removes a member (admin role). The last owner cannot leave.
func (s *Service) RemoveMember(ctx context.Context, userID, projectID, targetID uuid.UUID) error {
	if _, err := s.Authorize(ctx, userID, projectID, domain.RoleAdmin); err != nil {
		return err
	}
	target, err := s.db.Queries.GetProjectMember(ctx, sqlcgen.GetProjectMemberParams{ProjectID: projectID, UserID: targetID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound.WithMessage("member not found")
		}
		return err
	}
	if target.Role == sqlcgen.MemberRoleOwner {
		n, err := s.db.Queries.CountProjectOwners(ctx, projectID)
		if err != nil {
			return err
		}
		if n <= 1 {
			return domain.ErrConflict.WithMessage("a project needs at least one owner")
		}
	}
	return s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		if err := q.RemoveProjectMember(ctx, sqlcgen.RemoveProjectMemberParams{ProjectID: projectID, UserID: targetID}); err != nil {
			return err
		}
		audit.Record(ctx, q, &projectID, "member.remove", "user", targetID.String(), nil)
		return nil
	})
}

// Slugify derives a URL-safe slug from a name.
func Slugify(name string) string {
	var b strings.Builder
	lastDash := true
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == 'ä':
			b.WriteString("a")
			lastDash = false
		case r == 'ö':
			b.WriteString("o")
			lastDash = false
		case r == 'ü':
			b.WriteString("u")
			lastDash = false
		case r == 'ň':
			b.WriteString("n")
			lastDash = false
		case r == 'ş':
			b.WriteString("s")
			lastDash = false
		case r == 'ž':
			b.WriteString("z")
			lastDash = false
		case r == 'ç':
			b.WriteString("c")
			lastDash = false
		case r == 'ý':
			b.WriteString("y")
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// ValidateCIDRs checks a list of IPs or CIDR ranges.
func ValidateCIDRs(list []string) error {
	for _, s := range list {
		if net.ParseIP(s) != nil {
			continue
		}
		if _, _, err := net.ParseCIDR(s); err != nil {
			return fmt.Errorf("invalid ip or cidr %q", s)
		}
	}
	return nil
}

// IPAllowed reports whether ip matches any entry (plain IP or CIDR). An
// empty list allows everything.
func IPAllowed(list []string, ip string) bool {
	if len(list) == 0 {
		return true
	}
	addr := net.ParseIP(ip)
	if addr == nil {
		return false
	}
	for _, s := range list {
		if a := net.ParseIP(s); a != nil {
			if a.Equal(addr) {
				return true
			}
			continue
		}
		if _, n, err := net.ParseCIDR(s); err == nil && n.Contains(addr) {
			return true
		}
	}
	return false
}

// ---- API keys ----

// APIKeyInput is the payload for CreateAPIKey.
type APIKeyInput struct {
	Name             string
	Live             bool
	Scopes           []domain.APIKeyScope
	IPAllowlist      []string
	ExpiresAt        *time.Time
	RequireSignature bool
}

// CreatedAPIKey carries the plaintext key, which is returned exactly once.
type CreatedAPIKey struct {
	Key       sqlcgen.ApiKey
	Plaintext string
}

// CreateAPIKey generates a key (admin role). The plaintext is
// prefix + 43 url-safe base64 chars and is never stored.
func (s *Service) CreateAPIKey(ctx context.Context, userID, projectID uuid.UUID, in APIKeyInput) (*CreatedAPIKey, error) {
	if _, err := s.Authorize(ctx, userID, projectID, domain.RoleAdmin); err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Name) == "" {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"name": "required"})
	}
	if len(in.Scopes) == 0 {
		in.Scopes = domain.AllScopes
	}
	scopes := make([]string, 0, len(in.Scopes))
	for _, sc := range in.Scopes {
		if !sc.Valid() {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"scopes": "unknown scope " + string(sc)})
		}
		scopes = append(scopes, string(sc))
	}
	if err := ValidateCIDRs(in.IPAllowlist); err != nil {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"ip_allowlist": err.Error()})
	}
	if in.IPAllowlist == nil {
		in.IPAllowlist = []string{}
	}
	prefix := domain.APIKeyPrefixTest
	if in.Live {
		prefix = domain.APIKeyPrefixLive
	}
	secret, err := crypto.RandomToken(32)
	if err != nil {
		return nil, err
	}
	plaintext := prefix + secret
	var key sqlcgen.ApiKey
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		key, err = q.CreateAPIKey(ctx, sqlcgen.CreateAPIKeyParams{
			ProjectID: projectID, Name: strings.TrimSpace(in.Name), Prefix: prefix, Hint: plaintext[len(plaintext)-4:],
			KeyHash: crypto.SHA256(plaintext), Scopes: scopes, IpAllowlist: in.IPAllowlist,
			ExpiresAt: in.ExpiresAt, RequireSignature: in.RequireSignature,
		})
		if err != nil {
			return err
		}
		audit.Record(ctx, q, &projectID, "api_key.create", "api_key", key.ID.String(), map[string]any{"name": key.Name, "prefix": prefix, "scopes": scopes})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &CreatedAPIKey{Key: key, Plaintext: plaintext}, nil
}

// ListAPIKeys returns keys without hashes (developer role).
func (s *Service) ListAPIKeys(ctx context.Context, userID, projectID uuid.UUID) ([]sqlcgen.ApiKey, error) {
	if _, err := s.Authorize(ctx, userID, projectID, domain.RoleDeveloper); err != nil {
		return nil, err
	}
	return s.db.Queries.ListAPIKeys(ctx, projectID)
}

// RevokeAPIKey disables a key permanently (admin role).
func (s *Service) RevokeAPIKey(ctx context.Context, userID, projectID, keyID uuid.UUID) error {
	if _, err := s.Authorize(ctx, userID, projectID, domain.RoleAdmin); err != nil {
		return err
	}
	return s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		n, err := q.RevokeAPIKey(ctx, sqlcgen.RevokeAPIKeyParams{ID: keyID, ProjectID: projectID})
		if err != nil {
			return err
		}
		if n == 0 {
			return domain.ErrNotFound.WithMessage("api key not found or already revoked")
		}
		audit.Record(ctx, q, &projectID, "api_key.revoke", "api_key", keyID.String(), nil)
		return nil
	})
}

// Caller is the authenticated client of a public API request.
type Caller struct {
	Project sqlcgen.Project
	Key     sqlcgen.ApiKey
}

// IsTest reports whether the request used a test key.
func (c Caller) IsTest() bool { return c.Key.Prefix == domain.APIKeyPrefixTest }

// HasScope reports whether the key grants s.
func (c Caller) HasScope(s domain.APIKeyScope) bool {
	for _, have := range c.Key.Scopes {
		if have == string(s) {
			return true
		}
	}
	return false
}

// SignatureInput carries the headers and request parts needed to verify an
// HMAC request signature.
type SignatureInput struct {
	Signature string // X-Signature (hex)
	Timestamp string // X-Timestamp (unix seconds)
	Method    string
	Path      string
	Body      []byte
	Tolerance time.Duration
	Now       time.Time
}

// AuthenticateAPIKey resolves a plaintext key to its project, checking
// revocation, expiry, IP allowlists (key and project) and, when present or
// required, the HMAC signature. It also records last_used_at.
func (s *Service) AuthenticateAPIKey(ctx context.Context, plaintext, ip string, sig SignatureInput) (*Caller, error) {
	plaintext = strings.TrimSpace(plaintext)
	if plaintext == "" || (!strings.HasPrefix(plaintext, domain.APIKeyPrefixLive) && !strings.HasPrefix(plaintext, domain.APIKeyPrefixTest)) {
		return nil, domain.ErrUnauthorized.WithMessage("missing or malformed api key")
	}
	key, err := s.db.Queries.GetAPIKeyByHash(ctx, crypto.SHA256(plaintext))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUnauthorized.WithMessage("invalid api key")
		}
		return nil, err
	}
	now := sig.Now
	if now.IsZero() {
		now = time.Now()
	}
	if key.RevokedAt != nil || (key.ExpiresAt != nil && now.After(*key.ExpiresAt)) {
		return nil, domain.ErrUnauthorized.WithMessage("api key revoked or expired")
	}
	project, err := s.db.Queries.GetProject(ctx, key.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.Status != sqlcgen.ProjectStatusActive {
		return nil, domain.ErrUnauthorized.WithMessage("project is " + string(project.Status))
	}
	if !IPAllowed(key.IpAllowlist, ip) || !IPAllowed(project.AllowedIps, ip) {
		return nil, domain.ErrUnauthorized.WithMessage("ip address not allowed")
	}
	if key.RequireSignature || sig.Signature != "" || sig.Timestamp != "" {
		if err := VerifySignature(plaintext, sig, now); err != nil {
			return nil, err
		}
	}
	_ = s.db.Queries.TouchAPIKey(ctx, key.ID)
	return &Caller{Project: project, Key: key}, nil
}

// SignatureMessage builds the canonical string that is signed:
//
//	timestamp "\n" METHOD "\n" path "\n" hex(sha256(body))
func SignatureMessage(timestamp, method, path string, body []byte) []byte {
	bodyHash := fmt.Sprintf("%x", crypto.SHA256(string(body)))
	return []byte(timestamp + "\n" + strings.ToUpper(method) + "\n" + path + "\n" + bodyHash)
}

// Sign computes the X-Signature value for a request; SDKs mirror this.
func Sign(apiKey, timestamp, method, path string, body []byte) string {
	return crypto.HMACSHA256Hex([]byte(apiKey), SignatureMessage(timestamp, method, path, body))
}

// VerifySignature checks X-Timestamp freshness and the HMAC.
func VerifySignature(apiKey string, sig SignatureInput, now time.Time) error {
	if sig.Signature == "" || sig.Timestamp == "" {
		return domain.ErrUnauthorized.WithMessage("request signature required (X-Signature, X-Timestamp)")
	}
	var ts int64
	if _, err := fmt.Sscanf(sig.Timestamp, "%d", &ts); err != nil {
		return domain.ErrUnauthorized.WithMessage("X-Timestamp must be unix seconds")
	}
	tol := sig.Tolerance
	if tol <= 0 {
		tol = 5 * time.Minute
	}
	if d := now.Sub(time.Unix(ts, 0)); d > tol || d < -tol {
		return domain.ErrUnauthorized.WithMessage("X-Timestamp outside the allowed window")
	}
	if !crypto.VerifyHMACSHA256Hex([]byte(apiKey), SignatureMessage(sig.Timestamp, sig.Method, sig.Path, sig.Body), sig.Signature) {
		return domain.ErrUnauthorized.WithMessage("invalid request signature")
	}
	return nil
}

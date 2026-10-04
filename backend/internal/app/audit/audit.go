// Package audit records admin actions into audit_logs.
package audit

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
)

// Actor identifies who performs an action; carried in the request context.
type Actor struct {
	UserID    uuid.UUID
	IP        string
	UserAgent string
}

type ctxKey struct{}

// WithActor stores the actor in ctx.
func WithActor(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, ctxKey{}, a)
}

// ActorFrom returns the actor stored in ctx, if any.
func ActorFrom(ctx context.Context) (Actor, bool) {
	a, ok := ctx.Value(ctxKey{}).(Actor)
	return a, ok
}

// Record writes one audit row. Failures are logged, never returned: an
// audit write must not fail the business operation it describes.
func Record(ctx context.Context, q *sqlcgen.Queries, projectID *uuid.UUID, action, entityType, entityID string, changes map[string]any) {
	var userID *uuid.UUID
	ip, ua := "", ""
	if a, ok := ActorFrom(ctx); ok {
		id := a.UserID
		userID, ip, ua = &id, a.IP, a.UserAgent
	}
	payload, err := json.Marshal(changes)
	if err != nil || changes == nil {
		payload = []byte("{}")
	}
	if _, err := q.CreateAuditLog(ctx, sqlcgen.CreateAuditLogParams{
		ProjectID: projectID, UserID: userID, Action: action, EntityType: entityType,
		EntityID: entityID, Changes: payload, Ip: ip, UserAgent: ua,
	}); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("action", action).Msg("audit log write failed")
	}
}

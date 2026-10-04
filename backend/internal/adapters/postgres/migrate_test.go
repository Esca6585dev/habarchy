package postgres_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/pgtest"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

func TestMigrationsAndCoreQueries(t *testing.T) {
	db := pgtest.Open(t)
	ctx := context.Background()
	q := db.Queries

	v, err := db.MigrationVersion(ctx)
	if err != nil || v < 7 {
		t.Fatalf("migration version %d, err %v", v, err)
	}

	// Partitions for this month and next must exist after migrating.
	var parts int
	err = db.Pool.QueryRow(ctx, `SELECT count(*) FROM pg_inherits WHERE inhparent = 'messages'::regclass`).Scan(&parts)
	if err != nil || parts < 3 { // current, next, default
		t.Fatalf("expected >= 3 message partitions, got %d (%v)", parts, err)
	}
	name, err := q.EnsureMessagesPartition(ctx, time.Now().AddDate(0, 5, 0))
	if err != nil || name == "" {
		t.Fatalf("EnsureMessagesPartition: %q %v", name, err)
	}
	// Idempotent.
	if again, err := q.EnsureMessagesPartition(ctx, time.Now().AddDate(0, 5, 0)); err != nil || again != name {
		t.Fatalf("second EnsureMessagesPartition: %q %v", again, err)
	}

	user, err := q.CreateUser(ctx, sqlcgen.CreateUserParams{Email: "Admin@Example.com", PasswordHash: "x", FullName: "Admin"})
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "admin@example.com" {
		t.Fatalf("email not lower-cased: %s", user.Email)
	}
	if _, err := q.CreateUser(ctx, sqlcgen.CreateUserParams{Email: "ADMIN@example.com", PasswordHash: "x"}); !postgres.IsUniqueViolation(err) {
		t.Fatalf("duplicate email should violate unique index, got %v", err)
	}

	project, err := q.CreateProject(ctx, sqlcgen.CreateProjectParams{Name: "Demo", Slug: "demo", DefaultLocale: "tk"})
	if err != nil {
		t.Fatal(err)
	}
	if len(project.AutoChannelOrder) != 3 || project.AutoChannelOrder[0] != sqlcgen.ChannelPush {
		t.Fatalf("auto channel default: %v", project.AutoChannelOrder)
	}
	if _, err := q.UpsertProjectMember(ctx, sqlcgen.UpsertProjectMemberParams{ProjectID: project.ID, UserID: user.ID, Role: sqlcgen.MemberRoleOwner}); err != nil {
		t.Fatal(err)
	}
	rows, err := q.ListProjectsForUser(ctx, user.ID)
	if err != nil || len(rows) != 1 || rows[0].Role != sqlcgen.MemberRoleOwner || rows[0].Project.Slug != "demo" {
		t.Fatalf("ListProjectsForUser: %+v %v", rows, err)
	}

	key, err := q.CreateAPIKey(ctx, sqlcgen.CreateAPIKeyParams{
		ProjectID: project.ID, Name: "default", Prefix: "hb_test_", Hint: "abcd",
		KeyHash: []byte("hash"), Scopes: []string{"messages:send"}, IpAllowlist: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := q.GetAPIKeyByHash(ctx, []byte("hash")); err != nil || got.ID != key.ID {
		t.Fatalf("GetAPIKeyByHash: %v", err)
	}
	if n, err := q.RevokeAPIKey(ctx, sqlcgen.RevokeAPIKeyParams{ID: key.ID, ProjectID: project.ID}); err != nil || n != 1 {
		t.Fatalf("RevokeAPIKey: %d %v", n, err)
	}
	if n, _ := q.RevokeAPIKey(ctx, sqlcgen.RevokeAPIKeyParams{ID: key.ID, ProjectID: project.ID}); n != 0 {
		t.Fatal("revoking twice must affect 0 rows")
	}

	tpl, err := q.CreateTemplate(ctx, sqlcgen.CreateTemplateParams{
		ProjectID: project.ID, Key: "otp", Channel: sqlcgen.ChannelSms, Locale: "tk",
		Body: "Kod: {{.code}}", RequiredVars: []string{"code"},
	})
	if err != nil {
		t.Fatal(err)
	}
	found, err := q.FindTemplate(ctx, sqlcgen.FindTemplateParams{ProjectID: project.ID, Key: "otp", Channel: sqlcgen.ChannelSms, Locale: "tk"})
	if err != nil || found.ID != tpl.ID || found.Version != 1 {
		t.Fatalf("FindTemplate: %+v %v", found, err)
	}
	updated, err := q.UpdateTemplate(ctx, sqlcgen.UpdateTemplateParams{ID: tpl.ID, ProjectID: project.ID, Body: "Siziň koduňyz: {{.code}}", RequiredVars: []string{"code"}, IsActive: true})
	if err != nil || updated.Version != 2 {
		t.Fatalf("UpdateTemplate version: %d %v", updated.Version, err)
	}

	contact, err := q.CreateContact(ctx, sqlcgen.CreateContactParams{
		ProjectID: project.ID, ExternalID: "u-1", Phone: "+99365123456", Email: "a@b.tm", Locale: "tk",
		Tags: []string{"vip"}, Attributes: json.RawMessage(`{"city":"Aşgabat"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	tag := "vip"
	list, err := q.ListContacts(ctx, sqlcgen.ListContactsParams{ProjectID: project.ID, Tag: &tag, RowLimit: 10})
	if err != nil || len(list) != 1 {
		t.Fatalf("ListContacts by tag: %d %v", len(list), err)
	}
	dev, err := q.UpsertDevice(ctx, sqlcgen.UpsertDeviceParams{ProjectID: project.ID, ContactID: &contact.ID, Platform: sqlcgen.DevicePlatformAndroid, FcmToken: "tok-1", AppVersion: "1.0"})
	if err != nil {
		t.Fatal(err)
	}
	dev2, err := q.UpsertDevice(ctx, sqlcgen.UpsertDeviceParams{ProjectID: project.ID, Platform: sqlcgen.DevicePlatformAndroid, FcmToken: "tok-1", AppVersion: "1.1"})
	if err != nil || dev2.ID != dev.ID || dev2.ContactID == nil || *dev2.ContactID != contact.ID || dev2.AppVersion != "1.1" {
		t.Fatalf("UpsertDevice should update in place and keep contact: %+v %v", dev2, err)
	}

	// Messages: insert, idempotency, list with cursor, state transitions.
	idem := "req-1"
	msgID := ids.New()
	msg, err := q.CreateMessage(ctx, sqlcgen.CreateMessageParams{
		ID: msgID, ProjectID: project.ID, Channel: sqlcgen.ChannelSms, ToAddress: "+99365123456",
		ContactID: &contact.ID, TemplateKey: "otp", RenderedBody: "Kod: 1234",
		Status: sqlcgen.MessageStatusQueued, Priority: sqlcgen.MessagePriorityNormal,
		Metadata: json.RawMessage(`{}`), IdempotencyKey: &idem,
	})
	if err != nil {
		t.Fatal(err)
	}
	reserved, err := q.ReserveIdempotencyKey(ctx, sqlcgen.ReserveIdempotencyKeyParams{ProjectID: project.ID, IdempotencyKey: idem, MessageID: msgID})
	if err != nil || reserved != msgID {
		t.Fatalf("ReserveIdempotencyKey: %v %v", reserved, err)
	}
	if _, err := q.ReserveIdempotencyKey(ctx, sqlcgen.ReserveIdempotencyKeyParams{ProjectID: project.ID, IdempotencyKey: idem, MessageID: ids.New()}); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("second reservation must return no rows, got %v", err)
	}
	if existing, err := q.GetIdempotencyKey(ctx, sqlcgen.GetIdempotencyKeyParams{ProjectID: project.ID, IdempotencyKey: idem}); err != nil || existing != msgID {
		t.Fatalf("GetIdempotencyKey: %v %v", existing, err)
	}

	second, err := q.CreateMessage(ctx, sqlcgen.CreateMessageParams{
		ID: ids.New(), ProjectID: project.ID, Channel: sqlcgen.ChannelEmail, ToAddress: "a@b.tm",
		RenderedBody: "hi", Status: sqlcgen.MessageStatusQueued, Priority: sqlcgen.MessagePriorityHigh, Metadata: json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	page, err := q.ListMessages(ctx, sqlcgen.ListMessagesParams{ProjectID: project.ID, RowLimit: 1})
	if err != nil || len(page) != 1 || page[0].ID != second.ID {
		t.Fatalf("ListMessages first page: %v %v", page, err)
	}
	page, err = q.ListMessages(ctx, sqlcgen.ListMessagesParams{ProjectID: project.ID, RowLimit: 10, CursorID: &second.ID})
	if err != nil || len(page) != 1 || page[0].ID != msg.ID {
		t.Fatalf("ListMessages cursor page: %v %v", page, err)
	}
	ch := sqlcgen.NullChannel{Channel: sqlcgen.ChannelEmail, Valid: true}
	page, err = q.ListMessages(ctx, sqlcgen.ListMessagesParams{ProjectID: project.ID, RowLimit: 10, Channel: ch})
	if err != nil || len(page) != 1 || page[0].Channel != sqlcgen.ChannelEmail {
		t.Fatalf("ListMessages by channel: %v %v", page, err)
	}

	proc, err := q.MarkMessageProcessing(ctx, msg.ID)
	if err != nil || proc.Status != sqlcgen.MessageStatusProcessing || proc.Attempts != 1 {
		t.Fatalf("MarkMessageProcessing: %+v %v", proc, err)
	}
	if err := q.MarkMessageSent(ctx, sqlcgen.MarkMessageSentParams{ID: msg.ID, ProviderMessageID: "ext-1", Currency: "TMT"}); err != nil {
		t.Fatal(err)
	}
	if n, err := q.MarkMessageDelivered(ctx, msg.ID); err != nil || n != 1 {
		t.Fatalf("MarkMessageDelivered: %d %v", n, err)
	}
	if n, _ := q.CancelMessage(ctx, sqlcgen.CancelMessageParams{ID: msg.ID, ProjectID: project.ID}); n != 0 {
		t.Fatal("delivered message must not be cancellable")
	}
	if n, err := q.CancelMessage(ctx, sqlcgen.CancelMessageParams{ID: second.ID, ProjectID: project.ID}); err != nil || n != 1 {
		t.Fatalf("CancelMessage queued: %d %v", n, err)
	}
	if _, err := q.CreateMessageEvent(ctx, sqlcgen.CreateMessageEventParams{MessageID: msg.ID, Type: sqlcgen.EventTypeDelivered, Payload: json.RawMessage(`{"dlr":"DELIVRD"}`)}); err != nil {
		t.Fatal(err)
	}
	events, err := q.ListMessageEvents(ctx, msg.ID)
	if err != nil || len(events) != 1 {
		t.Fatalf("ListMessageEvents: %d %v", len(events), err)
	}
	count, err := q.CountMessagesSince(ctx, sqlcgen.CountMessagesSinceParams{ProjectID: project.ID, Since: time.Now().Add(-time.Hour)})
	if err != nil || count != 1 { // cancelled one excluded
		t.Fatalf("CountMessagesSince: %d %v", count, err)
	}

	// Usage aggregation reads straight from messages.
	day := time.Now().UTC().Truncate(24 * time.Hour)
	agg, err := q.AggregateUsageForDay(ctx, sqlcgen.AggregateUsageForDayParams{DayStart: day, DayEnd: day.Add(24 * time.Hour)})
	if err != nil || len(agg) != 2 {
		t.Fatalf("AggregateUsageForDay: %+v %v", agg, err)
	}
	for _, a := range agg {
		if err := q.UpsertUsageDaily(ctx, sqlcgen.UpsertUsageDailyParams{
			ProjectID: a.ProjectID, Day: day, Channel: a.Channel, Queued: a.Queued, Sent: a.Sent,
			Delivered: a.Delivered, Failed: a.Failed, CostMicros: a.CostMicros, Currency: a.Currency,
		}); err != nil {
			t.Fatal(err)
		}
	}
	usage, err := q.ListUsageDaily(ctx, sqlcgen.ListUsageDailyParams{ProjectID: project.ID, FromDay: day, ToDay: day})
	if err != nil || len(usage) != 2 {
		t.Fatalf("ListUsageDaily: %d %v", len(usage), err)
	}

	// Cascade: deleting the project removes everything that belongs to it.
	if err := q.DeleteProject(ctx, project.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := q.GetMessage(ctx, sqlcgen.GetMessageParams{ID: msg.ID, ProjectID: project.ID}); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("message should be gone after project delete, got %v", err)
	}
}

func TestWithTxRollsBack(t *testing.T) {
	db := pgtest.Open(t)
	ctx := context.Background()
	err := db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		if _, err := q.CreateProject(ctx, sqlcgen.CreateProjectParams{Name: "tx", Slug: "tx", DefaultLocale: "tk"}); err != nil {
			return err
		}
		return errors.New("abort")
	})
	if err == nil || err.Error() != "abort" {
		t.Fatalf("expected abort error, got %v", err)
	}
	if _, err := db.Queries.GetProjectBySlug(ctx, "tx"); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("project should have been rolled back, got %v", err)
	}
}

// Package chat is the internal team chat: channels (public, private
// groups, direct messages), membership, messages with image attachments
// and live delivery over Redis pub/sub.
package chat

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// Service holds the chat use cases.
type Service struct {
	db     *postgres.DB
	events *events.Publisher
}

// New creates the service. A nil publisher disables live delivery.
func New(db *postgres.DB, pub *events.Publisher) *Service { return &Service{db: db, events: pub} }

// Channel is the API shape of a chat channel for one viewer.
type Channel struct {
	ID          uuid.UUID  `json:"id"`
	Kind        string     `json:"kind"` // public | private | direct
	Name        string     `json:"name"`
	Topic       string     `json:"topic"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Unread      int64      `json:"unread"`
	LastBody    string     `json:"last_body,omitempty"`
	LastAt      *time.Time `json:"last_at,omitempty"`
	LastHasFile bool       `json:"last_has_file,omitempty"`
	// For direct channels: the other person.
	PeerID     *uuid.UUID `json:"peer_id,omitempty"`
	PeerName   string     `json:"peer_name,omitempty"`
	PeerAvatar *uuid.UUID `json:"peer_avatar,omitempty"`
}

// Member is a channel participant.
type Member struct {
	ID       uuid.UUID  `json:"id"`
	FullName string     `json:"full_name"`
	Email    string     `json:"email"`
	AvatarID *uuid.UUID `json:"avatar_id,omitempty"`
	Bio      string     `json:"bio,omitempty"`
	Role     string     `json:"role"`
}

// Message is a chat message.
type Message struct {
	ID             uuid.UUID  `json:"id"`
	ChannelID      uuid.UUID  `json:"channel_id"`
	UserID         *uuid.UUID `json:"user_id,omitempty"`
	AuthorName     string     `json:"author_name"`
	AuthorAvatar   *uuid.UUID `json:"author_avatar,omitempty"`
	Body           string     `json:"body"`
	AttachmentID   *uuid.UUID `json:"attachment_id,omitempty"`
	AttachmentType string     `json:"attachment_type,omitempty"`
	AttachmentName string     `json:"attachment_name,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	EditedAt       *time.Time `json:"edited_at,omitempty"`
}

// dmKey returns the stable sorted key for a direct channel between two users.
func dmKey(a, b uuid.UUID) string {
	s := []string{a.String(), b.String()}
	sort.Strings(s)
	return s[0] + ":" + s[1]
}

// ListChannels returns the viewer's channels (auto-joining public ones) with
// unread counts, last message and the peer for direct channels.
func (s *Service) ListChannels(ctx context.Context, me uuid.UUID) ([]Channel, error) {
	if err := s.db.Queries.AutoJoinPublicChannels(ctx, me); err != nil {
		return nil, err
	}
	rows, err := s.db.Queries.ListMyChannelRows(ctx, me)
	if err != nil {
		return nil, err
	}
	unread, err := s.db.Queries.ChannelUnreadCounts(ctx, me)
	if err != nil {
		return nil, err
	}
	last, err := s.db.Queries.ChannelLastMessages(ctx, me)
	if err != nil {
		return nil, err
	}
	peers, err := s.db.Queries.DirectPeers(ctx, me)
	if err != nil {
		return nil, err
	}
	unreadBy := map[uuid.UUID]int64{}
	for _, u := range unread {
		unreadBy[u.ChannelID] = u.Unread
	}
	lastBy := map[uuid.UUID]sqlcgen.ChannelLastMessagesRow{}
	for _, l := range last {
		lastBy[l.ChannelID] = l
	}
	peerBy := map[uuid.UUID]sqlcgen.DirectPeersRow{}
	for _, p := range peers {
		peerBy[p.ChannelID] = p
	}
	out := make([]Channel, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		ch := Channel{ID: r.ID, Kind: string(r.Kind), Name: r.Name, Topic: r.Topic, CreatedBy: r.CreatedBy, CreatedAt: r.CreatedAt, Unread: unreadBy[r.ID]}
		if l, ok := lastBy[r.ID]; ok {
			ch.LastBody, ch.LastHasFile = l.Body, l.HasAttachment
			at := l.CreatedAt
			ch.LastAt = &at
		}
		if p, ok := peerBy[r.ID]; ok {
			pid := p.UserID
			ch.PeerID, ch.PeerName, ch.PeerAvatar = &pid, p.FullName, p.AvatarID
		}
		out = append(out, ch)
	}
	sort.SliceStable(out, func(i, j int) bool {
		ti, tj := out[i].CreatedAt, out[j].CreatedAt
		if out[i].LastAt != nil {
			ti = *out[i].LastAt
		}
		if out[j].LastAt != nil {
			tj = *out[j].LastAt
		}
		return ti.After(tj)
	})
	return out, nil
}

// CreateChannel makes a public or private (group) channel. The creator is
// the owner; given members are added. Public channels are visible to all.
func (s *Service) CreateChannel(ctx context.Context, me uuid.UUID, kind, name, topic string, memberIDs []uuid.UUID) (*sqlcgen.ChatChannel, error) {
	kind = strings.TrimSpace(kind)
	if kind != "public" && kind != "private" {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"kind": "public or private"})
	}
	if name = strings.TrimSpace(name); name == "" || len(name) > 100 {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"name": "required, max 100"})
	}
	var ch sqlcgen.ChatChannel
	err := s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		ch, err = q.CreateChannel(ctx, sqlcgen.CreateChannelParams{ID: ids.New(), Kind: sqlcgen.ChatChannelKind(kind), Name: name, Topic: strings.TrimSpace(topic), CreatedBy: &me})
		if err != nil {
			return err
		}
		if err := q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{ChannelID: ch.ID, UserID: me, Role: "owner"}); err != nil {
			return err
		}
		for _, uid := range memberIDs {
			if uid == me {
				continue
			}
			if err := q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{ChannelID: ch.ID, UserID: uid, Role: "member"}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.notifyChannel(ctx, &ch, memberIDs)
	return &ch, nil
}

// OpenDirect returns (creating if needed) the direct channel between me and
// other.
func (s *Service) OpenDirect(ctx context.Context, me, other uuid.UUID) (*sqlcgen.ChatChannel, error) {
	if me == other {
		return nil, domain.ErrValidation.WithMessage("cannot open a direct chat with yourself")
	}
	key := dmKey(me, other)
	if ch, err := s.db.Queries.GetDirectChannel(ctx, &key); err == nil {
		return &ch, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	var ch sqlcgen.ChatChannel
	err := s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		ch, err = q.CreateChannel(ctx, sqlcgen.CreateChannelParams{ID: ids.New(), Kind: sqlcgen.ChatChannelKindDirect, CreatedBy: &me, DmKey: &key})
		if err != nil {
			return err
		}
		if err := q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{ChannelID: ch.ID, UserID: me, Role: "member"}); err != nil {
			return err
		}
		return q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{ChannelID: ch.ID, UserID: other, Role: "member"})
	})
	if err != nil {
		if postgres.IsUniqueViolation(err) { // raced with another open
			ch, err = s.db.Queries.GetDirectChannel(ctx, &key)
			return &ch, err
		}
		return nil, err
	}
	return &ch, nil
}

// access checks that the viewer may see a channel: a member, or any user for
// a public channel. Returns the channel.
func (s *Service) access(ctx context.Context, me, channelID uuid.UUID) (*sqlcgen.ChatChannel, error) {
	ch, err := s.db.Queries.GetChannel(ctx, channelID)
	if err != nil {
		return nil, domain.ErrNotFound.WithMessage("channel not found")
	}
	if ch.Kind == sqlcgen.ChatChannelKindPublic {
		return &ch, nil
	}
	if _, err := s.db.Queries.GetChatMember(ctx, sqlcgen.GetChatMemberParams{ChannelID: channelID, UserID: me}); err != nil {
		return nil, domain.ErrForbiddenScope.WithMessage("not a member of this channel")
	}
	return &ch, nil
}

// Members lists a channel's participants.
func (s *Service) Members(ctx context.Context, me, channelID uuid.UUID) ([]Member, error) {
	if _, err := s.access(ctx, me, channelID); err != nil {
		return nil, err
	}
	rows, err := s.db.Queries.ListChannelMembers(ctx, channelID)
	if err != nil {
		return nil, err
	}
	out := make([]Member, 0, len(rows))
	for _, r := range rows {
		out = append(out, Member{ID: r.ID, FullName: r.FullName, Email: r.Email, AvatarID: r.AvatarID, Bio: r.Bio, Role: r.Role})
	}
	return out, nil
}

// memberIDs returns the user ids of a channel's members.
func (s *Service) memberIDs(ctx context.Context, channelID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := s.db.Queries.ListChannelMembers(ctx, channelID)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	return ids, nil
}

// AddMembers adds users to a private group (not direct).
func (s *Service) AddMembers(ctx context.Context, me, channelID uuid.UUID, userIDs []uuid.UUID) error {
	ch, err := s.access(ctx, me, channelID)
	if err != nil {
		return err
	}
	if ch.Kind == sqlcgen.ChatChannelKindDirect {
		return domain.ErrValidation.WithMessage("cannot add members to a direct chat")
	}
	for _, uid := range userIDs {
		if err := s.db.Queries.AddChatMember(ctx, sqlcgen.AddChatMemberParams{ChannelID: channelID, UserID: uid, Role: "member"}); err != nil {
			return err
		}
	}
	s.notifyChannel(ctx, ch, userIDs)
	return nil
}

// Leave removes the viewer from a channel.
func (s *Service) Leave(ctx context.Context, me, channelID uuid.UUID) error {
	_, err := s.db.Queries.RemoveChatMember(ctx, sqlcgen.RemoveChatMemberParams{ChannelID: channelID, UserID: me})
	return err
}

// Post sends a message (text and/or an already-uploaded attachment).
func (s *Service) Post(ctx context.Context, me, channelID uuid.UUID, author string, authorAvatar *uuid.UUID, body string, attachmentID *uuid.UUID) (*Message, error) {
	ch, err := s.access(ctx, me, channelID)
	if err != nil {
		return nil, err
	}
	body = strings.TrimSpace(body)
	if body == "" && attachmentID == nil {
		return nil, domain.ErrValidation.WithMessage("empty message")
	}
	if len(body) > 8000 {
		return nil, domain.ErrValidation.WithMessage("message too long (max 8000 characters)")
	}
	// Public channels create membership on first post so the author tracks reads.
	_ = s.db.Queries.AddChatMember(ctx, sqlcgen.AddChatMemberParams{ChannelID: channelID, UserID: me, Role: "member"})
	row, err := s.db.Queries.CreateChatMessage(ctx, sqlcgen.CreateChatMessageParams{ID: ids.New(), ChannelID: channelID, UserID: &me, Body: body, AttachmentID: attachmentID})
	if err != nil {
		return nil, err
	}
	_ = s.db.Queries.TouchChannel(ctx, channelID)
	_ = s.db.Queries.MarkChannelRead(ctx, sqlcgen.MarkChannelReadParams{ChannelID: channelID, UserID: me})
	msg := &Message{ID: row.ID, ChannelID: channelID, UserID: &me, AuthorName: author, AuthorAvatar: authorAvatar, Body: body, AttachmentID: attachmentID, CreatedAt: row.CreatedAt}
	s.deliver(ctx, ch, "message", msg)
	return msg, nil
}

// Messages lists a page of a channel's messages (newest first); pass before
// to page backwards.
func (s *Service) Messages(ctx context.Context, me, channelID uuid.UUID, before *uuid.UUID, limit int32) ([]Message, error) {
	if _, err := s.access(ctx, me, channelID); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := s.db.Queries.ListChatMessages(ctx, sqlcgen.ListChatMessagesParams{ChannelID: channelID, Before: before, RowLimit: limit})
	if err != nil {
		return nil, err
	}
	out := make([]Message, 0, len(rows))
	for _, r := range rows {
		m := Message{ID: r.ID, ChannelID: r.ChannelID, UserID: r.UserID, Body: r.Body, AttachmentID: r.AttachmentID, CreatedAt: r.CreatedAt, EditedAt: r.EditedAt}
		if r.AuthorName != nil {
			m.AuthorName = *r.AuthorName
		}
		m.AuthorAvatar = r.AuthorAvatar
		if r.AttachmentType != nil {
			m.AttachmentType = *r.AttachmentType
		}
		if r.AttachmentName != nil {
			m.AttachmentName = *r.AttachmentName
		}
		out = append(out, m)
	}
	return out, nil
}

// DeleteMessage soft-deletes the viewer's own message.
func (s *Service) DeleteMessage(ctx context.Context, me, channelID, messageID uuid.UUID) error {
	ch, err := s.access(ctx, me, channelID)
	if err != nil {
		return err
	}
	n, err := s.db.Queries.SoftDeleteChatMessage(ctx, sqlcgen.SoftDeleteChatMessageParams{ID: messageID, UserID: &me})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound.WithMessage("message not found or not yours")
	}
	s.deliver(ctx, ch, "message.deleted", map[string]any{"id": messageID, "channel_id": channelID})
	return nil
}

// MarkRead advances the viewer's read marker.
func (s *Service) MarkRead(ctx context.Context, me, channelID uuid.UUID) error {
	if _, err := s.access(ctx, me, channelID); err != nil {
		return err
	}
	return s.db.Queries.MarkChannelRead(ctx, sqlcgen.MarkChannelReadParams{ChannelID: channelID, UserID: me})
}

// deliver fans an event out: public → broadcast, otherwise to each member.
func (s *Service) deliver(ctx context.Context, ch *sqlcgen.ChatChannel, typ string, payload any) {
	if s.events == nil {
		return
	}
	b, _ := json.Marshal(payload)
	ev := events.ChatEvent{Type: typ, ChannelID: ch.ID, Payload: b}
	if ch.Kind == sqlcgen.ChatChannelKindPublic {
		s.events.PublishChatBroadcast(ctx, ev)
		return
	}
	mem, err := s.memberIDs(ctx, ch.ID)
	if err != nil {
		return
	}
	s.events.PublishChatToUsers(ctx, mem, ev)
}

// notifyChannel tells users a channel was created / they were added.
func (s *Service) notifyChannel(ctx context.Context, ch *sqlcgen.ChatChannel, extra []uuid.UUID) {
	if s.events == nil {
		return
	}
	b, _ := json.Marshal(map[string]any{"id": ch.ID, "kind": string(ch.Kind), "name": ch.Name})
	ev := events.ChatEvent{Type: "channel", ChannelID: ch.ID, Payload: b}
	if ch.Kind == sqlcgen.ChatChannelKindPublic {
		s.events.PublishChatBroadcast(ctx, ev)
		return
	}
	mem, _ := s.memberIDs(ctx, ch.ID)
	mem = append(mem, extra...)
	s.events.PublishChatToUsers(ctx, mem, ev)
}

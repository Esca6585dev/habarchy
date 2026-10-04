package admin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// ---- profile ----

type profileRequest struct {
	FullName string     `json:"full_name" validate:"required,max=100"`
	Bio      string     `json:"bio" validate:"max=500"`
	AvatarID *uuid.UUID `json:"avatar_id"`
}

func (h *Handlers) updateProfile(c *fiber.Ctx) error {
	var req profileRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	u, err := h.Auth.UpdateSelfProfile(c.UserContext(), middleware.UserID(c), req.FullName, req.Bio, req.AvatarID)
	if err != nil {
		return err
	}
	return httpx.OK(c, toUser(u))
}

// ---- attachments (avatars, chat images) ----

func (h *Handlers) uploadAttachment(c *fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return domain.ErrValidation.WithDetails(map[string]any{"file": "multipart field 'file' required (image)"})
	}
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	buf := make([]byte, fh.Size)
	if _, err := f.Read(buf); err != nil && fh.Size > 0 {
		return err
	}
	meta, err := h.Attachments.SaveImage(c.UserContext(), middleware.UserID(c), fh.Filename, buf)
	if err != nil {
		return err
	}
	return httpx.Created(c, meta)
}

func (h *Handlers) getAttachment(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "attachment_id")
	if err != nil {
		return err
	}
	a, err := h.Attachments.Get(c.UserContext(), id)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, a.ContentType)
	c.Set(fiber.HeaderCacheControl, "private, max-age=86400")
	return c.Send(a.Data)
}

// ---- chat ----

func (h *Handlers) chatChannels(c *fiber.Ctx) error {
	list, err := h.Chat.ListChannels(c.UserContext(), middleware.UserID(c))
	if err != nil {
		return err
	}
	return httpx.OK(c, list)
}

type createChannelRequest struct {
	Kind    string      `json:"kind" validate:"required,oneof=public private"`
	Name    string      `json:"name" validate:"required,max=100"`
	Topic   string      `json:"topic" validate:"max=300"`
	Members []uuid.UUID `json:"members"`
}

func (h *Handlers) createChatChannel(c *fiber.Ctx) error {
	var req createChannelRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	ch, err := h.Chat.CreateChannel(c.UserContext(), middleware.UserID(c), req.Kind, req.Name, req.Topic, req.Members)
	if err != nil {
		return err
	}
	return httpx.Created(c, fiber.Map{"id": ch.ID, "kind": ch.Kind, "name": ch.Name, "topic": ch.Topic, "created_at": ch.CreatedAt})
}

type directRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (h *Handlers) openDirect(c *fiber.Ctx) error {
	var req directRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	ch, err := h.Chat.OpenDirect(c.UserContext(), middleware.UserID(c), req.UserID)
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"id": ch.ID, "kind": ch.Kind})
}

func (h *Handlers) chatMembers(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	m, err := h.Chat.Members(c.UserContext(), middleware.UserID(c), id)
	if err != nil {
		return err
	}
	return httpx.OK(c, m)
}

type membersRequestChat struct {
	UserIDs []uuid.UUID `json:"user_ids" validate:"required,min=1"`
}

func (h *Handlers) addChatMembers(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	var req membersRequestChat
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if err := h.Chat.AddMembers(c.UserContext(), middleware.UserID(c), id, req.UserIDs); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) leaveChatChannel(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	if err := h.Chat.Leave(c.UserContext(), middleware.UserID(c), id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) chatMessages(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	var before *uuid.UUID
	if b := c.Query("before"); b != "" {
		bid, err := uuid.Parse(b)
		if err != nil {
			return domain.ErrValidation.WithDetails(map[string]any{"before": "must be a uuid"})
		}
		before = &bid
	}
	msgs, err := h.Chat.Messages(c.UserContext(), middleware.UserID(c), id, before, int32(c.QueryInt("limit", 50))) //nolint:gosec // bounded in service
	if err != nil {
		return err
	}
	return httpx.OK(c, msgs)
}

type postMessageRequest struct {
	Body         string     `json:"body" validate:"max=8000"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
}

func (h *Handlers) postChatMessage(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	var req postMessageRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	me, err := h.Auth.User(c.UserContext(), middleware.UserID(c))
	if err != nil {
		return err
	}
	msg, err := h.Chat.Post(c.UserContext(), me.ID, id, me.FullName, me.AvatarID, req.Body, req.AttachmentID)
	if err != nil {
		return err
	}
	return httpx.Created(c, msg)
}

func (h *Handlers) deleteChatMessage(c *fiber.Ctx) error {
	cid, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	mid, err := httpx.ParamUUID(c, "message_id")
	if err != nil {
		return err
	}
	if err := h.Chat.DeleteMessage(c.UserContext(), middleware.UserID(c), cid, mid); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) markChatRead(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "channel_id")
	if err != nil {
		return err
	}
	if err := h.Chat.MarkRead(c.UserContext(), middleware.UserID(c), id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

// chatStream is the per-user SSE feed of chat events.
func (h *Handlers) chatStream(c *fiber.Ctx) error {
	if h.Redis == nil {
		return domain.ErrInternal.WithMessage("live chat unavailable")
	}
	userID := middleware.UserID(c)
	c.Set(fiber.HeaderContentType, "text/event-stream")
	c.Set(fiber.HeaderCacheControl, "no-cache")
	c.Set(fiber.HeaderConnection, "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	ctx, cancel := contextWithCancel(c)
	ch := events.SubscribeChat(ctx, h.Redis, userID)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer cancel()
		_, _ = fmt.Fprint(w, "event: ready\ndata: {}\n\n")
		if err := w.Flush(); err != nil {
			return
		}
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					return
				}
				b, _ := json.Marshal(ev)
				_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, b)
				if err := w.Flush(); err != nil {
					return
				}
			case <-ticker.C:
				_, _ = fmt.Fprint(w, ": keep-alive\n\n")
				if err := w.Flush(); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}))
	return nil
}

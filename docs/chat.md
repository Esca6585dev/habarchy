# Team chat

An internal chat for the Habarchy instance's admin users (the people who log
into the panel / app), not the project's message recipients. It has profiles
with photos, public channels, private groups and direct messages, text and
image messages, unread counts and live delivery.

## Concepts

- **Profile** — each user has a name, an optional bio and an avatar (image,
  max 5 MB). Edit at Profile (web) / More → Profile (app);
  `PUT /api/admin/me/profile`, images via `POST /api/admin/attachments`.
- **Channels** — three kinds:
  - **public**: visible to everyone in the instance; every user is auto-joined.
  - **private**: a named group; only invited members see it.
  - **direct**: a 1:1 conversation, opened by picking a user (idempotent).
- **Messages** — text and/or one image attachment; you can delete your own.
  History pages backwards with `before=<message id>`.
- **Unread** — a per-user read marker; the channel list shows unread counts,
  the last message and, for direct chats, the other person.

## Live delivery

Messages are delivered over Redis pub/sub: public-channel messages go to a
broadcast topic, private/direct messages to each member's personal topic. The
web admin consumes `GET /api/admin/chat/stream` (SSE); the mobile app polls
every few seconds. Both need Redis (already required by the gateway).

## API (admin JWT)

| Call | Purpose |
|------|---------|
| `GET /api/admin/chat/channels` | my channels with unread, last message, DM peer |
| `POST /api/admin/chat/channels` | create `{kind: public\|private, name, members}` |
| `POST /api/admin/chat/direct` | open/get a DM `{user_id}` |
| `GET /api/admin/chat/channels/{id}/messages?before=&limit=` | history (newest first) |
| `POST /api/admin/chat/channels/{id}/messages` | send `{body, attachment_id}` |
| `DELETE /api/admin/chat/channels/{id}/messages/{mid}` | delete your message |
| `POST /api/admin/chat/channels/{id}/read` | mark read |
| `GET/POST /api/admin/chat/channels/{id}/members`, `POST .../leave` | membership |
| `GET /api/admin/chat/stream` | SSE feed |
| `POST /api/admin/attachments`, `GET /api/admin/attachments/{id}` | images |
| `PUT /api/admin/me/profile` | name, bio, avatar |

Attachments are stored in Postgres (capped, image types only) and served to
any signed-in instance user; ids are UUID v7 and not guessable.

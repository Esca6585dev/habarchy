package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/chat"
)

// loginAs creates a user and returns its bearer token + id.
func (f *flowEnv) loginAs(t *testing.T, email string) (string, string) {
	t.Helper()
	u, err := f.auth.CreateUser(context.Background(), email, "chat-pass-1234", "User "+email[:4])
	if err != nil {
		t.Fatal(err)
	}
	type loginOut struct {
		Tokens auth.Tokens `json:"tokens"`
	}
	lo := decode[loginOut](t, f.expect(f.do("POST", "/api/admin/auth/login", map[string]any{"email": email, "password": "chat-pass-1234"}, nil), 200).Data)
	return lo.Tokens.AccessToken, u.ID.String()
}

func (f *flowEnv) uploadAvatar(t *testing.T, tok string) string {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, _ := mw.CreateFormFile("file", "a.png")
	// 1x1 transparent PNG.
	png := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89,
		0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00, 0x01,
		0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}
	_, _ = fw.Write(png)
	_ = mw.Close()
	req := httptest.NewRequest("POST", "/api/admin/attachments", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+tok)
	res, err := f.app.Test(req, 10_000)
	if err != nil {
		t.Fatal(err)
	}
	var env struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&env)
	if res.StatusCode != 201 || env.Data.ID == "" {
		t.Fatalf("avatar upload: %d", res.StatusCode)
	}
	return env.Data.ID
}

// TestChat: profile + avatar, public channel broadcast, direct message, unread
// counts, image attachment, message delete, membership enforcement.
func TestChat(t *testing.T) {
	f := newFlowEnv(t)
	aliceTok := f.adminTok // from bootstrap
	bobTok, bobID := f.loginAs(t, "bob@habarchy.tm")
	_, carolID := f.loginAs(t, "carol@habarchy.tm")
	A, B := bearer(aliceTok), bearer(bobTok)

	// Alice sets a profile with an avatar.
	avatar := f.uploadAvatar(t, aliceTok)
	prof := f.expect(f.do("PUT", "/api/admin/me/profile", map[string]any{"full_name": "Alice A", "bio": "hi", "avatar_id": avatar}, A), 200)
	var me admin.UserResponse
	_ = json.Unmarshal(prof.Data, &me)
	if me.FullName != "Alice A" || me.Bio != "hi" || me.AvatarID == nil {
		t.Fatalf("profile: %+v", me)
	}
	// The avatar is served.
	if st := f.do("GET", "/api/admin/attachments/"+avatar, nil, A).Status; st != 200 {
		t.Fatalf("avatar fetch: %d", st)
	}

	// Alice creates a public channel; Bob sees it (auto-join) and posts.
	ch := decode[chatChannelLite](t, f.expect(f.do("POST", "/api/admin/chat/channels", map[string]any{"kind": "public", "name": "general"}, A), 201).Data)
	var bobChannels []chat.Channel
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels", nil, B), 200).Data, &bobChannels)
	if len(bobChannels) != 1 || bobChannels[0].ID.String() != ch.ID || bobChannels[0].Kind != "public" {
		t.Fatalf("bob channels: %+v", bobChannels)
	}
	// Alice posts; Bob has 1 unread, Alice 0.
	f.expect(f.do("POST", "/api/admin/chat/channels/"+ch.ID+"/messages", map[string]any{"body": "Salam hemmä"}, A), 201)
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels", nil, B), 200).Data, &bobChannels)
	if bobChannels[0].Unread != 1 || bobChannels[0].LastBody != "Salam hemmä" {
		t.Fatalf("bob unread: %+v", bobChannels[0])
	}
	// Bob reads, then has 0 unread.
	f.expect(f.do("POST", "/api/admin/chat/channels/"+ch.ID+"/read", nil, B), 204)
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels", nil, B), 200).Data, &bobChannels)
	if bobChannels[0].Unread != 0 {
		t.Fatalf("bob unread after read: %d", bobChannels[0].Unread)
	}

	// Direct chat Alice↔Bob (idempotent).
	dm1 := decode[chatChannelLite](t, f.expect(f.do("POST", "/api/admin/chat/direct", map[string]any{"user_id": bobID}, A), 200).Data)
	dm2 := decode[chatChannelLite](t, f.expect(f.do("POST", "/api/admin/chat/direct", map[string]any{"user_id": me.ID.String()}, B), 200).Data)
	if dm1.ID != dm2.ID || dm1.Kind != "direct" {
		t.Fatalf("dm not shared: %s %s", dm1.ID, dm2.ID)
	}
	// Carol is not a member → cannot read the DM.
	carolTok, _ := f.loginAs(t, "carol2@habarchy.tm")
	_ = carolID
	if st := f.do("GET", "/api/admin/chat/channels/"+dm1.ID+"/messages", nil, bearer(carolTok)).Status; st != 403 {
		t.Fatalf("carol DM access: %d", st)
	}
	// Alice sends Bob an image message; Bob lists the DM and sees it.
	img := f.uploadAvatar(t, aliceTok)
	posted := decode[chat.Message](t, f.expect(f.do("POST", "/api/admin/chat/channels/"+dm1.ID+"/messages", map[string]any{"body": "gör", "attachment_id": img}, A), 201).Data)
	var msgs []chat.Message
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels/"+dm1.ID+"/messages", nil, B), 200).Data, &msgs)
	if len(msgs) != 1 || msgs[0].Body != "gör" || msgs[0].AttachmentID == nil || msgs[0].AuthorName != "Alice A" {
		t.Fatalf("dm messages: %+v", msgs)
	}
	// Bob now has the DM with an unread.
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels", nil, B), 200).Data, &bobChannels)
	var dmFound bool
	for _, c := range bobChannels {
		if c.ID.String() == dm1.ID {
			dmFound = true
			if c.Unread != 1 || c.PeerName != "Alice A" || c.PeerID == nil {
				t.Fatalf("bob dm row: %+v", c)
			}
		}
	}
	if !dmFound {
		t.Fatalf("bob missing dm: %+v", bobChannels)
	}
	// Bob cannot delete Alice's message; Alice can.
	if st := f.do("DELETE", "/api/admin/chat/channels/"+dm1.ID+"/messages/"+posted.ID.String(), nil, B).Status; st != 404 {
		t.Fatalf("bob delete alice msg: %d", st)
	}
	f.expect(f.do("DELETE", "/api/admin/chat/channels/"+dm1.ID+"/messages/"+posted.ID.String(), nil, A), 204)
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels/"+dm1.ID+"/messages", nil, A), 200).Data, &msgs)
	if len(msgs) != 0 {
		t.Fatalf("after delete: %+v", msgs)
	}

	// Private group with Bob; Carol excluded.
	grp := decode[chatChannelLite](t, f.expect(f.do("POST", "/api/admin/chat/channels", map[string]any{"kind": "private", "name": "secret", "members": []string{bobID}}, A), 201).Data)
	var members []chat.Member
	_ = json.Unmarshal(f.expect(f.do("GET", "/api/admin/chat/channels/"+grp.ID+"/members", nil, A), 200).Data, &members)
	if len(members) != 2 {
		t.Fatalf("group members: %+v", members)
	}
	if st := f.do("GET", "/api/admin/chat/channels/"+grp.ID+"/messages", nil, bearer(carolTok)).Status; st != 403 {
		t.Fatalf("carol group access: %d", st)
	}
}

type chatChannelLite struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

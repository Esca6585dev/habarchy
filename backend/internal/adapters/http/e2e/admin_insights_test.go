package e2e

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/public"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/app/stats"
)

// TestAdminInsights drives the admin read models after real traffic.
func TestAdminInsights(t *testing.T) {
	f := newFlowEnv(t)
	key := f.key()
	tok := bearer(f.adminTok)
	pid := f.project.ID.String()
	pbase := "/api/admin/projects/" + pid

	// Live stream: run the app on a real socket so SSE can be read.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = f.app.Listener(ln) }()
	t.Cleanup(func() { _ = f.app.Shutdown() })
	streamURL := "http://" + ln.Addr().String() + "/api/admin/stream?project_id=" + pid + "&access_token=" + f.adminTok
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "text/event-stream") {
		t.Fatalf("stream: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	sse := bufio.NewReader(res.Body)
	readEvent := func() (string, events.Event) {
		var typ string
		var ev events.Event
		for {
			line, err := sse.ReadString('\n')
			if err != nil {
				t.Fatalf("sse read: %v", err)
			}
			line = strings.TrimRight(line, "\n")
			switch {
			case strings.HasPrefix(line, "event: "):
				typ = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				_ = json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &ev)
			case line == "":
				if typ != "" {
					return typ, ev
				}
			}
		}
	}
	if typ, _ := readEvent(); typ != "ready" {
		t.Fatalf("first event %q", typ)
	}
	// Unauthenticated stream is rejected.
	bad, _ := http.Get("http://" + ln.Addr().String() + "/api/admin/stream") //nolint:noctx // test
	if bad.StatusCode != 401 {
		t.Fatalf("stream without token: %d", bad.StatusCode)
	}
	bad.Body.Close()

	// ---- traffic: one ok, one permanent failure ----
	type accepted struct {
		ID uuid.UUID `json:"id"`
	}
	ok := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365123456", "template": "otp", "data": map[string]any{"code": "1", "minutes": 5}}, key), 202).Data)
	f.waitStatus(t, ok.ID.String(), "sent")
	seen := map[string]bool{}
	for i := 0; i < 2; i++ {
		typ, ev := readEvent()
		seen[typ] = true
		if ev.MessageID == nil || *ev.MessageID != ok.ID || ev.ProjectID != f.project.ID {
			t.Fatalf("event %s: %+v", typ, ev)
		}
	}
	if !seen["message.queued"] || !seen["message.sent"] {
		t.Fatalf("sse events: %v", seen)
	}
	bad2 := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365000999", "body": "x"}, key), 202).Data)
	f.waitStatus(t, bad2.ID.String(), "failed")
	f.hooks.wait(t, "message.failed")

	// ---- admin message log ----
	list := f.expect(f.do("GET", pbase+"/messages?limit=10", nil, tok), 200)
	var rows []admin.MessageRow
	_ = json.Unmarshal(list.Data, &rows)
	if len(rows) != 2 || rows[0].ID != bad2.ID {
		t.Fatalf("log: %+v", rows)
	}
	list = f.expect(f.do("GET", pbase+"/messages?status=failed&search=000999", nil, tok), 200)
	_ = json.Unmarshal(list.Data, &rows)
	if len(rows) != 1 || rows[0].ErrorCode != "provider_rejected" {
		t.Fatalf("filtered log: %+v", rows)
	}
	f.expect(f.do("GET", pbase+"/messages?status=nope", nil, tok), 400)

	detail := f.expect(f.do("GET", pbase+"/messages/"+ok.ID.String(), nil, tok), 200)
	var d struct {
		Message  admin.MessageRow       `json:"message"`
		Events   []public.EventResponse `json:"events"`
		Webhooks []admin.WebhookRow     `json:"webhooks"`
	}
	_ = json.Unmarshal(detail.Data, &d)
	if len(d.Events) < 2 || len(d.Webhooks) < 1 || d.Webhooks[0].Status != "delivered" || d.Webhooks[0].ResponseCode == nil || *d.Webhooks[0].ResponseCode != 200 {
		t.Fatalf("detail: events=%d webhooks=%+v", len(d.Events), d.Webhooks)
	}
	// raw provider response is in the sent event payload
	foundRaw := false
	for _, e := range d.Events {
		if e.Type == "sent" && strings.Contains(string(e.Payload), `"raw"`) {
			foundRaw = true
		}
	}
	if !foundRaw {
		t.Fatalf("sent event lacks raw provider response: %+v", d.Events)
	}

	// ---- resend failed -> goes through again (gateway still rejects 999) ----
	f.expect(f.do("POST", pbase+"/messages/"+ok.ID.String()+"/resend", nil, tok), 409) // sent, not resendable
	re := decode[admin.MessageRow](t, f.expect(f.do("POST", pbase+"/messages/"+bad2.ID.String()+"/resend", nil, tok), 202).Data)
	if re.Status != "queued" || re.Attempts != 0 {
		t.Fatalf("resend: %+v", re)
	}
	f.waitStatus(t, bad2.ID.String(), "failed")

	// ---- webhook log + resend ----
	hooks := f.expect(f.do("GET", pbase+"/webhooks", nil, tok), 200)
	var wrows []admin.WebhookRow
	_ = json.Unmarshal(hooks.Data, &wrows)
	if len(wrows) < 3 {
		t.Fatalf("webhook log: %d", len(wrows))
	}
	before := len(f.hooks.waitN(t, "message.sent", 1))
	rs := decode[admin.WebhookRow](t, f.expect(f.do("POST", pbase+"/webhooks/"+d.Webhooks[0].ID.String()+"/resend", nil, tok), 202).Data)
	if rs.Status != "pending" || rs.Attempts != 0 {
		t.Fatalf("webhook resend: %+v", rs)
	}
	f.hooks.waitN(t, "message.sent", before+1)
	f.expect(f.do("GET", pbase+"/webhooks/"+uuid.NewString(), nil, tok), 404)

	// ---- dashboard + latency ----
	dash := decode[stats.Dashboard](t, f.expect(f.do("GET", pbase+"/dashboard?days=7", nil, tok), 200).Data)
	if dash.Totals.Total != 2 || dash.Totals.Sent != 1 || dash.Totals.Failed != 1 || dash.Latency.Samples != 1 || dash.Latency.P95SentSec <= 0 {
		t.Fatalf("dashboard: %+v", dash.Totals)
	}
	if len(dash.Daily) != 1 || dash.Daily[0].Channel != "sms" || len(dash.Failures) != 1 {
		t.Fatalf("daily: %+v failures=%d", dash.Daily, len(dash.Failures))
	}
	f.expect(f.do("GET", pbase+"/dashboard?from=2026-01-10&to=2026-01-01", nil, tok), 400)

	// ---- usage: empty until aggregated, then populated; csv export; public endpoint ----
	usage := f.expect(f.do("GET", pbase+"/usage", nil, tok), 200)
	var urows []stats.UsageRow
	_ = json.Unmarshal(usage.Data, &urows)
	if len(urows) != 0 {
		t.Fatalf("usage before aggregation: %+v", urows)
	}
	st := stats.New(f.db, nil)
	if n, err := st.AggregateUsage(context.Background(), time.Now().UTC()); err != nil || n != 1 {
		t.Fatalf("aggregate: %d %v", n, err)
	}
	usage = f.expect(f.do("GET", pbase+"/usage?group_by=channel", nil, tok), 200)
	_ = json.Unmarshal(usage.Data, &urows)
	if len(urows) != 1 || urows[0].Channel != "sms" || urows[0].Sent != 1 || urows[0].Failed != 1 {
		t.Fatalf("usage: %+v", urows)
	}
	csvRes := f.do("GET", pbase+"/usage?format=csv", nil, tok)
	if csvRes.Status != 200 || !strings.HasPrefix(string(csvRes.Raw), "day,channel,queued") || !strings.Contains(string(csvRes.Raw), ",sms,") {
		t.Fatalf("csv: %d %s", csvRes.Status, csvRes.Raw)
	}
	pub := f.expect(f.do("GET", "/api/v1/usage?group_by=day", nil, key), 200)
	_ = json.Unmarshal(pub.Data, &urows)
	if len(urows) != 1 || urows[0].Day == "" {
		t.Fatalf("public usage: %+v", urows)
	}

	// ---- health, audit, users, overview, contacts/devices views ----
	health := decode[stats.Health](t, f.expect(f.do("GET", pbase+"/health", nil, tok), 200).Data)
	if len(health.Providers) != 1 || health.Providers[0].OkCount != 1 || health.Providers[0].FailCount != 1 || health.Providers[0].Status != "failing" {
		t.Fatalf("health providers: %+v", health.Providers)
	}
	if len(health.Queues) == 0 {
		t.Fatalf("health queues empty")
	}
	audit := f.expect(f.do("GET", pbase+"/audit-logs", nil, tok), 200)
	if !strings.Contains(string(audit.Data), `"project.create"`) || !strings.Contains(string(audit.Data), `"provider.create"`) {
		t.Fatalf("audit: %s", truncateStr(string(audit.Data), 300))
	}
	f.expect(f.do("POST", "/api/admin/users", map[string]any{"email": "viewer@habarchy.tm", "password": "viewer-pass-1"}, tok), 201)
	f.expect(f.do("POST", "/api/admin/users", map[string]any{"email": "viewer@habarchy.tm", "password": "viewer-pass-1"}, tok), 409)
	users := f.expect(f.do("GET", "/api/admin/users", nil, tok), 200)
	if !strings.Contains(string(users.Data), "viewer@habarchy.tm") {
		t.Fatalf("users: %s", users.Data)
	}
	ov := f.expect(f.do("GET", "/api/admin/overview", nil, tok), 200)
	if !strings.Contains(string(ov.Data), `"projects":1`) {
		t.Fatalf("overview: %s", ov.Data)
	}
	f.expect(f.do("GET", pbase+"/contacts", nil, tok), 200)
	f.expect(f.do("GET", pbase+"/devices", nil, tok), 200)

	// Viewer role can read the log but not resend.
	f.expect(f.do("PUT", pbase+"/members", map[string]any{"email": "viewer@habarchy.tm", "role": "viewer"}, tok), 200)
	vlogin := f.expect(f.do("POST", "/api/admin/auth/login", map[string]any{"email": "viewer@habarchy.tm", "password": "viewer-pass-1"}, nil), 200)
	var lo struct {
		Tokens struct {
			Access string `json:"access_token"`
		} `json:"tokens"`
	}
	_ = json.Unmarshal(vlogin.Data, &lo)
	viewer := bearer(lo.Tokens.Access)
	f.expect(f.do("GET", pbase+"/messages", nil, viewer), 200)
	f.expect(f.do("POST", pbase+"/messages/"+bad2.ID.String()+"/resend", nil, viewer), 403)
	f.expect(f.do("GET", pbase+"/audit-logs", nil, viewer), 403)

	// Metrics endpoint is not mounted in tests; the registry still counts.
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

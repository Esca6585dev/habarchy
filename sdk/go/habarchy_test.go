package habarchy_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	habarchy "github.com/Esca6585dev/habarchy/sdk/go"
)

func TestSendMessageSignedAndEnvelope(t *testing.T) {
	const key = "hb_live_testkey"
	fixed := time.Unix(1_700_000_000, 0)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if r.Header.Get("X-Api-Key") != key || r.Header.Get("X-Timestamp") != "1700000000" {
			t.Errorf("headers %v", r.Header)
		}
		// Independent re-implementation of the server-side check.
		sum := sha256.Sum256(body)
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte("1700000000\nPOST\n" + r.URL.Path + "\n" + hex.EncodeToString(sum[:])))
		if r.Header.Get("X-Signature") != hex.EncodeToString(mac.Sum(nil)) {
			t.Errorf("bad signature")
		}
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["to"] != "+99365123456" || req["channel"] != "sms" {
			t.Errorf("body %s", body)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"data":{"id":"m1","status":"queued","channel":"sms","to":"+99365123456"},"meta":null,"error":null}`))
	}))
	defer srv.Close()

	c := habarchy.New(srv.URL+"/", key, habarchy.WithSigning())
	c.Now = func() time.Time { return fixed }
	acc, err := c.SendMessage(context.Background(), habarchy.SendMessageRequest{Channel: "sms", To: habarchy.To("+99365123456"), Template: "otp", Data: map[string]any{"code": "1234"}})
	if err != nil || acc.ID != "m1" || acc.Status != "queued" || acc.Duplicate {
		t.Fatalf("%+v %v", acc, err)
	}
}

func TestRecipientShapesAndErrors(t *testing.T) {
	b, _ := json.Marshal(habarchy.ToExternal("user-42"))
	if string(b) != `{"external_id":"user-42"}` {
		t.Fatalf("external: %s", b)
	}
	b, _ = json.Marshal(habarchy.To("a@b.tm"))
	if string(b) != `"a@b.tm"` {
		t.Fatalf("address: %s", b)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/otp/verify"):
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"data":null,"error":{"code":"unauthorized","message":"invalid api key"}}`))
		case strings.HasSuffix(r.URL.Path, "/messages/batch"):
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"data":{"id":"b1","status":"queued","total":2},"meta":{"accepted":2,"rejected":0}}`))
		case strings.Contains(r.URL.Path, "/messages/m1"):
			_, _ = w.Write([]byte(`{"data":{"message":{"id":"m1","status":"delivered","body":"hi"},"events":[{"type":"sent"}]}}`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		}
	}))
	defer srv.Close()
	c := habarchy.New(srv.URL, "hb_test_x")
	ctx := context.Background()

	_, err := c.VerifyOTP(ctx, "+993", "1234")
	if !habarchy.IsCode(err, "unauthorized") {
		t.Fatalf("want api error, got %v", err)
	}
	var apiErr *habarchy.Error
	if !strings.Contains(err.Error(), "401") {
		t.Fatalf("status missing: %v", err)
	}
	_ = apiErr

	batch, err := c.SendBatch(ctx, habarchy.SendBatchRequest{Channel: "sms", Body: "x", Recipients: []habarchy.BatchRecipient{{To: habarchy.To("+1")}, {To: habarchy.To("+2")}}})
	if err != nil || batch.ID != "b1" || batch.Accepted != 2 || batch.Total != 2 {
		t.Fatalf("%+v %v", batch, err)
	}
	d, err := c.GetMessage(ctx, "m1")
	if err != nil || d.Message.Status != "delivered" || len(d.Events) != 1 {
		t.Fatalf("%+v %v", d, err)
	}
	_, err = c.RegisterDevice(ctx, habarchy.RegisterDeviceRequest{Token: "t", Platform: "android"})
	if !habarchy.IsCode(err, "http_error") {
		t.Fatalf("non-json error: %v", err)
	}
}

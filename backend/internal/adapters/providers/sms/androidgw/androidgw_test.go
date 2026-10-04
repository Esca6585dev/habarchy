package androidgw_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/androidgw"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

type fakeStore struct {
	mu       sync.Mutex
	rows     map[uuid.UUID]androidgw.Outcome
	expired  int
	onCreate func(id uuid.UUID)
}

func newStore() *fakeStore { return &fakeStore{rows: map[uuid.UUID]androidgw.Outcome{}} }

func (f *fakeStore) Enqueue(_ context.Context, _, _ uuid.UUID, _, _ string, _ int, _ time.Time) (uuid.UUID, error) {
	id := uuid.New()
	f.mu.Lock()
	f.rows[id] = androidgw.Outcome{Status: "pending"}
	f.mu.Unlock()
	if f.onCreate != nil {
		f.onCreate(id)
	}
	return id, nil
}
func (f *fakeStore) Outcome(_ context.Context, id uuid.UUID) (androidgw.Outcome, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rows[id], nil
}
func (f *fakeStore) Expire(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.expired++
	f.rows[id] = androidgw.Outcome{Status: "expired"}
	return nil
}
func (f *fakeStore) set(id uuid.UUID, o androidgw.Outcome) {
	f.mu.Lock()
	f.rows[id] = o
	f.mu.Unlock()
}

func TestParseConfig(t *testing.T) {
	if _, err := androidgw.ParseConfig(json.RawMessage(`{"gateway_key":"short"}`)); err == nil {
		t.Fatal("short key accepted")
	}
	if _, err := androidgw.ParseConfig(json.RawMessage(`{"gateway_key":"0123456789abcdef","sim_slot":2}`)); err == nil {
		t.Fatal("bad sim slot accepted")
	}
	cfg, err := androidgw.ParseConfig(json.RawMessage(`{"gateway_key":"0123456789abcdef"}`))
	if err != nil || cfg.TimeoutSec != 45 || cfg.SimSlot != 0 {
		t.Fatalf("defaults: %+v %v", cfg, err)
	}
}

func TestSendWaitsForPhone(t *testing.T) {
	st := newStore()
	st.onCreate = func(id uuid.UUID) {
		go func() { time.Sleep(30 * time.Millisecond); st.set(id, androidgw.Outcome{Status: "sent", Parts: 2}) }()
	}
	p := androidgw.New(uuid.New(), androidgw.Config{GatewayKey: "0123456789abcdef", TimeoutSec: 10}, st)
	p.Poll = 10 * time.Millisecond
	res, err := p.Send(context.Background(), ports.SMSMessage{To: "+99365123456", Text: "hi", Ref: uuid.New().String()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := uuid.Parse(res.ProviderMessageID); err != nil || res.Raw["parts"] != 2 {
		t.Fatalf("result %+v", res)
	}
}

func TestSendPhoneFailure(t *testing.T) {
	st := newStore()
	st.onCreate = func(id uuid.UUID) {
		st.set(id, androidgw.Outcome{Status: "failed", ErrorCode: "null_pdu", ErrorMessage: "bad pdu"})
	}
	p := androidgw.New(uuid.New(), androidgw.Config{GatewayKey: "0123456789abcdef", TimeoutSec: 10}, st)
	p.Poll = 5 * time.Millisecond
	_, err := p.Send(context.Background(), ports.SMSMessage{To: "+99365123456", Text: "hi", Ref: uuid.New().String()})
	var pe *ports.ProviderError
	if !errors.As(err, &pe) || pe.Retryable || pe.Code != "gateway_null_pdu" {
		t.Fatalf("got %v", err)
	}
	st2 := newStore()
	st2.onCreate = func(id uuid.UUID) { st2.set(id, androidgw.Outcome{Status: "failed", ErrorCode: "no_service"}) }
	p2 := androidgw.New(uuid.New(), androidgw.Config{GatewayKey: "0123456789abcdef", TimeoutSec: 10}, st2)
	p2.Poll = 5 * time.Millisecond
	_, err = p2.Send(context.Background(), ports.SMSMessage{To: "+99365123456", Text: "hi", Ref: uuid.New().String()})
	if !errors.As(err, &pe) || !pe.Retryable {
		t.Fatalf("no_service should be retryable: %v", err)
	}
}

func TestSendCancelledExpiresRow(t *testing.T) {
	st := newStore()
	p := androidgw.New(uuid.New(), androidgw.Config{GatewayKey: "0123456789abcdef", TimeoutSec: 10}, st)
	p.Poll = 5 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	_, err := p.Send(ctx, ports.SMSMessage{To: "+99365123456", Text: "hi", Ref: uuid.New().String()})
	var pe *ports.ProviderError
	if !errors.As(err, &pe) || !pe.Retryable || st.expired != 1 {
		t.Fatalf("got %v expired=%d", err, st.expired)
	}
	if _, err := p.Send(context.Background(), ports.SMSMessage{To: "+1", Text: "x", Ref: "not-a-uuid"}); err == nil {
		t.Fatal("bad ref accepted")
	}
}

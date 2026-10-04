package redis

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/config"
)

func testClient(t *testing.T) *Client {
	t.Helper()
	url := os.Getenv("HABARCHY_TEST_REDIS_URL")
	if url == "" {
		t.Skip("HABARCHY_TEST_REDIS_URL not set")
	}
	c, err := Connect(context.Background(), config.Redis{URL: url})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func TestIdempotency(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	key := "t:" + t.Name() + time.Now().String()
	if err := c.ReserveIdempotency(ctx, key, "m1", time.Minute); err != nil {
		t.Fatal(err)
	}
	err := c.ReserveIdempotency(ctx, key, "m2", time.Minute)
	var coll *ErrIdempotencyCollision
	if !errors.As(err, &coll) || coll.Existing != "m1" {
		t.Fatalf("expected collision with m1, got %v", err)
	}
	_ = c.ReleaseIdempotency(ctx, key)
	if err := c.ReserveIdempotency(ctx, key, "m3", time.Minute); err != nil {
		t.Fatalf("after release: %v", err)
	}
	_ = c.ReleaseIdempotency(ctx, key)
}

func TestTokenBucket(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	key := "t:" + t.Name() + time.Now().String()
	allowed := 0
	for i := 0; i < 5; i++ {
		ok, _, err := c.Allow(ctx, key, 10, 3)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("burst 3 should allow exactly 3 immediately, got %d", allowed)
	}
	start := time.Now()
	if err := c.Wait(ctx, key, 10, 3); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) < 50*time.Millisecond {
		t.Fatalf("Wait returned too early: %v", time.Since(start))
	}
	if ok, _, _ := c.Allow(ctx, key, 0, 1); !ok {
		t.Fatal("rate 0 means unlimited")
	}
}

func TestCountersAndOTP(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	key := "t:" + t.Name() + time.Now().String()
	v, created, err := c.Incr(ctx, key, 1, time.Minute)
	if err != nil || v != 1 || !created {
		t.Fatalf("first incr: %d %v %v", v, created, err)
	}
	v, created, _ = c.Incr(ctx, key, 2, time.Minute)
	if v != 3 || created {
		t.Fatalf("second incr: %d %v", v, created)
	}
	_ = c.Del(ctx, key)

	addr := "+99365123456" + time.Now().String()
	h := HashOTP(addr, "123456")
	if h == HashOTP("+99365000000", "123456") {
		t.Fatal("hash must depend on address")
	}
	if err := c.PutOTP(ctx, "p1", addr, h, time.Minute); err != nil {
		t.Fatal(err)
	}
	rec, ok, err := c.GetOTP(ctx, "p1", addr)
	if err != nil || !ok || rec.Hash != h || rec.Attempts != 0 {
		t.Fatalf("get otp: %+v %v %v", rec, ok, err)
	}
	if n, _ := c.IncrOTPAttempts(ctx, "p1", addr); n != 1 {
		t.Fatalf("attempts = %d", n)
	}
	_ = c.DelOTP(ctx, "p1", addr)
	if _, ok, _ := c.GetOTP(ctx, "p1", addr); ok {
		t.Fatal("otp should be gone")
	}
}

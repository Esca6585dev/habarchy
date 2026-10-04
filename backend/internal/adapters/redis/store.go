package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// ---- Idempotency ----

// ErrIdempotencyCollision is returned by ReserveIdempotency when the key
// already exists; Existing holds the stored value.
type ErrIdempotencyCollision struct{ Existing string }

func (e *ErrIdempotencyCollision) Error() string { return "idempotency key already used" }

// ReserveIdempotency stores value under key for ttl unless it exists.
func (c *Client) ReserveIdempotency(ctx context.Context, key, value string, ttl time.Duration) error {
	ok, err := c.c.SetNX(ctx, "idem:"+key, value, ttl).Result()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	existing, err := c.c.Get(ctx, "idem:"+key).Result()
	if err != nil && !errors.Is(err, goredis.Nil) {
		return err
	}
	return &ErrIdempotencyCollision{Existing: existing}
}

// ReleaseIdempotency removes a reservation (when the request failed after
// reserving).
func (c *Client) ReleaseIdempotency(ctx context.Context, key string) error {
	return c.c.Del(ctx, "idem:"+key).Err()
}

// ---- Token bucket rate limiter ----

// tokenBucket atomically refills and takes one token. KEYS[1] = bucket,
// ARGV[1] = rate per second, ARGV[2] = burst, ARGV[3] = now (ms).
// Returns {allowed(0/1), wait_ms}.
var tokenBucket = goredis.NewScript(`
local key   = KEYS[1]
local rate  = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local now   = tonumber(ARGV[3])
local data  = redis.call('HMGET', key, 'tokens', 'ts')
local tokens = tonumber(data[1])
local ts     = tonumber(data[2])
if tokens == nil then tokens = burst; ts = now end
local elapsed = math.max(0, now - ts) / 1000.0
tokens = math.min(burst, tokens + elapsed * rate)
local allowed = 0
local wait = 0
if tokens >= 1 then
  tokens = tokens - 1
  allowed = 1
else
  wait = math.ceil((1 - tokens) / rate * 1000)
end
redis.call('HSET', key, 'tokens', tokens, 'ts', now)
redis.call('PEXPIRE', key, math.ceil(burst / rate * 1000) + 1000)
return {allowed, wait}
`)

// Allow takes one token from bucket key at rate tokens/sec with the given
// burst. It returns whether the call may proceed and, if not, how long to
// wait for the next token.
func (c *Client) Allow(ctx context.Context, key string, rate float64, burst int) (bool, time.Duration, error) {
	if rate <= 0 {
		return true, 0, nil
	}
	if burst < 1 {
		burst = 1
	}
	res, err := tokenBucket.Run(ctx, c.c, []string{"rl:" + key}, rate, burst, time.Now().UnixMilli()).Int64Slice()
	if err != nil || len(res) != 2 {
		return false, 0, fmt.Errorf("redis: rate limit: %w", err)
	}
	return res[0] == 1, time.Duration(res[1]) * time.Millisecond, nil
}

// Wait blocks until a token is available or ctx ends. Used by workers to
// pace provider calls instead of failing them.
func (c *Client) Wait(ctx context.Context, key string, rate float64, burst int) error {
	for {
		ok, wait, err := c.Allow(ctx, key, rate, burst)
		if err != nil || ok {
			return err
		}
		if wait < 5*time.Millisecond {
			wait = 5 * time.Millisecond
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// ---- Counters (quotas, OTP rate limits) ----

// Incr increments key by n and sets ttl when the key is new. It returns the
// new value and whether the key was created by this call.
func (c *Client) Incr(ctx context.Context, key string, n int64, ttl time.Duration) (int64, bool, error) {
	pipe := c.c.TxPipeline()
	incr := pipe.IncrBy(ctx, key, n)
	pttl := pipe.PTTL(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, false, err
	}
	created := false
	if pttl.Val() < 0 { // -1: no expiry yet -> freshly created
		created = true
		_ = c.c.Expire(ctx, key, ttl).Err()
	}
	return incr.Val(), created, nil
}

// SetIfAbsent stores value with ttl only if key is missing.
func (c *Client) SetIfAbsent(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return c.c.SetNX(ctx, key, value, ttl).Result()
}

// GetInt returns an integer value or 0 when missing.
func (c *Client) GetInt(ctx context.Context, key string) (int64, error) {
	v, err := c.c.Get(ctx, key).Int64()
	if errors.Is(err, goredis.Nil) {
		return 0, nil
	}
	return v, err
}

// Del removes keys.
func (c *Client) Del(ctx context.Context, keys ...string) error { return c.c.Del(ctx, keys...).Err() }

// ---- OTP store ----

// OTPRecord is what verify needs; the code itself is never stored.
type OTPRecord struct {
	Hash      string
	Attempts  int64
	ExpiresAt time.Time
}

// HashOTP hashes a code with a per-address salt so equal codes for
// different recipients do not share a hash.
func HashOTP(address, code string) string {
	sum := sha256.Sum256([]byte(address + ":" + code))
	return hex.EncodeToString(sum[:])
}

// PutOTP stores the hashed code for address with ttl, replacing any
// previous code.
func (c *Client) PutOTP(ctx context.Context, scope, address, hash string, ttl time.Duration) error {
	key := otpKey(scope, address)
	pipe := c.c.TxPipeline()
	pipe.Del(ctx, key)
	pipe.HSet(ctx, key, "hash", hash, "attempts", 0, "exp", time.Now().Add(ttl).Unix())
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

// GetOTP loads the record; ok is false when none exists.
func (c *Client) GetOTP(ctx context.Context, scope, address string) (*OTPRecord, bool, error) {
	vals, err := c.c.HGetAll(ctx, otpKey(scope, address)).Result()
	if err != nil {
		return nil, false, err
	}
	if len(vals) == 0 {
		return nil, false, nil
	}
	var rec OTPRecord
	rec.Hash = vals["hash"]
	rec.Attempts, _ = strconv.ParseInt(vals["attempts"], 10, 64) // zero on parse failure is fine
	exp, _ := strconv.ParseInt(vals["exp"], 10, 64)
	rec.ExpiresAt = time.Unix(exp, 0)
	return &rec, true, nil
}

// IncrOTPAttempts counts a failed verification and returns the new total.
func (c *Client) IncrOTPAttempts(ctx context.Context, scope, address string) (int64, error) {
	return c.c.HIncrBy(ctx, otpKey(scope, address), "attempts", 1).Result()
}

// DelOTP removes the code after success or lockout.
func (c *Client) DelOTP(ctx context.Context, scope, address string) error {
	return c.c.Del(ctx, otpKey(scope, address)).Err()
}

func otpKey(scope, address string) string { return "otp:" + scope + ":" + address }

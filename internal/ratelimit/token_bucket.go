package ratelimit

import (
	"context"
	"errors"
	"math"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const tokenBucketScript = `
local rate = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local ttl = tonumber(ARGV[3])

local nowData = redis.call("TIME")
local now = (nowData[1] * 1000) + math.floor(nowData[2] / 1000)

local data = redis.call("HMGET", KEYS[1], "tokens", "ts")
local tokens = tonumber(data[1])
local ts = tonumber(data[2])

if tokens == nil then
  tokens = burst
  ts = now
else
  local delta = now - ts
  if delta < 0 then
    delta = 0
  end
  local refill = (delta / 1000) * rate
  tokens = math.min(burst, tokens + refill)
  ts = now
end

local allowed = 0
if tokens >= 1 then
  allowed = 1
  tokens = tokens - 1
end

redis.call("HMSET", KEYS[1], "tokens", tokens, "ts", ts)
redis.call("PEXPIRE", KEYS[1], ttl)
return allowed
`

type TokenBucket struct {
	client *redis.Client
	script *redis.Script
}

func NewTokenBucket(client *redis.Client) *TokenBucket {
	if client == nil {
		return nil
	}
	return &TokenBucket{
		client: client,
		script: redis.NewScript(tokenBucketScript),
	}
}

func (t *TokenBucket) Allow(ctx context.Context, key string, rate float64, burst int) (bool, error) {
	if t == nil || t.client == nil {
		return false, errors.New("rate limiter not configured")
	}
	if key == "" {
		return false, errors.New("rate limiter key is empty")
	}
	if rate <= 0 {
		return false, errors.New("rate limiter rate must be positive")
	}
	if burst <= 0 {
		return false, errors.New("rate limiter burst must be positive")
	}

	ttl := defaultBucketTTL(rate, burst)
	res, err := t.script.Run(
		ctx,
		t.client,
		[]string{key},
		rate,
		burst,
		int64(ttl/time.Millisecond),
	).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func defaultBucketTTL(rate float64, burst int) time.Duration {
	if rate <= 0 || burst <= 0 {
		return time.Second
	}
	seconds := math.Ceil((float64(burst) / rate) * 2)
	if seconds < 1 {
		seconds = 1
	}
	return time.Duration(seconds) * time.Second
}

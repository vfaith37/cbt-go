package security

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redis       *redis.Client
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(redis *redis.Client, maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:       redis,
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	ctx := context.Background()
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	pipe := rl.redis.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
	pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, rl.window)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false
	}

	requests := cmds[2].(*redis.IntCmd).Val()
	return requests <= int64(rl.maxRequests)
}

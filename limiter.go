package tools

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RateLimiter struct {
	Rdb *redis.Client
	Ctx context.Context
}

func NewRL(redisCli *redis.Client, ctx context.Context) *RateLimiter {
	return &RateLimiter{
		Rdb: redisCli,
		Ctx: ctx,
	}
}

func (rl *RateLimiter) Limit(key string, limit int, window int64) bool {
	cr, err := rl.Rdb.Get(rl.Ctx, key).Int()
	if err != nil && err != redis.Nil {
		logrus.Error(err)
		return false
	}

	if cr >= limit {
		return false
	}

	if err := rl.Rdb.Incr(rl.Ctx, key).Err(); err != nil {
		logrus.Error(err)
		return false
	}

	if cr == 0 {
		rl.Rdb.Expire(rl.Ctx, key, time.Duration(window)*time.Second)
	}
	return true
}

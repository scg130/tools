package rate_limit

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ilylx/gconv"
	"github.com/sirupsen/logrus"
)

type CountRateLimiter struct {
	Rdb *redis.Client
	Ctx context.Context
}

func NewCRL(redisCli *redis.Client, ctx context.Context) *CountRateLimiter {
	return &CountRateLimiter{
		Rdb: redisCli,
		Ctx: ctx,
	}
}

func (crl *CountRateLimiter) Allow(key string, limit int64, window int64) bool {
	pipe := crl.Rdb.Pipeline()
	cr, err := crl.Rdb.Get(crl.Ctx, key).Int64()
	if err != nil && err != redis.Nil {
		logrus.Error(err)
		return false
	}
	if gconv.Int64(cr) >= limit {
		return false
	}
	pipe.Incr(crl.Ctx, key)
	pipe.Expire(crl.Ctx, key, time.Duration(window)*time.Second)
	_, err = pipe.Exec(crl.Rdb.Context())
	if err != nil && err != redis.Nil {
		logrus.Error(err)
		return false
	}
	return true
}

package rate_limit

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ilylx/gconv"
	"github.com/sirupsen/logrus"
)

type CountRateLimiter struct {
	Rdb    *redis.Client
	Key    string
	Window time.Duration
	Limit  int64
}

func NewCRL(key string, window time.Duration, limit int64, redisCli *redis.Client) *CountRateLimiter {
	return &CountRateLimiter{
		Rdb:    redisCli,
		Key:    key,
		Window: window,
		Limit:  limit,
	}
}

func (crl *CountRateLimiter) Allow() (bool, error) {
	pipe := crl.Rdb.Pipeline()

	pipe.Get(crl.Rdb.Context(), crl.Key)
	pipe.Incr(crl.Rdb.Context(), crl.Key)
	cmds, err := pipe.Exec(crl.Rdb.Context())
	if err != nil && err != redis.Nil {
		logrus.Error(err)
		return false, err
	}

	count, _ := cmds[1].(*redis.IntCmd).Result()
	if gconv.Int64(count) > crl.Limit {
		return false, err
	}
	crl.Rdb.Expire(crl.Rdb.Context(), crl.Key, crl.Window)
	return true, nil
}

func (crl *CountRateLimiter) Done() {
	pipe := crl.Rdb.Pipeline()

	pipe.Decr(crl.Rdb.Context(), crl.Key)
	pipe.Expire(crl.Rdb.Context(), crl.Key, crl.Window)
	_, err := pipe.Exec(crl.Rdb.Context())
	if err != nil && err != redis.Nil {
		logrus.Error(err)
	}
}

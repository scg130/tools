package rate_limit

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type WindowRateLimiter struct {
	client   *redis.Client
	rate     int           // 在窗口内可以允许的请求数
	interval time.Duration // 时间窗口大小
	key      string        // Redis键名
}

func NewWindowRateLimiter(key string, interval time.Duration, rate int, client *redis.Client) *WindowRateLimiter {
	return &WindowRateLimiter{
		client:   client,
		rate:     rate,
		interval: interval,
		key:      key,
	}
}

func (limiter *WindowRateLimiter) Allow() (bool, error) {
	now := time.Now()

	pipe := limiter.client.Pipeline()

	pipe.ZRemRangeByScore(limiter.client.Context(), limiter.key, "-inf", strconv.FormatInt(now.Add(-limiter.interval).Unix(), 10))
	pipe.ZCard(limiter.client.Context(), limiter.key)
	_, err := pipe.ZAdd(limiter.client.Context(), limiter.key, &redis.Z{
		Score:  float64(now.Unix()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	}).Result()
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	pipe.Expire(limiter.client.Context(), limiter.key, limiter.interval)
	cmds, err := pipe.Exec(limiter.client.Context())
	if err != nil {
		logrus.Error(err)
		return false, err
	}

	count, _ := cmds[1].(*redis.IntCmd).Result()
	return int(count) < limiter.rate, nil
}

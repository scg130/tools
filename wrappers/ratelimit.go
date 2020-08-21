package wrappers

import (
	"github.com/micro/go-micro/v2/client"
	"go.uber.org/ratelimit"

	"context"
)

type clientWrapper struct {
	ratelimit.Limiter
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return c.Client.Call(ctx, req, rsp, opts...)
}

//rate  每秒可以通过多少可以请求
func NewRateLimitClientWrapper(rate int, opts ...ratelimit.Option) client.Wrapper {
	r := ratelimit.New(rate, opts...)
	return func(c client.Client) client.Client {
		return &clientWrapper{r, c}
	}
}

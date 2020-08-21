package wrappers

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

type TestWrapper struct {
	client.Client
}

func (tw *TestWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	conf := hystrix.CommandConfig{
		Timeout:                12,  //超时时间设置  单位毫秒
		MaxConcurrentRequests:  100, //最大请求数
		ErrorPercentThreshold:  4,
		SleepWindow:            2000, //过多长时间，熔断器再次检测是否开启。单位毫秒
		RequestVolumeThreshold: 5,    //请求阈值  熔断器是否打开首先要满足这个条件；这里的设置表示至少有5个请求才进行ErrorPercentThreshold错误百分比计算
	}
	command := req.Service() + req.Endpoint()
	hystrix.ConfigureCommand(command, conf)
	return hystrix.Do(command, func() error {
		return tw.Client.Call(ctx, req, rsp)
	}, func(e error) error {
		rsp = map[string]interface{}{
			"code": -1,
			"msg":  "timeout",
			"data": nil,
		}
		return nil
	})
}

func NewHystrixWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &TestWrapper{c}
	}
}

package tools

import (
	"context"
	"fmt"
	"io"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/server"
	opcplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/ratelimit"
)

const (
	REGISTRY_ADDR   = "127.0.0.1:2379"
	TRACER_ADDR     = "127.0.0.1:6831"
	TRACER_SRV_NAME = "tracer"
)

type selfTracer struct {
	T    opentracing.Tracer
	IC   io.Closer
	flag bool
}

var st selfTracer

func Reg() registry.Registry {
	return etcd.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{REGISTRY_ADDR}
	})
}

func Tracer() selfTracer {
	if st.flag {
		return st
	}

	t, ic, err := NewTracer(TRACER_SRV_NAME, TRACER_ADDR)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(t)
	st = selfTracer{
		t,
		ic,
		true,
	}
	return st
}

func NewService(name string) micro.Service {
	reg := Reg()
	sr := selector.NewSelector(func(op *selector.Options) {
		op.Registry = reg
		op.Strategy = selector.RoundRobin
	})
	return micro.NewService(
		micro.Version("latest"),
		micro.Name(name),
		micro.Registry(reg),
		micro.WrapHandler(opcplugin.NewHandlerWrapper(opentracing.GlobalTracer()), func(h server.HandlerFunc) server.HandlerFunc {
			return func(ctx context.Context, req server.Request, rsp interface{}) error {
				fmt.Println("server wrapper")
				h(ctx, req, rsp)
				return nil
			}
		}),
		micro.Selector(sr),
	)

}

func GetMicroClient(serviceName string) client.Client {
	reg := Reg()
	t := Tracer()
	srv := micro.NewService(
		micro.Registry(reg),
		micro.Name(serviceName),
		micro.WrapClient(
			opcplugin.NewClientWrapper(t.T),
			func(c client.Client) client.Client {
				r := ratelimit.New(100)
				return &struct {
					client.Client
					ratelimit.Limiter
				}{
					c, r,
				}
			},
		),
	)
	srv.Init()
	return srv.Client()
}

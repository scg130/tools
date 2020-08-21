package tools

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/server"
)

const (
	REGISTRY_ADDR = "127.0.0.1:2379"
)

func Reg() registry.Registry {
	return etcd.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{REGISTRY_ADDR}
	})
}

func NewService(name string, handlers ...server.HandlerWrapper) micro.Service {
	reg := Reg()
	sr := selector.NewSelector(func(op *selector.Options) {
		op.Registry = reg
		op.Strategy = selector.RoundRobin
	})
	return micro.NewService(
		micro.Version("latest"),
		micro.Name(name),
		micro.Registry(reg),
		micro.WrapHandler(handlers...),
		micro.Selector(sr),
	)

}

func GetMicroClient(serviceName string, wrappers ...client.Wrapper) client.Client {
	reg := Reg()
	srv := micro.NewService(
		micro.Registry(reg),
		micro.Name(serviceName),
		micro.WrapClient(
			wrappers...,
		),
	)
	srv.Init()
	return srv.Client()
}

package tools

import (
	"fmt"
	"os"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/server"
)

func Reg() registry.Registry {
	host := os.Getenv("ETCD_HOST")
	port := os.Getenv("ETCD_PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "2379"
	}
	registrAddr := fmt.Sprintf("%s:%s", host, port)
	return etcd.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{registrAddr}
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
		micro.RegisterInterval(time.Second*10),
		micro.RegisterTTL(15*time.Second),
		micro.WrapHandler(handlers...),
		micro.Selector(sr),
	)

}

func GetMicroClient(serviceName string, wrappers ...client.Wrapper) client.Client {
	reg := Reg()
	srv := micro.NewService(
		micro.Registry(reg),
		micro.RegisterInterval(time.Second*10),
		micro.RegisterTTL(15*time.Second),
		micro.Name(serviceName),
		micro.WrapClient(
			wrappers...,
		),
	)
	srv.Init()
	return srv.Client()
}

package handlers

import (
	"github.com/micro/go-micro/v2/server"
	opcplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/scg130/tools/wrappers"
	"os"
	"fmt"
)

func NewOpentracing(srvName string) server.HandlerWrapper {
	host := os.Getenv("TRACER_HOST")
	if host == "" {
		panic("tracerAddr is invalid")
	}
	tracerAddr := fmt.Sprintf("%s:5775",host)
	t,_,err := wrappers.NewTracer(srvName,tracerAddr)
	if err != nil {
		panic(err)
	}
	return opcplugin.NewHandlerWrapper(t)
}

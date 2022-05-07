package handlers

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/v2/server"
	opcplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/scg130/tools/wrappers"
)

func NewOpentracing(srvName string) server.HandlerWrapper {
	host := os.Getenv("TRACER_HOST")
	port := os.Getenv("TRACER_PORT")
	if host == "" {
		panic("tracerAddr is invalid")
	}
	tracerAddr := fmt.Sprintf("%s:%s", host, port)
	t, _, err := wrappers.NewTracer(srvName, tracerAddr)
	if err != nil {
		panic(err)
	}
	return opcplugin.NewHandlerWrapper(t)
}

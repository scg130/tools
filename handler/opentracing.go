package handler

import (
	"github.com/micro/go-micro/v2/api/server"
	"github.com/opentracing/opentracing-go"
)

func NewOpentracing() server.HandlerWrapper {
	return opcplugin.NewHandlerWrapper(opentracing.GlobalTracer())
}

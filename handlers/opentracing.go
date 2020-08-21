package handlers

import (
	"github.com/micro/go-micro/v2/server"
	opcplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func NewOpentracing() server.HandlerWrapper {
	return opcplugin.NewHandlerWrapper(opentracing.GlobalTracer())
}

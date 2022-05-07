package wrappers

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/micro/go-micro/v2/client"
	opcplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

const (
	TRACER_SRV_NAME = "tracer"
)

func NewTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}

type selfTracer struct {
	T    opentracing.Tracer
	IC   io.Closer
	flag bool
}

var st selfTracer

func NewTracerWrapper() client.Wrapper {
	if st.flag {
		return opcplugin.NewClientWrapper(st.T)
	}
	tracerAddr := ""
	host := os.Getenv("TRACER_HOST")
	port := os.Getenv("TRACER_PORT")
	if host == "" {
		panic("tracerAddr is invalid")
	}
	tracerAddr = fmt.Sprintf("%s:%s", host, port)
	t, ic, err := NewTracer(TRACER_SRV_NAME, tracerAddr)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(t)
	st = selfTracer{
		t,
		ic,
		true,
	}
	return opcplugin.NewClientWrapper(st.T)
}

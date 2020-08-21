package wrappers

import (
	"io"
	"time"

	"github.com/micro/go-micro/v2/client"
	opcplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

const (
	TRACER_ADDR     = "127.0.0.1:6831"
	TRACER_SRV_NAME = "tracer"
)

func newTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
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

	t, ic, err := newTracer(TRACER_SRV_NAME, TRACER_ADDR)
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

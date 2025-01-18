package tracer

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type nopTracerImpl struct {
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

func NewNopTracer(p TracerParam) TracerResult {
	return TracerResult{
		Tracer: &nopTracerImpl{
			tracer: noop.NewTracerProvider().Tracer(p.Config.ServiceIdentifier),
			propagator: propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		},
	}
}

func (t *nopTracerImpl) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName)
}

func (t *nopTracerImpl) StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	return t.tracer.Start(ctx, callerName)
}

func (t *nopTracerImpl) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	t.propagator.Inject(ctx, carrier)
}

func (t *nopTracerImpl) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return t.propagator.Extract(ctx, carrier)
}

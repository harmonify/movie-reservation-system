package tracer

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type jaegerTracerImpl struct {
	provider      trace.TracerProvider
	propagator    propagation.TextMapPropagator
	defaultTracer trace.Tracer
}

func NewJaegerTracer(cfg *TracerConfig, lc fx.Lifecycle) (Tracer, error) {
	exporter := otlptrace.NewUnstarted(
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		),
	)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceIdentifier),
			attribute.String("service.environment", cfg.Env),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not set OTel resources: %v\n", err)
	}

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(provider)

	// Set an error handler for trace exports
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		fmt.Printf("OTel trace export error: %v\n", err)
	}))

	t := &jaegerTracerImpl{
		propagator:    propagator,
		provider:      provider,
		defaultTracer: provider.Tracer(cfg.ServiceIdentifier),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := exporter.Start(ctx); err != nil {
				fmt.Printf("Failed to connect to OTel exporter: %v\n", err)
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err = exporter.Shutdown(ctx); err != nil {
				fmt.Printf("Failed to shutdown OTel exporter: %v\n", err)
				return err
			}
			return nil
		},
	})

	return t, nil
}

func (t *jaegerTracerImpl) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return t.defaultTracer.Start(ctx, spanName)
}

func (s *jaegerTracerImpl) StartSpanWithCaller(ctx context.Context, skip ...int) (context.Context, trace.Span) {
	finalSkip := 1
	if len(skip) > 0 && skip[0] > 0 {
		finalSkip = skip[0]
	}
	pc, _, _, _ := runtime.Caller(finalSkip)
	callerName := runtime.FuncForPC(pc).Name()
	ctx, span := s.Start(ctx, callerName)
	return ctx, span
}

func (s *jaegerTracerImpl) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func (s *jaegerTracerImpl) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

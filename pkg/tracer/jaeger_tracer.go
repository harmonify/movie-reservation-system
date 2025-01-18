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
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

func NewJaegerTracer(p TracerParam) TracerResult {
	exporter := otlptrace.NewUnstarted(
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(p.Config.OtelHost),
		),
	)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", p.Config.ServiceIdentifier),
			attribute.String("service.environment", p.Config.Env),
		),
	)
	if err != nil {
		fmt.Printf("Could not set tracer resources: %v\n", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(provider)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	// Set an error handler for trace exports
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		fmt.Println("OpenTelemetry trace export error", err.Error())
	}))

	t := &jaegerTracerImpl{
		tracer:     provider.Tracer(p.Config.ServiceIdentifier),
		propagator: propagator,
	}

	p.Lifecycle.Append(fx.Hook{
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

	return TracerResult{
		Tracer: t,
	}
}

func (t *jaegerTracerImpl) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName)
}

func (s *jaegerTracerImpl) StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()

	// segments := strings.Split(callerName, ".")
	// spanName := segments[len(segments)-1]

	ctx, span := s.Start(ctx, callerName)
	return ctx, span
}

func (s *jaegerTracerImpl) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	s.propagator.Inject(ctx, carrier)
}

func (s *jaegerTracerImpl) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return s.propagator.Extract(ctx, carrier)
}

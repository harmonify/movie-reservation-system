package tracer

import (
	"context"
	"fmt"
	"runtime"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type jaegerTracerImpl struct {
	exporter *otlptrace.Exporter
	config   *config.Config
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
			attribute.String("service.name", p.Config.AppName),
			attribute.String("service.environment", p.Config.Env),
		),
	)
	if err != nil {
		fmt.Printf("Could not set tracer resources: %v\n", err)
	}

	provider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resources),
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
		exporter: exporter,
		config:   p.Config,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := t.exporter.Start(ctx); err != nil {
				fmt.Printf("Failed to connect to OTel exporter: %v\n", err)
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err = t.exporter.Shutdown(ctx); err != nil {
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
	return otel.Tracer(t.config.AppName).Start(ctx, spanName)
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
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func (s *jaegerTracerImpl) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

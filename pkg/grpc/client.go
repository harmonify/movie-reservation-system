package grpc_pkg

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats/opentelemetry"
)

type GrpcClient struct {
	Conn   *grpc.ClientConn
	logger logger.Logger
}

// see https://github.com/grpc/grpc/blob/master/doc/service_config.md to know more about service config
var defaultServiceConfig = `{
    "loadBalancingConfig": [{"round_robin":{}}],
    "methodConfig": [{
        "name": [{"service": "grpc.examples.echo.Echo"}],
        "retryPolicy": {
            "MaxAttempts": 4,
            "InitialBackoff": ".01s",
            "MaxBackoff": ".01s",
            "BackoffMultiplier": 1.0,
            "RetryableStatusCodes": [ "UNAVAILABLE" ]
        }
    }]
}`

type GrpcClientParam struct {
	fx.In
	fx.Lifecycle

	Logger logger.Logger
}

type GrpcClientConfig struct {
	Address string `validate:"required"`
}

func NewGrpcClient(p GrpcClientParam, cfg *GrpcClientConfig, opts ...grpc.DialOption) (*GrpcClient, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	exporter, err := prometheus.New()
	if err != nil {
		p.Logger.Error(fmt.Sprintf("Failed to start prometheus exporter: %v", err))
		return nil, err
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))

	// Create a connection with timeout and other options
	finalOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(defaultServiceConfig),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		opentelemetry.DialOption(opentelemetry.Options{MetricsOptions: opentelemetry.MetricsOptions{MeterProvider: provider}}),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	finalOpts = append(finalOpts, opts...)
	conn, err := grpc.NewClient(cfg.Address, finalOpts...)
	if err != nil {
		return nil, err
	}

	// Add lifecycle hooks for the client connection
	p.Lifecycle.Append(fx.StartStopHook(
		func(ctx context.Context) error {
			// p.Logger.WithCtx(ctx).Info("Starting gRPC client connection...")
			return nil
		},
		func(ctx context.Context) error {
			p.Logger.WithCtx(ctx).Warn("Closing gRPC client connection...")
			return conn.Close()
		},
	))

	return &GrpcClient{
		Conn: conn,
	}, nil
}

package grpc_pkg

import (
	"context"
	"fmt"
	"net"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats/opentelemetry"
)

type GrpcServerParam struct {
	fx.In
	fx.Lifecycle

	Config *config.Config
	Logger logger.Logger
}

type GrpcServerResult struct {
	fx.Out

	GrpcServer *GrpcServer
}

type GrpcServer struct {
	is_started bool

	Server *grpc.Server
	cfg    *config.Config
	logger logger.Logger
}

func NewGrpcServer(
	p GrpcServerParam,
) GrpcServerResult {
	exporter, err := prometheus.New()
	if err != nil {
		p.Logger.Error(fmt.Sprintf("Failed to start prometheus exporter: %v", err))
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))

	server := grpc.NewServer(
		grpc.MaxConcurrentStreams(100),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		opentelemetry.ServerOption(opentelemetry.Options{MetricsOptions: opentelemetry.MetricsOptions{MeterProvider: provider}}),
	)

	g := &GrpcServer{
		Server: server,
		cfg:    p.Config,
		logger: p.Logger,
	}

	result := GrpcServerResult{
		GrpcServer: g,
	}

	p.Lifecycle.Append(fx.StartStopHook(
		func(ctx context.Context) error {
			return g.Start(ctx)
		},
		func(ctx context.Context) error {
			g.Shutdown(ctx)
			return nil
		},
	))

	return result
}

func (g *GrpcServer) Start(ctx context.Context) error {
	if g.is_started {
		g.logger.WithCtx(ctx).Warn(fmt.Sprintf(">> gRPC server is already running on port: %s", g.cfg.GrpcPort))
		return nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.cfg.GrpcPort))
	if err != nil {
		g.logger.WithCtx(ctx).Error(fmt.Sprintf(">> gRPC server failed to listen on port %s. error: %s", g.cfg.GrpcPort, err.Error()))
		return err
	}

	err = g.Server.Serve(listener)
	if err != nil {
		g.logger.WithCtx(ctx).Error(fmt.Sprintf(">> gRPC server failed to start. error: %s", err.Error()))
		return err
	}

	g.logger.WithCtx(ctx).Info(fmt.Sprintf(">> gRPC server is running on port: %s", g.cfg.GrpcPort))
	return nil
}

func (g *GrpcServer) Shutdown(ctx context.Context) {
	g.logger.WithCtx(ctx).Info(">> gRPC server shutting down...")
	g.Server.GracefulStop()
	g.logger.WithCtx(ctx).Info(">> gRPC server is shutdown")
}

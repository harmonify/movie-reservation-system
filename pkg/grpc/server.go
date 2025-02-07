package grpc_pkg

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
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

	Logger logger.Logger
}

type GrpcServerResult struct {
	fx.Out

	GrpcServer *GrpcServer
}

type GrpcServer struct {
	started bool
	mu      sync.RWMutex

	Server *grpc.Server
	cfg    *GrpcServerConfig
	logger logger.Logger
}

type GrpcServerConfig struct {
	GrpcPort int `validate:"required,numeric,min=1024,max=65535"`
}

func NewGrpcServer(
	p GrpcServerParam,
	cfg *GrpcServerConfig,
) (GrpcServerResult, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return GrpcServerResult{}, err
	}

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
		cfg:    cfg,
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

	return result, nil
}

func (g *GrpcServer) Start(ctx context.Context) error {
	if g.getStarted() {
		g.logger.WithCtx(ctx).Warn(fmt.Sprintf(">> gRPC server is already running on port: %d", g.cfg.GrpcPort))
		return nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.cfg.GrpcPort))
	if err != nil {
		g.logger.WithCtx(ctx).Error(fmt.Sprintf(">> gRPC server failed to listen on port %d. error: %v", g.cfg.GrpcPort, err.Error()))
		return err
	}

	// Start gRPC server in a goroutine
	// Wait for 1 second to see if the server is running
	// If the server is not running, return an error

	go func() {
		g.setStarted(true)
		if err := g.Server.Serve(listener); err != nil {
			g.setStarted(false)
			g.logger.WithCtx(ctx).Error(fmt.Sprintf(">> gRPC server failed to shutdown gracefully. error: %s", err.Error()))
		}
	}()

	time.Sleep(1 * time.Second)
	if g.getStarted() {
		g.logger.WithCtx(ctx).Info(fmt.Sprintf(">> gRPC server is running on port: %d", g.cfg.GrpcPort))
		return nil
	} else {
		g.logger.WithCtx(ctx).Error(fmt.Sprintf(">> gRPC server failed to start on port: %d", g.cfg.GrpcPort))
		return fmt.Errorf("gRPC server failed to start on port: %d", g.cfg.GrpcPort)
	}
}

func (g *GrpcServer) Shutdown(ctx context.Context) {
	g.logger.WithCtx(ctx).Info(">> gRPC server shutting down...")
	g.Server.GracefulStop()
	g.logger.WithCtx(ctx).Info(">> gRPC server is shut down")
}

func (g *GrpcServer) getStarted() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.started
}

func (g *GrpcServer) setStarted(started bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.started = started
}

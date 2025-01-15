package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcServerParam struct {
	fx.In

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
) (GrpcServerResult, error) {
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	g := &GrpcServer{
		Server: server,
		cfg:    p.Config,
		logger: p.Logger,
	}

	result := GrpcServerResult{
		GrpcServer: g,
	}

	err := g.Start()

	return result, err
}

func (g *GrpcServer) Start() error {
	if g.is_started {
		g.logger.Warn(fmt.Sprintf(">> gRPC server is already running on port: %s", g.cfg.GrpcPort))
		return nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.cfg.GrpcPort))
	if err != nil {
		g.logger.Error(fmt.Sprintf(">> gRPC server failed to listen on port %s. error: %s", g.cfg.GrpcPort, err.Error()))
		return err
	}

	err = g.Server.Serve(listener)
	if err != nil {
		g.logger.Error(fmt.Sprintf(">> gRPC server failed to start. error: %s", err.Error()))
		return err
	}

	g.logger.Info(fmt.Sprintf(">> gRPC server is running on port: %s", g.cfg.GrpcPort))
	return nil
}

func (g *GrpcServer) Shutdown() {
	sh := make(chan os.Signal, 2)
	signal.Notify(sh, os.Interrupt, syscall.SIGTERM)
	<-sh

	log.Println(">> gRPC server shutting down...")
	g.Server.GracefulStop()
	log.Println(">> gRPC server is shutdown")
}

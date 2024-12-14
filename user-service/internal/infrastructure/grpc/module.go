package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	config "github.com/harmonify/movie-reservation-system/pkg/configs"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var GRPCModule = fx.Module("grpc", fx.Provide(NewGRPCServer))

type GRPCServer struct {
	Server *grpc.Server
	cfg    *config.Config
	logger logger_shared.Logger
}

func NewGRPCServer(
	cfg *config.Config,
	logger logger_shared.Logger,
) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	return &GRPCServer{
		Server: server,
		cfg:    cfg,
		logger: logger,
	}
}

func (g *GRPCServer) Start() {
	var err error

	defer func() {
		if err != nil {
			g.logger.Error(fmt.Sprintf(">> Failed to start gRPC server on port : %s with error: %s", g.cfg.GRPCPort, err.Error()))
		}
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.cfg.GRPCPort))
	if err != nil {
		return
	}

	g.logger.Info(fmt.Sprintf(">> gRPC Server run on port: %s", g.cfg.GRPCPort))

	err = g.Server.Serve(listener)
	if err != nil {
		return
	}
}

func (g *GRPCServer) Shutdown() {
	sh := make(chan os.Signal, 2)
	signal.Notify(sh, os.Interrupt, syscall.SIGTERM)
	<-sh

	log.Println(">> Shutdown GRPC Server...")
	g.Server.GracefulStop()
	log.Println(">> GRPC Server already shutdown")
}

package grpc_driver

import (
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	"go.uber.org/fx"
)

var (
	GrpcModule = fx.Module(
		"grpc-driver",
		fx.Provide(
			func(p grpc_pkg.GrpcServerParam, cfg *config.MovieServiceConfig) (grpc_pkg.GrpcServerResult, error) {
				return grpc_pkg.NewGrpcServer(p, &grpc_pkg.GrpcServerConfig{
					GrpcPort: cfg.GrpcPort,
				})
			},
			NewMovieServiceServer,
		),
		fx.Invoke(RegisterMovieServiceServer),
	)
)

package grpc

import (
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/config"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"go.uber.org/fx"
)

type TheaterServiceClientParam struct {
	fx.In
	Client *grpc_pkg.GrpcClient
}

func NewTheaterServiceClient(p grpc_pkg.GrpcClientParam, cfg *config.MovieSearchServiceConfig) (theater_proto.TheaterServiceClient, error) {
	client, err := grpc_pkg.NewGrpcClient(p, &grpc_pkg.GrpcClientConfig{
		Address: cfg.GrpcTheaterServiceUrl,
	})
	if err != nil {
		return nil, err
	}
	return theater_proto.NewTheaterServiceClient(client.Conn), nil
}

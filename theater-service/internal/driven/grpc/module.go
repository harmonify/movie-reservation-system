package grpc

import (
	grpc_repository "github.com/harmonify/movie-reservation-system/theater-service/internal/driven/grpc/repository"
	"go.uber.org/fx"
)

var DrivenGrpcModule = fx.Module(
	"driven-grpc",
	grpc_repository.DrivenGrpcRepositoryModule,
	fx.Provide(
		NewAuthServiceClient,
		NewMovieServiceClient,
	),
)

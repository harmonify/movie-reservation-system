package grpc_repository

import "go.uber.org/fx"

var DrivenGrpcRepositoryModule = fx.Module(
	"driven-grpc-repository",
	fx.Provide(NewMovieGrpcRepository),
)

package grpc

import (
	"go.uber.org/fx"
)

var DrivenGrpcModule = fx.Module(
	"driven-grpc",
	fx.Provide(
		NewTheaterServiceClient,
	),
)

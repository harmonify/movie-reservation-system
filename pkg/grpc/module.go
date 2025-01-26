package grpc_pkg

import "go.uber.org/fx"

var GrpcModule = fx.Module(
	"grpc",
	fx.Provide(NewGrpcServer),
)


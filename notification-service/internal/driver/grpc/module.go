package grpc_driver

import (
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"go.uber.org/fx"
)

var (
	GrpcModule = fx.Module(
		"grpc-driver",
		fx.Provide(
			func(p grpc_pkg.GrpcServerParam, cfg *config.NotificationServiceConfig) (grpc_pkg.GrpcServerResult, error) {
				return grpc_pkg.NewGrpcServer(p, &grpc_pkg.GrpcServerConfig{
					GrpcPort: cfg.GrpcPort,
				})
			},
			NewNotificationServiceServer,
		),
		fx.Invoke(RegisterNotificationServiceServer),
	)
)

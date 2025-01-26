package grpc_driver

import (
	// email_grpc "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/grpc/email"
	"go.uber.org/fx"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
)

var (
	GrpcModule = fx.Module(
		"grpc-driver",
		grpc_pkg.GrpcModule,
		fx.Invoke(RegisterNotificationServiceServer),
		fx.Invoke(BootstrapGrpcServer),
	)
)

func BootstrapGrpcServer(h *grpc_pkg.GrpcServer) {}

package pkg

import (
	notification_proto "github.com/harmonify/movie-reservation-system/notification-service/internal/driven/proto/notification"
	"github.com/harmonify/movie-reservation-system/pkg/grpc"
)

func NewNotificationServiceGrpcClient(p grpc_pkg.GrpcClientParam) (notification_proto.NotificationServiceClient, error) {
	client, err := grpc_pkg.NewGrpcClient(p, &grpc_pkg.GrpcClientConfig{
		Address: p.Config.NotificationServiceGrpcAddress,
	})
	if err != nil {
		return nil, err
	}
	return notification_proto.NewNotificationServiceClient(client.Conn), nil
}

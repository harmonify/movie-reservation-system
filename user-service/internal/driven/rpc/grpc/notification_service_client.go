package grpc

import (
	"context"

	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	"go.uber.org/fx"
)

type NotificationServiceClientParam struct {
	fx.In
	Client *grpc_pkg.GrpcClient
}

type notificationServiceClient struct {
	client notification_proto.NotificationServiceClient
}

func NewNotificationServiceClient(p grpc_pkg.GrpcClientParam, cfg *config.UserServiceConfig) (shared.NotificationProvider, error) {
	client, err := grpc_pkg.NewGrpcClient(p, &grpc_pkg.GrpcClientConfig{
		Address: cfg.GrpcNotificationServiceUrl,
	})
	if err != nil {
		return nil, err
	}
	return &notificationServiceClient{
		client: notification_proto.NewNotificationServiceClient(client.Conn),
	}, nil
}

func (n *notificationServiceClient) SendEmail(ctx context.Context, p *notification_proto.SendEmailRequest) error {
	_, err := n.client.SendEmail(ctx, p)
	return err
}

func (n *notificationServiceClient) SendSms(ctx context.Context, p *notification_proto.SendSmsRequest) error {
	_, err := n.client.SendSms(ctx, p)
	return err
}

func (n *notificationServiceClient) BulkSendSms(ctx context.Context, p *notification_proto.BulkSendSmsRequest) error {
	_, err := n.client.BulkSendSms(ctx, p)
	return err
}

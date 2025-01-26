package rpc

import (
	"context"

	notification_pkg "github.com/harmonify/movie-reservation-system/notification-service/pkg"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type NotificationServiceClientParam struct {
	fx.In
	Client *grpc_pkg.GrpcClient
}

type notificationServiceClient struct {
	notificationService notification_proto.NotificationServiceClient
}

func NewNotificationServiceClient(p NotificationServiceClientParam) shared.NotificationProvider {
	return &notificationServiceClient{
		notificationService: notification_pkg.NewNotificationServiceClient(p.Client.Conn),
	}
}

func (n *notificationServiceClient) SendEmail(ctx context.Context, p *notification_proto.SendEmailRequest) error {
	_, err := n.notificationService.SendEmail(ctx, p)
	return err
}

func (n *notificationServiceClient) SendSms(ctx context.Context, p *notification_proto.SendSmsRequest) error {
	_, err := n.notificationService.SendSms(ctx, p)
	return err
}

func (n *notificationServiceClient) BulkSendSms(ctx context.Context, p *notification_proto.BulkSendSmsRequest) error {
	_, err := n.notificationService.BulkSendSms(ctx, p, grpc.WaitForReady(true))
	return err
}

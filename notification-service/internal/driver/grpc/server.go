package grpc_driver

import (
	"context"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/services"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/templates"
	notification_proto "github.com/harmonify/movie-reservation-system/notification-service/internal/driven/proto/notification"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func RegisterNotificationServiceServer(
	server *grpc_pkg.GrpcServer,
	handler notification_proto.NotificationServiceServer,
) {
	notification_proto.RegisterNotificationServiceServer(server.Server, handler)
}

type NotificationServiceServerParam struct {
	fx.In
	Logger               logger.Logger
	Tracer               tracer.Tracer
	ErrorMapper          error_pkg.ErrorMapper
	SmsService           services.SmsService
	EmailTemplateService services.EmailTemplateService
	EmailService         services.EmailService
	Util                 *util.Util
}

type notificationServiceServerImpl struct {
	notification_proto.UnimplementedNotificationServiceServer // Embedding for compatibility
	logger                                                    logger.Logger
	tracer                                                    tracer.Tracer
	errorMapper                                               error_pkg.ErrorMapper
	smsService                                                services.SmsService
	emailTemplateService                                      services.EmailTemplateService
	emailService                                              services.EmailService
	util                                                      *util.Util
}

func NewNotificationServiceServer(
	p NotificationServiceServerParam,
) notification_proto.NotificationServiceServer {
	return &notificationServiceServerImpl{
		UnimplementedNotificationServiceServer: notification_proto.UnimplementedNotificationServiceServer{},
		logger:                                 p.Logger,
		tracer:                                 p.Tracer,
		errorMapper:                            p.ErrorMapper,
		smsService:                             p.SmsService,
		emailTemplateService:                   p.EmailTemplateService,
		emailService:                           p.EmailService,
		util:                                   p.Util,
	}
}

func (s *notificationServiceServerImpl) SendEmail(
	ctx context.Context,
	req *notification_proto.SendEmailRequest,
) (*notification_proto.SendEmailResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// TODO: Notify listeners (useful for testing)
	// for _, listener := range s.listeners {
	// 	listener.OnEvent(ctx, event)
	// }

	if len(req.Recipients) <= 0 {
		s.logger.WithCtx(ctx).Error("Empty recipients")
		return nil, s.errorMapper.ToGrpcError(shared.EmptyRecipientError)
	}
	if req.Subject == "" {
		s.logger.WithCtx(ctx).Error("Empty subject")
		return nil, s.errorMapper.ToGrpcError(shared.EmptySubjectError)
	}
	if req.TemplateId == "" {
		s.logger.WithCtx(ctx).Error("Empty template id")
		return nil, s.errorMapper.ToGrpcError(shared.EmptyTemplateError)
	}

	tmplPath := templates.MapEmailTemplateIdToPath(req.GetTemplateId())
	if tmplPath == "" {
		s.logger.WithCtx(ctx).Error("Invalid email template id", zap.String("template_id", req.GetTemplateId()))
		return nil, s.errorMapper.ToGrpcError(shared.InvalidTemplateIdError)
	}

	tmplData, err := anypb.UnmarshalNew(req.GetTemplateData(), proto.UnmarshalOptions{})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to unmarshal email template data proto", zap.Error(err))
		return nil, s.errorMapper.ToGrpcError(shared.InvalidTemplateDataError)
	}

	tmplDataMap, err := s.util.StructUtil.ConvertProtoToMap(ctx, tmplData)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to convert email template data proto to a map", zap.Error(err))
		return nil, s.errorMapper.ToGrpcError(error_pkg.InternalServerError)
	}

	content, err := s.emailTemplateService.Render(ctx, tmplPath, tmplDataMap)
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	emailId, err := s.emailService.Send(ctx, shared.EmailMessage{
		Recipients: req.GetRecipients(),
		Subject:    req.GetSubject(),
		Body:       content,
		Type:       shared.EmailTypeHtml,
	})
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return &notification_proto.SendEmailResponse{EmailId: emailId}, nil
}

func (s *notificationServiceServerImpl) SendSms(
	ctx context.Context,
	req *notification_proto.SendSmsRequest,
) (*notification_proto.SendSmsResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// TODO: Notify listeners
	// for _, listener := range s.listeners {
	// 	listener.OnEvent(ctx, event)
	// }

	if req.Recipient == "" {
		s.logger.WithCtx(ctx).Error("Empty recipient")
		return nil, s.errorMapper.ToGrpcError(shared.EmptyRecipientError)
	}
	if req.Body == "" {
		s.logger.WithCtx(ctx).Error("Empty body")
		return nil, s.errorMapper.ToGrpcError(shared.EmptyBodyError)
	}

	smsId, err := s.smsService.Send(ctx, shared.SmsMessage{
		Recipient: req.GetRecipient(),
		Body:      req.GetBody(),
	})
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return &notification_proto.SendSmsResponse{SmsId: smsId}, nil
}

func (s *notificationServiceServerImpl) BulkSendSms(ctx context.Context, req *notification_proto.BulkSendSmsRequest) (*notification_proto.BulkSendSmsResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// TODO: Notify listeners
	// for _, listener := range s.listeners {
	// 	listener.OnEvent(ctx, event)
	// }

	if len(req.Recipients) <= 0 {
		s.logger.WithCtx(ctx).Error("Empty recipients")
		return nil, s.errorMapper.ToGrpcError(shared.EmptyRecipientError)
	}
	if req.Body == "" {
		s.logger.WithCtx(ctx).Error("Empty body")
		return nil, s.errorMapper.ToGrpcError(shared.EmptyBodyError)
	}

	smsIds, err := s.smsService.BulkSend(ctx, shared.BulkSmsMessage{
		Recipients: req.GetRecipients(),
		Body:       req.GetBody(),
	})
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return &notification_proto.BulkSendSmsResponse{SmsIds: smsIds}, nil
}

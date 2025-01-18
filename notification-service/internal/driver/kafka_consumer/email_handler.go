package kafkaconsumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/services"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	notification_proto "github.com/harmonify/movie-reservation-system/notification-service/internal/driven/proto/notification"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type EmailRouteParam struct {
	fx.In

	Logger               logger.Logger
	Tracer               tracer.Tracer
	EmailService         services.EmailService
	EmailTemplateService services.EmailTemplateService
}

type emailRouteImpl struct {
	logger               logger.Logger
	tracer               tracer.Tracer
	emailService         services.EmailService
	emailTemplateService services.EmailTemplateService
}

func NewEmailRoute(p EmailRouteParam) kafka.Route {
	return &emailRouteImpl{
		logger:               p.Logger,
		tracer:               p.Tracer,
		emailService:         p.EmailService,
		emailTemplateService: p.EmailTemplateService,
	}
}

func (r *emailRouteImpl) Match(topic string) bool {
	return topic == shared.EmailTopicV1_0_0
}

func (r *emailRouteImpl) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	val := &notification_proto.Email{}
	if err := proto.Unmarshal(message.Value, val); err != nil {
		return err
	}

	var tmplData *notification_proto.EmailVerificationTemplateData
	if err := anypb.UnmarshalTo(val.GetTemplateData(), tmplData, proto.UnmarshalOptions{}); err != nil {
		return err
	}

	content, err := r.emailTemplateService.Render(ctx, shared.EmailVerificationTemplatePath, tmplData)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to render verification link email template", zap.Error(err), zap.Strings("recipients", val.GetRecipients()), zap.String("subject", val.GetSubject()))
		return err
	}

	err = r.emailService.Send(ctx, shared.EmailMessage{
		Recipients: val.GetRecipients(),
		Subject:    val.GetSubject(),
		Body:       content,
	})
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to send verification link email", zap.Error(err), zap.Strings("recipients", val.GetRecipients()), zap.String("subject", val.GetSubject()))
		return err
	}

	return nil
}

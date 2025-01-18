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
)

type SmsRouteParam struct {
	fx.In

	Logger     logger.Logger
	Tracer     tracer.Tracer
	SmsService services.SmsService
}

type smsRouteImpl struct {
	logger     logger.Logger
	tracer     tracer.Tracer
	smsService services.SmsService
}

func NewSmsRoute(p SmsRouteParam) kafka.Route {
	return &smsRouteImpl{
		logger:     p.Logger,
		tracer:     p.Tracer,
		smsService: p.SmsService,
	}
}

func (r *smsRouteImpl) Match(topic string) bool {
	return topic == shared.SmsTopicV1_0_0
}

func (r *smsRouteImpl) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	val := &notification_proto.Sms{}
	if err := proto.Unmarshal(message.Value, val); err != nil {
		return err
	}

	err := r.smsService.Send(ctx, shared.SmsMessage{
		Recipient: val.GetRecipient(),
		Body:      val.GetBody(),
	})
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to send verification link sms", zap.Error(err), zap.String("recipient", val.GetRecipient()), zap.String("body", val.GetBody()))
		return err
	}

	return nil
}

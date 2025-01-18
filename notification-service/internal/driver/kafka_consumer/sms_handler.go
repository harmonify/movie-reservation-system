package kafka_consumer

import (
	"context"
	"reflect"

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
	listeners  []kafka.EventListener
	logger     logger.Logger
	tracer     tracer.Tracer
	smsService services.SmsService
}

func NewSmsRoute(p SmsRouteParam) kafka.Route {
	return &smsRouteImpl{
		listeners:  []kafka.EventListener{},
		logger:     p.Logger,
		tracer:     p.Tracer,
		smsService: p.SmsService,
	}
}

func (r *smsRouteImpl) Decode(ctx context.Context, value []byte) (interface{}, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	val := &notification_proto.Sms{}
	if err := proto.Unmarshal(value, val); err != nil {
		r.logger.WithCtx(ctx).Error("Failed to decode sms proto", zap.Error(err))
		return nil, kafka.ErrDecodeFailed
	}
	return val, nil
}

func (r *smsRouteImpl) Match(ctx context.Context, event *kafka.Event) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return event.Topic == shared.SmsTopicV1_0_0, nil
}

func (r *smsRouteImpl) Handle(ctx context.Context, event *kafka.Event) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Notify listeners
	for _, listener := range r.listeners {
		listener.OnEvent(ctx, event)
	}

	val, ok := event.Value.(*notification_proto.Sms)
	if !ok {
		r.logger.WithCtx(ctx).Error(
			"Failed to assert correct event type",
			zap.String("want", reflect.TypeFor[*notification_proto.Email]().Name()),
			zap.String("got", reflect.TypeOf(val).Name()),
		)
		return kafka.ErrInvalidValueType
	}

	err := r.smsService.Send(ctx, shared.SmsMessage{
		Recipient: val.GetRecipient(),
		Body:      val.GetBody(),
	})
	if err != nil {
		r.logger.WithCtx(ctx).Error(
			"Failed to send verification link sms",
			zap.Error(err),
			zap.String("recipient", val.GetRecipient()),
			zap.String("body", val.GetBody()),
		)
		return err
	}

	return nil
}

func (r *smsRouteImpl) AddEventListener(listener kafka.EventListener) {
	r.listeners = append(r.listeners, listener)
}

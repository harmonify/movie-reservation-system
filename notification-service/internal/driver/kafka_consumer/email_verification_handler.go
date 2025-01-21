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
	listeners            []kafka.EventListener
	logger               logger.Logger
	tracer               tracer.Tracer
	emailService         services.EmailService
	emailTemplateService services.EmailTemplateService
}

func NewEmailVerificationRoute(p EmailRouteParam) kafka.Route {
	return &emailRouteImpl{
		listeners:            []kafka.EventListener{},
		logger:               p.Logger,
		tracer:               p.Tracer,
		emailService:         p.EmailService,
		emailTemplateService: p.EmailTemplateService,
	}
}

func (r *emailRouteImpl) Decode(ctx context.Context, value []byte) (interface{}, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	val := &notification_proto.Email{}
	if err := proto.Unmarshal(value, val); err != nil {
		r.logger.WithCtx(ctx).Error("Failed to decode email proto", zap.Error(err))
		return nil, kafka.ErrDecodeFailed
	}
	return val, nil
}

func (r *emailRouteImpl) Match(ctx context.Context, event *kafka.Event) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if event.Topic != shared.EmailTopicV1_0_0 {
		return false, nil
	}

	val, ok := event.Value.(*notification_proto.Email)
	if !ok {
		r.logger.WithCtx(ctx).Error(
			"Failed to assert correct event type",
			zap.String("want", reflect.TypeFor[*notification_proto.Email]().Name()),
			zap.String("got", reflect.TypeOf(val).Name()),
		)
		return false, kafka.ErrInvalidValueType
	}

	return val.TemplateId == shared.EmailVerificationTemplateId.String(), nil
}

func (r *emailRouteImpl) Handle(ctx context.Context, event *kafka.Event) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Notify listeners
	for _, listener := range r.listeners {
		listener.OnEvent(ctx, event)
	}

	val, ok := event.Value.(*notification_proto.Email)
	if !ok {
		r.logger.WithCtx(ctx).Error(
			"Failed to assert correct event type",
			zap.String("want", reflect.TypeFor[*notification_proto.Email]().Name()),
			zap.String("got", reflect.TypeOf(val).Name()),
		)
		return kafka.ErrInvalidValueType
	}

	var tmplData *notification_proto.EmailVerificationTemplateData
	if err := anypb.UnmarshalTo(val.GetTemplateData(), tmplData, proto.UnmarshalOptions{}); err != nil {
		r.logger.WithCtx(ctx).Error("Failed to decode email template data proto", zap.Error(err))
		return err
	}

	content, err := r.emailTemplateService.Render(ctx, shared.EmailVerificationTemplatePath, tmplData)
	if err != nil {
		r.logger.WithCtx(ctx).Error(
			"Failed to render verification link email template",
			zap.Error(err),
			zap.Strings("recipients", val.GetRecipients()),
			zap.String("subject", val.GetSubject()),
		)
		return err
	}

	err = r.emailService.Send(ctx, shared.EmailMessage{
		Recipients: val.GetRecipients(),
		Subject:    val.GetSubject(),
		Body:       content,
	})
	if err != nil {
		r.logger.WithCtx(ctx).Error(
			"Failed to send verification link email",
			zap.Error(err),
			zap.Strings("recipients", val.GetRecipients()),
			zap.String("subject", val.GetSubject()),
		)
		return err
	}

	return nil
}

func (r *emailRouteImpl) AddEventListener(listener kafka.EventListener) {
	r.listeners = append(r.listeners, listener)
}

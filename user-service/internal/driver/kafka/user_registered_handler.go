package kafka_driver

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	watermill_pkg "github.com/harmonify/movie-reservation-system/pkg/kafka/watermill"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	user_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/user"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type UserRegisteredRouteParam struct {
	fx.In

	Logger     logger.Logger
	Tracer     tracer.Tracer
	OtpService otp_service.OtpService
}

type userRegisteredRouteImpl struct {
	listeners  []watermill_pkg.MessageListener
	logger     logger.Logger
	tracer     tracer.Tracer
	otpService otp_service.OtpService
}

func NewUserRegisteredRoute(p UserRegisteredRouteParam) watermill_pkg.Route {
	return &userRegisteredRouteImpl{
		listeners:  []watermill_pkg.MessageListener{},
		logger:     p.Logger,
		tracer:     p.Tracer,
		otpService: p.OtpService,
	}
}

func (r *userRegisteredRouteImpl) Identifier() string {
	return "user-registered-handler"
}

func (r *userRegisteredRouteImpl) Register(router *message.Router, subscriber message.Subscriber) error {
	router.AddNoPublisherHandler(
		r.Identifier(),
		shared.PublicUserRegisteredV1.String(),
		subscriber,
		r.Handle,
	)
	return nil
}

func (r *userRegisteredRouteImpl) Match(ctx context.Context, event *kafka.Event) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return event.Topic == shared.PublicUserRegisteredV1.String(), nil
}

func (r *userRegisteredRouteImpl) Handle(message *message.Message) error {
	ctx, span := r.tracer.StartSpanWithCaller(message.Context())
	defer span.End()

	// Notify listeners
	for _, listener := range r.listeners {
		listener.OnMessage(ctx, message)
	}

	val := &user_proto.UserRegistered{}
	if err := proto.Unmarshal(message.Payload, val); err != nil {
		r.logger.WithCtx(ctx).Error("Failed to decode UserRegistered proto", zap.Error(err))
		return kafka.ErrMalformedMessage
	}

	r.logger.WithCtx(ctx).Debug("Received UserRegistered event", zap.Any("user_registered", val))

	// Send email verification link
	err := r.otpService.SendEmailVerificationLink(ctx, otp_service.SendEmailVerificationLinkParam{
		Name:  fmt.Sprintf("%s %s", val.GetFirstName(), val.GetLastName()),
		Email: val.GetEmail(),
	})
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to send email verification link", zap.Error(err))
		return err
	}

	return nil
}

func (r *userRegisteredRouteImpl) AddListener(listener watermill_pkg.MessageListener) {
	r.listeners = append(r.listeners, listener)
}

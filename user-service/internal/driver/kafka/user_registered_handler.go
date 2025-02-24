package kafka_driver

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill/message"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	watermill_pkg "github.com/harmonify/movie-reservation-system/pkg/kafka/watermill"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	user_proto "github.com/harmonify/movie-reservation-system/pkg/proto/user"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type UserRegisteredRouteParam struct {
	fx.In

	Config     *config.UserServiceConfig
	Logger     logger.Logger
	Tracer     tracer.Tracer
	OtpService otp_service.OtpService
}

type userRegisteredRouteImpl struct {
	listeners  []watermill_pkg.MessageListener
	config     *config.UserServiceConfig
	logger     logger.Logger
	tracer     tracer.Tracer
	otpService otp_service.OtpService
}

func NewUserRegisteredRoute(p UserRegisteredRouteParam) watermill_pkg.Route {
	return &userRegisteredRouteImpl{
		listeners:  []watermill_pkg.MessageListener{},
		config:     p.Config,
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
		r.config.KafkaTopicUserRegisteredV1,
		subscriber,
		r.handle,
	)
	return nil
}

func (r *userRegisteredRouteImpl) handle(message *message.Message) error {
	ctx, span := r.tracer.StartSpanWithCaller(message.Context())
	defer span.End()

	r.logger.WithCtx(ctx).Debug("Received UserRegistered event")

	// Notify listeners
	for _, listener := range r.listeners {
		listener.OnMessage(ctx, message)
	}

	val := &user_proto.UserRegistered{}
	if err := proto.Unmarshal(message.Payload, val); err != nil {
		r.logger.WithCtx(ctx).Error("Failed to decode UserRegistered proto", zap.Error(err))
		return kafka.ErrMalformedMessage
	}

	r.logger.WithCtx(ctx).Debug("UserRegistered event payload", zap.Any("user_registered", val))

	err := r.otpService.SendSignupEmail(ctx, otp_service.SendSignupEmailParam{
		UUID:      val.GetUuid(),
		Email:     val.GetEmail(),
		FirstName: val.GetFirstName(),
		LastName:  val.GetLastName(),
	})
	var ed *error_pkg.ErrorWithDetails
	if err != nil {
		if errors.As(err, &ed) {
			if ed.Code == otp_service.OtpAlreadySentError.Code {
				r.logger.WithCtx(ctx).Info("Email verification link already sent. Skipping sending email verification link")
			} else {
				r.logger.WithCtx(ctx).Error("Failed to send email verification link", zap.Error(err), zap.Object("ed", ed))
				return err
			}
		} else {
			r.logger.WithCtx(ctx).Error("Failed to send email verification link", zap.Error(err))
			return err
		}
	}

	return nil
}

func (r *userRegisteredRouteImpl) AddListener(listener watermill_pkg.MessageListener) {
	r.listeners = append(r.listeners, listener)
}

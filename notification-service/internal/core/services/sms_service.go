package services

import (
	"context"
	"errors"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	SmsService interface {
		// Send sends sms to a recipient
		Send(ctx context.Context, msg shared.SmsMessage) (string, error)
		// BulkSend sends sms to many recipients
		BulkSend(ctx context.Context, msg shared.BulkSmsMessage) ([]string, error)
	}

	SmsServiceParam struct {
		fx.In

		SmsProvider shared.SmsProvider
		Logger      logger.Logger
		Tracer      tracer.Tracer
	}

	SmsServiceResult struct {
		fx.Out

		SmsService SmsService
	}

	smsServiceImpl struct {
		logger                      logger.Logger
		tracer                      tracer.Tracer
		smsProvider                 shared.SmsProvider
		smsProviderSendExecutor     failsafe.Executor[string]
		smsProviderBulkSendExecutor failsafe.Executor[[]string]
	}
)

func NewSmsService(p SmsServiceParam) SmsServiceResult {
	s := &smsServiceImpl{
		smsProvider:                 p.SmsProvider,
		logger:                      p.Logger,
		tracer:                      p.Tracer,
		smsProviderSendExecutor:     buildSendSmsExecutor[string](),
		smsProviderBulkSendExecutor: buildSendSmsExecutor[[]string](),
	}

	return SmsServiceResult{
		SmsService: s,
	}
}

func (s *smsServiceImpl) Send(ctx context.Context, message shared.SmsMessage) (string, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	s.logger.WithCtx(ctx).Debug("Send sms")

	smsId, err := s.smsProviderSendExecutor.WithContext(ctx).Get(func() (string, error) {
		return s.smsProvider.Send(ctx, message)
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Send sms failed", zap.Error(err), zap.Any("message", message))
		return "", err
	}

	s.logger.WithCtx(ctx).Info("Send sms success", zap.String("sms_id", smsId), zap.Any("message", message))
	return smsId, nil
}

func (s *smsServiceImpl) BulkSend(ctx context.Context, message shared.BulkSmsMessage) ([]string, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	s.logger.WithCtx(ctx).Debug("Bulk send sms")

	smsIds, err := s.smsProviderBulkSendExecutor.WithContext(ctx).Get(func() ([]string, error) {
		return s.smsProvider.BulkSend(ctx, message)
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Bulk send sms failed", zap.Error(err), zap.Any("message", message))
		return nil, err
	}

	s.logger.WithCtx(ctx).Info("Bulk send sms success", zap.Strings("sms_ids", smsIds), zap.Any("message", message))

	return smsIds, nil
}

func buildSendSmsExecutor[T any]() failsafe.Executor[T] {
	return failsafe.NewExecutor(
		retrypolicy.Builder[T]().
			WithBackoff(100*time.Millisecond, time.Second).
			WithJitterFactor(0.2).
			WithMaxRetries(3).
			Build(),
		circuitbreaker.Builder[T]().
			WithDelayFunc(func(exec failsafe.ExecutionAttempt[T]) time.Duration {
				err := exec.LastError()
				if err == nil {
					return 0
				}

				var ed *error_pkg.ErrorWithDetails
				if errors.As(err, &ed) {
					if ed.Code == error_pkg.BadGatewayError.Code {
						return 5 * time.Second
					} else if ed.Code == error_pkg.RateLimitExceededError.Code {
						data, ok := ed.Data.(*error_pkg.RateLimitExceededErrorData)
						if ok {
							return (time.Duration(data.RetryAfter) * time.Second) + (5 * time.Second)
						}
					}
				}

				return 30 * time.Second
			}).
			Build(),
		timeout.With[T](5*time.Second),
	)
}

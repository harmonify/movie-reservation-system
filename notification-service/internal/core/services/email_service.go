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
	EmailService interface {
		// Send sends email
		Send(ctx context.Context, msg shared.EmailMessage) (string, error)
	}

	EmailServiceParam struct {
		fx.In

		EmailProvider shared.EmailProvider
		Logger        logger.Logger
		Tracer        tracer.Tracer
	}

	EmailServiceResult struct {
		fx.Out

		EmailService EmailService
	}

	emailServiceImpl struct {
		logger                    logger.Logger
		tracer                    tracer.Tracer
		emailProvider             shared.EmailProvider
		emailProviderSendExecutor failsafe.Executor[string]
	}
)

func NewEmailService(p EmailServiceParam) EmailServiceResult {
	s := &emailServiceImpl{
		logger:        p.Logger,
		tracer:        p.Tracer,
		emailProvider: p.EmailProvider,
		emailProviderSendExecutor: failsafe.NewExecutor(
			retrypolicy.Builder[string]().
				WithBackoff(100*time.Millisecond, time.Second).
				WithJitterFactor(0.2).
				WithMaxRetries(3).
				Build(),
			circuitbreaker.Builder[string]().
				WithDelayFunc(func(exec failsafe.ExecutionAttempt[string]) time.Duration {
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
			timeout.With[string](5*time.Second),
		),
	}

	return EmailServiceResult{
		EmailService: s,
	}
}

func (s *emailServiceImpl) Send(ctx context.Context, message shared.EmailMessage) (string, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	s.logger.WithCtx(ctx).Debug("Send email")

	emailId, err := s.emailProviderSendExecutor.WithContext(ctx).Get(func() (string, error) {
		return s.emailProvider.Send(ctx, message)
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Send email failed", zap.Error(err), zap.Any("message", message))
		return "", err
	}

	s.logger.WithCtx(ctx).Info("Send email success", zap.String("email_id", emailId), zap.Any("message", message))
	return emailId, nil
}

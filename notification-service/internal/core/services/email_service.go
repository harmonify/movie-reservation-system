package services

import (
	"context"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
)

type (
	EmailService interface {
		// Send sends email
		Send(ctx context.Context, msg shared.EmailMessage) error
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
				HandleErrorTypes(&shared.RateLimitError{}).
				WithDelayFunc(func(exec failsafe.ExecutionAttempt[string]) time.Duration {
					err := exec.LastError()
					switch e := (err).(type) {
					case *shared.RateLimitError:
						return (time.Duration(e.RetryAfter) * time.Second) + (5 * time.Second)
					default:
						return 0
					}
				}).
				Build(),
			timeout.With[string](5*time.Second),
		),
	}

	return EmailServiceResult{
		EmailService: s,
	}
}

func (s *emailServiceImpl) Send(ctx context.Context, message shared.EmailMessage) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return s.emailProviderSendExecutor.Run(func() error {
		_, err := s.emailProvider.Send(ctx, message)
		return err
	})
}

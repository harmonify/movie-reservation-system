package services

import (
	"context"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	EmailService interface {
		// Send sends email
		Send(ctx context.Context, msg shared.EmailMessage) error
	}

	EmailServiceParam struct {
		fx.In

		MailProvider   shared.EmailProvider
		Logger         logger.Logger
		Tracer         tracer.Tracer
	}

	EmailServiceResult struct {
		fx.Out

		EmailService EmailService
	}

	emailServiceImpl struct {
		mailProvider shared.EmailProvider
		logger       logger.Logger
		tracer       tracer.Tracer
	}
)

func NewEmailService(p EmailServiceParam) EmailServiceResult {
	s := &emailServiceImpl{
		mailProvider: p.MailProvider,
		logger:       p.Logger,
		tracer:       p.Tracer,
	}

	return EmailServiceResult{
		EmailService: s,
	}
}

func (s *emailServiceImpl) Send(ctx context.Context, message shared.EmailMessage) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	emailMsg, emailId, err := s.mailProvider.Send(ctx, message)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send email", zap.Error(err), zap.String("email_message", emailMsg), zap.String("email_id", emailId), zap.Any("message", message))
		return err
	}

	return nil
}

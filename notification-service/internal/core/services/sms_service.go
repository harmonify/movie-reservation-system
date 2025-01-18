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
	SmsService interface {
		// Send sends sms
		Send(ctx context.Context, msg shared.SmsMessage) error
		// Send sends sms to many recipients
		BulkSend(ctx context.Context, msg shared.BulkSmsMessage) error
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
		smsProvider shared.SmsProvider
		logger      logger.Logger
		tracer      tracer.Tracer
	}
)

func NewSmsService(p SmsServiceParam) SmsServiceResult {
	s := &smsServiceImpl{
		smsProvider: p.SmsProvider,
		logger:      p.Logger,
		tracer:      p.Tracer,
	}

	return SmsServiceResult{
		SmsService: s,
	}
}

func (s *smsServiceImpl) Send(ctx context.Context, message shared.SmsMessage) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	smsId, err := s.smsProvider.Send(ctx, message)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send sms", zap.Error(err), zap.String("sms_id", smsId), zap.Any("message", message))
		return err
	}

	return nil
}

func (s *smsServiceImpl) BulkSend(ctx context.Context, message shared.BulkSmsMessage) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	smsIds, err := s.smsProvider.BulkSend(ctx, message)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send sms", zap.Error(err), zap.Strings("sms_ids", smsIds), zap.Any("message", message))
		return err
	}

	return nil
}

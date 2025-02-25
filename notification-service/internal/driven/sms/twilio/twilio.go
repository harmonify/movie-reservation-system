package twilio

import (
	"context"
	"errors"
	"fmt"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	SmsProviderParam struct {
		fx.In

		Config *config.NotificationServiceConfig
		Logger logger.Logger
		Tracer tracer.Tracer
		Util   *util.Util
	}

	SmsProviderResult struct {
		fx.Out

		SmsProvider shared.SmsProvider
	}

	twilioSmsProviderImpl struct {
		client *twilio.RestClient
		config *config.NotificationServiceConfig
		logger logger.Logger
		tracer tracer.Tracer
		util   *util.Util
	}
)

func NewTwilioSmsProvider(p SmsProviderParam) (SmsProviderResult, error) {
	if p.Config.TwilioAccountSid == "" {
		return SmsProviderResult{}, fmt.Errorf("TwilioAccountSid is empty")
	}
	if p.Config.TwilioAuthToken == "" {
		return SmsProviderResult{}, fmt.Errorf("TwilioAuthToken is empty")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: p.Config.TwilioAccountSid,
		Password: p.Config.TwilioAuthToken,
	})

	return SmsProviderResult{
		SmsProvider: &twilioSmsProviderImpl{
			client: client,
			config: p.Config,
			logger: p.Logger,
			tracer: p.Tracer,
			util:   p.Util,
		},
	}, nil
}

func (s *twilioSmsProviderImpl) Send(ctx context.Context, message shared.SmsMessage) (id string, err error) {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if message.Recipient == "" {
		return "", shared.NewInvalidPhoneNumberError(message.Recipient)
	}
	formattedPhoneNumber := s.util.FormatterUtil.FormatPhoneNumberToE164(message.Recipient, s.config.AppDefaultCountryDialCode)

	// Create Twilio message
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(formattedPhoneNumber)
	params.SetFrom(s.config.AppName)
	params.SetBody(message.Body)

	res, err := s.client.Api.CreateMessage(params)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send SMS message", zap.Error(err), zap.Any("params", params))
		return "", fmt.Errorf("error sending SMS message: %v", err)
	}

	s.logger.WithCtx(ctx).Info("Successfully send a message", zap.Any("response", res), zap.Any("params", params))
	return *res.Sid, nil

}

func (s *twilioSmsProviderImpl) BulkSend(ctx context.Context, message shared.BulkSmsMessage) ([]string, error) {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	message_ids := make([]string, 0, len(message.Recipients))
	var finalErr error

	for _, recipient := range message.Recipients {
		if recipient == "" {
			finalErr = errors.Join(finalErr, shared.NewInvalidPhoneNumberError(recipient))
			continue
		}

		formattedPhoneNumber := s.util.FormatterUtil.FormatPhoneNumberToE164(recipient, s.config.AppDefaultCountryDialCode)

		// Create Twilio message
		params := &twilioApi.CreateMessageParams{}
		params.SetTo(formattedPhoneNumber)
		params.SetFrom(s.config.AppName)
		params.SetBody(message.Body)

		resp, err := s.client.Api.CreateMessage(params)
		if err != nil {
			s.logger.WithCtx(ctx).Error("Failed to send SMS message", zap.Error(err), zap.Any("params", params))
			finalErr = errors.Join(finalErr, fmt.Errorf("error sending SMS message: %v", err))
			continue
		}

		s.logger.WithCtx(ctx).Info("Successfully send a message", zap.Any("response", resp), zap.Any("params", params))
		message_ids = append(message_ids, *resp.Sid)
	}

	return message_ids, nil
}

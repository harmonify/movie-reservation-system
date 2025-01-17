package messaging

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

type twilioMessagerImpl struct {
	client *twilio.RestClient
	config *config.Config
	logger logger.Logger
	tracer tracer.Tracer
	util   *util.Util
}

func NewTwilioMessager(p MessagerParam) (MessagerResult, error) {
	if p.Config.TwilioAccountSid == "" {
		return MessagerResult{}, fmt.Errorf("TwilioAccountSid is empty")
	}
	if p.Config.TwilioAuthToken == "" {
		return MessagerResult{}, fmt.Errorf("TwilioAuthToken is empty")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: p.Config.TwilioAccountSid,
		Password: p.Config.TwilioAuthToken,
	})

	return MessagerResult{
		Messager: &twilioMessagerImpl{
			client: client,
			config: p.Config,
			logger: p.Logger,
			tracer: p.Tracer,
			util:   p.Util,
		},
	}, nil
}

func (s *twilioMessagerImpl) Send(ctx context.Context, message Message) (id string, err error) {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if message.To == "" {
		return "", ErrInvalidPhoneNumber
	}
	formattedPhoneNumber := s.util.FormatterUtil.FormatPhoneNumberToE164(message.To, s.config.AppDefaultCountryDialCode)

	// Create Twilio message
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(formattedPhoneNumber)
	params.SetFrom(s.config.AppName)
	params.SetBody(message.Body)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		s.logger.Error("Failed to send SMS message", zap.Error(err), zap.Any("params", params))
		return "", fmt.Errorf("error sending SMS message: %v", err)
	}

	s.logger.Info("Successfully send a message", zap.Any("response", resp), zap.Any("params", params))
	return *resp.Sid, nil

}

func (s *twilioMessagerImpl) SendOTP(ctx context.Context, otpCode string, phoneNumber string) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if phoneNumber == "" {
		return ErrInvalidPhoneNumber
	}

	//Format phone number to E.164(the international telephone numbering)
	formattedPhoneNumber := s.util.FormatterUtil.FormatPhoneNumberToE164(phoneNumber, s.config.AppDefaultCountryDialCode)

	// Twilio message body with the OTP
	body := fmt.Sprintf("Your verification code for %s account is: %s", s.config.AppName, otpCode)

	// Create Twilio message
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(formattedPhoneNumber)
	params.SetFrom(s.config.AppName)
	params.SetBody(body)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		s.logger.Error("Failed to send SMS message", zap.Error(err), zap.Any("params", params))
		return fmt.Errorf("error sending OTP: %v", err)
	}

	s.logger.Info("Successfully send a message", zap.Any("response", resp), zap.Any("params", params))
	return nil
}

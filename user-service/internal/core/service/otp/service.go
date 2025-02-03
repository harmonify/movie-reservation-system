package otp_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

type (
	OtpService interface {
		SendEmailVerificationLink(ctx context.Context, p SendEmailVerificationLinkParam) error
		VerifyEmail(ctx context.Context, p VerifyEmailParam) error
		SendPhoneOtp(ctx context.Context, p SendPhoneOtpParam) error
		VerifyPhoneOtp(ctx context.Context, p VerifyPhoneOtpParam) error
	}

	OtpServiceParam struct {
		fx.In

		Config               *config.UserServiceConfig
		Logger               logger.Logger
		Tracer               tracer.Tracer
		NotificationProvider shared.NotificationProvider
		OtpStorage           shared.OtpStorage
		Util                 *util.Util
	}

	OtpServiceResult struct {
		fx.Out

		OtpService OtpService
	}

	otpServiceImpl struct {
		config               *config.UserServiceConfig
		logger               logger.Logger
		tracer               tracer.Tracer
		notificationProvider shared.NotificationProvider
		otpStorage           shared.OtpStorage
		util                 *util.Util

		EmailVerificationLinkTTL uint // in seconds
		PhoneOtpTTL              uint // in seconds
	}
)

func NewOtpService(p OtpServiceParam) OtpServiceResult {
	return OtpServiceResult{
		OtpService: &otpServiceImpl{
			config:               p.Config,
			logger:               p.Logger,
			tracer:               p.Tracer,
			notificationProvider: p.NotificationProvider,
			otpStorage:           p.OtpStorage,
			util:                 p.Util,

			EmailVerificationLinkTTL: 24 * 60 * 60, // 24 hours
			PhoneOtpTTL:              15 * 60,      // 15 minutes
		},
	}
}

// Construct email verification link to frontend app.
// Note: Frontend app should handle the parameters to
// make request to POST /profile/email/verify
func (s *otpServiceImpl) constructEmailVerificationLink(email, token string) string {
	return fmt.Sprintf("%s/profile/email/verify?email=%s&token=%s", s.config.FrontEndUrl, email, token)
}

func (s *otpServiceImpl) SendEmailVerificationLink(ctx context.Context, p SendEmailVerificationLinkParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	savedToken, err := s.otpStorage.GetEmailVerificationToken(ctx, p.Email)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) && ed.Code != VerificationTokenNotFoundError.Code {
		s.logger.WithCtx(ctx).Error("Failed to get existing verification token", zap.Error(err))
		return error_pkg.InternalServerError
	}
	if savedToken != "" {
		return OtpAlreadySentError
	}

	token, err := s.util.GeneratorUtil.GenerateRandomHex(32)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate verification token", zap.Error(err))
		return error_pkg.InternalServerError
	}

	templateData, err := anypb.New(&notification_proto.EmailVerificationTemplateData{
		Name:             p.Name,
		VerificationLink: s.constructEmailVerificationLink(p.Email, token),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to marshal email template data into protobuf", zap.Error(err))
		return error_pkg.InternalServerError
	}

	err = s.notificationProvider.SendEmail(ctx, &notification_proto.SendEmailRequest{
		Recipients:   []string{p.Email},
		Subject:      "Account verification",
		TemplateId:   shared.EmailVerificationTemplateId.String(),
		TemplateData: templateData,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send email", zap.Error(err))
		return SendVerificationLinkFailedError
	}

	err = s.otpStorage.SaveEmailVerificationToken(ctx, shared.SaveEmailVerificationTokenParam{
		Email: p.Email,
		Token: token,
		TTL:   time.Second * time.Duration(s.EmailVerificationLinkTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save email verification token", zap.Error(err))
		return error_pkg.InternalServerError
	}

	return nil
}

func (s *otpServiceImpl) VerifyEmail(ctx context.Context, p VerifyEmailParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	token, err := s.otpStorage.GetEmailVerificationToken(ctx, p.Email)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) {
		if ed.Code == VerificationTokenNotFoundError.Code {
			return VerificationTokenNotFoundError
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing verification token", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}

	if token == "" {
		return VerificationTokenNotFoundError
	}

	if token != p.Token {
		return VerificationTokenInvalidError
	}

	return nil
}

func (s *otpServiceImpl) SendPhoneOtp(ctx context.Context, p SendPhoneOtpParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	savedOtp, err := s.otpStorage.GetPhoneOtp(ctx, p.PhoneNumber)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) && ed.Code != VerificationTokenNotFoundError.Code {
		s.logger.WithCtx(ctx).Error("Failed to get existing OTP", zap.Error(err))
		return error_pkg.InternalServerError
	}
	if savedOtp != "" {
		return OtpAlreadySentError
	}

	otp, err := s.util.GeneratorUtil.GenerateRandomNumber(6)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate OTP", zap.Error(err))
		return error_pkg.InternalServerError
	}

	err = s.notificationProvider.SendSms(ctx, &notification_proto.SendSmsRequest{
		Recipient: p.PhoneNumber,
		Body:      fmt.Sprintf("Your verification code for %s is %s", s.config.AppName, otp),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send OTP", zap.Error(err))
		return SendOtpFailedError
	}

	err = s.otpStorage.SavePhoneOtp(ctx, shared.SavePhoneOtpParam{
		PhoneNumber: p.PhoneNumber,
		Otp:         otp,
		TTL:         time.Second * time.Duration(s.PhoneOtpTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save OTP", zap.Error(err))
		return error_pkg.InternalServerError
	}

	return nil
}

func (s *otpServiceImpl) VerifyPhoneOtp(ctx context.Context, p VerifyPhoneOtpParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	attempt, err := s.otpStorage.GetPhoneOtpAttempt(ctx, p.PhoneNumber)
	if err != nil {
		var ed *error_pkg.ErrorWithDetails
		if errors.As(err, &ed) && ed.Code == VerificationTokenNotFoundError.Code {
			return ed
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing phone OTP attempt", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}
	if attempt >= 3 {
		s.logger.Info("User attempted to verify phone OTP too many times", zap.Int("attempt", attempt))
		return OtpTooManyAttemptError
	}

	otp, err := s.otpStorage.GetPhoneOtp(ctx, p.PhoneNumber)
	if err != nil {
		var ed *error_pkg.ErrorWithDetails
		if errors.As(err, &ed) && ed.Code == VerificationTokenNotFoundError.Code {
			return ed
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing phone OTP", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}

	err = s.otpStorage.IncrementPhoneOtpAttempt(ctx, p.PhoneNumber)
	if err != nil {
		s.logger.Error("Failed to increment user phone OTP attempt", zap.Error(err))
		return error_pkg.InternalServerError
	}

	if p.Otp != otp {
		return OtpInvalidError
	}

	return nil
}

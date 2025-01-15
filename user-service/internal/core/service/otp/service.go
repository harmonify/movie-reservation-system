package otp_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/mail"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"

	"github.com/harmonify/movie-reservation-system/pkg/messaging"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

		Config     *config.Config
		Logger     logger.Logger
		Tracer     tracer.Tracer
		Mailer     mail.Mailer
		Messager   messaging.Messager
		OtpStorage shared_service.OtpStorage
		Util       *util.Util
	}

	OtpServiceResult struct {
		fx.Out

		OtpService OtpService
	}

	otpServiceImpl struct {
		config     *config.Config
		logger     logger.Logger
		tracer     tracer.Tracer
		mailer     mail.Mailer
		messager   messaging.Messager
		otpStorage shared_service.OtpStorage
		util       *util.Util

		EmailVerificationLinkTTL uint // in seconds
		PhoneOtpTTL              uint // in seconds
	}
)

func NewOtpService(p OtpServiceParam) OtpServiceResult {
	return OtpServiceResult{
		OtpService: &otpServiceImpl{
			config:     p.Config,
			logger:     p.Logger,
			tracer:     p.Tracer,
			mailer:     p.Mailer,
			messager:   p.Messager,
			otpStorage: p.OtpStorage,
			util:       p.Util,

			EmailVerificationLinkTTL: 24 * 60 * 60, // 24 hours
			PhoneOtpTTL:              15 * 60,      // 15 minutes
		},
	}
}

// Construct email verification link to frontend app.
// Note: Frontend app should handle the parameters to
// make request to POST /email/verify
func (s *otpServiceImpl) constructEmailVerificationLink(email, token string) string {
	return fmt.Sprintf("https://%s/email/verify?email=%s&token=%s", s.config.FrontEndUrl, email, token)
}

func (s *otpServiceImpl) SendEmailVerificationLink(ctx context.Context, p SendEmailVerificationLinkParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	savedToken, err := s.otpStorage.GetEmailVerificationToken(ctx, p.Email)
	if err != nil && !errors.Is(err, error_constant.ErrNotFound) {
		s.logger.WithCtx(ctx).Error("Failed to get existing verification token", zap.Error(err))
		return error_constant.ErrInternalServerError
	}
	if savedToken != "" {
		return ErrOtpAlreadyExist
	}

	token, err := s.util.GeneratorUtil.GenerateRandomBase64(32)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate verification token", zap.Error(err))
		return error_constant.ErrInternalServerError
	}

	_, err = s.mailer.Send(ctx, mail.Message{
		To:      []string{p.Email},
		Subject: "Account verification",
		Body: fmt.Sprintf(`<!DOCTYPE html>
<html>
<body>
  <p>Howdy,</p>
  <p>Thank you for choosing %s! Please confirm your email address by clicking the link below. We'll communicate important updates with you from time to time via email, so it's essential that we have an up-to-date email address on file.</p>
  <p><a href="%s" style="color: #1a73e8; text-decoration: none;">Verify your email address</a></p>
  <p>If you did not sign up for a %s account, you can simply disregard this email.</p>
  <p>Hope you have a nice day!</p>
  <p>The %s Team</p>
</body>
</html>`,
			s.config.AppName,
			s.constructEmailVerificationLink(p.Email, token),
			s.config.AppName,
			s.config.AppName,
		),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send email", zap.Error(err))
		return ErrSendVerificationLinkFailed
	}

	err = s.otpStorage.SaveEmailVerificationToken(ctx, shared_service.SaveEmailVerificationTokenParam{
		Email: p.Email,
		Token: token,
		TTL:   time.Second * time.Duration(s.EmailVerificationLinkTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save email verification token", zap.Error(err))
		return error_constant.ErrInternalServerError
	}

	return nil
}

func (s *otpServiceImpl) VerifyEmail(ctx context.Context, p VerifyEmailParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	token, err := s.otpStorage.GetEmailVerificationToken(ctx, p.Email)
	if err != nil {
		if errors.Is(err, error_constant.ErrNotFound) {
			return ErrVerificationTokenNotFound
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing verification token", zap.Error(err))
			return error_constant.ErrInternalServerError
		}
	}

	if token == "" {
		return ErrVerificationTokenNotFound
	}

	if token != p.Token {
		return ErrVerificationTokenInvalid
	}

	return nil
}

func (s *otpServiceImpl) SendPhoneOtp(ctx context.Context, p SendPhoneOtpParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	savedOtp, err := s.otpStorage.GetPhoneOtp(ctx, p.PhoneNumber)
	if err != nil && !errors.Is(err, error_constant.ErrNotFound) {
		s.logger.WithCtx(ctx).Error("Failed to get existing OTP", zap.Error(err))
		return error_constant.ErrInternalServerError
	}
	if savedOtp != "" {
		return ErrOtpAlreadyExist
	}

	otp, err := s.util.GeneratorUtil.GenerateRandomNumber(6)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate OTP", zap.Error(err))
		return error_constant.ErrInternalServerError
	}

	_, err = s.messager.Send(ctx, messaging.Message{
		To:   p.PhoneNumber,
		Body: fmt.Sprintf("Your verification code for %s is %s", s.config.AppName, otp),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send OTP", zap.Error(err))
		return ErrSendOtpFailed
	}

	err = s.otpStorage.SavePhoneOtp(ctx, shared_service.SavePhoneOtpParam{
		PhoneNumber: p.PhoneNumber,
		Otp:         otp,
		TTL:         time.Second * time.Duration(s.PhoneOtpTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save OTP", zap.Error(err))
		return error_constant.ErrInternalServerError
	}

	return nil
}

func (s *otpServiceImpl) VerifyPhoneOtp(ctx context.Context, p VerifyPhoneOtpParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	attempt, err := s.otpStorage.GetPhoneOtpAttempt(ctx, p.PhoneNumber)
	if err != nil {
		if errors.Is(err, error_constant.ErrNotFound) {
			return ErrOtpNotFound
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get user existing phone OTP attempt", zap.Error(err))
			return error_constant.ErrInternalServerError
		}
	}
	if attempt >= 3 {
		s.logger.Info("User attempted to verify phone OTP too many times", zap.Int("attempt", attempt))
		return ErrOtpTooManyAttempt
	}

	otp, err := s.otpStorage.GetPhoneOtp(ctx, p.PhoneNumber)
	if err != nil {
		if errors.Is(err, error_constant.ErrNotFound) {
			return ErrOtpNotFound
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing phone OTP", zap.Error(err))
			return error_constant.ErrInternalServerError
		}
	}

	err = s.otpStorage.IncrementPhoneOtpAttempt(ctx, p.PhoneNumber)
	if err != nil {
		s.logger.Error("Failed to increment user phone OTP attempt", zap.Error(err))
		return error_constant.ErrInternalServerError
	}

	if p.Otp != otp {
		return ErrOtpInvalid
	}

	return nil
}

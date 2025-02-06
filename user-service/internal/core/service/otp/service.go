package otp_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

type (
	OtpService interface {
		SendSignupEmail(ctx context.Context, p SendSignupEmailParam) error
		SendVerificationEmail(ctx context.Context, p SendVerificationEmailParam) error
		VerifyEmail(ctx context.Context, p VerifyEmailParam) error
		SendPhoneNumberVerificationOtp(ctx context.Context, p SendPhoneNumberVerificationOtpParam) error
		VerifyPhoneNumber(ctx context.Context, p VerifyPhoneNumberParam) error
	}

	OtpServiceParam struct {
		fx.In

		Config               *config.UserServiceConfig
		Logger               logger.Logger
		Tracer               tracer.Tracer
		Util                 *util.Util
		OtpCache             shared.OtpCache
		UserStorage          shared.UserStorage
		NotificationProvider shared.NotificationProvider
	}

	OtpServiceResult struct {
		fx.Out

		OtpService OtpService
	}

	otpServiceImpl struct {
		config               *config.UserServiceConfig
		logger               logger.Logger
		tracer               tracer.Tracer
		util                 *util.Util
		otpCache             shared.OtpCache
		userStorage          shared.UserStorage
		notificationProvider shared.NotificationProvider

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
			util:                 p.Util,
			otpCache:             p.OtpCache,
			userStorage:          p.UserStorage,
			notificationProvider: p.NotificationProvider,

			EmailVerificationLinkTTL: 24 * 60 * 60, // 24 hours
			PhoneOtpTTL:              15 * 60,      // 15 minutes
		},
	}
}

// Construct email verification link to frontend app.
// Note: Frontend app should handle the parameters to
// make request to POST /profile/email/verify
func (s *otpServiceImpl) constructEmailVerificationLink(email, code string) string {
	return fmt.Sprintf("%s/profile/email/verify?email=%s&code=%s", s.config.FrontEndUrl, email, code)
}

func (s *otpServiceImpl) SendSignupEmail(ctx context.Context, p SendSignupEmailParam) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	savedToken, err := s.otpCache.GetEmailVerificationCode(ctx, p.Email)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) && ed.Code != VerificationTokenNotFoundError.Code {
		s.logger.WithCtx(ctx).Error("Failed to get existing verification code", zap.Error(err))
		return error_pkg.InternalServerError
	}
	if savedToken != "" {
		return nil
	}

	code, err := s.util.GeneratorUtil.GenerateRandomHex(32)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate verification code", zap.Error(err))
		return error_pkg.InternalServerError
	}

	templateData, err := anypb.New(&notification_proto.SignupEmailVerificationTemplateData{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Url:       s.constructEmailVerificationLink(p.Email, code),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to marshal email template data into protobuf", zap.Error(err))
		return error_pkg.InternalServerError
	}

	err = s.notificationProvider.SendEmail(ctx, &notification_proto.SendEmailRequest{
		Recipients:   []string{p.Email},
		Subject:      "Welcome to " + s.config.AppName,
		TemplateId:   shared.SignupEmailTemplateId.String(),
		TemplateData: templateData,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send email", zap.Error(err))
		return SendVerificationLinkFailedError
	}

	err = s.otpCache.SaveEmailVerificationCode(ctx, shared.SaveEmailVerificationCodeParam{
		Email: p.Email,
		Code:  code,
		TTL:   time.Second * time.Duration(s.EmailVerificationLinkTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save email verification code", zap.Error(err))
		return error_pkg.InternalServerError
	}

	return nil
}

func (s *otpServiceImpl) constructUpdateEmailVerificationLink(email, code string) string {
	return fmt.Sprintf("%s/profile/email/update/verify?email=%s&code=%s", s.config.FrontEndUrl, email, code)
}

func (s *otpServiceImpl) SendVerificationEmail(ctx context.Context, p SendVerificationEmailParam) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	user, err := s.userStorage.FindUser(ctx, entity.FindUser{UUID: sql.NullString{String: p.UUID, Valid: true}})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to find user", zap.Error(err))
		return error_pkg.NotFoundError
	}

	savedToken, err := s.otpCache.GetEmailVerificationCode(ctx, user.Email)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) && ed.Code != VerificationTokenNotFoundError.Code {
		s.logger.WithCtx(ctx).Error("Failed to get existing verification code", zap.Error(err))
		return error_pkg.InternalServerError
	}
	if savedToken != "" {
		return VerificationLinkAlreadyExistError
	}

	code, err := s.util.GeneratorUtil.GenerateRandomHex(32)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate verification code", zap.Error(err))
		return error_pkg.InternalServerError
	}

	templateData, err := anypb.New(&notification_proto.SignupEmailVerificationTemplateData{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Url:       s.constructEmailVerificationLink(user.Email, code),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to marshal email template data into protobuf", zap.Error(err))
		return error_pkg.InternalServerError
	}

	err = s.notificationProvider.SendEmail(ctx, &notification_proto.SendEmailRequest{
		Recipients:   []string{user.Email},
		Subject:      "Email Verification",
		TemplateId:   shared.VerificationEmailTemplateId.String(),
		TemplateData: templateData,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send email", zap.Error(err))
		return SendVerificationLinkFailedError
	}

	err = s.otpCache.SaveEmailVerificationCode(ctx, shared.SaveEmailVerificationCodeParam{
		Email: user.Email,
		Code:  code,
		TTL:   time.Second * time.Duration(s.EmailVerificationLinkTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save email verification code", zap.Error(err))
		return error_pkg.InternalServerError
	}

	return nil
}

func (s *otpServiceImpl) VerifyEmail(ctx context.Context, p VerifyEmailParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	user, err := s.userStorage.FindUser(ctx, entity.FindUser{
		UUID: sql.NullString{String: p.UUID, Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to find user", zap.Error(err))
		return err
	}

	code, err := s.otpCache.GetEmailVerificationCode(ctx, user.Email)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) {
		if ed.Code == VerificationTokenNotFoundError.Code {
			return VerificationTokenNotFoundError
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing verification code", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}

	if code == "" {
		return VerificationTokenNotFoundError
	}

	if code != p.Code {
		return VerificationTokenInvalidError
	}

	go func() {
		_, err = s.otpCache.DeleteEmailVerificationCode(ctx, user.Email)
		if err != nil {
			s.logger.WithCtx(ctx).Warn("Failed to delete user email verification code", zap.Error(err))
		}
	}()

	_, err = s.userStorage.UpdateUser(
		ctx,
		entity.FindUser{
			UUID:  sql.NullString{String: p.UUID, Valid: true},
			Email: sql.NullString{String: user.Email, Valid: true},
		},
		entity.UpdateUser{
			IsEmailVerified: sql.NullBool{Bool: true, Valid: true},
		},
	)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to update user email verification status", zap.Error(err))
	}

	return err
}

func (s *otpServiceImpl) SendPhoneNumberVerificationOtp(ctx context.Context, p SendPhoneNumberVerificationOtpParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var user *entity.UserPhoneNumber
	err := s.userStorage.FindUserWithResult(ctx, entity.FindUser{UUID: sql.NullString{String: p.UUID, Valid: true}}, &user)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to find user", zap.Error(err))
		return error_pkg.NotFoundError
	}

	savedOtp, err := s.otpCache.GetPhoneNumberVerificationOtp(ctx, user.PhoneNumber)
	var ed *error_pkg.ErrorWithDetails
	if err != nil && errors.As(err, &ed) && ed.Code != OtpNotFoundError.Code {
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
		Recipient: user.PhoneNumber,
		Body:      fmt.Sprintf("Your verification code for %s is %s", s.config.AppName, otp),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send OTP", zap.Error(err))
		return SendOtpFailedError
	}

	err = s.otpCache.SavePhoneNumberVerificationOtp(ctx, shared.SavePhoneNumberVerificationOtpParam{
		PhoneNumber: user.PhoneNumber,
		Otp:         otp,
		TTL:         time.Second * time.Duration(s.PhoneOtpTTL),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save OTP", zap.Error(err))
		return error_pkg.InternalServerError
	}

	return nil
}

func (s *otpServiceImpl) VerifyPhoneNumber(ctx context.Context, p VerifyPhoneNumberParam) error {
	_, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var user *entity.UserPhoneNumber
	err := s.userStorage.FindUserWithResult(
		ctx,
		entity.FindUser{
			UUID: sql.NullString{String: p.UUID, Valid: true},
		},
		&user,
	)
	if err != nil {
		return err
	}

	attempt, err := s.otpCache.GetPhoneNumberVerificationAttempt(ctx, user.PhoneNumber)
	if err != nil {
		var ed *error_pkg.ErrorWithDetails
		if errors.As(err, &ed) {
			if ed.Code != OtpNotFoundError.Code {
				s.logger.WithCtx(ctx).Error("Failed to get existing phone OTP attempt", zap.Object("error", ed))
				return ed
			}
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing phone OTP attempt", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}
	if attempt >= 3 {
		s.logger.Info("User attempted to verify phone OTP too many times", zap.Int("attempt", attempt))
		return OtpTooManyAttemptError
	}

	otp, err := s.otpCache.GetPhoneNumberVerificationOtp(ctx, user.PhoneNumber)
	if err != nil {
		var ed *error_pkg.ErrorWithDetails
		if errors.As(err, &ed) && ed.Code == OtpNotFoundError.Code {
			return ed
		} else {
			s.logger.WithCtx(ctx).Error("Failed to get existing phone OTP", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}

	err = s.otpCache.IncrementPhoneNumberVerificationAttempt(ctx, user.PhoneNumber)
	if err != nil {
		s.logger.Error("Failed to increment user phone OTP attempt", zap.Error(err))
		return error_pkg.InternalServerError
	}

	if otp == "" {
		return OtpNotFoundError
	}

	if p.Otp != otp {
		return OtpInvalidError
	}

	go func() {
		_, err = s.otpCache.DeletePhoneNumberVerificationAttempt(ctx, user.PhoneNumber)
		if err != nil {
			s.logger.WithCtx(ctx).Warn("Failed to delete user phone OTP attempt", zap.Error(err))
		}

		_, err = s.otpCache.DeletePhoneNumberVerificationOtp(ctx, user.PhoneNumber)
		if err != nil {
			s.logger.WithCtx(ctx).Warn("Failed to delete user phone OTP", zap.Error(err))
		}
	}()

	_, err = s.userStorage.UpdateUser(
		ctx,
		entity.FindUser{
			UUID:        sql.NullString{String: p.UUID, Valid: true},
			PhoneNumber: sql.NullString{String: user.PhoneNumber, Valid: true},
		},
		entity.UpdateUser{
			IsPhoneNumberVerified: sql.NullBool{Bool: true, Valid: true},
		},
	)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to update user phone number verification status", zap.Error(err))
	}

	return err
}

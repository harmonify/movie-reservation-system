package otp_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		OtpCacheV2           shared.OtpCacheV2
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
		otpCacheV2           shared.OtpCacheV2
		userStorage          shared.UserStorage
		notificationProvider shared.NotificationProvider
	}
)

func NewOtpService(p OtpServiceParam) OtpServiceResult {
	return OtpServiceResult{
		OtpService: &otpServiceImpl{
			config:               p.Config,
			logger:               p.Logger,
			tracer:               p.Tracer,
			util:                 p.Util,
			otpCacheV2:           p.OtpCacheV2,
			userStorage:          p.UserStorage,
			notificationProvider: p.NotificationProvider,
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

	otp, err := s.otpCacheV2.GetOtp(ctx, p.UUID, shared.EmailVerificationOtpType)
	var ed *error_pkg.ErrorWithDetails
	if err != nil {
		if errors.As(err, &ed) {
			if ed.Code != OtpNotFoundError.Code {
				s.logger.WithCtx(ctx).Error("failed to get existing verification code", zap.Object("error", ed))
				return ed
			}
		} else {
			s.logger.WithCtx(ctx).Error("failed to get existing verification code", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}
	if otp != nil {
		s.logger.WithCtx(ctx).Info("Email verification link already sent. Skipping sending email verification link")
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

	err = s.otpCacheV2.SaveOtp(ctx, p.UUID, shared.EmailVerificationOtpType, code)
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
		s.logger.WithCtx(ctx).Error("failed to find user", zap.Error(err))
		return error_pkg.NotFoundError
	}

	otp, err := s.otpCacheV2.GetOtp(ctx, user.UUID, shared.EmailVerificationOtpType)
	var ed *error_pkg.ErrorWithDetails
	if err != nil {
		if errors.As(err, &ed) {
			if ed.Code != OtpNotFoundError.Code {
				s.logger.WithCtx(ctx).Error("failed to get existing verification code", zap.Object("error", ed))
				return ed
			}
		} else {
			s.logger.WithCtx(ctx).Error("failed to get existing verification code", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}
	if otp != nil {
		s.logger.WithCtx(ctx).Info("verification email code already sent. Skipping sending verification email code")
		return VerificationLinkAlreadySentError
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

	err = s.otpCacheV2.SaveOtp(ctx, user.UUID, shared.EmailVerificationOtpType, code)
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

	otp, err := s.otpCacheV2.IncrementOtpAttempt(ctx, p.UUID, shared.EmailVerificationOtpType)
	if err != nil {
		var ed *error_pkg.ErrorWithDetails
		if errors.As(err, &ed) {
			s.logger.WithCtx(ctx).Info("failed to increment email verification attempt", zap.Object("error", ed))
			if ed.Code == OtpNotFoundError.Code {
				return VerificationTokenNotFoundError
			} else {
				return ed
			}
		} else {
			s.logger.WithCtx(ctx).Error("failed to increment email verification attempt", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}

	if otp == nil {
		return VerificationTokenNotFoundError
	}

	if otp.Attempts >= shared.EmailVerificationOtpType.MaxAttempts+1 {
		// +1 because we increment the attempt before checking the OTP code validity
		return TooManyVerificationAttemptError
	}

	if otp.Code != p.Code {
		return IncorrectVerificationCodeError
	}

	go func() {
		_, err = s.otpCacheV2.DeleteOtp(ctx, p.UUID, shared.EmailVerificationOtpType)
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

	otp, err := s.otpCacheV2.GetOtp(ctx, p.UUID, shared.PhoneNumberVerificationOtpType)
	var ed *error_pkg.ErrorWithDetails
	if err != nil {
		if errors.As(err, &ed) {
			if ed.Code != OtpNotFoundError.Code {
				s.logger.WithCtx(ctx).Error("failed to get existing otp code", zap.Object("error", ed))
				return ed
			}
		} else {
			s.logger.WithCtx(ctx).Error("failed to get existing otp code", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}
	if otp != nil {
		return OtpAlreadySentError
	}

	code, err := s.util.GeneratorUtil.GenerateRandomNumber(6)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate OTP", zap.Error(err))
		return error_pkg.InternalServerError
	}

	err = s.notificationProvider.SendSms(ctx, &notification_proto.SendSmsRequest{
		Recipient: user.PhoneNumber,
		Body:      fmt.Sprintf("Your verification code for %s is %s", s.config.AppName, code),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send OTP", zap.Error(err))
		return SendPhoneOtpFailedError
	}

	err = s.otpCacheV2.SaveOtp(ctx, p.UUID, shared.PhoneNumberVerificationOtpType, code)
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

	otp, err := s.otpCacheV2.IncrementOtpAttempt(ctx, p.UUID, shared.PhoneNumberVerificationOtpType)
	if err != nil {
		var ed *error_pkg.ErrorWithDetails
		if errors.As(err, &ed) {
			if ed.Code != OtpNotFoundError.Code {
				s.logger.WithCtx(ctx).Error("failed to increment phone OTP attempt", zap.Object("error", ed))
			}
			return ed
		} else {
			s.logger.WithCtx(ctx).Error("failed to increment phone OTP attempt", zap.Error(err))
			return error_pkg.InternalServerError
		}
	}

	if otp == nil {
		return OtpNotFoundError
	}

	if otp.Attempts >= shared.PhoneNumberVerificationOtpType.MaxAttempts+1 {
		// +1 because we increment the attempt before checking the OTP code validity
		return TooManyOtpAttemptError
	}

	if otp.Code != p.Otp {
		return IncorrectOtpError
	}

	go func() {
		_, err := s.otpCacheV2.DeleteOtp(ctx, p.UUID, shared.PhoneNumberVerificationOtpType)
		if err != nil {
			s.logger.WithCtx(ctx).Warn("Failed to delete user phone OTP attempt", zap.Error(err))
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

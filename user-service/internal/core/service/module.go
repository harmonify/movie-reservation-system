package service

import (
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	token_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/token"
	user_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/user"
	"go.uber.org/fx"
)

var (
	AuthServiceModule = fx.Module(
		"auth-service",
		fx.Provide(
			auth_service.NewAuthService,
		),
		fx.Invoke(func(errorMapper error_pkg.ErrorMapper) {
			errorMapper.RegisterErrors(
				auth_service.DuplicateUsernameError,
				auth_service.DuplicateEmailError,
				auth_service.DuplicatePhoneNumberError,
				auth_service.InvalidUsernameError,
				auth_service.InvalidEmailError,
				auth_service.InvalidPhoneNumberError,
				auth_service.UnverifiedEmailError,
				auth_service.UnverifiedPhoneNumberError,
				auth_service.AccountNotFoundError,
				auth_service.IncorrectPasswordError,
				auth_service.InvalidRefreshTokenError,
				auth_service.RefreshTokenAlreadyExpiredError,
			)
		}),
	)

	OtpServiceModule = fx.Module(
		"otp-service",
		fx.Provide(
			otp_service.NewOtpService,
		),
		fx.Invoke(func(errorMapper error_pkg.ErrorMapper) {
			errorMapper.RegisterErrors(
				otp_service.SendVerificationLinkFailedError,
				otp_service.VerificationTokenNotFoundError,
				otp_service.VerificationLinkAlreadyExistError,
				otp_service.VerificationTokenInvalidError,
				otp_service.SendOtpFailedError,
				otp_service.OtpNotFoundError,
				otp_service.OtpAlreadyExistError,
				otp_service.OtpInvalidError,
				otp_service.OtpTooManyAttemptError,
			)
		}),
	)

	TokenServiceModule = fx.Module(
		"token-service",
		fx.Provide(
			token_service.NewTokenService,
		),
		fx.Invoke(func(errorMapper error_pkg.ErrorMapper) {
			errorMapper.RegisterErrors(
				token_service.SessionInvalidError,
				token_service.SessionRevokedError,
				token_service.SessionExpiredError,
			)
		}),
	)

	UserServiceModule = fx.Module(
		"user-service",
		fx.Provide(
			user_service.NewUserService,
		),
	)

	ServiceModule = fx.Module(
		"service",
		AuthServiceModule,
		OtpServiceModule,
		TokenServiceModule,
		UserServiceModule,
	)
)

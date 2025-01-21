package service

import (
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
	)

	OtpServiceModule = fx.Module(
		"otp-service",
		fx.Provide(
			otp_service.NewOtpService,
		),
	)

	TokenServiceModule = fx.Module(
		"token-service",
		fx.Provide(
			token_service.NewTokenService,
		),
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

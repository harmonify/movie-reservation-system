package auth_service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	token_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/token"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/mail"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	AuthService interface {
		Register(ctx context.Context, p RegisterParam) error
		VerifyEmail(ctx context.Context, p VerifyEmailParam) error
		Login(ctx context.Context, p LoginParam) (*LoginResult, error)
		GetToken(ctx context.Context, p GetTokenParam) (*GetTokenResult, error)
		Logout(ctx context.Context, p LogoutParam) error
	}

	AuthServiceParam struct {
		fx.In

		Logger             logger.Logger
		Tracer             tracer.Tracer
		Database           *database.Database
		UserStorage        shared_service.UserStorage
		UserKeyStorage     shared_service.UserKeyStorage
		UserSessionStorage shared_service.UserSessionStorage
		RbacStorage        shared_service.RbacStorage
		Mailer             mail.Mailer
		Util               *util.Util
		Config             *config.Config
		TokenService       token_service.TokenService
		OtpService         otp_service.OtpService
	}

	AuthServiceResult struct {
		fx.Out

		AuthService AuthService
	}

	authServiceImpl struct {
		logger             logger.Logger
		tracer             tracer.Tracer
		database           *database.Database
		userStorage        shared_service.UserStorage
		userKeyStorage     shared_service.UserKeyStorage
		userSessionStorage shared_service.UserSessionStorage
		rbacStorage        shared_service.RbacStorage
		mailer             mail.Mailer
		util               *util.Util
		config             *config.Config
		tokenService       token_service.TokenService
		otpService         otp_service.OtpService
	}
)

func NewAuthService(p AuthServiceParam) AuthServiceResult {
	return AuthServiceResult{
		AuthService: &authServiceImpl{
			logger:             p.Logger,
			tracer:             p.Tracer,
			database:           p.Database,
			userStorage:        p.UserStorage,
			userSessionStorage: p.UserSessionStorage,
			userKeyStorage:     p.UserKeyStorage,
			rbacStorage:        p.RbacStorage,
			mailer:             p.Mailer,
			util:               p.Util,
			config:             p.Config,
			tokenService:       p.TokenService,
			otpService:         p.OtpService,
		},
	}
}

func (s *authServiceImpl) Register(ctx context.Context, p RegisterParam) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Hash user password
	hashedPassword, err := s.util.EncryptionUtil.Argon2Hasher.Hash(p.Password)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to hash user password", zap.Error(err))
		return err
	}

	// Generate user key
	userKey, err := s.tokenService.GenerateUserKey(ctx)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate user key", zap.Error(err))
		return err
	}

	// Start transaction
	err = s.database.Transaction(func(tx *database.Transaction) error {
		// Save user record
		user, err := s.userStorage.WithTx(tx).SaveUser(ctx, entity.SaveUser{
			Username:    p.Username,
			Password:    hashedPassword,
			Email:       p.Email,
			PhoneNumber: p.PhoneNumber,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
		})
		s.logger.WithCtx(ctx).Debug("User record", zap.Any("param", p), zap.Any("user", user))

		if err != nil {
			s.logger.WithCtx(ctx).Error("Failed to save user record", zap.Error(err))
			return err
		}

		// Save encryption key pair as user key record
		_, err = s.userKeyStorage.WithTx(tx).SaveUserKey(ctx, entity.SaveUserKey{
			UserUUID:   user.UUID,
			PublicKey:  userKey.PublicKey,
			PrivateKey: userKey.PrivateKey,
		})
		if err != nil {
			s.logger.WithCtx(ctx).Error("Failed to save user key record", zap.Error(err))
			return err
		}

		// Grant user role
		_, err = s.rbacStorage.WithTx(tx).GrantRole(ctx, shared_service.GrantRoleParam{
			UUID: user.UUID.String(),
			Role: shared_service.RoleUser,
		})
		if err != nil {
			s.logger.WithCtx(ctx).Error("Failed to grant user role", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save records", zap.Error(err))
		return err
	}

	// Send email verification link
	err = s.otpService.SendEmailVerificationLink(ctx, otp_service.SendEmailVerificationLinkParam{
		Email: p.Email,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to send email verification link", zap.Error(err))
		return err
	}

	return err
}

func (s *authServiceImpl) VerifyEmail(ctx context.Context, p VerifyEmailParam) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	err := s.otpService.VerifyEmail(ctx, otp_service.VerifyEmailParam{
		Email: p.Email,
		Token: p.Token,
	})

	if err != nil {
		s.logger.Error("Failed to verify user email", zap.Error(err))
		return err
	}

	s.userStorage.UpdateUser(
		ctx,
		entity.FindUser{
			Email: sql.NullString{String: p.Email, Valid: true},
		},
		entity.UpdateUser{
			IsEmailVerified: sql.NullBool{Bool: true, Valid: true},
		},
	)

	return nil
}

func (s *authServiceImpl) Login(ctx context.Context, p LoginParam) (*LoginResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Get user record
	user, err := s.userStorage.FindUser(ctx, entity.FindUser{Username: sql.NullString{String: p.Username, Valid: true}})
	if err != nil {
		var terr *database.RecordNotFoundError
		if errors.As(err, &terr) {
			return nil, ErrAccountNotFound
		}
		s.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}
	// Get user key record
	userKey, err := s.userKeyStorage.FindUserKey(ctx, entity.FindUserKey{
		UserUUID: sql.NullString{String: user.UUID.String(), Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	// Compare user hashed password with password param.
	match, err := s.util.EncryptionUtil.Argon2Hasher.Compare(user.Password, p.Password)
	if err != nil {
		s.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	} else if !match {
		s.logger.WithCtx(ctx).Info("Password didn't match")
		return nil, ErrIncorrectPassword
	}

	// Generate and encrypt user session
	accessToken, err := s.tokenService.GenerateAccessToken(ctx, token_service.GenerateAccessTokenParam{
		UUID:        user.UUID.String(),
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		PrivateKey:  userKey.PrivateKey,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	// Save user session record
	_, err = s.userSessionStorage.SaveSession(ctx, entity.SaveUserSession{
		UserUUID:     user.UUID.String(),
		RefreshToken: refreshToken.HashedRefreshToken,
		IpAddress:    sql.NullString{String: p.IpAddress, Valid: true},
		UserAgent:    sql.NullString{String: p.UserAgent, Valid: true},
		ExpiredAt:    refreshToken.RefreshTokenExpiredAt,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save user session", zap.Error(err))
		return nil, err
	}

	return &LoginResult{
		AccessToken:           accessToken.AccessToken,
		AccessTokenDuration:   accessToken.AccessTokenDuration,
		RefreshToken:          refreshToken.RefreshToken,
		RefreshTokenExpiredAt: refreshToken.RefreshTokenExpiredAt,
	}, nil
}

func (s *authServiceImpl) GetToken(ctx context.Context, p GetTokenParam) (*GetTokenResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	verifyResult, err := s.tokenService.VerifyRefreshToken(ctx, token_service.VerifyRefreshTokenParam{
		RefreshToken: p.RefreshToken,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to verify refresh token", zap.Error(err))
		return nil, err
	}

	// Get user record
	user, err := s.userStorage.FindUser(ctx, entity.FindUser{UUID: sql.NullString{String: verifyResult.User.UserUUID, Valid: true}})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	// Get user key record
	userKey, err := s.userKeyStorage.FindUserKey(ctx, entity.FindUserKey{
		UserUUID: sql.NullString{String: user.UUID.String(), Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get user key", zap.Error(err))
		return nil, err
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateAccessToken(ctx, token_service.GenerateAccessTokenParam{
		UUID:        user.UUID.String(),
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		PrivateKey:  userKey.PrivateKey,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	return &GetTokenResult{
		AccessToken:         accessToken.AccessToken,
		AccessTokenDuration: accessToken.AccessTokenDuration,
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, p LogoutParam) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	hashedRefreshToken, err := s.util.EncryptionUtil.SHA256Hasher.Hash(p.RefreshToken)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to hash refresh token with SHA256", zap.Error(err))
		return err
	}

	// Revoke user session if exists
	err = s.userSessionStorage.RevokeSession(ctx, hashedRefreshToken)
	if err != nil {
		s.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	return nil
}

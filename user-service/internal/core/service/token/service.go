package token_service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	TokenService interface {
		GenerateUserKey(ctx context.Context) (*GenerateUserKeyResult, error)
		GenerateAccessToken(ctx context.Context, p GenerateAccessTokenParam) (*GenerateAccessTokenResult, error)
		GenerateRefreshToken(ctx context.Context) (*GenerateRefreshTokenResult, error)
		VerifyRefreshToken(ctx context.Context, p VerifyRefreshTokenParam) (*VerifyRefreshTokenResult, error)
	}

	TokenServiceParam struct {
		fx.In

		Logger             logger.Logger
		Tracer             tracer.Tracer
		Config             *config.UserServiceConfig
		Util               *util.Util
		UserSessionStorage shared.UserSessionStorage
	}

	TokenServiceResult struct {
		fx.Out

		TokenService
	}

	tokenServiceImpl struct {
		logger             logger.Logger
		tracer             tracer.Tracer
		config             *config.UserServiceConfig
		util               *util.Util
		userSessionStorage shared.UserSessionStorage

		AccessTokenDuration  int // in seconds
		RefreshTokenDuration int // in seconds
	}
)

func NewTokenService(p TokenServiceParam) TokenServiceResult {
	return TokenServiceResult{
		TokenService: &tokenServiceImpl{
			logger:             p.Logger,
			tracer:             p.Tracer,
			config:             p.Config,
			util:               p.Util,
			userSessionStorage: p.UserSessionStorage,

			AccessTokenDuration:  15 * 60,           // 15 minutes
			RefreshTokenDuration: 30 * 24 * 60 * 60, // 30 days
		},
	}
}

func (s *tokenServiceImpl) GenerateUserKey(ctx context.Context) (*GenerateUserKeyResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Generate RSA key pair for user session encryption
	keyPair, err := s.util.EncryptionUtil.RSAEncryption.Generate()
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate user encryption key pair", zap.Error(err))
		return nil, err
	}

	// Encode the public key to base64
	encodedPublicKey := base64.RawStdEncoding.EncodeToString(keyPair.PublicKey)

	// Encrypt the private key to be securely stored
	encryptedPrivateKey, err := s.util.EncryptionUtil.AESEncryption.Encrypt(string(keyPair.PrivateKey))
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to encrypt user's RSA private key", zap.Error(err))
		return nil, err
	}

	return &GenerateUserKeyResult{
		PublicKey:  encodedPublicKey,
		PrivateKey: encryptedPrivateKey,
	}, nil
}

func (s *tokenServiceImpl) GenerateAccessToken(ctx context.Context, p GenerateAccessTokenParam) (*GenerateAccessTokenResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	decryptedPrivateKey, err := s.util.EncryptionUtil.AESEncryption.Decrypt(p.PrivateKey)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to decrypt user private key", zap.Error(err))
	}
	accessToken, err := s.util.JWTUtil.JWTSign(ctx, jwt_util.JWTSignParam{
		ExpInSeconds: s.AccessTokenDuration,
		PrivateKey:   []byte(decryptedPrivateKey),
		BodyPayload:  p.BodyPayload,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to sign JWT", zap.Error(err))
	}

	return &GenerateAccessTokenResult{
		AccessToken:         accessToken,
		AccessTokenDuration: s.AccessTokenDuration,
	}, nil
}

func (s *tokenServiceImpl) GenerateRefreshToken(ctx context.Context) (*GenerateRefreshTokenResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	refreshToken, err := s.util.GeneratorUtil.GenerateRandomBase64(32)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	hashedRefreshToken, err := s.util.EncryptionUtil.SHA256Hasher.Hash(refreshToken)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to hash refresh token with SHA256", zap.Error(err))
		return nil, err
	}

	return &GenerateRefreshTokenResult{
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: time.Now().Add(time.Second * time.Duration(s.RefreshTokenDuration)),
		HashedRefreshToken:    hashedRefreshToken,
	}, nil
}

func (s *tokenServiceImpl) VerifyRefreshToken(ctx context.Context, p VerifyRefreshTokenParam) (*VerifyRefreshTokenResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	hashedRefreshToken, err := s.util.EncryptionUtil.SHA256Hasher.Hash(p.RefreshToken)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to hash refresh token with SHA256", zap.Error(err))
		return nil, err
	}

	session, err := s.userSessionStorage.GetSession(ctx, entity.GetUserSession{
		RefreshToken: sql.NullString{String: hashedRefreshToken, Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get user session from the storage", zap.Error(err))
		return nil, SessionInvalidError
	}

	if session == nil {
		s.logger.WithCtx(ctx).Info("Session not found")
		return nil, SessionInvalidError
	}

	if session.IsRevoked {
		s.logger.WithCtx(ctx).Info("Session is revoked")
		return nil, SessionRevokedError
	}

	if time.Now().After(session.ExpiredAt) {
		s.logger.WithCtx(ctx).Info("Session is expired")
		return nil, SessionExpiredError
	}

	return &VerifyRefreshTokenResult{
		User: *session,
	}, nil
}

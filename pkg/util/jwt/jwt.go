package jwt_util

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	JwtUtil interface {
		JWTSign(ctx context.Context, payload JWTSignParam) (string, error)
		JWTVerify(ctx context.Context, token string) (*JWTBodyPayload, error)
	}

	JwtUtilParam struct {
		fx.In

		Logger     logger.Logger
		Tracer     tracer.Tracer
		Encryption *encryption.Encryption
	}

	JwtUtilResult struct {
		fx.Out

		JwtUtil JwtUtil
	}

	jwtUtilImpl struct {
		logger     logger.Logger
		tracer     tracer.Tracer
		encryption *encryption.Encryption
		config     *JwtUtilConfig
	}

	JwtUtilConfig struct {
		ServiceIdentifier      string `validate:"required"` // Identifier used when signing JWT iss claim
		JwtAudienceIdentifiers string `validate:"required"` // Comma separated list of identifiers used when signing JWT aud claim
		JwtIssuerIdentifier    string `validate:"required"` // Identifier used when verifying JWT iss claim
	}

	JWTCustomClaims struct {
		Data JWTBodyPayload `json:"data"`
		jwt.RegisteredClaims
	}

	JWTBodyPayload struct {
		UUID        string   `json:"uuid"`
		Permissions []string `json:"permissions"`
	}

	JWTSignParam struct {
		ExpInSeconds int
		PrivateKey   []byte // in PEM format
		BodyPayload  JWTBodyPayload
	}
)

func NewJwtUtil(
	p JwtUtilParam,
	cfg *JwtUtilConfig,
) (JwtUtilResult, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return JwtUtilResult{}, err
	}
	return JwtUtilResult{
		JwtUtil: &jwtUtilImpl{
			logger:     p.Logger,
			tracer:     p.Tracer,
			encryption: p.Encryption,
			config: &JwtUtilConfig{
				ServiceIdentifier:      cfg.ServiceIdentifier,
				JwtAudienceIdentifiers: cfg.JwtAudienceIdentifiers,
				JwtIssuerIdentifier:    cfg.JwtIssuerIdentifier,
			},
		},
	}, nil
}

func (i *jwtUtilImpl) JWTSign(ctx context.Context, payload JWTSignParam) (string, error) {
	ctx, span := i.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(payload.PrivateKey)
	if err != nil {
		i.logger.WithCtx(ctx).Error("failed to parse RSA private key", zap.Error(err))
		return "", err
	}

	// Define time expiration
	now := time.Now()

	// Claim Property
	claims := JWTCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.config.ServiceIdentifier,
			Subject:   payload.BodyPayload.UUID,
			Audience:  strings.Split(i.config.JwtAudienceIdentifiers, ","),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(payload.ExpInSeconds))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			// ID: "", // not needed, apparently to prevent replay attack but it will cause the token to be one-time use, open an issue if this is false
		},
		Data: payload.BodyPayload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = string(i.encryption.RSAEncryption.EncodePublicKey(&privKey.PublicKey))

	// Sign the JWT
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		i.logger.WithCtx(ctx).Error("failed to sign JWT", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (i *jwtUtilImpl) JWTVerify(ctx context.Context, tokenString string) (*JWTBodyPayload, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&JWTCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			publicKey, ok := token.Header["kid"].(string)
			if !ok || publicKey == "" {
				i.logger.WithCtx(ctx).Info("public key is empty")
				return nil, fmt.Errorf("public key is empty")
			}

			// Decode PEM block
			block, _ := pem.Decode([]byte(publicKey))
			if block == nil {
				i.logger.WithCtx(ctx).Info("failed to decode public key into PEM")
				return nil, fmt.Errorf("failed to decode public key into PEM")
			}

			// Parse RSA public key
			rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				i.logger.WithCtx(ctx).Info("failed to parse RSA public key from PEM", zap.Error(err))
				return nil, err
			}

			return rsaPublicKey, nil
		},
		jwt.WithAudience(i.config.ServiceIdentifier),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(i.config.JwtIssuerIdentifier),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
	if err != nil {
		i.logger.WithCtx(ctx).Info("failed to parse JWT", zap.Error(err))
		return nil, error_pkg.InvalidJwtError
	}

	if !parsedToken.Valid {
		i.logger.WithCtx(ctx).Info("invalid JWT", zap.Any("parsed_token", parsedToken))
		return nil, error_pkg.InvalidJwtError
	}

	claims, ok := parsedToken.Claims.(*JWTCustomClaims)
	if !ok {
		i.logger.WithCtx(ctx).Info("failed to assert correct JWT claims type")
		return nil, error_pkg.InvalidJwtClaimsError
	}

	return &claims.Data, nil
}

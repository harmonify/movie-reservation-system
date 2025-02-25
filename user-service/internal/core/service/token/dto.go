package token_service

import (
	"time"

	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type (
	GenerateUserKeyResult struct {
		PublicKey  string // encoded in base64
		PrivateKey string // encrypted with AES
	}

	GenerateAccessTokenParam struct {
		PrivateKey  string
		BodyPayload jwt_util.JWTBodyPayload
	}

	GenerateAccessTokenResult struct {
		AccessToken         string
		AccessTokenDuration int // in seconds
	}

	GenerateRefreshTokenResult struct {
		RefreshToken          string
		RefreshTokenExpiredAt time.Time
		HashedRefreshToken    string
	}

	VerifyRefreshTokenParam struct {
		RefreshToken string
	}

	VerifyRefreshTokenResult struct {
		User entity.UserSession
	}
)

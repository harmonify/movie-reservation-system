package token_service

import (
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type (
	GenerateUserKeyResult struct {
		PublicKey  string // encoded in base64
		PrivateKey string // encrypted with AES
	}

	GenerateAccessTokenParam struct {
		UUID        string
		Username    string
		Email       string
		PhoneNumber string
		PrivateKey  string
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

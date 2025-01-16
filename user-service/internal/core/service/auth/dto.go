package auth_service

import "time"

type (
	RegisterParam struct {
		Username    string
		Email       string
		Password    string
		PhoneNumber string
		FirstName   string
		LastName    string
	}

	VerifyEmailParam struct {
		Email string
		Token string
	}

	LoginParam struct {
		Username  string
		Password  string
		UserAgent string // track
		IpAddress string // track
	}

	LoginResult struct {
		AccessToken           string    `json:"accessToken"`
		AccessTokenDuration   int       `json:"accessTokenDuration"` // in seconds
		RefreshToken          string    `json:"refreshToken"`
		RefreshTokenExpiredAt time.Time `json:"refreshTokenExpiredAt"`
	}

	GetTokenParam struct {
		RefreshToken string
	}

	GetTokenResult struct {
		AccessToken         string `json:"accessToken"`
		AccessTokenDuration int    `json:"accessTokenDuration"` // in seconds
	}

	LogoutParam struct {
		RefreshToken string
		TraceId      string
	}
)

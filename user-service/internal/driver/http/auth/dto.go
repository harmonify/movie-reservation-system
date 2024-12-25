package auth_rest

import (
	http_interface "github.com/harmonify/movie-reservation-system/user-service/lib/http/interface"
)

type (
	RegisterUserReq struct {
		PostUserRegisterReq
		http_interface.HeadersExtension
	}

	PostUserRegisterReq struct {
		Username    string `json:"username" validate:"required"`
		Password    string `json:"password" validate:"required,min=8"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
	}

	PostVerifyEmailReq struct {
		Email string `form:"email" json:"email" validate:"required,email"`
		Token string `form:"token" json:"token" validate:"required"`
	}

	PostUserLoginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	PostUserLoginRes struct {
		AccessToken         string `json:"accessToken"`
		AccessTokenDuration int    `json:"accessTokenDuration"` // in seconds
	}

	GetTokenReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	GetTokenRes struct {
		AccessToken         string `json:"accessToken"`
		AccessTokenDuration int    `json:"accessTokenDuration"` // in seconds
	}

	PostUserLogoutReq struct {
		RefreshToken string
	}
)

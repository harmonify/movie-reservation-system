package auth_rest

import (
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	http_interface "github.com/harmonify/movie-reservation-system/user-service/lib/http/interface"
)

type (
	RegisterUserReq struct {
		PostUserRegisterReq
		http_interface.HeadersExtension
	}

	RegisterUserRes struct {
		FullName    string
		Email       string
		PhoneNumber string
	}

	PostUserRegisterReq struct {
		Username    string `json:"username" validate:"required"`
		Password    string `json:"password" validate:"required,min=8"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		FullName    string `json:"full_name" validate:"required"`
	}

	PostUserRegisterRes struct {
		auth_service.LoginResult
	}

	PostVerifyEmailReq struct {
		Email string `form:"email" json:"email" validate:"required,email"`
	}

	PostUserLoginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	PostUserLoginRes struct {
		auth_service.LoginResult
	}

	GetTokenReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	GetTokenRes struct {
		auth_service.LoginResult
	}
)

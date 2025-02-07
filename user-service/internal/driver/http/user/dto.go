package user_rest

import (
	user_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/user"
)

type (
	GetUserRes user_service.GetUserResult

	PatchUserReq struct {
		Email       string `form:"email" json:"email" validate:"omitempty,email"`
		PhoneNumber string `form:"phone_number" json:"phone_number" validate:"omitempty,e164"`
		Username    string `form:"username" json:"username" validate:"omitempty,alphanum,min=3,max=20"`
		FirstName   string `form:"first_name" json:"first_name" validate:"omitempty,alpha,max=50"`
		LastName    string `form:"last_name" json:"last_name" validate:"omitempty,alpha,max=50"`
	}

	PatchUserRes user_service.UpdateUserResult

	VerifyEmailReq struct {
		Code string `form:"code" json:"code" validate:"required"`
	}

	VerifyPhoneNumberReq struct {
		Otp string `form:"otp" json:"otp" validate:"required,len=6"`
	}
)

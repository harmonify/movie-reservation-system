package user_rest

import (
	"database/sql"
	"time"
)

type (
	GetUserReq struct{}

	GetUserRes struct {
		UUID                  string       `json:"uuid"`
		Username              string       `json:"username"`
		Email                 string       `json:"email"`
		PhoneNumber           string       `json:"phone_number"`
		FirstName             string       `json:"first_name"`
		LastName              string       `json:"last_name"`
		IsEmailVerified       bool         `json:"is_email_verified"`
		IsPhoneNumberVerified bool         `json:"is_phone_number_verified"`
		CreatedAt             time.Time    `json:"created_at"`
		UpdatedAt             time.Time    `json:"updated_at"`
		DeletedAt             sql.NullTime `json:"deleted_at"`
	}

	PatchUserReq struct {
		Username  string `form:"username" json:"username"`
		FirstName string `form:"first_name" json:"first_name"`
		LastName  string `form:"last_name" json:"last_name"`
	}

	PatchUserRes struct {
		UUID                  string       `json:"uuid"`
		Username              string       `json:"username"`
		Email                 string       `json:"email"`
		PhoneNumber           string       `json:"phone_number"`
		FirstName             string       `json:"first_name"`
		LastName              string       `json:"last_name"`
		IsEmailVerified       bool         `json:"is_email_verified"`
		IsPhoneNumberVerified bool         `json:"is_phone_number_verified"`
		CreatedAt             time.Time    `json:"created_at"`
		UpdatedAt             time.Time    `json:"updated_at"`
		DeletedAt             sql.NullTime `json:"deleted_at"`
	}

	GetVerifyUpdateEmailReq struct{}

	PostVerifyUserEmailReq struct {
		Token    string `form:"token" json:"token" validate:"required"`
		NewEmail string `form:"new_email" json:"new_email" validate:"required,email"`
	}

	GetVerifyUpdatePhoneNumberReq struct{}

	PostVerifyUserPhoneNumberReq struct {
		Token          string `form:"token" json:"token" validate:"required"`
		NewPhoneNumber string `form:"new_phone_number" json:"new_phone_number" validate:"required,phone_number"`
	}
)

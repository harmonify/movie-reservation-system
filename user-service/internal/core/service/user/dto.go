package user_service

import (
	"time"
)

type (
	GetUserParam struct {
		UUID string
	}

	GetUserResult struct {
		UUID                  string     `json:"uuid"`
		Username              string     `json:"username"`
		Email                 string     `json:"email"`
		PhoneNumber           string     `json:"phone_number"`
		FirstName             string     `json:"first_name"`
		LastName              string     `json:"last_name"`
		IsEmailVerified       bool       `json:"is_email_verified"`
		IsPhoneNumberVerified bool       `json:"is_phone_number_verified"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             time.Time  `json:"updated_at"`
		DeletedAt             *time.Time `json:"deleted_at"`
	}

	UpdateUserParam struct {
		UUID        string
		Username    string
		FirstName   string
		LastName    string
		Email       string
		PhoneNumber string
	}

	UpdateUserResult struct {
		UUID                  string     `json:"uuid"`
		Username              string     `json:"username"`
		Email                 string     `json:"email"`
		PhoneNumber           string     `json:"phone_number"`
		FirstName             string     `json:"first_name"`
		LastName              string     `json:"last_name"`
		IsEmailVerified       bool       `json:"is_email_verified"`
		IsPhoneNumberVerified bool       `json:"is_phone_number_verified"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             time.Time  `json:"updated_at"`
		DeletedAt             *time.Time `json:"deleted_at"`
	}

	UpdateEmailParam struct {
		UUID     string
		Token    string
		NewEmail string
	}

	UpdatePhoneNumberParam struct {
		UUID           string
		Token          string
		NewPhoneNumber string
	}
)

package user_service

import (
	"database/sql"
	"time"
)

type (
	GetUserParam struct {
		UUID string
	}

	GetUserResult struct {
		UUID                  string
		Username              string
		Email                 string
		PhoneNumber           string
		FirstName             string
		LastName              string
		IsEmailVerified       bool
		IsPhoneNumberVerified bool
		CreatedAt             time.Time
		UpdatedAt             time.Time
		DeletedAt             sql.NullTime
	}

	UpdateUserParam struct {
		UUID      string
		Username  string
		FirstName string
		LastName  string
	}

	UpdateUserResult struct {
		UUID                  string
		Username              string
		Email                 string
		PhoneNumber           string
		FirstName             string
		LastName              string
		IsEmailVerified       bool
		IsPhoneNumberVerified bool
		CreatedAt             time.Time
		UpdatedAt             time.Time
		DeletedAt             sql.NullTime
	}

	GetUpdateEmailVerificationParam struct {
		Email string
	}

	VerifyUpdateEmailParam struct {
		Email    string
		Token    string
		NewEmail string
	}

	GetUpdatePhoneNumberVerificationParam struct {
		PhoneNumber string
	}

	VerifyUpdatePhoneNumberParam struct {
		PhoneNumber    string
		Token          string
		NewPhoneNumber string
	}
)

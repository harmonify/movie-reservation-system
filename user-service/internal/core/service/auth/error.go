package auth_service

import (
	"errors"
	"net/http"

	http_constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
)

var (
	DuplicateEmail       = "DUPLICATE_EMAIL"
	DuplicatePhoneNumber = "DUPLICATE_PHONE_NUMBER"
	InvalidEmail         = "INVALID_EMAIL"
	InvalidPhoneNumber   = "INVALID_PHONE_NUMBER"
	UnverifiedEmail      = "UNVERIFIED_EMAIL"

	ErrDuplicateEmail       = errors.New(DuplicateEmail)
	ErrDuplicatePhoneNumber = errors.New(DuplicatePhoneNumber)
	ErrInvalidEmail         = errors.New(InvalidEmail)
	ErrInvalidPhoneNumber   = errors.New(InvalidPhoneNumber)
	ErrUnverifiedEmail      = errors.New(UnverifiedEmail)

	AuthServiceErrorMap = http_constant.CustomHttpErrorMap{
		DuplicateEmail: {
			HttpCode: http.StatusBadRequest,
			Message:  "The email address is already registered. Please use a different email.",
		},
		DuplicatePhoneNumber: {
			HttpCode: http.StatusBadRequest,
			Message:  "The phone number is already registered. Please use a different phone number.",
		},
		InvalidEmail: {
			HttpCode: http.StatusBadRequest,
			Message:  "The provided email address is invalid. Please enter a valid email.",
		},
		InvalidPhoneNumber: {
			HttpCode: http.StatusBadRequest,
			Message:  "The provided phone number is invalid. Please enter a valid phone number.",
		},
		UnverifiedEmail: {
			HttpCode: http.StatusBadRequest,
			Message:  "Please check your email inbox to verify the account and try again.",
		},
	}
)

package auth_service

import (
	"errors"
	"net/http"

	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
)

var (
	DuplicateUsername          = "DUPLICATE_USERNAME"
	DuplicateEmail             = "DUPLICATE_EMAIL"
	DuplicatePhoneNumber       = "DUPLICATE_PHONE_NUMBER"
	InvalidUsername            = "INVALID_USERNAME"
	InvalidEmail               = "INVALID_EMAIL"
	InvalidPhoneNumber         = "INVALID_PHONE_NUMBER"
	UnverifiedEmail            = "UNVERIFIED_EMAIL"
	UnverifiedPhoneNumber      = "UNVERIFIED_PHONE_NUMBER"
	AccountNotFound            = "ACCOUNT_NOT_FOUND"
	IncorrectPassword          = "INCORRECT_PASSWORD"
	InvalidRefreshToken        = "INVALID_REFRESH_TOKEN"
	RefreshTokenAlreadyExpired = "REFRESH_TOKEN_ALREADY_EXPIRED"

	ErrDuplicateUsername          = errors.New(DuplicateUsername)
	ErrDuplicateEmail             = errors.New(DuplicateEmail)
	ErrDuplicatePhoneNumber       = errors.New(DuplicatePhoneNumber)
	ErrInvalidUsername            = errors.New(InvalidUsername)
	ErrInvalidEmail               = errors.New(InvalidEmail)
	ErrInvalidPhoneNumber         = errors.New(InvalidPhoneNumber)
	ErrUnverifiedEmail            = errors.New(UnverifiedEmail)
	ErrUnverifiedPhoneNumber      = errors.New(UnverifiedPhoneNumber)
	ErrAccountNotFound            = errors.New(AccountNotFound)
	ErrIncorrectPassword          = errors.New(IncorrectPassword)
	ErrInvalidRefreshToken        = errors.New(InvalidRefreshToken)
	ErrRefreshTokenAlreadyExpired = errors.New(RefreshTokenAlreadyExpired)

	AuthServiceErrorMap = error_constant.CustomErrorMap{
		DuplicateUsername: {
			HttpCode: http.StatusBadRequest,
			Message:  "The username is already registered. Please use a different username.",
		},
		DuplicateEmail: {
			HttpCode: http.StatusBadRequest,
			Message:  "The email address is already registered. Please use a different email.",
		},
		DuplicatePhoneNumber: {
			HttpCode: http.StatusBadRequest,
			Message:  "The phone number is already registered. Please use a different phone number.",
		},
		InvalidUsername: {
			HttpCode: http.StatusBadRequest,
			Message:  "The provided username address is invalid. Please enter a valid username.",
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
		UnverifiedPhoneNumber: {
			HttpCode: http.StatusBadRequest,
			Message:  "Please verify your phone number and try again.",
		},
		AccountNotFound: {
			HttpCode: http.StatusNotFound,
			Message:  "The account you're trying to access is not found. Please register an account or check the username you've entered.",
		},
		IncorrectPassword: {
			HttpCode: http.StatusForbidden,
			Message:  "The password you've entered is incorrect. Please try again.",
		},
		InvalidRefreshToken: {
			HttpCode: http.StatusUnauthorized,
			Message:  "Your session is expired. Please login again.",
		},
		RefreshTokenAlreadyExpired: {
			HttpCode: http.StatusBadRequest,
			Message:  "Your session is already expired.",
		},
	}
)

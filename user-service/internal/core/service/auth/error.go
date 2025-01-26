package auth_service

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

var (
	DuplicateUsernameError = &error_pkg.ErrorWithDetails{
		Code:     "DUPLICATE_USERNAME",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.AlreadyExists,
		Message:  "The username is already registered. Please use a different username.",
	}

	DuplicateEmailError = &error_pkg.ErrorWithDetails{
		Code:     "DUPLICATE_EMAIL",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.AlreadyExists,
		Message:  "The email address is already registered. Please use a different email.",
	}

	DuplicatePhoneNumberError = &error_pkg.ErrorWithDetails{
		Code:     "DUPLICATE_PHONE_NUMBER",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.AlreadyExists,
		Message:  "The phone number is already registered. Please use a different phone number.",
	}

	InvalidUsernameError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_USERNAME",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
		Message:  "The provided username address is invalid. Please enter a valid username.",
	}

	InvalidEmailError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_EMAIL",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
		Message:  "The provided email address is invalid. Please enter a valid email.",
	}

	InvalidPhoneNumberError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_PHONE_NUMBER",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
		Message:  "The provided phone number is invalid. Please enter a valid phone number.",
	}

	UnverifiedEmailError = &error_pkg.ErrorWithDetails{
		Code:     "UNVERIFIED_EMAIL",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.PermissionDenied,
		Message:  "Please check your email inbox to verify the account and try again.",
	}

	UnverifiedPhoneNumberError = &error_pkg.ErrorWithDetails{
		Code:     "UNVERIFIED_PHONE_NUMBER",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.PermissionDenied,
		Message:  "Please verify your phone number and try again.",
	}

	AccountNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "ACCOUNT_NOT_FOUND",
		HttpCode: http.StatusNotFound,
		GrpcCode: codes.NotFound,
		Message:  "The account you're trying to access is not found. Please register an account or check the username you've entered.",
	}

	IncorrectPasswordError = &error_pkg.ErrorWithDetails{
		Code:     "INCORRECT_PASSWORD",
		HttpCode: http.StatusForbidden,
		GrpcCode: codes.PermissionDenied,
		Message:  "The password you've entered is incorrect. Please try again.",
	}

	InvalidRefreshTokenError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_REFRESH_TOKEN",
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your session is expired. Please login again.",
	}

	RefreshTokenAlreadyExpiredError = &error_pkg.ErrorWithDetails{
		Code:     "REFRESH_TOKEN_ALREADY_EXPIRED",
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your session is already expired.",
	}
)

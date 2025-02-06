package user_service

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

var (
	UsernameAlreadyExistsError = &error_pkg.ErrorWithDetails{
		Code:     error_pkg.ErrorCode("USERNAME_ALREADY_EXISTS"),
		Message:  "Username is already taken by another user",
		HttpCode: http.StatusConflict,
		GrpcCode: codes.AlreadyExists,
	}

	EmailAlreadyExistsError = &error_pkg.ErrorWithDetails{
		Code:     error_pkg.ErrorCode("EMAIL_ALREADY_EXISTS"),
		Message:  "Email is already used by another user",
		HttpCode: http.StatusConflict,
		GrpcCode: codes.AlreadyExists,
	}

	PhoneNumberAlreadyExistsError = &error_pkg.ErrorWithDetails{
		Code:     error_pkg.ErrorCode("PHONE_NUMBER_ALREADY_EXISTS"),
		Message:  "Phone number is already used by another user",
		HttpCode: http.StatusConflict,
		GrpcCode: codes.AlreadyExists,
	}
)

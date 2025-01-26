package token_service

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

var (
	SessionInvalidError = &error_pkg.ErrorWithDetails{
		Code:     "SESSION_INVALID",
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Invalid session. Please log in again.",
	}

	SessionRevokedError = &error_pkg.ErrorWithDetails{
		Code:     "SESSION_REVOKED",
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "You have been logged out from this device.",
	}

	SessionExpiredError = &error_pkg.ErrorWithDetails{
		Code:     "SESSION_EXPIRED",
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your session has expired. Please log in again.",
	}
)

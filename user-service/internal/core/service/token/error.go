package token_service

import (
	"errors"
	"net/http"

	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	grpc_constant "github.com/harmonify/movie-reservation-system/user-service/lib/grpc/constant"
)

var (
	SessionInvalid = "SESSION_INVALID"
	SessionRevoked = "SESSION_REVOKED"
	SessionExpired = "SESSION_EXPIRED"

	ErrSessionInvalid = errors.New(SessionInvalid)
	ErrSessionRevoked = errors.New(SessionRevoked)
	ErrSessionExpired = errors.New(SessionExpired)

	OtpServiceErrorMap = error_constant.CustomErrorMap{
		SessionInvalid: {
			HttpCode: http.StatusUnauthorized,
			GrpcCode: grpc_constant.GrpcUnauthenticated,
			Message:  "Invalid session. Please log in again.",
		},
		SessionRevoked: {
			HttpCode: http.StatusUnauthorized,
			GrpcCode: grpc_constant.GrpcUnauthenticated,
			Message:  "You have been logged out from this device.",
		},
		SessionExpired: {
			HttpCode: http.StatusUnauthorized,
			GrpcCode: grpc_constant.GrpcUnauthenticated,
			Message:  "Your session has expired. Please log in again.",
		},
	}
)

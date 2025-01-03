package error_constant

import (
	"errors"
	"net/http"

	grpc_constant "github.com/harmonify/movie-reservation-system/user-service/lib/grpc/constant"
)

var (
	InvalidRequestBody      = "INVALID_REQUEST_BODY"
	Unauthorized            = "UNAUTHORIZED"
	InvalidJwt              = "INVALID_JWT"
	InvalidJwtClaims        = "INVALID_JWT_CLAIMS"
	InvalidJwtSigningMethod = "INVALID_JWT_SIGNING_METHOD"
	Forbidden               = "FORBIDDEN"
	NotFound                = "NOT_FOUND"
	UnprocessableEntity     = "UNPROCESSABLE_ENTITY"
	RateLimitExceeded       = "RATE_LIMIT_EXCEEDED"
	InternalServerError     = "INTERNAL_SERVER_ERROR"
	DefaultError            = ""
	ServiceUnavailable      = "SERVICE_UNAVAILABLE"
)

var (
	ErrInvalidRequestBody      = errors.New(InvalidRequestBody)
	ErrUnauthorized            = errors.New(Unauthorized)
	ErrInvalidJwt              = errors.New(InvalidJwt)
	ErrInvalidJwtClaims        = errors.New(InvalidJwtClaims)
	ErrInvalidJwtSigningMethod = errors.New(InvalidJwtSigningMethod)
	ErrForbidden               = errors.New(Forbidden)
	ErrNotFound                = errors.New(NotFound)
	ErrUnprocessableEntity     = errors.New(UnprocessableEntity)
	ErrRateLimitExceeded       = errors.New(RateLimitExceeded)
	ErrInternalServerError     = errors.New(InternalServerError)
	ErrServiceUnavailable      = errors.New(ServiceUnavailable)
)

type CustomError struct {
	HttpCode int
	GrpcCode int
	Message  string
	Errors   map[string]interface{}
}

type CustomErrorMap = map[string]CustomError

var (
	DefaultCustomErrorMap = CustomErrorMap{
		InvalidRequestBody: {
			HttpCode: http.StatusBadRequest,
			GrpcCode: grpc_constant.GrpcInvalidArgument,
			Message:  "Please ensure you have filled all the information required and try again. If the problem persists, please contact our technical support.",
		},

		Unauthorized: {
			HttpCode: http.StatusUnauthorized,
			GrpcCode: grpc_constant.GrpcUnauthenticated,
			Message:  "Your request is unauthorized. Please ensure you have the correct credentials and try again.",
		},

		Forbidden: {
			HttpCode: http.StatusForbidden,
			GrpcCode: grpc_constant.GrpcPermissionDenied,
			Message:  "You are forbidden to access this resource.",
		},

		NotFound: {
			HttpCode: http.StatusNotFound,
			GrpcCode: grpc_constant.GrpcNotFound,
			Message:  "Resource not found.",
		},

		UnprocessableEntity: {
			HttpCode: http.StatusUnprocessableEntity,
			GrpcCode: grpc_constant.GrpcInvalidArgument,
			Message:  "We are unable to process your request. Please contact our technical support.",
		},

		RateLimitExceeded: {
			HttpCode: http.StatusTooManyRequests,
			GrpcCode: grpc_constant.GrpcResourceExhausted,
			Message:  "You have exceeded the allowed rate limit for this operation. Please try again later.",
		},

		InternalServerError: {
			HttpCode: http.StatusInternalServerError,
			GrpcCode: grpc_constant.GrpcInternal,
			Message:  "Something went wrong from our side. Please try again later.",
		},

		DefaultError: {
			HttpCode: http.StatusInternalServerError,
			GrpcCode: grpc_constant.GrpcInternal,
			Message:  "Something went wrong from our side. Please try again later.",
		},

		ServiceUnavailable: {
			HttpCode: http.StatusServiceUnavailable,
			GrpcCode: grpc_constant.GrpcUnavailable,
			Message:  "The server is currently unable to handle the request. Please try again later.",
		},
	}
)

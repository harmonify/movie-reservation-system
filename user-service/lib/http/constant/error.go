package http_constant

import (
	"errors"
	"net/http"
)

var (
	InvalidRequestBody      = "INVALID_REQUEST_BODY"
	Unauthorized            = "UNAUTHORIZED"
	InvalidJwtSigningMethod = "INVALID_JWT_SIGNING_METHOD"
	InvalidJwtFormat        = "INVALID_JWT_FORMAT"
	InvalidJwt              = "INVALID_JWT"
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
	ErrInvalidJwtFormat        = errors.New(InvalidJwtFormat)
	ErrInvalidJwtSigningMethod = errors.New(InvalidJwtSigningMethod)
	ErrInvalidJwt              = errors.New(InvalidJwt)
	ErrForbidden               = errors.New(Forbidden)
	ErrNotFound                = errors.New(NotFound)
	ErrUnprocessableEntity     = errors.New(UnprocessableEntity)
	ErrRateLimitExceeded       = errors.New(RateLimitExceeded)
	ErrInternalServerError     = errors.New(InternalServerError)
	ErrServiceUnavailable      = errors.New(ServiceUnavailable)
)

type CustomHttpError struct {
	HttpCode int
	Message  string
	Errors   map[string]interface{}
}

type CustomHttpErrorMap = map[string]CustomHttpError

var (
	DefaultCustomHttpErrorMap = CustomHttpErrorMap{
		InvalidRequestBody: {
			HttpCode: http.StatusBadRequest,
			Message:  "Please ensure you have filled all the information required and try again. If the problem persists, please contact our technical support.",
		},

		Unauthorized: {
			HttpCode: http.StatusUnauthorized,
			Message:  "Your request is unauthorized. Please ensure you have the correct credentials and try again.",
		},

		Forbidden: {
			HttpCode: http.StatusForbidden,
			Message:  "You are forbidden to access this resource.",
		},

		NotFound: {
			HttpCode: http.StatusNotFound,
			Message:  "Resource not found.",
		},

		UnprocessableEntity: {
			HttpCode: http.StatusUnprocessableEntity,
			Message:  "We are unable to process your request. Please contact our technical support.",
		},

		RateLimitExceeded: {
			HttpCode: http.StatusTooManyRequests,
			Message:  "You have exceeded the allowed rate limit for this operation. Please try again later.",
		},

		InternalServerError: {
			HttpCode: http.StatusInternalServerError,
			Message:  "Something went wrong from our side. Please try again later.",
		},

		DefaultError: {
			HttpCode: http.StatusInternalServerError,
			Message:  "Something went wrong from our side. Please try again later.",
		},

		ServiceUnavailable: {
			HttpCode: http.StatusServiceUnavailable,
			Message:  "The server is currently unable to handle the request. Please try again later.",
		},
	}
)

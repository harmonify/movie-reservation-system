package error_pkg

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	DefaultError = &ErrorWithDetails{
		Code:     ErrorCode("ERROR"),
		HttpCode: http.StatusInternalServerError,
		GrpcCode: codes.Internal,
		Message:  "Something went wrong from our side. Please try again later.",
	}

	InternalServerError = &ErrorWithDetails{
		Code:     ErrorCode("INTERNAL_SERVER_ERROR"),
		HttpCode: http.StatusInternalServerError,
		GrpcCode: codes.Internal,
		Message:  "Something went wrong from our side. Please try again later.",
	}

	InvalidRequestBodyError = &ErrorWithDetails{
		Code:     ErrorCode("INVALID_REQUEST_BODY_ERROR"),
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
		Message:  "Please ensure you have filled all the required information correctly and try again. If the problem persists, please contact our technical support.",
	}

	UnauthorizedError = &ErrorWithDetails{
		Code:     ErrorCode("UNAUTHORIZED_ERROR"),
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your request is unauthorized. Please ensure you have the correct credentials and try again.",
	}

	InvalidAuthorizationHeaderError = &ErrorWithDetails{
		Code:     ErrorCode("INVALID_AUTHORIZATION_HEADER_ERROR"),
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your request is malformed. Please ensure you have the correct credentials and try again.",
	}

	InvalidJwtError = &ErrorWithDetails{
		Code:     ErrorCode("INVALID_JWT_ERROR"),
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your request is unauthorized. Please ensure you have the correct credentials and try again.",
	}

	InvalidJwtClaimsError = &ErrorWithDetails{
		Code:     ErrorCode("INVALID_JWT_CLAIMS_ERROR"),
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your request is unauthorized. Please ensure you have the correct credentials and try again.",
	}

	InvalidJwtSigningMethodError = &ErrorWithDetails{
		Code:     ErrorCode("INVALID_JWT_SIGNING_METHOD_ERROR"),
		HttpCode: http.StatusUnauthorized,
		GrpcCode: codes.Unauthenticated,
		Message:  "Your request is unauthorized. Please ensure you have the correct credentials and try again.",
	}

	ForbiddenError = &ErrorWithDetails{
		Code:     ErrorCode("FORBIDDEN_ERROR"),
		HttpCode: http.StatusForbidden,
		GrpcCode: codes.PermissionDenied,
		Message:  "You are forbidden to access this resource.",
	}

	NotFoundError = &ErrorWithDetails{
		Code:     ErrorCode("NOT_FOUND_ERROR"),
		HttpCode: http.StatusNotFound,
		GrpcCode: codes.NotFound,
		Message:  "Resource not found.",
	}

	UnprocessableEntityError = &ErrorWithDetails{
		Code:     ErrorCode("UNPROCESSABLE_ENTITY_ERROR"),
		HttpCode: http.StatusUnprocessableEntity,
		GrpcCode: codes.InvalidArgument,
		Message:  "We are unable to process your request. Please contact our technical support.",
	}

	RateLimitExceededError = &ErrorWithDetails{
		Code:     ErrorCode("RATE_LIMIT_EXCEEDED_ERROR"),
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "You have exceeded the allowed rate limit for this operation. Please try again later.",
	}

	BadGatewayError = &ErrorWithDetails{
		Code:     ErrorCode("BAD_GATEWAY_ERROR"),
		HttpCode: http.StatusBadGateway,
		GrpcCode: codes.Unavailable,
		Message:  "The server is currently unable to handle the request. Please try again later.",
	}

	ServiceUnavailableError = &ErrorWithDetails{
		Code:     ErrorCode("SERVICE_UNAVAILABLE_ERROR"),
		HttpCode: http.StatusServiceUnavailable,
		GrpcCode: codes.Unavailable,
		Message:  "The server is currently unable to handle the request. Please try again later.",
	}
)

type RateLimitExceededErrorData struct {
	RetryAfter int `json:"retry_after"`
}

func NewRateLimitExceededError(retryAfter int) *ErrorWithDetails {
	return &ErrorWithDetails{
		Code:     RateLimitExceededError.Code,
		HttpCode: RateLimitExceededError.HttpCode,
		GrpcCode: RateLimitExceededError.GrpcCode,
		Message:  RateLimitExceededError.Message,
		Data:     &RateLimitExceededErrorData{RetryAfter: retryAfter},
	}
}

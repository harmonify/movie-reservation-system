package otp_service

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

var (
	SendVerificationLinkFailedError = &error_pkg.ErrorWithDetails{
		Code:     "SEND_VERIFICATION_LINK_FAILED",
		HttpCode: http.StatusBadGateway,
		GrpcCode: codes.Unavailable,
		Message:  "Failed to send a verification link to your email. If issue persists, please contact our technical support and try again later",
	}

	VerificationTokenNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "VERIFICATION_TOKEN_NOT_FOUND",
		HttpCode: http.StatusNotFound,
		GrpcCode: codes.NotFound,
		Message:  "Your verification link may be expired. Please try to request a new verification link.",
	}

	VerificationLinkAlreadyExistError = &error_pkg.ErrorWithDetails{
		Code:     "VERIFICATION_LINK_ALREADY_EXIST",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "A verification link is already sent to your inbox. Please check your inbox and try again later.",
	}

	VerificationTokenInvalidError = &error_pkg.ErrorWithDetails{
		Code:     "VERIFICATION_TOKEN_INVALID",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
		Message:  "Failed to verify your email. Please try to request a new verification link.",
	}

	SendOtpFailedError = &error_pkg.ErrorWithDetails{
		Code:     "SEND_OTP_FAILED",
		HttpCode: http.StatusBadGateway,
		GrpcCode: codes.Unavailable,
		Message:  "Failed to send an OTP to your phone number. If issue persists, please contact our technical support and try again later",
	}

	OtpNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "OTP_NOT_FOUND",
		HttpCode: http.StatusNotFound,
		GrpcCode: codes.NotFound,
		Message:  "Your OTP may be expired. Please try to request a new OTP.",
	}

	OtpAlreadySentError = &error_pkg.ErrorWithDetails{
		Code:     "OTP_ALREADY_SENT",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "OTP is already sent to your inbox. Please try again later.",
	}

	OtpInvalidError = &error_pkg.ErrorWithDetails{
		Code:     "OTP_INVALID",
		HttpCode: http.StatusForbidden,
		GrpcCode: codes.PermissionDenied,
		Message:  "Incorrect OTP. Please try again.",
	}

	OtpTooManyAttemptError = &error_pkg.ErrorWithDetails{
		Code:     "OTP_TOO_MANY_ATTEMPT",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "You have attempted OTP too many times. Please try again later.",
	}
)

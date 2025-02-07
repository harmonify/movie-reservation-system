package otp_service

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

// Generic OTP errors
var (
	OtpNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "OTP_NOT_FOUND",
		HttpCode: http.StatusNotFound,
		GrpcCode: codes.NotFound,
		Message:  "Your OTP may be expired. Please try to request a new OTP.",
	}
)

// Verification email errors
var (
	SendVerificationLinkFailedError = &error_pkg.ErrorWithDetails{
		Code:     "SEND_VERIFICATION_LINK_FAILED",
		HttpCode: http.StatusBadGateway,
		GrpcCode: codes.Unavailable,
		Message:  "Failed to send a verification link to your email. If issue persists, please contact our technical support and try again later",
	}

	VerificationLinkAlreadySentError = &error_pkg.ErrorWithDetails{
		Code:     "VERIFICATION_LINK_ALREADY_SENT",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "A verification link is already sent to your inbox. Please check your inbox and try again later.",
	}

	IncorrectVerificationCodeError = &error_pkg.ErrorWithDetails{
		Code:     "INCORRECT_VERIFICATION_CODE",
		HttpCode: http.StatusForbidden,
		GrpcCode: codes.PermissionDenied,
		Message:  "Incorrect verification code. Please try to request a new verification link.",
	}

	TooManyVerificationAttemptError = &error_pkg.ErrorWithDetails{
		Code:     "TOO_MANY_VERIFICATION_ATTEMPT",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "You have attempted verification process too many times. Please try again later.",
	}

	VerificationTokenNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "VERIFICATION_TOKEN_NOT_FOUND",
		HttpCode: http.StatusNotFound,
		GrpcCode: codes.NotFound,
		Message:  "Your verification link may be expired. Please try to request a new verification link.",
	}
)

// Phone OTP errors
var (
	SendPhoneOtpFailedError = &error_pkg.ErrorWithDetails{
		Code:     "SEND_OTP_FAILED",
		HttpCode: http.StatusBadGateway,
		GrpcCode: codes.Unavailable,
		Message:  "Failed to send an OTP to your phone number. If issue persists, please contact our technical support and try again later",
	}

	OtpAlreadySentError = &error_pkg.ErrorWithDetails{
		Code:     "OTP_ALREADY_SENT",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "OTP is already sent to your inbox. Please try again later.",
	}

	IncorrectOtpError = &error_pkg.ErrorWithDetails{
		Code:     "INCORRECT_OTP",
		HttpCode: http.StatusForbidden,
		GrpcCode: codes.PermissionDenied,
		Message:  "Incorrect OTP. Please try again.",
	}

	TooManyOtpAttemptError = &error_pkg.ErrorWithDetails{
		Code:     "TOO_MANY_OTP_ATTEMPT",
		HttpCode: http.StatusTooManyRequests,
		GrpcCode: codes.ResourceExhausted,
		Message:  "You have attempted OTP too many times. Please try again later.",
	}
)

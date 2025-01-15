package otp_service

import (
	"errors"
	"net/http"

	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
)

var (
	SendVerificationLinkFailed   = "SEND_VERIFICATION_LINK_FAILED"
	VerificationTokenNotFound    = "VERIFICATION_TOKEN_NOT_FOUND"
	VerificationLinkAlreadyExist = "VERIFICATION_LINK_ALREADY_EXIST"
	VerificationTokenInvalid     = "VERIFICATION_TOKEN_INVALID"
	SendOtpFailed                = "SEND_OTP_FAILED"
	OtpNotFound                  = "OTP_NOT_FOUND"
	OtpAlreadyExist              = "OTP_ALREADY_EXIST"
	OtpInvalid                   = "OTP_INVALID"
	OtpTooManyAttempt            = "OTP_TOO_MANY_ATTEMPT"

	ErrSendVerificationLinkFailed   = errors.New(SendVerificationLinkFailed)
	ErrVerificationTokenNotFound    = errors.New(VerificationTokenNotFound)
	ErrVerificationLinkAlreadyExist = errors.New(VerificationLinkAlreadyExist)
	ErrVerificationTokenInvalid     = errors.New(VerificationTokenInvalid)
	ErrSendOtpFailed                = errors.New(SendOtpFailed)
	ErrOtpNotFound                  = errors.New(OtpNotFound)
	ErrOtpAlreadyExist              = errors.New(OtpAlreadyExist)
	ErrOtpInvalid                   = errors.New(OtpInvalid)
	ErrOtpTooManyAttempt            = errors.New(OtpTooManyAttempt)

	OtpServiceErrorMap = error_constant.CustomErrorMap{
		SendVerificationLinkFailed: {
			HttpCode: http.StatusBadGateway,
			Message:  "Failed to send a verification link to your email. If issue persists, please contact our technical support and try again later",
		},
		VerificationTokenNotFound: {
			HttpCode: http.StatusNotFound,
			Message:  "Your verification link may be expired. Please try to request a new verification link.",
		},
		VerificationLinkAlreadyExist: {
			HttpCode: http.StatusTooManyRequests,
			Message:  "A verification link is already sent to your inbox. Please check your inbox and try again later.",
		},
		VerificationTokenInvalid: {
			HttpCode: http.StatusForbidden,
			Message:  "Failed to verify your email. Please try to request a new verification link.",
		},
		SendOtpFailed: {
			HttpCode: http.StatusBadGateway,
			Message:  "Failed to send an OTP to your phone number. If issue persists, please contact our technical support and try again later",
		},
		OtpNotFound: {
			HttpCode: http.StatusNotFound,
			Message:  "Your OTP may be expired. Please try to request a new OTP.",
		},
		OtpAlreadyExist: {
			HttpCode: http.StatusTooManyRequests,
			Message:  "OTP is already sent to your inbox. Please try again later.",
		},
		OtpInvalid: {
			HttpCode: http.StatusForbidden,
			Message:  "Incorrect OTP. Please try again.",
		},
		OtpTooManyAttempt: {
			HttpCode: http.StatusTooManyRequests,
			Message:  "You have attempted OTP too many times. Please try again later.",
		},
	}
)

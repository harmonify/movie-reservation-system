package otp_service

import (
	"errors"
	"net/http"

	http_constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
)

var (
	KeyNotExist                  = "KEY_NOT_EXIST"
	VerificationLinkAlreadyExist = "VERIFICATION_LINK_ALREADY_EXIST"
	OtpAlreadyExist              = "OTP_ALREADY_EXIST"

	ErrKeyNotExist                  = errors.New(KeyNotExist)
	ErrVerificationLinkAlreadyExist = errors.New(VerificationLinkAlreadyExist)
	ErrOtpAlreadyExist              = errors.New(OtpAlreadyExist)

	OtpServiceErrorMap = http_constant.CustomHttpErrorMap{
		VerificationLinkAlreadyExist: {
			HttpCode: http.StatusBadRequest,
			Message:  "The verification link is already sent to your inbox. Please check your inbox and try again later.",
		},
		OtpAlreadyExist: {
			HttpCode: http.StatusBadRequest,
			Message:  "The OTP is already sent to your inbox. Please try again later.",
		},
	}
)

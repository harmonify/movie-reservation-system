package shared

import "time"

type (
	SaveEmailVerificationCodeParam struct {
		Email string
		Code  string
		TTL   time.Duration
	}

	SavePhoneNumberVerificationOtpParam struct {
		PhoneNumber string
		Otp         string
		TTL         time.Duration
	}
)

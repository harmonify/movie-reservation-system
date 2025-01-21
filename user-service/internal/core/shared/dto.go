package shared

import "time"

type (
	SaveEmailVerificationTokenParam struct {
		Email string
		Token string
		TTL   time.Duration
	}

	SavePhoneOtpParam struct {
		PhoneNumber string
		Otp         string
		TTL         time.Duration
	}
)

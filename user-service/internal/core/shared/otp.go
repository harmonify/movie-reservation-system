package shared

import "time"

var (
	EmailVerificationOtpType = OtpType{
		Name:        "email-verification",
		TTL:         24 * time.Hour,
		MaxAttempts: 3,
	}

	PhoneNumberVerificationOtpType = OtpType{
		Name:        "phone-number-verification",
		TTL:         5 * time.Minute,
		MaxAttempts: 3,
	}
)

type (
	// V1 types

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

	// V2 types

	OtpType struct {
		Name        string
		TTL         time.Duration
		MaxAttempts int
	}

	Otp struct {
		Code     string
		Attempts int
	}
)

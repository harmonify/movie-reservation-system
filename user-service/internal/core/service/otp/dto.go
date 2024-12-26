package otp_service

type (
	GetEmailVerificationTokenParam struct {
		Email string
		TTL   uint32 // in seconds
	}

	SendEmailVerificationLinkParam struct {
		Email string
	}

	VerifyEmailParam struct {
		Email string
		Token string
	}

	SendPhoneOtpParam struct {
		PhoneNumber string
	}

	VerifyPhoneOtpParam struct {
		PhoneNumber string
		Otp         string
	}
)

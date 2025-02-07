package otp_service

type (
	GetEmailVerificationTokenParam struct {
		Email string
		TTL   uint32 // in seconds
	}

	SendSignupEmailParam struct {
		UUID      string
		Email     string
		FirstName string
		LastName  string
	}

	VerifyEmailParam struct {
		UUID string
		Code string
	}

	SendPhoneNumberVerificationOtpParam struct {
		UUID string
	}

	VerifyPhoneNumberParam struct {
		UUID string
		Otp  string
	}

	SendVerificationEmailParam struct {
		UUID string
	}
)

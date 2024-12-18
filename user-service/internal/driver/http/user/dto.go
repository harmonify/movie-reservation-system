package user_rest

type (
	GetQueryInfoReq struct {
		Email       string `form:"email" json:"email"`
		PhoneNumber string `form:"phone_number" json:"phone_number"`
	}

	GetQueryInfoRes struct {
		Registered           bool `json:"isRegistered"`
		OtpRequired          bool `json:"otpRequired"`
		TwoFactorRequired    bool `json:"twoFactorRequired"`
		EmailVerified        bool `json:"emailVerified"`
		PhoneNumberVerified  bool `json:"phoneNumberVerified"`
		IsPINSetup           bool `json:"isPINSetup"`
		IsFinishedOnboarding bool `json:"isFinishedOnboarding"`
		IsGuest              bool `json:"isGuest"`
	}

	PatchUserReq struct {
		Email       string `form:"email" json:"email" validate:"omitempty,email"`
		PhoneNumber string `form:"phoneNumber" json:"phoneNumber" validate:"omitempty,phone_number"`
	}

	PatchUserRes struct {
		Email       string `form:"email" json:"email"`
		PhoneNumber string `form:"phoneNumber" json:"phoneNumber"`
	}
)

package auth_rest

type (
	PostRegisterReq struct {
		Username    string `json:"username" validate:"required"`
		Password    string `json:"password" validate:"required,min=8"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
	}

	PostVerifyEmailReq struct {
		Email string `form:"email" json:"email" validate:"required,email"`
		Token string `form:"token" json:"token" validate:"required"`
	}

	PostLoginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	PostLoginRes struct {
		AccessToken         string `json:"accessToken"`
		AccessTokenDuration int    `json:"accessTokenDuration"` // in seconds
	}

	GetTokenReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	GetTokenRes struct {
		AccessToken         string `json:"accessToken"`
		AccessTokenDuration int    `json:"accessTokenDuration"` // in seconds
	}

	PostUserLogoutReq struct {
		RefreshToken string
	}
)

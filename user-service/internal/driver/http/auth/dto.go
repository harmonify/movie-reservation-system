package auth_rest

type (
	PostRegisterReq struct {
		Username    string `json:"username" validate:"required,alphanum,min=3,max=20"`
		Password    string `json:"password" validate:"required,min=8"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required,e164"`
		FirstName   string `json:"first_name" validate:"required,alpha_space,max=50"`
		LastName    string `json:"last_name" validate:"required,alpha_space,max=50"`
	}

	PostVerifyEmailReq struct {
		Email string `form:"email" json:"email" validate:"required,email"`
		Token string `form:"token" json:"token" validate:"required"`
	}

	PostLoginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	PostLoginRes struct {
		AccessToken         string `json:"access_token"`
		AccessTokenDuration int    `json:"access_token_duration"` // in seconds
	}

	GetTokenRes struct {
		AccessToken         string `json:"access_token"`
		AccessTokenDuration int    `json:"access_token_duration"` // in seconds
	}

	PostUserLogoutReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)

package auth_service

type (
	RegisterParam struct {
		Username    string
		Password    string
		Email       string
		PhoneNumber string
		FullName    string
	}

	LoginRecord struct {
		UUID             string
		RefreshTokenHash string
		Revoked          NullBool
		IpAddress        string
		UserAgent        string
	}

	LoginParam struct {
		Username string
		Password string
	}

	LoginResult struct {
		AccessToken           string
		RefreshToken          string
		AccessTokenDuration   int   // in seconds
		RefreshTokenExpiredAt int64 // epoch
	}

	GetTokenParam struct {
		RefreshToken string
	}

	LogoutParam struct {
		RefreshTokens []string
	}
)

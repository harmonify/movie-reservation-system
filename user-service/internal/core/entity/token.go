package entity

type Token struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenDuration   int   // in seconds
	RefreshTokenExpiredAt int64 // epoch
}

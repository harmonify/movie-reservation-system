package factory

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
)

type UserSessionFactory interface {
	CreateUserSession(p CreateUserSessionParam) (session *model.UserSession, hashedRefreshToken string)
}

func NewUserSessionFactory() UserSessionFactory {
	return &userSessionFactoryImpl{}
}

type userSessionFactoryImpl struct {
}

type CreateUserSessionParam struct {
	UserUUID         string
	HashRefreshToken bool
}

func (f *userSessionFactoryImpl) CreateUserSession(p CreateUserSessionParam) (*model.UserSession, string) {
	expiredAt := time.Now().Add(30 * 24 * time.Hour)                    // 30 days
	hashedRefreshToken := "dKnLXT92KEHhVAhvPUuEDzzwal/ZExRKbgLa3YUxVbo" // ANfKOUKiJwys5Lyg49JTy3SbboHJmMO9SbRdbZ9jKA0
	refreshToken := "ANfKOUKiJwys5Lyg49JTy3SbboHJmMO9SbRdbZ9jKA0"
	if p.HashRefreshToken {
		refreshToken = hashedRefreshToken
	}
	return &model.UserSession{
		UserUUID:     p.UserUUID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiredAt:    expiredAt,
		IpAddress:    sql.NullString{String: "", Valid: false},
		UserAgent:    sql.NullString{String: "", Valid: false},
	}, hashedRefreshToken
}

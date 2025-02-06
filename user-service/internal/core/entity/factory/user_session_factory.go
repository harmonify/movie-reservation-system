package entityfactory

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type UserSessionRaw struct {
	RefreshToken string
}

type UserSessionFactory interface {
	GenerateUserSession(user *entity.User) (session *entity.UserSession, raw *UserSessionRaw)
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

func (f *userSessionFactoryImpl) GenerateUserSession(user *entity.User) (*entity.UserSession, *UserSessionRaw) {

	return &entity.UserSession{
			UserUUID:     user.UUID,
			TraceID:      user.TraceID,
			RefreshToken: "dKnLXT92KEHhVAhvPUuEDzzwal/ZExRKbgLa3YUxVbo", // hashed ANfKOUKiJwys5Lyg49JTy3SbboHJmMO9SbRdbZ9jKA0
			IsRevoked:    false,
			ExpiredAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
			IpAddress:    sql.NullString{String: "", Valid: false},
			UserAgent:    sql.NullString{String: "", Valid: false},
		},
		&UserSessionRaw{
			RefreshToken: "ANfKOUKiJwys5Lyg49JTy3SbboHJmMO9SbRdbZ9jKA0",
		}
}

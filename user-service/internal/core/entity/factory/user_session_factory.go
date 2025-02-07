package entityfactory

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type UserSessionRaw struct {
	RefreshToken string
}

type UserSessionFactory interface {
	GenerateUserSession(user *entity.User) (session *entity.UserSession, raw *UserSessionRaw, err error)
}

func NewUserSessionFactory(sha256Hasher encryption.SHA256Hasher) UserSessionFactory {
	return &userSessionFactoryImpl{
		sha256Hasher: sha256Hasher,
	}
}

type userSessionFactoryImpl struct {
	sha256Hasher encryption.SHA256Hasher
}

func (f *userSessionFactoryImpl) GenerateUserSession(user *entity.User) (*entity.UserSession, *UserSessionRaw, error) {
	refreshToken := "ANfKOUKiJwys5Lyg49JTy3SbboHJmMO9SbRdbZ9jKA0"
	hashed, err := f.sha256Hasher.Hash(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	return &entity.UserSession{
			UserUUID:     user.UUID,
			TraceID:      user.TraceID,
			RefreshToken: hashed,
			IsRevoked:    false,
			ExpiredAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
			IpAddress:    sql.NullString{String: "", Valid: false},
			UserAgent:    sql.NullString{String: "", Valid: false},
		},
		&UserSessionRaw{
			RefreshToken: refreshToken,
		},
		nil
}

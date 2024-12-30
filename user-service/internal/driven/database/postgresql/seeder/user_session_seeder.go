package seeder

import (
	"database/sql"
	"errors"
	"time"

	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"go.uber.org/fx"
)

var (
	expiredAt, _    = time.Parse(time.RFC3339, "2025-01-30T03:02:30Z07:00")
	TestUserSession = model.UserSession{
		UserUUID:     "868b606b-26d5-4c8d-ba45-9587919e059f",
		RefreshToken: "S1QVIM/yacA4Rw56O8BRqnjPZZ7kawuHhiS9mK0+XPY",
		IsRevoked:    false,
		ExpiredAt:    expiredAt,
		IpAddress:    sql.NullString{String: "", Valid: false},
		UserAgent:    sql.NullString{String: "", Valid: false},
	}
)

type UserSessionSeeder interface {
	CreateUserSession(user model.User) (*model.UserSession, error)
	DeleteUserSession(user model.User) error
}

type UserSessionSeederParam struct {
	fx.In

	PostgresqlErrorTranslator database.PostgresqlErrorTranslator
	Database                  *database.Database
	UserSessionStorage        shared_service.UserSessionStorage
}

func NewUserSessionSeeder(p UserSessionSeederParam) UserSessionSeeder {
	return &userSessionSeederImpl{
		translator:         p.PostgresqlErrorTranslator,
		database:           p.Database,
		userSessionStorage: p.UserSessionStorage,
	}
}

type userSessionSeederImpl struct {
	translator         database.PostgresqlErrorTranslator
	database           *database.Database
	userSessionStorage shared_service.UserSessionStorage
}

func (s *userSessionSeederImpl) CreateUserSession(user model.User) (*model.UserSession, error) {
	var notFoundError *database.RecordNotFoundError

	newUserSession := model.UserSession{}
	err := s.database.DB.Unscoped().Where(model.UserSession{UserUUID: user.UUID.String()}).Assign(TestUserSession).FirstOrCreate(&newUserSession).Error
	err = s.translator.Translate(err)
	if err != nil && !errors.As(err, &notFoundError) {
		return nil, err
	}

	return &newUserSession, nil
}

func (s *userSessionSeederImpl) DeleteUserSession(user model.User) error {
	err := s.database.DB.Unscoped().Where(model.UserSession{UserUUID: user.UUID.String()}).Delete(&model.UserSession{}).Error
	err = s.translator.Translate(err)
	if err != nil {
		var terr *database.RecordNotFoundError
		if errors.As(err, &terr) {
			return nil
		}
		return err
	}
	return nil
}

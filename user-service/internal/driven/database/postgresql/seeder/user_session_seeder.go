package seeder

import (
	"errors"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"go.uber.org/fx"
)

type UserSessionSeeder interface {
	SaveUserSession(session model.UserSession) (*model.UserSession, error)
	DeleteAllUserSessionsByUUID(uuidString string) error
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

func (s *userSessionSeederImpl) SaveUserSession(session model.UserSession) (*model.UserSession, error) {
	var notFoundError *database.RecordNotFoundError

	newUserSession := model.UserSession{}
	err := s.database.DB.Unscoped().Where(model.UserSession{UserUUID: session.UserUUID}).Assign(session).FirstOrCreate(&newUserSession).Error
	err = s.translator.Translate(err)
	if err != nil && !errors.As(err, &notFoundError) {
		return nil, err
	}

	return &newUserSession, nil
}

func (s *userSessionSeederImpl) DeleteAllUserSessionsByUUID(uuidString string) error {
	err := s.database.DB.Unscoped().Where(model.UserSession{UserUUID: uuidString}).Delete(&model.UserSession{}).Error
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

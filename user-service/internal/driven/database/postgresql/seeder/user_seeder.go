package seeder

import (
	"errors"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UserSeeder interface {
	SaveUser(user model.User) (*model.User, error)
	DeleteUser(username string) error
}

type UserSeederParam struct {
	fx.In

	PostgresqlErrorTranslator database.PostgresqlErrorTranslator
	Database                  *database.Database
	UserStorage               shared_service.UserStorage
	UserKeySeeder             UserKeySeeder
	UserSessionSeeder         UserSessionSeeder
}

func NewUserSeeder(p UserSeederParam) UserSeeder {
	return &userSeederImpl{
		translator:        p.PostgresqlErrorTranslator,
		database:          p.Database,
		userStorage:       p.UserStorage,
		userKeySeeder:     p.UserKeySeeder,
		userSessionSeeder: p.UserSessionSeeder,
	}
}

type userSeederImpl struct {
	translator        database.PostgresqlErrorTranslator
	database          *database.Database
	userStorage       shared_service.UserStorage
	userKeySeeder     UserKeySeeder
	userSessionSeeder UserSessionSeeder
}

// SaveUser saves a user to the database
// TODO: Including the user's keys and sessions, if provided
// This function is idempotent, meaning that it will not return an error if the user already exists
func (s *userSeederImpl) SaveUser(user model.User) (*model.User, error) {
	var terr *database.DuplicatedKeyError

	newUser := model.User{}

	err := s.database.DB.Unscoped().Where(model.User{Username: user.Username}).Assign(user).FirstOrCreate(&newUser).Error
	err = s.translator.Translate(err)
	if err != nil && !errors.As(err, &terr) {
		return &newUser, err
	}

	err = s.database.DB.Model(&newUser).Update("deleted_at", gorm.Expr("NULL")).Error
	err = s.translator.Translate(err)
	if err != nil && !errors.As(err, &terr) {
		return &newUser, err
	}

	if _, err := s.userKeySeeder.CreateUserKey(user); err != nil {
		return &newUser, err
	}

	return &newUser, nil
}

// DeleteUser deletes a user from the database
// Including the user's keys and sessions
// This function is idempotent, meaning that it will not return an error if the user does not exist
func (s *userSeederImpl) DeleteUser(username string) error {
	user := model.User{}

	err := s.database.DB.Unscoped().Where(&model.User{Username: username}).First(&user).Error
	err = s.translator.Translate(err)
	if err != nil {
		var terr *database.RecordNotFoundError
		if errors.As(err, &terr) {
			return nil
		}
		return err
	}

	err = s.userKeySeeder.DeleteUserKey(user)
	if err != nil {
		return err
	}

	err = s.userSessionSeeder.DeleteAllUserSessionsByUUID(user.UUID.String())
	if err != nil {
		return err
	}

	err = s.database.DB.Unscoped().Delete(&user).Error
	err = s.translator.Translate(err)
	if err != nil {
		return err
	}
	return nil
}

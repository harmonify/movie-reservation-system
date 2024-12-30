package seeder

import (
	"errors"
	"time"

	"github.com/google/uuid"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	TestUser = model.User{
		UUID:                  uuid.MustParse("868b606b-26d5-4c8d-ba45-9587919e059f"),
		Username:              "user1234",
		Password:              "user1234",
		Email:                 "user1234@example.com",
		PhoneNumber:           "+6281234567890",
		FirstName:             "Example",
		LastName:              "User",
		IsEmailVerified:       true,
		IsPhoneNumberVerified: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             gorm.DeletedAt{Time: time.Time{}, Valid: false},
	}
)

type UserSeeder interface {
	CreateUser(user model.User) (*model.User, error)
	DeleteUser(uuidString string) error
	CreateTestUser() (*model.User, error)
	DeleteTestUser() error
}

type UserSeederParam struct {
	fx.In

	PostgresqlErrorTranslator database.PostgresqlErrorTranslator
	Database                  *database.Database
	UserStorage               shared_service.UserStorage
	UserKeySeeder             UserKeySeeder
}

func NewUserSeeder(p UserSeederParam) UserSeeder {
	return &userSeederImpl{
		translator:    p.PostgresqlErrorTranslator,
		database:      p.Database,
		userStorage:   p.UserStorage,
		userKeySeeder: p.UserKeySeeder,
	}
}

type userSeederImpl struct {
	translator    database.PostgresqlErrorTranslator
	database      *database.Database
	userStorage   shared_service.UserStorage
	userKeySeeder UserKeySeeder
}

func (s *userSeederImpl) CreateTestUser() (*model.User, error) {
	return s.CreateUser(TestUser)
}

func (s *userSeederImpl) CreateUser(user model.User) (*model.User, error) {
	var terr *database.DuplicatedKeyError

	newUser := model.User{}

	err := s.database.DB.Unscoped().Where(model.User{UUID: user.UUID}).Assign(user).FirstOrCreate(&newUser).Error
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

func (s *userSeederImpl) DeleteTestUser() error {
	return s.DeleteUser(TestUser.UUID.String())
}

func (s *userSeederImpl) DeleteUser(uuidString string) error {
	user := model.User{}

	err := s.database.DB.Unscoped().Where(&model.User{UUID: uuid.MustParse(uuidString)}).First(&user).Error
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

	err = s.database.DB.Unscoped().Delete(&user).Error
	err = s.translator.Translate(err)
	if err != nil {
		return err
	}
	return nil
}

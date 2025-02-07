package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/database"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestUserRepository(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}
	os.Setenv("ENV", config_pkg.EnvironmentTest)
	suite.Run(t, new(UserRepositoryTestSuite))
}

type (
	saveUserTestConfig struct {
		Data *entity.SaveUser
	}
	saveUserTestExpectation struct {
		Data  *entity.User
		Error error
	}

	findUserTestConfig struct {
		Data *entity.User
		Find *entity.FindUser
	}
	findUserTestExpectation struct {
		Data  *entity.User
		Error error
	}

	updateUserTestConfig struct {
		Find   *entity.FindUser
		Update *entity.UpdateUser
	}
	updateUserTestExpectation struct {
		Data  *entity.User
		Error error
	}

	softDeleteUserTestConfig struct {
		Find entity.FindUser
	}
	softDeleteUserTestExpectation struct {
		Error error
	}
)

type UserRepositoryTestSuite struct {
	suite.Suite
	app         *fx.App
	db          *database.Database
	userSeeder  seeder.UserSeeder
	userStorage shared.UserStorage
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		entityfactory.UserEntityFactoryModule,
		seeder.DrivenPostgresqlSeederModule,
		fx.Invoke(func(
			db *database.Database,
			userSeeder seeder.UserSeeder,
			userStorage shared.UserStorage,
		) {
			s.db = db
			s.userStorage = userStorage
			s.userSeeder = userSeeder
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *UserRepositoryTestSuite) TestUserRepository_Basic() {
	var newUser *entity.User
	var err error

	ctx := context.Background()

	userSaveConfig := entity.SaveUser{
		TraceID:     uuid.New().String(),
		Username:    "saveuser123",
		Password:    "saveuser123",
		Email:       "saveuser123@example.com",
		PhoneNumber: "+6281230599852",
		FirstName:   "Save",
		LastName:    "User123",
	}

	userUpdateConfig := entity.UpdateUser{
		Username:              sql.NullString{String: "saveuser456", Valid: true},
		Password:              sql.NullString{String: "saveuser456", Valid: true},
		Email:                 sql.NullString{String: "saveuser456@example.com", Valid: true},
		PhoneNumber:           sql.NullString{String: "+6281230599853", Valid: true},
		FirstName:             sql.NullString{String: "SaveUpdated", Valid: true},
		LastName:              sql.NullString{String: "User456", Valid: true},
		IsEmailVerified:       sql.NullBool{Bool: true, Valid: true},
		IsPhoneNumberVerified: sql.NullBool{Bool: true, Valid: true},
	}

	_ = s.db.DB.Unscoped().Delete(&model.User{}, "username IN (?)", []string{userSaveConfig.Username, userUpdateConfig.Username.String})

	userBeforeSave, err := s.userStorage.FindUser(ctx, entity.FindUser{Username: sql.NullString{String: userSaveConfig.Username, Valid: true}})
	s.Require().Nil(userBeforeSave)
	s.Require().Error(err)
	s.Require().ErrorAs(err, &database.RecordNotFoundError{})

	newUser, err = s.userStorage.SaveUser(ctx, userSaveConfig)
	s.Require().Nil(err)
	s.Assert().NotEmpty(newUser.UUID)
	s.Assert().Equal(userSaveConfig.Username, newUser.Username)
	s.Assert().Equal(userSaveConfig.Password, newUser.Password)
	s.Assert().Equal(userSaveConfig.Email, newUser.Email)
	s.Assert().Equal(userSaveConfig.PhoneNumber, newUser.PhoneNumber)
	s.Assert().Equal(userSaveConfig.FirstName, newUser.FirstName)
	s.Assert().Equal(userSaveConfig.LastName, newUser.LastName)
	s.Assert().Equal(false, newUser.IsEmailVerified)
	s.Assert().Equal(false, newUser.IsPhoneNumberVerified)
	now := time.Now().Unix()
	s.Assert().GreaterOrEqual(newUser.CreatedAt.Unix(), now)
	s.Assert().GreaterOrEqual(newUser.UpdatedAt.Unix(), now)
	s.Assert().Empty(newUser.DeletedAt.Time)
	s.Require().False(newUser.DeletedAt.Valid)

	userAfterSave, err := s.userStorage.FindUser(ctx, entity.FindUser{Username: sql.NullString{String: userSaveConfig.Username, Valid: true}})
	s.Require().Nil(err)
	s.Assert().Equal(newUser.Username, userAfterSave.Username)
	s.Assert().Equal(newUser.Password, userAfterSave.Password)
	s.Assert().Equal(newUser.Email, userAfterSave.Email)
	s.Assert().Equal(newUser.PhoneNumber, userAfterSave.PhoneNumber)
	s.Assert().Equal(newUser.FirstName, userAfterSave.FirstName)
	s.Assert().Equal(newUser.LastName, userAfterSave.LastName)
	s.Assert().Equal(newUser.IsEmailVerified, userAfterSave.IsEmailVerified)
	s.Assert().Equal(newUser.IsPhoneNumberVerified, userAfterSave.IsPhoneNumberVerified)
	s.Assert().Equal(newUser.CreatedAt, userAfterSave.CreatedAt)
	s.Assert().Equal(newUser.UpdatedAt, userAfterSave.UpdatedAt)
	s.Assert().Equal(newUser.DeletedAt, userAfterSave.DeletedAt)

	updatedUser, err := s.userStorage.UpdateUser(
		ctx,
		entity.FindUser{Username: sql.NullString{String: userSaveConfig.Username, Valid: true}},
		userUpdateConfig,
	)
	s.Require().Nil(err)
	s.Assert().Equal(newUser.UUID, updatedUser.UUID)
	s.Assert().Equal(userUpdateConfig.Username.String, updatedUser.Username)
	s.Assert().Equal(userUpdateConfig.Password.String, updatedUser.Password)
	s.Assert().Equal(userUpdateConfig.Email.String, updatedUser.Email)
	s.Assert().Equal(userUpdateConfig.PhoneNumber.String, updatedUser.PhoneNumber)
	s.Assert().Equal(userUpdateConfig.FirstName.String, updatedUser.FirstName)
	s.Assert().Equal(userUpdateConfig.LastName.String, updatedUser.LastName)
	s.Assert().Equal(userUpdateConfig.IsEmailVerified.Bool, updatedUser.IsEmailVerified)
	s.Assert().Equal(userUpdateConfig.IsPhoneNumberVerified.Bool, updatedUser.IsPhoneNumberVerified)
	s.Assert().Equal(newUser.CreatedAt, updatedUser.CreatedAt)
	s.Assert().Greater(updatedUser.UpdatedAt, newUser.UpdatedAt)
	s.Assert().Equal(newUser.DeletedAt, updatedUser.DeletedAt)

	userAfterUpdate, err := s.userStorage.FindUser(ctx, entity.FindUser{Username: sql.NullString{String: userUpdateConfig.Username.String, Valid: true}})
	s.Require().Nil(err)
	s.Assert().Equal(updatedUser.UUID, userAfterUpdate.UUID)
	s.Assert().Equal(updatedUser.Username, userAfterUpdate.Username)
	s.Assert().Equal(updatedUser.Password, userAfterUpdate.Password)
	s.Assert().Equal(updatedUser.Email, userAfterUpdate.Email)
	s.Assert().Equal(updatedUser.PhoneNumber, userAfterUpdate.PhoneNumber)
	s.Assert().Equal(updatedUser.FirstName, userAfterUpdate.FirstName)
	s.Assert().Equal(updatedUser.LastName, userAfterUpdate.LastName)
	s.Assert().Equal(updatedUser.IsEmailVerified, userAfterUpdate.IsEmailVerified)
	s.Assert().Equal(updatedUser.IsPhoneNumberVerified, userAfterUpdate.IsPhoneNumberVerified)
	s.Assert().Equal(updatedUser.CreatedAt, userAfterUpdate.CreatedAt)
	s.Assert().Equal(updatedUser.UpdatedAt, userAfterUpdate.UpdatedAt)
	s.Assert().Equal(updatedUser.DeletedAt, userAfterUpdate.DeletedAt)

	err = s.userStorage.SoftDeleteUser(ctx, entity.FindUser{Username: sql.NullString{String: userUpdateConfig.Username.String, Valid: true}})
	s.Require().Nil(err)

	userAfterDelete, err := s.userStorage.FindUser(ctx, entity.FindUser{Username: sql.NullString{String: userUpdateConfig.Username.String, Valid: true}})
	s.Require().Nil(userAfterDelete)
	s.Require().Error(err)
	s.Require().ErrorAs(err, &database.RecordNotFoundError{})

	// Teardown
	db := s.db.DB.Unscoped()
	err = s.userStorage.
		WithTx(&database.Transaction{DB: db}).
		SoftDeleteUser(ctx, entity.FindUser{UUID: sql.NullString{String: newUser.UUID, Valid: true}})
	s.Require().Nil(err, "Failed to teardown test user")
}

func (s *UserRepositoryTestSuite) TestUserRepository_FindUser() {
	testUser, err := s.userSeeder.CreateUser(context.Background())
	s.Require().Nil(err, "Failed to setup test user")
	defer func() {
		err := s.userSeeder.DeleteUser(context.Background(), entity.FindUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}})
		s.Require().Nil(err, "Failed to teardown test user")
	}()

	testCases := []test_interface.TestCase[findUserTestConfig, findUserTestExpectation]{
		{
			Description: "Should be able to find user by UUID",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					UUID: sql.NullString{
						String: testUser.User.UUID,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  testUser.User,
			},
		},
		{
			Description: "Should be able to find user by email",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					Email: sql.NullString{
						String: testUser.User.Email,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  testUser.User,
			},
		},
		{
			Description: "Should be able to find user by username",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					Username: sql.NullString{
						String: testUser.User.Username,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  testUser.User,
			},
		},
		{
			Description: "Should be able to find user by phone number",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					PhoneNumber: sql.NullString{
						String: testUser.User.PhoneNumber,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  testUser.User,
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			user, err := s.userStorage.FindUser(ctx, *testCase.Config.Find)

			s.Require().Equal(testCase.Expectation.Error, err)

			s.Assert().Equal(testCase.Expectation.Data.UUID, user.UUID)
			s.Assert().Equal(testCase.Expectation.Data.Username, user.Username)
			s.Assert().Equal(testCase.Expectation.Data.Email, user.Email)
			s.Assert().Equal(testCase.Expectation.Data.PhoneNumber, user.PhoneNumber)
			s.Assert().Equal(testCase.Expectation.Data.FirstName, user.FirstName)
			s.Assert().Equal(testCase.Expectation.Data.LastName, user.LastName)
			s.Assert().Equal(testCase.Expectation.Data.IsEmailVerified, user.IsEmailVerified)
			s.Assert().Equal(testCase.Expectation.Data.IsPhoneNumberVerified, user.IsPhoneNumberVerified)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}
		})
	}
}

package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/factory"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestUserRepository(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

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
		Find *entity.FindUser
	}
	softDeleteUserTestExpectation struct {
		Error error
	}
)

type UserRepositoryTestSuite struct {
	suite.Suite
	app         *fx.App
	userStorage shared.UserStorage
	testUser    *model.User
	userSeeder  seeder.UserSeeder
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		factory.DrivenPostgresqlFactoryModule,
		seeder.DrivenPostgresqlSeederModule,
		fx.Invoke(func(
			userStorage shared.UserStorage,
			userFactory factory.UserFactory,
			userSeeder seeder.UserSeeder,
		) {
			s.userStorage = userStorage
			s.testUser = userFactory.CreateTestUser(factory.CreateTestUserParam{HashPassword: true})
			s.userSeeder = userSeeder
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *UserRepositoryTestSuite) SetupSubTest() {
	_, err := s.userSeeder.SaveUser(*s.testUser)
	s.Require().Nil(err, "Failed to setup test user")
}

func (s *UserRepositoryTestSuite) TearDownSubTest() {
	err := s.userSeeder.DeleteUser(s.testUser.Username)
	s.Require().Nil(err, "Failed to teardown test user")
}

func (s *UserRepositoryTestSuite) TestUserRepository_SaveUser() {
	var newUser *entity.User
	var err error

	testCases := []test_interface.TestCase[saveUserTestConfig, saveUserTestExpectation]{
		{
			Description: "Should be able to find user by UUID",
			Config: saveUserTestConfig{
				Data: &entity.SaveUser{
					Username:    "saveuser123",
					Password:    "saveuser123",
					Email:       "saveuser123@example.com",
					PhoneNumber: "+6281230599852",
					FirstName:   "Save",
					LastName:    "User123",
				},
			},
			Expectation: saveUserTestExpectation{
				Error: nil,
				Data: &entity.User{
					Username:    "saveuser123",
					Password:    "saveuser123",
					Email:       "saveuser123@example.com",
					PhoneNumber: "+6281230599852",
					FirstName:   "Save",
					LastName:    "User123",
				},
			},
			AfterCall: func() {
				err := s.userSeeder.DeleteUser(newUser.Username)
				s.Require().Nil(err)
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			newUser, err = s.userStorage.SaveUser(ctx, *testCase.Config.Data)

			s.Require().Equal(testCase.Expectation.Error, err)

			s.Assert().NotEmpty(newUser.UUID.String())
			s.Assert().Equal(testCase.Expectation.Data.Username, newUser.Username)
			s.Assert().NotEmpty(newUser.Password)
			s.Assert().Equal(testCase.Expectation.Data.Email, newUser.Email)
			s.Assert().Equal(testCase.Expectation.Data.PhoneNumber, newUser.PhoneNumber)
			s.Assert().Equal(testCase.Expectation.Data.FirstName, newUser.FirstName)
			s.Assert().Equal(testCase.Expectation.Data.LastName, newUser.LastName)
			s.Assert().Equal(false, newUser.IsEmailVerified)
			s.Assert().Equal(false, newUser.IsPhoneNumberVerified)

			now := time.Now().Unix()
			s.Assert().GreaterOrEqual(newUser.CreatedAt.Unix(), now)
			s.Assert().GreaterOrEqual(newUser.UpdatedAt.Unix(), now)
			s.Assert().Empty(newUser.DeletedAt.Time)
			s.Assert().False(newUser.DeletedAt.Valid)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestUserRepository_FindUser() {
	testCases := []test_interface.TestCase[findUserTestConfig, findUserTestExpectation]{
		{
			Description: "Should be able to find user by UUID",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					UUID: sql.NullString{
						String: s.testUser.UUID.String(),
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  s.testUser.ToEntity(),
			},
		},
		{
			Description: "Should be able to find user by email",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					Email: sql.NullString{
						String: s.testUser.Email,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  s.testUser.ToEntity(),
			},
		},
		{
			Description: "Should be able to find user by username",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					Username: sql.NullString{
						String: s.testUser.Username,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  s.testUser.ToEntity(),
			},
		},
		{
			Description: "Should be able to find user by phone number",
			Config: findUserTestConfig{
				Find: &entity.FindUser{
					PhoneNumber: sql.NullString{
						String: s.testUser.PhoneNumber,
						Valid:  true,
					},
				},
			},
			Expectation: findUserTestExpectation{
				Error: nil,
				Data:  s.testUser.ToEntity(),
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

func (s *UserRepositoryTestSuite) TestUserRepository_UpdateUser() {
	expectedUser1 := s.testUser.ToEntity()
	expectedUser1.FirstName = "Testing"

	expectedUser2 := s.testUser.ToEntity()
	expectedUser2.Email = "test2@example.com"

	expectedUser3 := s.testUser.ToEntity()
	expectedUser3.PhoneNumber = "+62812345678922"
	expectedUser3.LastName = "User3"
	expectedUser3.IsEmailVerified = true
	expectedUser3.IsPhoneNumberVerified = true

	testCases := []test_interface.TestCase[updateUserTestConfig, updateUserTestExpectation]{
		{
			Description: "Should be able to update any field by UUID",
			Config: updateUserTestConfig{
				Find: &entity.FindUser{
					UUID: sql.NullString{
						String: s.testUser.UUID.String(),
						Valid:  true,
					},
				},
				Update: &entity.UpdateUser{
					FirstName: sql.NullString{String: "Testing", Valid: true},
				},
			},
			Expectation: updateUserTestExpectation{
				Error: nil,
				Data:  expectedUser1,
			},
		},
		{
			Description: "Should be able to update any field by email",
			Config: updateUserTestConfig{
				Find: &entity.FindUser{
					Email: sql.NullString{
						String: s.testUser.Email,
						Valid:  true,
					},
				},
				Update: &entity.UpdateUser{
					Email: sql.NullString{String: "test2@example.com", Valid: true},
				},
			},
			Expectation: updateUserTestExpectation{
				Error: nil,
				Data:  expectedUser2,
			},
		},
		{
			Description: "Should be able to update any field by username",
			Config: updateUserTestConfig{
				Find: &entity.FindUser{
					Username: sql.NullString{
						String: s.testUser.Username,
						Valid:  true,
					},
				},
				Update: &entity.UpdateUser{
					PhoneNumber:           sql.NullString{String: "+62812345678922", Valid: true},
					LastName:              sql.NullString{String: "User3", Valid: true},
					IsEmailVerified:       sql.NullBool{Bool: true, Valid: true},
					IsPhoneNumberVerified: sql.NullBool{Bool: true, Valid: true},
				},
			},
			Expectation: updateUserTestExpectation{
				Error: nil,
				Data:  expectedUser3,
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			updated, err := s.userStorage.UpdateUser(ctx, *testCase.Config.Find, *testCase.Config.Update)

			s.Require().Equal(testCase.Expectation.Error, err)

			s.Assert().Equal(testCase.Expectation.Data.UUID, updated.UUID)
			s.Assert().Equal(testCase.Expectation.Data.Username, updated.Username)
			s.Assert().Equal(testCase.Expectation.Data.Email, updated.Email)
			s.Assert().Equal(testCase.Expectation.Data.PhoneNumber, updated.PhoneNumber)
			s.Assert().Equal(testCase.Expectation.Data.FirstName, updated.FirstName)
			s.Assert().Equal(testCase.Expectation.Data.LastName, updated.LastName)
			s.Assert().Equal(testCase.Expectation.Data.IsEmailVerified, updated.IsEmailVerified)
			s.Assert().Equal(testCase.Expectation.Data.IsPhoneNumberVerified, updated.IsPhoneNumberVerified)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestUserRepository_SoftDeleteUser() {
	testCases := []test_interface.TestCase[softDeleteUserTestConfig, softDeleteUserTestExpectation]{
		{
			Description: "Should be able to delete with UUID",
			Config: softDeleteUserTestConfig{
				Find: &entity.FindUser{
					UUID: sql.NullString{
						String: s.testUser.UUID.String(),
						Valid:  true,
					},
				},
			},
			Expectation: softDeleteUserTestExpectation{
				Error: nil,
			},
		},
		{
			Description: "Should be able to delete with email",
			Config: softDeleteUserTestConfig{
				Find: &entity.FindUser{
					Email: sql.NullString{
						String: s.testUser.Email,
						Valid:  true,
					},
				},
			},
			Expectation: softDeleteUserTestExpectation{
				Error: nil,
			},
		},
		{
			Description: "Should be able to delete with username",
			Config: softDeleteUserTestConfig{
				Find: &entity.FindUser{
					Username: sql.NullString{
						String: s.testUser.Username,
						Valid:  true,
					},
				},
			},
			Expectation: softDeleteUserTestExpectation{
				Error: nil,
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			err := s.userStorage.SoftDeleteUser(ctx, *testCase.Config.Find)
			s.Require().Equal(testCase.Expectation.Error, err)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}
		})
	}
}

package seeder

import (
	"context"
	"database/sql"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"go.uber.org/fx"
)

type UserWithRelations struct {
	User            *entity.User
	UserRaw         *entityfactory.UserRaw
	UserKey         *entity.UserKey
	UserSessions    []*entity.UserSession
	UserSessionRaws []*entityfactory.UserSessionRaw
}

type UserSeeder interface {
	CreateUser(ctx context.Context) (*UserWithRelations, error)
	DeleteUser(ctx context.Context, findModel entity.FindUser) error
}

type UserSeederParam struct {
	fx.In

	PostgresqlErrorTranslator database.PostgresqlErrorTranslator
	Database                  *database.Database
	UserFactory               entityfactory.UserFactory
	UserKeyFactory            entityfactory.UserKeyFactory
	UserSessionFactory        entityfactory.UserSessionFactory
	UserStorage               shared.UserStorage
	UserKeyStorage            shared.UserKeyStorage
	UserSessionStorage        shared.UserSessionStorage
}

type userSeederImpl struct {
	translator         database.PostgresqlErrorTranslator
	database           *database.Database
	userFactory        entityfactory.UserFactory
	userKeyFactory     entityfactory.UserKeyFactory
	userSessionFactory entityfactory.UserSessionFactory
	userStorage        shared.UserStorage
	userKeyStorage     shared.UserKeyStorage
	userSessionStorage shared.UserSessionStorage
}

func NewUserSeeder(p UserSeederParam) UserSeeder {
	return &userSeederImpl{
		translator:         p.PostgresqlErrorTranslator,
		database:           p.Database,
		userFactory:        p.UserFactory,
		userKeyFactory:     p.UserKeyFactory,
		userSessionFactory: p.UserSessionFactory,
		userStorage:        p.UserStorage,
		userKeyStorage:     p.UserKeyStorage,
		userSessionStorage: p.UserSessionStorage,
	}
}

// CreateUser creates a user in the database
func (s *userSeederImpl) CreateUser(ctx context.Context) (*UserWithRelations, error) {
	user, userRaw, err := s.userFactory.GenerateUser()
	if err != nil {
		return nil, err
	}
	newUser, err := s.userStorage.SaveUser(ctx, entity.SaveUser{
		Username:    user.Username,
		TraceID:     user.TraceID,
		Password:    user.Password,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	})
	if err != nil {
		return nil, err
	}

	userKey := s.userKeyFactory.GenerateUserKey(user)
	newUserKey, err := s.userKeyStorage.SaveUserKey(ctx, entity.SaveUserKey{
		UserUUID:   newUser.UUID,
		PublicKey:  userKey.PublicKey,
		PrivateKey: userKey.PrivateKey,
	})
	if err != nil {
		return nil, err
	}

	userSession, userSessionRaw, err := s.userSessionFactory.GenerateUserSession(user)
	if err != nil {
		return nil, err
	}

	nu, err := s.userSessionStorage.SaveSession(ctx, entity.SaveUserSession{
		UserUUID:     newUser.UUID,
		TraceID:      userSession.TraceID,
		RefreshToken: userSession.RefreshToken,
		ExpiredAt:    userSession.ExpiredAt,
		IpAddress:    userSession.IpAddress,
		UserAgent:    userSession.UserAgent,
	})
	if err != nil {
		return nil, err
	}
	newUserSessions := []*entity.UserSession{nu}
	newUserSessionRaws := []*entityfactory.UserSessionRaw{userSessionRaw}

	return &UserWithRelations{
		User:            newUser,
		UserRaw:         userRaw,
		UserKey:         newUserKey,
		UserSessions:    newUserSessions,
		UserSessionRaws: newUserSessionRaws,
	}, nil
}

// DeleteUser deletes a user from the database
// Including the user's keys and sessions
// This function is idempotent, meaning that it will not return an error if the user does not exist
func (s *userSeederImpl) DeleteUser(ctx context.Context, findModel entity.FindUser) error {
	user, err := s.userStorage.FindUser(ctx, findModel)
	if err != nil {
		return err
	}

	db := s.database.DB.Unscoped() // Make it hard delete

	err = s.userSessionStorage.WithTx(&database.Transaction{DB: db}).SoftDeleteSession(ctx, entity.FindUserSession{
		UserUUID: sql.NullString{String: user.UUID, Valid: true},
	})
	if err != nil {
		return err
	}

	err = s.userKeyStorage.WithTx(&database.Transaction{DB: db}).SoftDeleteUserKey(ctx, entity.FindUserKey{
		UserUUID: sql.NullString{String: user.UUID, Valid: true},
	})
	if err != nil {
		return err
	}

	err = s.userStorage.WithTx(&database.Transaction{DB: db}).SoftDeleteUser(ctx, entity.FindUser{
		UUID: sql.NullString{String: user.UUID, Valid: true},
	})
	return err
}

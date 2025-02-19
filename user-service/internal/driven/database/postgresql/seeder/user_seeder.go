package seeder

import (
	"context"
	"database/sql"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
	entityseeder "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/seeder"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

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

func NewUserSeeder(p UserSeederParam) entityseeder.UserSeeder {
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
func (s *userSeederImpl) CreateUser(ctx context.Context) (*entityseeder.UserWithRelations, error) {
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

	return &entityseeder.UserWithRelations{
		User:            newUser,
		UserRaw:         userRaw,
		UserKey:         newUserKey,
		UserSessions:    newUserSessions,
		UserSessionRaws: newUserSessionRaws,
	}, nil
}

// CreateAdmin creates an admin user in the database
func (s *userSeederImpl) CreateAdmin(ctx context.Context) (*entityseeder.UserWithRelations, error) {
	user, userRaw, err := s.userFactory.GenerateUserV2()
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

	return &entityseeder.UserWithRelations{
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
func (s *userSeederImpl) DeleteUser(ctx context.Context, GetModel entity.GetUser) error {
	user, err := s.userStorage.GetUser(ctx, GetModel)
	if err != nil {
		return err
	}

	// Make it hard delete
	err = s.database.DB.Unscoped().Transaction(func(_tx *gorm.DB) error {
		tx := database.NewTransaction(_tx)

		err = s.userSessionStorage.WithTx(tx).SoftDeleteSession(ctx, entity.GetUserSession{
			UserUUID: sql.NullString{String: user.UUID, Valid: true},
		})
		if err != nil {
			return err
		}

		err = s.userKeyStorage.WithTx(tx).SoftDeleteUserKey(ctx, entity.GetUserKey{
			UserUUID: sql.NullString{String: user.UUID, Valid: true},
		})
		if err != nil {
			return err
		}

		err = s.userStorage.WithTx(tx).SoftDeleteUser(ctx, entity.GetUser{
			UUID: sql.NullString{String: user.UUID, Valid: true},
		})

		return err
	})

	return err
}

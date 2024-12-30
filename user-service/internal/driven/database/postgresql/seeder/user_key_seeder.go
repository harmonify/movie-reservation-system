package seeder

import (
	"errors"

	"github.com/google/uuid"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"go.uber.org/fx"
)

var (
	TestUserKey = model.UserKey{
		UserUUID:   uuid.MustParse("868b606b-26d5-4c8d-ba45-9587919e059f"),
		PublicKey:  "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJDZ0tDQVFFQW1KeDNpQmR1akNrMzdDS0hQb3ZHRm1hdnNuTUFMRjNJUUNpUnFSbVRNcEhHNmkzMTd1R3gKbGFtRUFQbTd0Zld1VnVueTJ2cnVjSktFWFNsaWhWMnJma0R2dlRVdjNyVXFwZEJiNlphRmhhcDZuamNSM1o0bwo5dGJyQ1Z2MXVPcG42WGQvRVMvVmdkUjREb01EYW9NNm5NUVAwWlBwTmdGTjcxVy95b1BPcEtjUlo0c2txcWJTCmMvVFQxZ2ZBR2lPQlJmSGwyZHdKOXZvK0ZQcktYSHQvbFVNSDBoNXpmV1JIalNRWEJhdU1odFlSMDBOMHA5T2UKcXh5R2tkYk94UVdlYWlMenRJdGZNUkJzTHJDQU5MN21zUTBHQ3EwTXh6UUhTVEF1ZmMxbUhMMCtoa0F3RHBaNQpDdHBLUGpHRjU0V1VHdWVPa2RGLzllS1J3YTQwWVhFajhRSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
		PrivateKey: "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.Iy1+Zux/wvz//kT2xW4fjbkawdwt9X4YJlpjpkCZ7oo=.a1VorD0gFYz9YL3D.15000",
	}
)

type UserKeySeeder interface {
	CreateUserKey(user model.User) (*model.UserKey, error)
	DeleteUserKey(user model.User) error
}

type UserKeySeederParam struct {
	fx.In

	PostgresqlErrorTranslator database.PostgresqlErrorTranslator
	Database                  *database.Database
	UserKeyStorage            shared_service.UserKeyStorage
}

func NewUserKeySeeder(p UserKeySeederParam) UserKeySeeder {
	return &userKeySeederImpl{
		translator:     p.PostgresqlErrorTranslator,
		database:       p.Database,
		userKeyStorage: p.UserKeyStorage,
	}
}

type userKeySeederImpl struct {
	translator     database.PostgresqlErrorTranslator
	database       *database.Database
	userKeyStorage shared_service.UserKeyStorage
}

func (s *userKeySeederImpl) CreateUserKey(user model.User) (*model.UserKey, error) {
	var notFoundError *database.RecordNotFoundError

	newUserKey := model.UserKey{}
	err := s.database.DB.Unscoped().Where(model.UserKey{UserUUID: user.UUID}).Assign(user).FirstOrCreate(&newUserKey).Error
	err = s.translator.Translate(err)
	if err != nil && !errors.As(err, &notFoundError) {
		return nil, err
	}

	return &newUserKey, nil
}

func (s *userKeySeederImpl) DeleteUserKey(user model.User) error {
	err := s.database.DB.Unscoped().Where(&model.UserKey{UserUUID: user.UUID}).Delete(&model.UserKey{}).Error
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

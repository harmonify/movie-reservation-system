package user_service

import (
	"context"
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UserService interface {
	GetUser(ctx context.Context, p GetUserParam) (*GetUserResult, error)
	UpdateUser(ctx context.Context, p UpdateUserParam) (*UpdateUserResult, error)
}

type UserServiceParam struct {
	fx.In
	logger.Logger
	tracer.Tracer
	shared.UserStorage
	shared.UserRoleStorage
}

type UserServiceResult struct {
	fx.Out

	UserService UserService
}

type userServiceImpl struct {
	logger          logger.Logger
	tracer          tracer.Tracer
	userStorage     shared.UserStorage
	userRoleStorage shared.UserRoleStorage
}

func NewUserService(p UserServiceParam) UserServiceResult {
	return UserServiceResult{
		UserService: &userServiceImpl{
			logger:          p.Logger,
			tracer:          p.Tracer,
			userStorage:     p.UserStorage,
			userRoleStorage: p.UserRoleStorage,
		},
	}
}

func (s *userServiceImpl) GetUser(ctx context.Context, p GetUserParam) (*GetUserResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	user, err := s.userStorage.GetUser(ctx, entity.GetUser{
		UUID: sql.NullString{String: p.UUID, Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to get user", zap.Error(err))
		return nil, err
	}

	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	roles, err := s.userRoleStorage.SearchUserRoles(ctx, entity.SearchUserRoles{
		UserUUID: p.UUID,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to search user roles", zap.Error(err))
		return nil, err
	}

	return &GetUserResult{
		UUID:                  p.UUID,
		Username:              user.Username,
		Email:                 user.Email,
		PhoneNumber:           user.PhoneNumber,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		IsEmailVerified:       user.IsEmailVerified,
		IsPhoneNumberVerified: user.IsPhoneNumberVerified,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		DeletedAt:             deletedAt,
		Roles:                 roles,
	}, err

}

func (s *userServiceImpl) UpdateUser(ctx context.Context, p UpdateUserParam) (*UpdateUserResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	user, err := s.userStorage.UpdateUser(
		ctx,
		entity.GetUser{
			UUID: sql.NullString{String: p.UUID, Valid: true},
		},
		entity.UpdateUser{
			Email:                 sql.NullString{String: p.Email, Valid: p.Email != ""},
			PhoneNumber:           sql.NullString{String: p.PhoneNumber, Valid: p.PhoneNumber != ""},
			Username:              sql.NullString{String: p.Username, Valid: p.Username != ""},
			FirstName:             sql.NullString{String: p.FirstName, Valid: p.FirstName != ""},
			LastName:              sql.NullString{String: p.LastName, Valid: p.LastName != ""},
			IsEmailVerified:       sql.NullBool{Bool: false, Valid: p.Email != ""},
			IsPhoneNumberVerified: sql.NullBool{Bool: false, Valid: p.PhoneNumber != ""},
		},
	)
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to update user", zap.Error(err))
		return nil, err
	}

	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return &UpdateUserResult{
		UUID:                  user.UUID,
		Username:              user.Username,
		Email:                 user.Email,
		PhoneNumber:           user.PhoneNumber,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		IsEmailVerified:       user.IsEmailVerified,
		IsPhoneNumberVerified: user.IsPhoneNumberVerified,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		DeletedAt:             deletedAt,
	}, nil
}

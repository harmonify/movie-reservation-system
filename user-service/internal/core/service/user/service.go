package user_service

import (
	"context"
	"database/sql"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"go.uber.org/fx"
)

type UserService interface {
	GetUser(ctx context.Context, p GetUserParam) (*GetUserResult, error)
	UpdateUser(ctx context.Context, p UpdateUserParam) (*UpdateUserResult, error)
	GetUpdateEmailVerification(ctx context.Context, p GetUpdateEmailVerificationParam) error
	VerifyUpdateEmail(ctx context.Context, p VerifyUpdateEmailParam) error
	GetUpdatePhoneNumberVerification(ctx context.Context, p GetUpdatePhoneNumberVerificationParam) error
	VerifyUpdatePhoneNumber(ctx context.Context, p VerifyUpdatePhoneNumberParam) error
}

type UserServiceParam struct {
	fx.In

	Logger      logger.Logger
	Tracer      tracer.Tracer
	UserStorage shared_service.UserStorage
}

type UserServiceResult struct {
	fx.Out

	UserService UserService
}

type userServiceImpl struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	userStorage shared_service.UserStorage
}

func NewUserService(p UserServiceParam) UserServiceResult {
	return UserServiceResult{
		UserService: &userServiceImpl{
			logger:      p.Logger,
			tracer:      p.Tracer,
			userStorage: p.UserStorage,
		},
	}
}

func (s *userServiceImpl) GetUser(ctx context.Context, p GetUserParam) (*GetUserResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	user, err := s.userStorage.FindUser(ctx, entity.FindUser{
		UUID: sql.NullString{String: p.UUID, Valid: true},
	})

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
		DeletedAt:             user.DeletedAt,
	}, err

}

func (s *userServiceImpl) UpdateUser(ctx context.Context, p UpdateUserParam) (*UpdateUserResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	user, err := s.userStorage.UpdateUser(
		ctx,
		entity.FindUser{
			UUID: sql.NullString{String: p.UUID, Valid: true},
		},
		entity.UpdateUser{
			Username:  sql.NullString{String: p.Username, Valid: true},
			FirstName: sql.NullString{String: p.FirstName, Valid: true},
			LastName:  sql.NullString{String: p.LastName, Valid: true},
		},
	)

	return &UpdateUserResult{
		UUID:                  user.UUID.String(),
		Username:              user.Username,
		Email:                 user.Email,
		PhoneNumber:           user.PhoneNumber,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		IsEmailVerified:       user.IsEmailVerified,
		IsPhoneNumberVerified: user.IsPhoneNumberVerified,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		DeletedAt:             user.DeletedAt,
	}, err
}

func (s *userServiceImpl) GetUpdateEmailVerification(ctx context.Context, p GetUpdateEmailVerificationParam) error {
	panic("unimplemented")
}

func (s *userServiceImpl) VerifyUpdateEmail(ctx context.Context, p VerifyUpdateEmailParam) error {
	panic("unimplemented")
}

func (s *userServiceImpl) GetUpdatePhoneNumberVerification(ctx context.Context, p GetUpdatePhoneNumberVerificationParam) error {
	panic("unimplemented")
}

func (s *userServiceImpl) VerifyUpdatePhoneNumber(ctx context.Context, p VerifyUpdatePhoneNumberParam) error {
	panic("unimplemented")
}

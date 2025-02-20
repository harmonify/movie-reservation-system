package entityseeder

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
)

type UserWithRelations struct {
	User            *entity.User
	UserRaw         *entityfactory.UserRaw
	UserRoles       []*entity.UserRole
	UserKey         *entity.UserKey
	UserSessions    []*entity.UserSession
	UserSessionRaws []*entityfactory.UserSessionRaw
}

type UserSeeder interface {
	CreateUser(ctx context.Context) (*UserWithRelations, error)
	CreateAdmin(ctx context.Context) (*UserWithRelations, error)
	DeleteUser(ctx context.Context, GetModel entity.GetUser) error
}

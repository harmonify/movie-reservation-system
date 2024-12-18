package user_service

import "context"

type UserService interface {
	GetUser(ctx context.Context)
	EditUser(ctx context.Context)
}

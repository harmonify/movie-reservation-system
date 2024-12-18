package auth_service

import (
	"context"

	"go.uber.org/fx"
)

type (
	AuthService interface {
		Register(ctx context.Context, p RegisterParam) (LoginResult, error)
		Login(ctx context.Context, p LoginParam) (LoginResult, error)
		Logout(ctx context.Context, p LogoutParam) error
	}

	AuthServiceParam struct {
		fx.In
	}

	AuthServiceResult struct {
		fx.Out
	}

	authServiceImpl struct {
	}
)

func NewAuthService() AuthService {
	return &authServiceImpl{}
}

func (a *authServiceImpl) Register(ctx context.Context, p RegisterParam) (LoginResult, error) {
	panic("unimplemented")
}

func (a *authServiceImpl) Login(ctx context.Context, p LoginParam) (LoginResult, error) {
	panic("unimplemented")
}

func (a *authServiceImpl) Logout(ctx context.Context, p LogoutParam) error {
	panic("unimplemented")
}

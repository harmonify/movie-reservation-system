// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	token_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/token"
	mock "github.com/stretchr/testify/mock"
)

// TokenService is an autogenerated mock type for the TokenService type
type TokenService struct {
	mock.Mock
}

type TokenService_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenService) EXPECT() *TokenService_Expecter {
	return &TokenService_Expecter{mock: &_m.Mock}
}

// GenerateAccessToken provides a mock function with given fields: ctx, p
func (_m *TokenService) GenerateAccessToken(ctx context.Context, p token_service.GenerateAccessTokenParam) (*token_service.GenerateAccessTokenResult, error) {
	ret := _m.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAccessToken")
	}

	var r0 *token_service.GenerateAccessTokenResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, token_service.GenerateAccessTokenParam) (*token_service.GenerateAccessTokenResult, error)); ok {
		return rf(ctx, p)
	}
	if rf, ok := ret.Get(0).(func(context.Context, token_service.GenerateAccessTokenParam) *token_service.GenerateAccessTokenResult); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token_service.GenerateAccessTokenResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, token_service.GenerateAccessTokenParam) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenService_GenerateAccessToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateAccessToken'
type TokenService_GenerateAccessToken_Call struct {
	*mock.Call
}

// GenerateAccessToken is a helper method to define mock.On call
//   - ctx context.Context
//   - p token_service.GenerateAccessTokenParam
func (_e *TokenService_Expecter) GenerateAccessToken(ctx interface{}, p interface{}) *TokenService_GenerateAccessToken_Call {
	return &TokenService_GenerateAccessToken_Call{Call: _e.mock.On("GenerateAccessToken", ctx, p)}
}

func (_c *TokenService_GenerateAccessToken_Call) Run(run func(ctx context.Context, p token_service.GenerateAccessTokenParam)) *TokenService_GenerateAccessToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(token_service.GenerateAccessTokenParam))
	})
	return _c
}

func (_c *TokenService_GenerateAccessToken_Call) Return(_a0 *token_service.GenerateAccessTokenResult, _a1 error) *TokenService_GenerateAccessToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenService_GenerateAccessToken_Call) RunAndReturn(run func(context.Context, token_service.GenerateAccessTokenParam) (*token_service.GenerateAccessTokenResult, error)) *TokenService_GenerateAccessToken_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateRefreshToken provides a mock function with given fields: ctx
func (_m *TokenService) GenerateRefreshToken(ctx context.Context) (*token_service.GenerateRefreshTokenResult, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GenerateRefreshToken")
	}

	var r0 *token_service.GenerateRefreshTokenResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*token_service.GenerateRefreshTokenResult, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *token_service.GenerateRefreshTokenResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token_service.GenerateRefreshTokenResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenService_GenerateRefreshToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateRefreshToken'
type TokenService_GenerateRefreshToken_Call struct {
	*mock.Call
}

// GenerateRefreshToken is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TokenService_Expecter) GenerateRefreshToken(ctx interface{}) *TokenService_GenerateRefreshToken_Call {
	return &TokenService_GenerateRefreshToken_Call{Call: _e.mock.On("GenerateRefreshToken", ctx)}
}

func (_c *TokenService_GenerateRefreshToken_Call) Run(run func(ctx context.Context)) *TokenService_GenerateRefreshToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TokenService_GenerateRefreshToken_Call) Return(_a0 *token_service.GenerateRefreshTokenResult, _a1 error) *TokenService_GenerateRefreshToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenService_GenerateRefreshToken_Call) RunAndReturn(run func(context.Context) (*token_service.GenerateRefreshTokenResult, error)) *TokenService_GenerateRefreshToken_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateUserKey provides a mock function with given fields: ctx
func (_m *TokenService) GenerateUserKey(ctx context.Context) (*token_service.GenerateUserKeyResult, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GenerateUserKey")
	}

	var r0 *token_service.GenerateUserKeyResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*token_service.GenerateUserKeyResult, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *token_service.GenerateUserKeyResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token_service.GenerateUserKeyResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenService_GenerateUserKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateUserKey'
type TokenService_GenerateUserKey_Call struct {
	*mock.Call
}

// GenerateUserKey is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TokenService_Expecter) GenerateUserKey(ctx interface{}) *TokenService_GenerateUserKey_Call {
	return &TokenService_GenerateUserKey_Call{Call: _e.mock.On("GenerateUserKey", ctx)}
}

func (_c *TokenService_GenerateUserKey_Call) Run(run func(ctx context.Context)) *TokenService_GenerateUserKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TokenService_GenerateUserKey_Call) Return(_a0 *token_service.GenerateUserKeyResult, _a1 error) *TokenService_GenerateUserKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenService_GenerateUserKey_Call) RunAndReturn(run func(context.Context) (*token_service.GenerateUserKeyResult, error)) *TokenService_GenerateUserKey_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyRefreshToken provides a mock function with given fields: ctx, p
func (_m *TokenService) VerifyRefreshToken(ctx context.Context, p token_service.VerifyRefreshTokenParam) (*token_service.VerifyRefreshTokenResult, error) {
	ret := _m.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for VerifyRefreshToken")
	}

	var r0 *token_service.VerifyRefreshTokenResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, token_service.VerifyRefreshTokenParam) (*token_service.VerifyRefreshTokenResult, error)); ok {
		return rf(ctx, p)
	}
	if rf, ok := ret.Get(0).(func(context.Context, token_service.VerifyRefreshTokenParam) *token_service.VerifyRefreshTokenResult); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token_service.VerifyRefreshTokenResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, token_service.VerifyRefreshTokenParam) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenService_VerifyRefreshToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyRefreshToken'
type TokenService_VerifyRefreshToken_Call struct {
	*mock.Call
}

// VerifyRefreshToken is a helper method to define mock.On call
//   - ctx context.Context
//   - p token_service.VerifyRefreshTokenParam
func (_e *TokenService_Expecter) VerifyRefreshToken(ctx interface{}, p interface{}) *TokenService_VerifyRefreshToken_Call {
	return &TokenService_VerifyRefreshToken_Call{Call: _e.mock.On("VerifyRefreshToken", ctx, p)}
}

func (_c *TokenService_VerifyRefreshToken_Call) Run(run func(ctx context.Context, p token_service.VerifyRefreshTokenParam)) *TokenService_VerifyRefreshToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(token_service.VerifyRefreshTokenParam))
	})
	return _c
}

func (_c *TokenService_VerifyRefreshToken_Call) Return(_a0 *token_service.VerifyRefreshTokenResult, _a1 error) *TokenService_VerifyRefreshToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenService_VerifyRefreshToken_Call) RunAndReturn(run func(context.Context, token_service.VerifyRefreshTokenParam) (*token_service.VerifyRefreshTokenResult, error)) *TokenService_VerifyRefreshToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewTokenService creates a new instance of TokenService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenService {
	mock := &TokenService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

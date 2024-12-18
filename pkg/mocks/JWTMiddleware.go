// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// JWTMiddleware is an autogenerated mock type for the JWTMiddleware type
type JWTMiddleware struct {
	mock.Mock
}

type JWTMiddleware_Expecter struct {
	mock *mock.Mock
}

func (_m *JWTMiddleware) EXPECT() *JWTMiddleware_Expecter {
	return &JWTMiddleware_Expecter{mock: &_m.Mock}
}

// AuthenticateUser provides a mock function with given fields: c
func (_m *JWTMiddleware) AuthenticateUser(c *gin.Context) {
	_m.Called(c)
}

// JWTMiddleware_AuthenticateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AuthenticateUser'
type JWTMiddleware_AuthenticateUser_Call struct {
	*mock.Call
}

// AuthenticateUser is a helper method to define mock.On call
//   - c *gin.Context
func (_e *JWTMiddleware_Expecter) AuthenticateUser(c interface{}) *JWTMiddleware_AuthenticateUser_Call {
	return &JWTMiddleware_AuthenticateUser_Call{Call: _e.mock.On("AuthenticateUser", c)}
}

func (_c *JWTMiddleware_AuthenticateUser_Call) Run(run func(c *gin.Context)) *JWTMiddleware_AuthenticateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *JWTMiddleware_AuthenticateUser_Call) Return() *JWTMiddleware_AuthenticateUser_Call {
	_c.Call.Return()
	return _c
}

func (_c *JWTMiddleware_AuthenticateUser_Call) RunAndReturn(run func(*gin.Context)) *JWTMiddleware_AuthenticateUser_Call {
	_c.Call.Return(run)
	return _c
}

// OptAuthenticateUser provides a mock function with given fields: c
func (_m *JWTMiddleware) OptAuthenticateUser(c *gin.Context) {
	_m.Called(c)
}

// JWTMiddleware_OptAuthenticateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OptAuthenticateUser'
type JWTMiddleware_OptAuthenticateUser_Call struct {
	*mock.Call
}

// OptAuthenticateUser is a helper method to define mock.On call
//   - c *gin.Context
func (_e *JWTMiddleware_Expecter) OptAuthenticateUser(c interface{}) *JWTMiddleware_OptAuthenticateUser_Call {
	return &JWTMiddleware_OptAuthenticateUser_Call{Call: _e.mock.On("OptAuthenticateUser", c)}
}

func (_c *JWTMiddleware_OptAuthenticateUser_Call) Run(run func(c *gin.Context)) *JWTMiddleware_OptAuthenticateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *JWTMiddleware_OptAuthenticateUser_Call) Return() *JWTMiddleware_OptAuthenticateUser_Call {
	_c.Call.Return()
	return _c
}

func (_c *JWTMiddleware_OptAuthenticateUser_Call) RunAndReturn(run func(*gin.Context)) *JWTMiddleware_OptAuthenticateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewJWTMiddleware creates a new instance of JWTMiddleware. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJWTMiddleware(t interface {
	mock.TestingT
	Cleanup(func())
}) *JWTMiddleware {
	mock := &JWTMiddleware{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

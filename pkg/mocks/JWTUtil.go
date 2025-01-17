// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	mock "github.com/stretchr/testify/mock"
)

// JwtUtil is an autogenerated mock type for the JwtUtil type
type JwtUtil struct {
	mock.Mock
}

type JwtUtil_Expecter struct {
	mock *mock.Mock
}

func (_m *JwtUtil) EXPECT() *JwtUtil_Expecter {
	return &JwtUtil_Expecter{mock: &_m.Mock}
}

// JWTSign provides a mock function with given fields: payload
func (_m *JwtUtil) JWTSign(payload jwt_util.JWTSignParam) (string, error) {
	ret := _m.Called(payload)

	if len(ret) == 0 {
		panic("no return value specified for JWTSign")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(jwt_util.JWTSignParam) (string, error)); ok {
		return rf(payload)
	}
	if rf, ok := ret.Get(0).(func(jwt_util.JWTSignParam) string); ok {
		r0 = rf(payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(jwt_util.JWTSignParam) error); ok {
		r1 = rf(payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// JwtUtil_JWTSign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'JWTSign'
type JwtUtil_JWTSign_Call struct {
	*mock.Call
}

// JWTSign is a helper method to define mock.On call
//   - payload jwt_util.JWTSignParam
func (_e *JwtUtil_Expecter) JWTSign(payload interface{}) *JwtUtil_JWTSign_Call {
	return &JwtUtil_JWTSign_Call{Call: _e.mock.On("JWTSign", payload)}
}

func (_c *JwtUtil_JWTSign_Call) Run(run func(payload jwt_util.JWTSignParam)) *JwtUtil_JWTSign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(jwt_util.JWTSignParam))
	})
	return _c
}

func (_c *JwtUtil_JWTSign_Call) Return(_a0 string, _a1 error) *JwtUtil_JWTSign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *JwtUtil_JWTSign_Call) RunAndReturn(run func(jwt_util.JWTSignParam) (string, error)) *JwtUtil_JWTSign_Call {
	_c.Call.Return(run)
	return _c
}

// JWTVerify provides a mock function with given fields: token
func (_m *JwtUtil) JWTVerify(token string) (*jwt_util.JWTBodyPayload, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for JWTVerify")
	}

	var r0 *jwt_util.JWTBodyPayload
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*jwt_util.JWTBodyPayload, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *jwt_util.JWTBodyPayload); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt_util.JWTBodyPayload)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// JwtUtil_JWTVerify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'JWTVerify'
type JwtUtil_JWTVerify_Call struct {
	*mock.Call
}

// JWTVerify is a helper method to define mock.On call
//   - token string
func (_e *JwtUtil_Expecter) JWTVerify(token interface{}) *JwtUtil_JWTVerify_Call {
	return &JwtUtil_JWTVerify_Call{Call: _e.mock.On("JWTVerify", token)}
}

func (_c *JwtUtil_JWTVerify_Call) Run(run func(token string)) *JwtUtil_JWTVerify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *JwtUtil_JWTVerify_Call) Return(_a0 *jwt_util.JWTBodyPayload, _a1 error) *JwtUtil_JWTVerify_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *JwtUtil_JWTVerify_Call) RunAndReturn(run func(string) (*jwt_util.JWTBodyPayload, error)) *JwtUtil_JWTVerify_Call {
	_c.Call.Return(run)
	return _c
}

// NewJwtUtil creates a new instance of JwtUtil. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJwtUtil(t interface {
	mock.TestingT
	Cleanup(func())
}) *JwtUtil {
	mock := &JwtUtil{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

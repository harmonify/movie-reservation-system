// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	entity "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"

	mock "github.com/stretchr/testify/mock"
)

// UserSessionFactory is an autogenerated mock type for the UserSessionFactory type
type UserSessionFactory struct {
	mock.Mock
}

type UserSessionFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *UserSessionFactory) EXPECT() *UserSessionFactory_Expecter {
	return &UserSessionFactory_Expecter{mock: &_m.Mock}
}

// GenerateUserSession provides a mock function with given fields: user
func (_m *UserSessionFactory) GenerateUserSession(user *entity.User) (*entity.UserSession, *entityfactory.UserSessionRaw, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for GenerateUserSession")
	}

	var r0 *entity.UserSession
	var r1 *entityfactory.UserSessionRaw
	var r2 error
	if rf, ok := ret.Get(0).(func(*entity.User) (*entity.UserSession, *entityfactory.UserSessionRaw, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*entity.User) *entity.UserSession); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.UserSession)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.User) *entityfactory.UserSessionRaw); ok {
		r1 = rf(user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entityfactory.UserSessionRaw)
		}
	}

	if rf, ok := ret.Get(2).(func(*entity.User) error); ok {
		r2 = rf(user)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UserSessionFactory_GenerateUserSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateUserSession'
type UserSessionFactory_GenerateUserSession_Call struct {
	*mock.Call
}

// GenerateUserSession is a helper method to define mock.On call
//   - user *entity.User
func (_e *UserSessionFactory_Expecter) GenerateUserSession(user interface{}) *UserSessionFactory_GenerateUserSession_Call {
	return &UserSessionFactory_GenerateUserSession_Call{Call: _e.mock.On("GenerateUserSession", user)}
}

func (_c *UserSessionFactory_GenerateUserSession_Call) Run(run func(user *entity.User)) *UserSessionFactory_GenerateUserSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entity.User))
	})
	return _c
}

func (_c *UserSessionFactory_GenerateUserSession_Call) Return(session *entity.UserSession, raw *entityfactory.UserSessionRaw, err error) *UserSessionFactory_GenerateUserSession_Call {
	_c.Call.Return(session, raw, err)
	return _c
}

func (_c *UserSessionFactory_GenerateUserSession_Call) RunAndReturn(run func(*entity.User) (*entity.UserSession, *entityfactory.UserSessionRaw, error)) *UserSessionFactory_GenerateUserSession_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserSessionFactory creates a new instance of UserSessionFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserSessionFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserSessionFactory {
	mock := &UserSessionFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

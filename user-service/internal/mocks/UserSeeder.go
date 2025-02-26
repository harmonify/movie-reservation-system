// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityseeder "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/seeder"

	mock "github.com/stretchr/testify/mock"
)

// UserSeeder is an autogenerated mock type for the UserSeeder type
type UserSeeder struct {
	mock.Mock
}

type UserSeeder_Expecter struct {
	mock *mock.Mock
}

func (_m *UserSeeder) EXPECT() *UserSeeder_Expecter {
	return &UserSeeder_Expecter{mock: &_m.Mock}
}

// CreateAdmin provides a mock function with given fields: ctx, username
func (_m *UserSeeder) CreateAdmin(ctx context.Context, username string) (*entityseeder.UserWithRelations, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for CreateAdmin")
	}

	var r0 *entityseeder.UserWithRelations
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entityseeder.UserWithRelations, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entityseeder.UserWithRelations); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entityseeder.UserWithRelations)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserSeeder_CreateAdmin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAdmin'
type UserSeeder_CreateAdmin_Call struct {
	*mock.Call
}

// CreateAdmin is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *UserSeeder_Expecter) CreateAdmin(ctx interface{}, username interface{}) *UserSeeder_CreateAdmin_Call {
	return &UserSeeder_CreateAdmin_Call{Call: _e.mock.On("CreateAdmin", ctx, username)}
}

func (_c *UserSeeder_CreateAdmin_Call) Run(run func(ctx context.Context, username string)) *UserSeeder_CreateAdmin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserSeeder_CreateAdmin_Call) Return(_a0 *entityseeder.UserWithRelations, _a1 error) *UserSeeder_CreateAdmin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserSeeder_CreateAdmin_Call) RunAndReturn(run func(context.Context, string) (*entityseeder.UserWithRelations, error)) *UserSeeder_CreateAdmin_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: ctx
func (_m *UserSeeder) CreateUser(ctx context.Context) (*entityseeder.UserWithRelations, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *entityseeder.UserWithRelations
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*entityseeder.UserWithRelations, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *entityseeder.UserWithRelations); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entityseeder.UserWithRelations)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserSeeder_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type UserSeeder_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserSeeder_Expecter) CreateUser(ctx interface{}) *UserSeeder_CreateUser_Call {
	return &UserSeeder_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx)}
}

func (_c *UserSeeder_CreateUser_Call) Run(run func(ctx context.Context)) *UserSeeder_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserSeeder_CreateUser_Call) Return(_a0 *entityseeder.UserWithRelations, _a1 error) *UserSeeder_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserSeeder_CreateUser_Call) RunAndReturn(run func(context.Context) (*entityseeder.UserWithRelations, error)) *UserSeeder_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUser provides a mock function with given fields: ctx, getModel
func (_m *UserSeeder) DeleteUser(ctx context.Context, getModel entity.GetUser) error {
	ret := _m.Called(ctx, getModel)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.GetUser) error); ok {
		r0 = rf(ctx, getModel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserSeeder_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type UserSeeder_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - ctx context.Context
//   - getModel entity.GetUser
func (_e *UserSeeder_Expecter) DeleteUser(ctx interface{}, getModel interface{}) *UserSeeder_DeleteUser_Call {
	return &UserSeeder_DeleteUser_Call{Call: _e.mock.On("DeleteUser", ctx, getModel)}
}

func (_c *UserSeeder_DeleteUser_Call) Run(run func(ctx context.Context, getModel entity.GetUser)) *UserSeeder_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.GetUser))
	})
	return _c
}

func (_c *UserSeeder_DeleteUser_Call) Return(_a0 error) *UserSeeder_DeleteUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserSeeder_DeleteUser_Call) RunAndReturn(run func(context.Context, entity.GetUser) error) *UserSeeder_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserSeeder creates a new instance of UserSeeder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserSeeder(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserSeeder {
	mock := &UserSeeder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

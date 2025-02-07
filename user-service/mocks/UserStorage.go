// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	database "github.com/harmonify/movie-reservation-system/pkg/database"
	entity "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"

	mock "github.com/stretchr/testify/mock"

	shared "github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
)

// UserStorage is an autogenerated mock type for the UserStorage type
type UserStorage struct {
	mock.Mock
}

type UserStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *UserStorage) EXPECT() *UserStorage_Expecter {
	return &UserStorage_Expecter{mock: &_m.Mock}
}

// FindUser provides a mock function with given fields: ctx, findModel
func (_m *UserStorage) FindUser(ctx context.Context, findModel entity.FindUser) (*entity.User, error) {
	ret := _m.Called(ctx, findModel)

	if len(ret) == 0 {
		panic("no return value specified for FindUser")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.FindUser) (*entity.User, error)); ok {
		return rf(ctx, findModel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.FindUser) *entity.User); ok {
		r0 = rf(ctx, findModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.FindUser) error); ok {
		r1 = rf(ctx, findModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserStorage_FindUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUser'
type UserStorage_FindUser_Call struct {
	*mock.Call
}

// FindUser is a helper method to define mock.On call
//   - ctx context.Context
//   - findModel entity.FindUser
func (_e *UserStorage_Expecter) FindUser(ctx interface{}, findModel interface{}) *UserStorage_FindUser_Call {
	return &UserStorage_FindUser_Call{Call: _e.mock.On("FindUser", ctx, findModel)}
}

func (_c *UserStorage_FindUser_Call) Run(run func(ctx context.Context, findModel entity.FindUser)) *UserStorage_FindUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.FindUser))
	})
	return _c
}

func (_c *UserStorage_FindUser_Call) Return(_a0 *entity.User, _a1 error) *UserStorage_FindUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserStorage_FindUser_Call) RunAndReturn(run func(context.Context, entity.FindUser) (*entity.User, error)) *UserStorage_FindUser_Call {
	_c.Call.Return(run)
	return _c
}

// FindUserWithResult provides a mock function with given fields: ctx, findModel, resultModel
func (_m *UserStorage) FindUserWithResult(ctx context.Context, findModel entity.FindUser, resultModel interface{}) error {
	ret := _m.Called(ctx, findModel, resultModel)

	if len(ret) == 0 {
		panic("no return value specified for FindUserWithResult")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.FindUser, interface{}) error); ok {
		r0 = rf(ctx, findModel, resultModel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserStorage_FindUserWithResult_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserWithResult'
type UserStorage_FindUserWithResult_Call struct {
	*mock.Call
}

// FindUserWithResult is a helper method to define mock.On call
//   - ctx context.Context
//   - findModel entity.FindUser
//   - resultModel interface{}
func (_e *UserStorage_Expecter) FindUserWithResult(ctx interface{}, findModel interface{}, resultModel interface{}) *UserStorage_FindUserWithResult_Call {
	return &UserStorage_FindUserWithResult_Call{Call: _e.mock.On("FindUserWithResult", ctx, findModel, resultModel)}
}

func (_c *UserStorage_FindUserWithResult_Call) Run(run func(ctx context.Context, findModel entity.FindUser, resultModel interface{})) *UserStorage_FindUserWithResult_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.FindUser), args[2].(interface{}))
	})
	return _c
}

func (_c *UserStorage_FindUserWithResult_Call) Return(_a0 error) *UserStorage_FindUserWithResult_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserStorage_FindUserWithResult_Call) RunAndReturn(run func(context.Context, entity.FindUser, interface{}) error) *UserStorage_FindUserWithResult_Call {
	_c.Call.Return(run)
	return _c
}

// SaveUser provides a mock function with given fields: ctx, createModel
func (_m *UserStorage) SaveUser(ctx context.Context, createModel entity.SaveUser) (*entity.User, error) {
	ret := _m.Called(ctx, createModel)

	if len(ret) == 0 {
		panic("no return value specified for SaveUser")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.SaveUser) (*entity.User, error)); ok {
		return rf(ctx, createModel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.SaveUser) *entity.User); ok {
		r0 = rf(ctx, createModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.SaveUser) error); ok {
		r1 = rf(ctx, createModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserStorage_SaveUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveUser'
type UserStorage_SaveUser_Call struct {
	*mock.Call
}

// SaveUser is a helper method to define mock.On call
//   - ctx context.Context
//   - createModel entity.SaveUser
func (_e *UserStorage_Expecter) SaveUser(ctx interface{}, createModel interface{}) *UserStorage_SaveUser_Call {
	return &UserStorage_SaveUser_Call{Call: _e.mock.On("SaveUser", ctx, createModel)}
}

func (_c *UserStorage_SaveUser_Call) Run(run func(ctx context.Context, createModel entity.SaveUser)) *UserStorage_SaveUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.SaveUser))
	})
	return _c
}

func (_c *UserStorage_SaveUser_Call) Return(_a0 *entity.User, _a1 error) *UserStorage_SaveUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserStorage_SaveUser_Call) RunAndReturn(run func(context.Context, entity.SaveUser) (*entity.User, error)) *UserStorage_SaveUser_Call {
	_c.Call.Return(run)
	return _c
}

// SoftDeleteUser provides a mock function with given fields: ctx, findModel
func (_m *UserStorage) SoftDeleteUser(ctx context.Context, findModel entity.FindUser) error {
	ret := _m.Called(ctx, findModel)

	if len(ret) == 0 {
		panic("no return value specified for SoftDeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.FindUser) error); ok {
		r0 = rf(ctx, findModel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserStorage_SoftDeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SoftDeleteUser'
type UserStorage_SoftDeleteUser_Call struct {
	*mock.Call
}

// SoftDeleteUser is a helper method to define mock.On call
//   - ctx context.Context
//   - findModel entity.FindUser
func (_e *UserStorage_Expecter) SoftDeleteUser(ctx interface{}, findModel interface{}) *UserStorage_SoftDeleteUser_Call {
	return &UserStorage_SoftDeleteUser_Call{Call: _e.mock.On("SoftDeleteUser", ctx, findModel)}
}

func (_c *UserStorage_SoftDeleteUser_Call) Run(run func(ctx context.Context, findModel entity.FindUser)) *UserStorage_SoftDeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.FindUser))
	})
	return _c
}

func (_c *UserStorage_SoftDeleteUser_Call) Return(_a0 error) *UserStorage_SoftDeleteUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserStorage_SoftDeleteUser_Call) RunAndReturn(run func(context.Context, entity.FindUser) error) *UserStorage_SoftDeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: ctx, findModel, updateModel
func (_m *UserStorage) UpdateUser(ctx context.Context, findModel entity.FindUser, updateModel entity.UpdateUser) (*entity.User, error) {
	ret := _m.Called(ctx, findModel, updateModel)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.FindUser, entity.UpdateUser) (*entity.User, error)); ok {
		return rf(ctx, findModel, updateModel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.FindUser, entity.UpdateUser) *entity.User); ok {
		r0 = rf(ctx, findModel, updateModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.FindUser, entity.UpdateUser) error); ok {
		r1 = rf(ctx, findModel, updateModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserStorage_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type UserStorage_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - findModel entity.FindUser
//   - updateModel entity.UpdateUser
func (_e *UserStorage_Expecter) UpdateUser(ctx interface{}, findModel interface{}, updateModel interface{}) *UserStorage_UpdateUser_Call {
	return &UserStorage_UpdateUser_Call{Call: _e.mock.On("UpdateUser", ctx, findModel, updateModel)}
}

func (_c *UserStorage_UpdateUser_Call) Run(run func(ctx context.Context, findModel entity.FindUser, updateModel entity.UpdateUser)) *UserStorage_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.FindUser), args[2].(entity.UpdateUser))
	})
	return _c
}

func (_c *UserStorage_UpdateUser_Call) Return(_a0 *entity.User, _a1 error) *UserStorage_UpdateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserStorage_UpdateUser_Call) RunAndReturn(run func(context.Context, entity.FindUser, entity.UpdateUser) (*entity.User, error)) *UserStorage_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// WithTx provides a mock function with given fields: tx
func (_m *UserStorage) WithTx(tx *database.Transaction) shared.UserStorage {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 shared.UserStorage
	if rf, ok := ret.Get(0).(func(*database.Transaction) shared.UserStorage); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(shared.UserStorage)
		}
	}

	return r0
}

// UserStorage_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type UserStorage_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - tx *database.Transaction
func (_e *UserStorage_Expecter) WithTx(tx interface{}) *UserStorage_WithTx_Call {
	return &UserStorage_WithTx_Call{Call: _e.mock.On("WithTx", tx)}
}

func (_c *UserStorage_WithTx_Call) Run(run func(tx *database.Transaction)) *UserStorage_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*database.Transaction))
	})
	return _c
}

func (_c *UserStorage_WithTx_Call) Return(_a0 shared.UserStorage) *UserStorage_WithTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserStorage_WithTx_Call) RunAndReturn(run func(*database.Transaction) shared.UserStorage) *UserStorage_WithTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserStorage creates a new instance of UserStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserStorage {
	mock := &UserStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

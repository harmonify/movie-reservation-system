// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	database "github.com/harmonify/movie-reservation-system/pkg/database"
	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

type Storage_Expecter struct {
	mock *mock.Mock
}

func (_m *Storage) EXPECT() *Storage_Expecter {
	return &Storage_Expecter{mock: &_m.Mock}
}

// Transaction provides a mock function with given fields: fc
func (_m *Storage) Transaction(fc func(*database.Transaction) error) error {
	ret := _m.Called(fc)

	if len(ret) == 0 {
		panic("no return value specified for Transaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(func(*database.Transaction) error) error); ok {
		r0 = rf(fc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storage_Transaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Transaction'
type Storage_Transaction_Call struct {
	*mock.Call
}

// Transaction is a helper method to define mock.On call
//   - fc func(*database.Transaction) error
func (_e *Storage_Expecter) Transaction(fc interface{}) *Storage_Transaction_Call {
	return &Storage_Transaction_Call{Call: _e.mock.On("Transaction", fc)}
}

func (_c *Storage_Transaction_Call) Run(run func(fc func(*database.Transaction) error)) *Storage_Transaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(*database.Transaction) error))
	})
	return _c
}

func (_c *Storage_Transaction_Call) Return(_a0 error) *Storage_Transaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Storage_Transaction_Call) RunAndReturn(run func(func(*database.Transaction) error) error) *Storage_Transaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

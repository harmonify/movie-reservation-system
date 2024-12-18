// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// BcryptHash is an autogenerated mock type for the BcryptHash type
type BcryptHash struct {
	mock.Mock
}

type BcryptHash_Expecter struct {
	mock *mock.Mock
}

func (_m *BcryptHash) EXPECT() *BcryptHash_Expecter {
	return &BcryptHash_Expecter{mock: &_m.Mock}
}

// Compare provides a mock function with given fields: hashedValue, currValue
func (_m *BcryptHash) Compare(hashedValue string, currValue string) (bool, error) {
	ret := _m.Called(hashedValue, currValue)

	if len(ret) == 0 {
		panic("no return value specified for Compare")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (bool, error)); ok {
		return rf(hashedValue, currValue)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(hashedValue, currValue)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(hashedValue, currValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BcryptHash_Compare_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Compare'
type BcryptHash_Compare_Call struct {
	*mock.Call
}

// Compare is a helper method to define mock.On call
//   - hashedValue string
//   - currValue string
func (_e *BcryptHash_Expecter) Compare(hashedValue interface{}, currValue interface{}) *BcryptHash_Compare_Call {
	return &BcryptHash_Compare_Call{Call: _e.mock.On("Compare", hashedValue, currValue)}
}

func (_c *BcryptHash_Compare_Call) Run(run func(hashedValue string, currValue string)) *BcryptHash_Compare_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *BcryptHash_Compare_Call) Return(_a0 bool, _a1 error) *BcryptHash_Compare_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BcryptHash_Compare_Call) RunAndReturn(run func(string, string) (bool, error)) *BcryptHash_Compare_Call {
	_c.Call.Return(run)
	return _c
}

// Hash provides a mock function with given fields: value
func (_m *BcryptHash) Hash(value string) (string, error) {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(value)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BcryptHash_Hash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Hash'
type BcryptHash_Hash_Call struct {
	*mock.Call
}

// Hash is a helper method to define mock.On call
//   - value string
func (_e *BcryptHash_Expecter) Hash(value interface{}) *BcryptHash_Hash_Call {
	return &BcryptHash_Hash_Call{Call: _e.mock.On("Hash", value)}
}

func (_c *BcryptHash_Hash_Call) Run(run func(value string)) *BcryptHash_Hash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *BcryptHash_Hash_Call) Return(_a0 string, _a1 error) *BcryptHash_Hash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BcryptHash_Hash_Call) RunAndReturn(run func(string) (string, error)) *BcryptHash_Hash_Call {
	_c.Call.Return(run)
	return _c
}

// NewBcryptHash creates a new instance of BcryptHash. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBcryptHash(t interface {
	mock.TestingT
	Cleanup(func())
}) *BcryptHash {
	mock := &BcryptHash{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

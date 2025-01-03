// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ValidationUtil is an autogenerated mock type for the ValidationUtil type
type ValidationUtil struct {
	mock.Mock
}

type ValidationUtil_Expecter struct {
	mock *mock.Mock
}

func (_m *ValidationUtil) EXPECT() *ValidationUtil_Expecter {
	return &ValidationUtil_Expecter{mock: &_m.Mock}
}

// ValidateE164PhoneNumber provides a mock function with given fields: value
func (_m *ValidationUtil) ValidateE164PhoneNumber(value string) bool {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for ValidateE164PhoneNumber")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ValidationUtil_ValidateE164PhoneNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateE164PhoneNumber'
type ValidationUtil_ValidateE164PhoneNumber_Call struct {
	*mock.Call
}

// ValidateE164PhoneNumber is a helper method to define mock.On call
//   - value string
func (_e *ValidationUtil_Expecter) ValidateE164PhoneNumber(value interface{}) *ValidationUtil_ValidateE164PhoneNumber_Call {
	return &ValidationUtil_ValidateE164PhoneNumber_Call{Call: _e.mock.On("ValidateE164PhoneNumber", value)}
}

func (_c *ValidationUtil_ValidateE164PhoneNumber_Call) Run(run func(value string)) *ValidationUtil_ValidateE164PhoneNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ValidationUtil_ValidateE164PhoneNumber_Call) Return(_a0 bool) *ValidationUtil_ValidateE164PhoneNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ValidationUtil_ValidateE164PhoneNumber_Call) RunAndReturn(run func(string) bool) *ValidationUtil_ValidateE164PhoneNumber_Call {
	_c.Call.Return(run)
	return _c
}

// ValidatePhoneNumber provides a mock function with given fields: value
func (_m *ValidationUtil) ValidatePhoneNumber(value string) bool {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for ValidatePhoneNumber")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ValidationUtil_ValidatePhoneNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidatePhoneNumber'
type ValidationUtil_ValidatePhoneNumber_Call struct {
	*mock.Call
}

// ValidatePhoneNumber is a helper method to define mock.On call
//   - value string
func (_e *ValidationUtil_Expecter) ValidatePhoneNumber(value interface{}) *ValidationUtil_ValidatePhoneNumber_Call {
	return &ValidationUtil_ValidatePhoneNumber_Call{Call: _e.mock.On("ValidatePhoneNumber", value)}
}

func (_c *ValidationUtil_ValidatePhoneNumber_Call) Run(run func(value string)) *ValidationUtil_ValidatePhoneNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ValidationUtil_ValidatePhoneNumber_Call) Return(_a0 bool) *ValidationUtil_ValidatePhoneNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ValidationUtil_ValidatePhoneNumber_Call) RunAndReturn(run func(string) bool) *ValidationUtil_ValidatePhoneNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewValidationUtil creates a new instance of ValidationUtil. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValidationUtil(t interface {
	mock.TestingT
	Cleanup(func())
}) *ValidationUtil {
	mock := &ValidationUtil{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

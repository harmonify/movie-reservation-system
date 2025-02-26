// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StructValidator is an autogenerated mock type for the StructValidator type
type StructValidator struct {
	mock.Mock
}

type StructValidator_Expecter struct {
	mock *mock.Mock
}

func (_m *StructValidator) EXPECT() *StructValidator_Expecter {
	return &StructValidator_Expecter{mock: &_m.Mock}
}

// ConstructValidationErrorFields provides a mock function with given fields: err
func (_m *StructValidator) ConstructValidationErrorFields(err error) []error {
	ret := _m.Called(err)

	if len(ret) == 0 {
		panic("no return value specified for ConstructValidationErrorFields")
	}

	var r0 []error
	if rf, ok := ret.Get(0).(func(error) []error); ok {
		r0 = rf(err)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]error)
		}
	}

	return r0
}

// StructValidator_ConstructValidationErrorFields_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConstructValidationErrorFields'
type StructValidator_ConstructValidationErrorFields_Call struct {
	*mock.Call
}

// ConstructValidationErrorFields is a helper method to define mock.On call
//   - err error
func (_e *StructValidator_Expecter) ConstructValidationErrorFields(err interface{}) *StructValidator_ConstructValidationErrorFields_Call {
	return &StructValidator_ConstructValidationErrorFields_Call{Call: _e.mock.On("ConstructValidationErrorFields", err)}
}

func (_c *StructValidator_ConstructValidationErrorFields_Call) Run(run func(err error)) *StructValidator_ConstructValidationErrorFields_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error))
	})
	return _c
}

func (_c *StructValidator_ConstructValidationErrorFields_Call) Return(_a0 []error) *StructValidator_ConstructValidationErrorFields_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StructValidator_ConstructValidationErrorFields_Call) RunAndReturn(run func(error) []error) *StructValidator_ConstructValidationErrorFields_Call {
	_c.Call.Return(run)
	return _c
}

// Validate provides a mock function with given fields: schema
func (_m *StructValidator) Validate(schema interface{}) (error, []error) {
	ret := _m.Called(schema)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	var r1 []error
	if rf, ok := ret.Get(0).(func(interface{}) (error, []error)); ok {
		return rf(schema)
	}
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(interface{}) []error); ok {
		r1 = rf(schema)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]error)
		}
	}

	return r0, r1
}

// StructValidator_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type StructValidator_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - schema interface{}
func (_e *StructValidator_Expecter) Validate(schema interface{}) *StructValidator_Validate_Call {
	return &StructValidator_Validate_Call{Call: _e.mock.On("Validate", schema)}
}

func (_c *StructValidator_Validate_Call) Run(run func(schema interface{})) *StructValidator_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *StructValidator_Validate_Call) Return(original error, errorFields []error) *StructValidator_Validate_Call {
	_c.Call.Return(original, errorFields)
	return _c
}

func (_c *StructValidator_Validate_Call) RunAndReturn(run func(interface{}) (error, []error)) *StructValidator_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// NewStructValidator creates a new instance of StructValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStructValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *StructValidator {
	mock := &StructValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

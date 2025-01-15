// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// HttpValidator is an autogenerated mock type for the HttpValidator type
type HttpValidator struct {
	mock.Mock
}

type HttpValidator_Expecter struct {
	mock *mock.Mock
}

func (_m *HttpValidator) EXPECT() *HttpValidator_Expecter {
	return &HttpValidator_Expecter{mock: &_m.Mock}
}

// ValidateRequestBody provides a mock function with given fields: c, schema
func (_m *HttpValidator) ValidateRequestBody(c *gin.Context, schema interface{}) error {
	ret := _m.Called(c, schema)

	if len(ret) == 0 {
		panic("no return value specified for ValidateRequestBody")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gin.Context, interface{}) error); ok {
		r0 = rf(c, schema)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HttpValidator_ValidateRequestBody_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateRequestBody'
type HttpValidator_ValidateRequestBody_Call struct {
	*mock.Call
}

// ValidateRequestBody is a helper method to define mock.On call
//   - c *gin.Context
//   - schema interface{}
func (_e *HttpValidator_Expecter) ValidateRequestBody(c interface{}, schema interface{}) *HttpValidator_ValidateRequestBody_Call {
	return &HttpValidator_ValidateRequestBody_Call{Call: _e.mock.On("ValidateRequestBody", c, schema)}
}

func (_c *HttpValidator_ValidateRequestBody_Call) Run(run func(c *gin.Context, schema interface{})) *HttpValidator_ValidateRequestBody_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context), args[1].(interface{}))
	})
	return _c
}

func (_c *HttpValidator_ValidateRequestBody_Call) Return(_a0 error) *HttpValidator_ValidateRequestBody_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpValidator_ValidateRequestBody_Call) RunAndReturn(run func(*gin.Context, interface{}) error) *HttpValidator_ValidateRequestBody_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateRequestQuery provides a mock function with given fields: c, schema
func (_m *HttpValidator) ValidateRequestQuery(c *gin.Context, schema interface{}) error {
	ret := _m.Called(c, schema)

	if len(ret) == 0 {
		panic("no return value specified for ValidateRequestQuery")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gin.Context, interface{}) error); ok {
		r0 = rf(c, schema)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HttpValidator_ValidateRequestQuery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateRequestQuery'
type HttpValidator_ValidateRequestQuery_Call struct {
	*mock.Call
}

// ValidateRequestQuery is a helper method to define mock.On call
//   - c *gin.Context
//   - schema interface{}
func (_e *HttpValidator_Expecter) ValidateRequestQuery(c interface{}, schema interface{}) *HttpValidator_ValidateRequestQuery_Call {
	return &HttpValidator_ValidateRequestQuery_Call{Call: _e.mock.On("ValidateRequestQuery", c, schema)}
}

func (_c *HttpValidator_ValidateRequestQuery_Call) Run(run func(c *gin.Context, schema interface{})) *HttpValidator_ValidateRequestQuery_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context), args[1].(interface{}))
	})
	return _c
}

func (_c *HttpValidator_ValidateRequestQuery_Call) Return(_a0 error) *HttpValidator_ValidateRequestQuery_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpValidator_ValidateRequestQuery_Call) RunAndReturn(run func(*gin.Context, interface{}) error) *HttpValidator_ValidateRequestQuery_Call {
	_c.Call.Return(run)
	return _c
}

// NewHttpValidator creates a new instance of HttpValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHttpValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *HttpValidator {
	mock := &HttpValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
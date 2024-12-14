// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"

	response "github.com/harmonify/movie-reservation-system/pkg/http/response"
)

// Response is an autogenerated mock type for the Response type
type Response struct {
	mock.Mock
}

type Response_Expecter struct {
	mock *mock.Mock
}

func (_m *Response) EXPECT() *Response_Expecter {
	return &Response_Expecter{mock: &_m.Mock}
}

// Build provides a mock function with given fields: ctx, httpCode, data, err
func (_m *Response) Build(ctx context.Context, httpCode int, data interface{}, err error) (int, response.BaseResponseSchema, error) {
	ret := _m.Called(ctx, httpCode, data, err)

	if len(ret) == 0 {
		panic("no return value specified for Build")
	}

	var r0 int
	var r1 response.BaseResponseSchema
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int, interface{}, error) (int, response.BaseResponseSchema, error)); ok {
		return rf(ctx, httpCode, data, err)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, interface{}, error) int); ok {
		r0 = rf(ctx, httpCode, data, err)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, interface{}, error) response.BaseResponseSchema); ok {
		r1 = rf(ctx, httpCode, data, err)
	} else {
		r1 = ret.Get(1).(response.BaseResponseSchema)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int, interface{}, error) error); ok {
		r2 = rf(ctx, httpCode, data, err)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Response_Build_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Build'
type Response_Build_Call struct {
	*mock.Call
}

// Build is a helper method to define mock.On call
//   - ctx context.Context
//   - httpCode int
//   - data interface{}
//   - err error
func (_e *Response_Expecter) Build(ctx interface{}, httpCode interface{}, data interface{}, err interface{}) *Response_Build_Call {
	return &Response_Build_Call{Call: _e.mock.On("Build", ctx, httpCode, data, err)}
}

func (_c *Response_Build_Call) Run(run func(ctx context.Context, httpCode int, data interface{}, err error)) *Response_Build_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(interface{}), args[3].(error))
	})
	return _c
}

func (_c *Response_Build_Call) Return(_a0 int, _a1 response.BaseResponseSchema, _a2 error) *Response_Build_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *Response_Build_Call) RunAndReturn(run func(context.Context, int, interface{}, error) (int, response.BaseResponseSchema, error)) *Response_Build_Call {
	_c.Call.Return(run)
	return _c
}

// BuildError provides a mock function with given fields: code, err
func (_m *Response) BuildError(code string, err error) response.ErrorHandler {
	ret := _m.Called(code, err)

	if len(ret) == 0 {
		panic("no return value specified for BuildError")
	}

	var r0 response.ErrorHandler
	if rf, ok := ret.Get(0).(func(string, error) response.ErrorHandler); ok {
		r0 = rf(code, err)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(response.ErrorHandler)
		}
	}

	return r0
}

// Response_BuildError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BuildError'
type Response_BuildError_Call struct {
	*mock.Call
}

// BuildError is a helper method to define mock.On call
//   - code string
//   - err error
func (_e *Response_Expecter) BuildError(code interface{}, err interface{}) *Response_BuildError_Call {
	return &Response_BuildError_Call{Call: _e.mock.On("BuildError", code, err)}
}

func (_c *Response_BuildError_Call) Run(run func(code string, err error)) *Response_BuildError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(error))
	})
	return _c
}

func (_c *Response_BuildError_Call) Return(_a0 response.ErrorHandler) *Response_BuildError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_BuildError_Call) RunAndReturn(run func(string, error) response.ErrorHandler) *Response_BuildError_Call {
	_c.Call.Return(run)
	return _c
}

// BuildValidationError provides a mock function with given fields: code, err, errorFields
func (_m *Response) BuildValidationError(code string, err error, errorFields interface{}) response.ErrorHandler {
	ret := _m.Called(code, err, errorFields)

	if len(ret) == 0 {
		panic("no return value specified for BuildValidationError")
	}

	var r0 response.ErrorHandler
	if rf, ok := ret.Get(0).(func(string, error, interface{}) response.ErrorHandler); ok {
		r0 = rf(code, err, errorFields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(response.ErrorHandler)
		}
	}

	return r0
}

// Response_BuildValidationError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BuildValidationError'
type Response_BuildValidationError_Call struct {
	*mock.Call
}

// BuildValidationError is a helper method to define mock.On call
//   - code string
//   - err error
//   - errorFields interface{}
func (_e *Response_Expecter) BuildValidationError(code interface{}, err interface{}, errorFields interface{}) *Response_BuildValidationError_Call {
	return &Response_BuildValidationError_Call{Call: _e.mock.On("BuildValidationError", code, err, errorFields)}
}

func (_c *Response_BuildValidationError_Call) Run(run func(code string, err error, errorFields interface{})) *Response_BuildValidationError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(error), args[2].(interface{}))
	})
	return _c
}

func (_c *Response_BuildValidationError_Call) Return(_a0 response.ErrorHandler) *Response_BuildValidationError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_BuildValidationError_Call) RunAndReturn(run func(string, error, interface{}) response.ErrorHandler) *Response_BuildValidationError_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: c, httpCode, data, err
func (_m *Response) Send(c *gin.Context, httpCode int, data interface{}, err error) (int, response.BaseResponseSchema, error) {
	ret := _m.Called(c, httpCode, data, err)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 int
	var r1 response.BaseResponseSchema
	var r2 error
	if rf, ok := ret.Get(0).(func(*gin.Context, int, interface{}, error) (int, response.BaseResponseSchema, error)); ok {
		return rf(c, httpCode, data, err)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context, int, interface{}, error) int); ok {
		r0 = rf(c, httpCode, data, err)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*gin.Context, int, interface{}, error) response.BaseResponseSchema); ok {
		r1 = rf(c, httpCode, data, err)
	} else {
		r1 = ret.Get(1).(response.BaseResponseSchema)
	}

	if rf, ok := ret.Get(2).(func(*gin.Context, int, interface{}, error) error); ok {
		r2 = rf(c, httpCode, data, err)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Response_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type Response_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - c *gin.Context
//   - httpCode int
//   - data interface{}
//   - err error
func (_e *Response_Expecter) Send(c interface{}, httpCode interface{}, data interface{}, err interface{}) *Response_Send_Call {
	return &Response_Send_Call{Call: _e.mock.On("Send", c, httpCode, data, err)}
}

func (_c *Response_Send_Call) Run(run func(c *gin.Context, httpCode int, data interface{}, err error)) *Response_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context), args[1].(int), args[2].(interface{}), args[3].(error))
	})
	return _c
}

func (_c *Response_Send_Call) Return(_a0 int, _a1 response.BaseResponseSchema, _a2 error) *Response_Send_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *Response_Send_Call) RunAndReturn(run func(*gin.Context, int, interface{}, error) (int, response.BaseResponseSchema, error)) *Response_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewResponse creates a new instance of Response. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResponse(t interface {
	mock.TestingT
	Cleanup(func())
}) *Response {
	mock := &Response{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

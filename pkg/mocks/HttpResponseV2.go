// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	gin "github.com/gin-gonic/gin"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"

	mock "github.com/stretchr/testify/mock"
)

// HttpResponseV2 is an autogenerated mock type for the HttpResponseV2 type
type HttpResponseV2 struct {
	mock.Mock
}

type HttpResponseV2_Expecter struct {
	mock *mock.Mock
}

func (_m *HttpResponseV2) EXPECT() *HttpResponseV2_Expecter {
	return &HttpResponseV2_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: c
func (_m *HttpResponseV2) Send(c *gin.Context) {
	_m.Called(c)
}

// HttpResponseV2_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type HttpResponseV2_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - c *gin.Context
func (_e *HttpResponseV2_Expecter) Send(c interface{}) *HttpResponseV2_Send_Call {
	return &HttpResponseV2_Send_Call{Call: _e.mock.On("Send", c)}
}

func (_c *HttpResponseV2_Send_Call) Run(run func(c *gin.Context)) *HttpResponseV2_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *HttpResponseV2_Send_Call) Return() *HttpResponseV2_Send_Call {
	_c.Call.Return()
	return _c
}

func (_c *HttpResponseV2_Send_Call) RunAndReturn(run func(*gin.Context)) *HttpResponseV2_Send_Call {
	_c.Call.Return(run)
	return _c
}

// WithCtx provides a mock function with given fields: ctx
func (_m *HttpResponseV2) WithCtx(ctx context.Context) http_pkg.HttpResponseV2 {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for WithCtx")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(context.Context) http_pkg.HttpResponseV2); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithCtx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithCtx'
type HttpResponseV2_WithCtx_Call struct {
	*mock.Call
}

// WithCtx is a helper method to define mock.On call
//   - ctx context.Context
func (_e *HttpResponseV2_Expecter) WithCtx(ctx interface{}) *HttpResponseV2_WithCtx_Call {
	return &HttpResponseV2_WithCtx_Call{Call: _e.mock.On("WithCtx", ctx)}
}

func (_c *HttpResponseV2_WithCtx_Call) Run(run func(ctx context.Context)) *HttpResponseV2_WithCtx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *HttpResponseV2_WithCtx_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithCtx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithCtx_Call) RunAndReturn(run func(context.Context) http_pkg.HttpResponseV2) *HttpResponseV2_WithCtx_Call {
	_c.Call.Return(run)
	return _c
}

// WithError provides a mock function with given fields: err
func (_m *HttpResponseV2) WithError(err error) http_pkg.HttpResponseV2 {
	ret := _m.Called(err)

	if len(ret) == 0 {
		panic("no return value specified for WithError")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(error) http_pkg.HttpResponseV2); ok {
		r0 = rf(err)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithError'
type HttpResponseV2_WithError_Call struct {
	*mock.Call
}

// WithError is a helper method to define mock.On call
//   - err error
func (_e *HttpResponseV2_Expecter) WithError(err interface{}) *HttpResponseV2_WithError_Call {
	return &HttpResponseV2_WithError_Call{Call: _e.mock.On("WithError", err)}
}

func (_c *HttpResponseV2_WithError_Call) Run(run func(err error)) *HttpResponseV2_WithError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error))
	})
	return _c
}

func (_c *HttpResponseV2_WithError_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithError_Call) RunAndReturn(run func(error) http_pkg.HttpResponseV2) *HttpResponseV2_WithError_Call {
	_c.Call.Return(run)
	return _c
}

// WithMetadata provides a mock function with given fields: metadata
func (_m *HttpResponseV2) WithMetadata(metadata map[string]interface{}) http_pkg.HttpResponseV2 {
	ret := _m.Called(metadata)

	if len(ret) == 0 {
		panic("no return value specified for WithMetadata")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(map[string]interface{}) http_pkg.HttpResponseV2); ok {
		r0 = rf(metadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithMetadata'
type HttpResponseV2_WithMetadata_Call struct {
	*mock.Call
}

// WithMetadata is a helper method to define mock.On call
//   - metadata map[string]interface{}
func (_e *HttpResponseV2_Expecter) WithMetadata(metadata interface{}) *HttpResponseV2_WithMetadata_Call {
	return &HttpResponseV2_WithMetadata_Call{Call: _e.mock.On("WithMetadata", metadata)}
}

func (_c *HttpResponseV2_WithMetadata_Call) Run(run func(metadata map[string]interface{})) *HttpResponseV2_WithMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]interface{}))
	})
	return _c
}

func (_c *HttpResponseV2_WithMetadata_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithMetadata_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithMetadata_Call) RunAndReturn(run func(map[string]interface{}) http_pkg.HttpResponseV2) *HttpResponseV2_WithMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// WithMetadataFromStruct provides a mock function with given fields: metadata
func (_m *HttpResponseV2) WithMetadataFromStruct(metadata interface{}) http_pkg.HttpResponseV2 {
	ret := _m.Called(metadata)

	if len(ret) == 0 {
		panic("no return value specified for WithMetadataFromStruct")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(interface{}) http_pkg.HttpResponseV2); ok {
		r0 = rf(metadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithMetadataFromStruct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithMetadataFromStruct'
type HttpResponseV2_WithMetadataFromStruct_Call struct {
	*mock.Call
}

// WithMetadataFromStruct is a helper method to define mock.On call
//   - metadata interface{}
func (_e *HttpResponseV2_Expecter) WithMetadataFromStruct(metadata interface{}) *HttpResponseV2_WithMetadataFromStruct_Call {
	return &HttpResponseV2_WithMetadataFromStruct_Call{Call: _e.mock.On("WithMetadataFromStruct", metadata)}
}

func (_c *HttpResponseV2_WithMetadataFromStruct_Call) Run(run func(metadata interface{})) *HttpResponseV2_WithMetadataFromStruct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *HttpResponseV2_WithMetadataFromStruct_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithMetadataFromStruct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithMetadataFromStruct_Call) RunAndReturn(run func(interface{}) http_pkg.HttpResponseV2) *HttpResponseV2_WithMetadataFromStruct_Call {
	_c.Call.Return(run)
	return _c
}

// WithPaginationMetadata provides a mock function with given fields: page, limit, total, totalPages
func (_m *HttpResponseV2) WithPaginationMetadata(page int, limit int, total int, totalPages int) http_pkg.HttpResponseV2 {
	ret := _m.Called(page, limit, total, totalPages)

	if len(ret) == 0 {
		panic("no return value specified for WithPaginationMetadata")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(int, int, int, int) http_pkg.HttpResponseV2); ok {
		r0 = rf(page, limit, total, totalPages)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithPaginationMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithPaginationMetadata'
type HttpResponseV2_WithPaginationMetadata_Call struct {
	*mock.Call
}

// WithPaginationMetadata is a helper method to define mock.On call
//   - page int
//   - limit int
//   - total int
//   - totalPages int
func (_e *HttpResponseV2_Expecter) WithPaginationMetadata(page interface{}, limit interface{}, total interface{}, totalPages interface{}) *HttpResponseV2_WithPaginationMetadata_Call {
	return &HttpResponseV2_WithPaginationMetadata_Call{Call: _e.mock.On("WithPaginationMetadata", page, limit, total, totalPages)}
}

func (_c *HttpResponseV2_WithPaginationMetadata_Call) Run(run func(page int, limit int, total int, totalPages int)) *HttpResponseV2_WithPaginationMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(int), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *HttpResponseV2_WithPaginationMetadata_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithPaginationMetadata_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithPaginationMetadata_Call) RunAndReturn(run func(int, int, int, int) http_pkg.HttpResponseV2) *HttpResponseV2_WithPaginationMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// WithResult provides a mock function with given fields: result
func (_m *HttpResponseV2) WithResult(result interface{}) http_pkg.HttpResponseV2 {
	ret := _m.Called(result)

	if len(ret) == 0 {
		panic("no return value specified for WithResult")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(interface{}) http_pkg.HttpResponseV2); ok {
		r0 = rf(result)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithResult_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithResult'
type HttpResponseV2_WithResult_Call struct {
	*mock.Call
}

// WithResult is a helper method to define mock.On call
//   - result interface{}
func (_e *HttpResponseV2_Expecter) WithResult(result interface{}) *HttpResponseV2_WithResult_Call {
	return &HttpResponseV2_WithResult_Call{Call: _e.mock.On("WithResult", result)}
}

func (_c *HttpResponseV2_WithResult_Call) Run(run func(result interface{})) *HttpResponseV2_WithResult_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *HttpResponseV2_WithResult_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithResult_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithResult_Call) RunAndReturn(run func(interface{}) http_pkg.HttpResponseV2) *HttpResponseV2_WithResult_Call {
	_c.Call.Return(run)
	return _c
}

// WithStatusCode provides a mock function with given fields: statusCode
func (_m *HttpResponseV2) WithStatusCode(statusCode int) http_pkg.HttpResponseV2 {
	ret := _m.Called(statusCode)

	if len(ret) == 0 {
		panic("no return value specified for WithStatusCode")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(int) http_pkg.HttpResponseV2); ok {
		r0 = rf(statusCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithStatusCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithStatusCode'
type HttpResponseV2_WithStatusCode_Call struct {
	*mock.Call
}

// WithStatusCode is a helper method to define mock.On call
//   - statusCode int
func (_e *HttpResponseV2_Expecter) WithStatusCode(statusCode interface{}) *HttpResponseV2_WithStatusCode_Call {
	return &HttpResponseV2_WithStatusCode_Call{Call: _e.mock.On("WithStatusCode", statusCode)}
}

func (_c *HttpResponseV2_WithStatusCode_Call) Run(run func(statusCode int)) *HttpResponseV2_WithStatusCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *HttpResponseV2_WithStatusCode_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithStatusCode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithStatusCode_Call) RunAndReturn(run func(int) http_pkg.HttpResponseV2) *HttpResponseV2_WithStatusCode_Call {
	_c.Call.Return(run)
	return _c
}

// WithTraceId provides a mock function with given fields: traceId
func (_m *HttpResponseV2) WithTraceId(traceId string) http_pkg.HttpResponseV2 {
	ret := _m.Called(traceId)

	if len(ret) == 0 {
		panic("no return value specified for WithTraceId")
	}

	var r0 http_pkg.HttpResponseV2
	if rf, ok := ret.Get(0).(func(string) http_pkg.HttpResponseV2); ok {
		r0 = rf(traceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http_pkg.HttpResponseV2)
		}
	}

	return r0
}

// HttpResponseV2_WithTraceId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTraceId'
type HttpResponseV2_WithTraceId_Call struct {
	*mock.Call
}

// WithTraceId is a helper method to define mock.On call
//   - traceId string
func (_e *HttpResponseV2_Expecter) WithTraceId(traceId interface{}) *HttpResponseV2_WithTraceId_Call {
	return &HttpResponseV2_WithTraceId_Call{Call: _e.mock.On("WithTraceId", traceId)}
}

func (_c *HttpResponseV2_WithTraceId_Call) Run(run func(traceId string)) *HttpResponseV2_WithTraceId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *HttpResponseV2_WithTraceId_Call) Return(_a0 http_pkg.HttpResponseV2) *HttpResponseV2_WithTraceId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpResponseV2_WithTraceId_Call) RunAndReturn(run func(string) http_pkg.HttpResponseV2) *HttpResponseV2_WithTraceId_Call {
	_c.Call.Return(run)
	return _c
}

// NewHttpResponseV2 creates a new instance of HttpResponseV2. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHttpResponseV2(t interface {
	mock.TestingT
	Cleanup(func())
}) *HttpResponseV2 {
	mock := &HttpResponseV2{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

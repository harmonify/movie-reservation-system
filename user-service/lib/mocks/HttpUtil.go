// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// HttpUtil is an autogenerated mock type for the HttpUtil type
type HttpUtil struct {
	mock.Mock
}

type HttpUtil_Expecter struct {
	mock *mock.Mock
}

func (_m *HttpUtil) EXPECT() *HttpUtil_Expecter {
	return &HttpUtil_Expecter{mock: &_m.Mock}
}

// GetUserIP provides a mock function with given fields: r
func (_m *HttpUtil) GetUserIP(r *http.Request) string {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for GetUserIP")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(*http.Request) string); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// HttpUtil_GetUserIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserIP'
type HttpUtil_GetUserIP_Call struct {
	*mock.Call
}

// GetUserIP is a helper method to define mock.On call
//   - r *http.Request
func (_e *HttpUtil_Expecter) GetUserIP(r interface{}) *HttpUtil_GetUserIP_Call {
	return &HttpUtil_GetUserIP_Call{Call: _e.mock.On("GetUserIP", r)}
}

func (_c *HttpUtil_GetUserIP_Call) Run(run func(r *http.Request)) *HttpUtil_GetUserIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *HttpUtil_GetUserIP_Call) Return(_a0 string) *HttpUtil_GetUserIP_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HttpUtil_GetUserIP_Call) RunAndReturn(run func(*http.Request) string) *HttpUtil_GetUserIP_Call {
	_c.Call.Return(run)
	return _c
}

// NewHttpUtil creates a new instance of HttpUtil. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHttpUtil(t interface {
	mock.TestingT
	Cleanup(func())
}) *HttpUtil {
	mock := &HttpUtil{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

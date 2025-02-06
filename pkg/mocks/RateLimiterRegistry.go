// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	ratelimiter "github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	mock "github.com/stretchr/testify/mock"
)

// RateLimiterRegistry is an autogenerated mock type for the RateLimiterRegistry type
type RateLimiterRegistry struct {
	mock.Mock
}

type RateLimiterRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *RateLimiterRegistry) EXPECT() *RateLimiterRegistry_Expecter {
	return &RateLimiterRegistry_Expecter{mock: &_m.Mock}
}

// GetHttpRequestRateLimiter provides a mock function with given fields: p, c
func (_m *RateLimiterRegistry) GetHttpRequestRateLimiter(p *ratelimiter.HttpRequestRateLimiterParam, c *ratelimiter.RateLimiterConfig) (ratelimiter.RateLimiter, error) {
	ret := _m.Called(p, c)

	if len(ret) == 0 {
		panic("no return value specified for GetHttpRequestRateLimiter")
	}

	var r0 ratelimiter.RateLimiter
	var r1 error
	if rf, ok := ret.Get(0).(func(*ratelimiter.HttpRequestRateLimiterParam, *ratelimiter.RateLimiterConfig) (ratelimiter.RateLimiter, error)); ok {
		return rf(p, c)
	}
	if rf, ok := ret.Get(0).(func(*ratelimiter.HttpRequestRateLimiterParam, *ratelimiter.RateLimiterConfig) ratelimiter.RateLimiter); ok {
		r0 = rf(p, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ratelimiter.RateLimiter)
		}
	}

	if rf, ok := ret.Get(1).(func(*ratelimiter.HttpRequestRateLimiterParam, *ratelimiter.RateLimiterConfig) error); ok {
		r1 = rf(p, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RateLimiterRegistry_GetHttpRequestRateLimiter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHttpRequestRateLimiter'
type RateLimiterRegistry_GetHttpRequestRateLimiter_Call struct {
	*mock.Call
}

// GetHttpRequestRateLimiter is a helper method to define mock.On call
//   - p *ratelimiter.HttpRequestRateLimiterParam
//   - c *ratelimiter.RateLimiterConfig
func (_e *RateLimiterRegistry_Expecter) GetHttpRequestRateLimiter(p interface{}, c interface{}) *RateLimiterRegistry_GetHttpRequestRateLimiter_Call {
	return &RateLimiterRegistry_GetHttpRequestRateLimiter_Call{Call: _e.mock.On("GetHttpRequestRateLimiter", p, c)}
}

func (_c *RateLimiterRegistry_GetHttpRequestRateLimiter_Call) Run(run func(p *ratelimiter.HttpRequestRateLimiterParam, c *ratelimiter.RateLimiterConfig)) *RateLimiterRegistry_GetHttpRequestRateLimiter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ratelimiter.HttpRequestRateLimiterParam), args[1].(*ratelimiter.RateLimiterConfig))
	})
	return _c
}

func (_c *RateLimiterRegistry_GetHttpRequestRateLimiter_Call) Return(_a0 ratelimiter.RateLimiter, _a1 error) *RateLimiterRegistry_GetHttpRequestRateLimiter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RateLimiterRegistry_GetHttpRequestRateLimiter_Call) RunAndReturn(run func(*ratelimiter.HttpRequestRateLimiterParam, *ratelimiter.RateLimiterConfig) (ratelimiter.RateLimiter, error)) *RateLimiterRegistry_GetHttpRequestRateLimiter_Call {
	_c.Call.Return(run)
	return _c
}

// Len provides a mock function with given fields:
func (_m *RateLimiterRegistry) Len() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Len")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// RateLimiterRegistry_Len_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Len'
type RateLimiterRegistry_Len_Call struct {
	*mock.Call
}

// Len is a helper method to define mock.On call
func (_e *RateLimiterRegistry_Expecter) Len() *RateLimiterRegistry_Len_Call {
	return &RateLimiterRegistry_Len_Call{Call: _e.mock.On("Len")}
}

func (_c *RateLimiterRegistry_Len_Call) Run(run func()) *RateLimiterRegistry_Len_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *RateLimiterRegistry_Len_Call) Return(_a0 int) *RateLimiterRegistry_Len_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RateLimiterRegistry_Len_Call) RunAndReturn(run func() int) *RateLimiterRegistry_Len_Call {
	_c.Call.Return(run)
	return _c
}

// NewRateLimiterRegistry creates a new instance of RateLimiterRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRateLimiterRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *RateLimiterRegistry {
	mock := &RateLimiterRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

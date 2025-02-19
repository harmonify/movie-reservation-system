// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	kafka "github.com/harmonify/movie-reservation-system/pkg/kafka"
	mock "github.com/stretchr/testify/mock"
)

// Route is an autogenerated mock type for the Route type
type Route struct {
	mock.Mock
}

type Route_Expecter struct {
	mock *mock.Mock
}

func (_m *Route) EXPECT() *Route_Expecter {
	return &Route_Expecter{mock: &_m.Mock}
}

// AddEventListener provides a mock function with given fields: listener
func (_m *Route) AddEventListener(listener kafka.EventListener) {
	_m.Called(listener)
}

// Route_AddEventListener_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddEventListener'
type Route_AddEventListener_Call struct {
	*mock.Call
}

// AddEventListener is a helper method to define mock.On call
//   - listener kafka.EventListener
func (_e *Route_Expecter) AddEventListener(listener interface{}) *Route_AddEventListener_Call {
	return &Route_AddEventListener_Call{Call: _e.mock.On("AddEventListener", listener)}
}

func (_c *Route_AddEventListener_Call) Run(run func(listener kafka.EventListener)) *Route_AddEventListener_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(kafka.EventListener))
	})
	return _c
}

func (_c *Route_AddEventListener_Call) Return() *Route_AddEventListener_Call {
	_c.Call.Return()
	return _c
}

func (_c *Route_AddEventListener_Call) RunAndReturn(run func(kafka.EventListener)) *Route_AddEventListener_Call {
	_c.Call.Return(run)
	return _c
}

// Handle provides a mock function with given fields: ctx, event
func (_m *Route) Handle(ctx context.Context, event *kafka.Event) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for Handle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *kafka.Event) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Route_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type Route_Handle_Call struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - ctx context.Context
//   - event *kafka.Event
func (_e *Route_Expecter) Handle(ctx interface{}, event interface{}) *Route_Handle_Call {
	return &Route_Handle_Call{Call: _e.mock.On("Handle", ctx, event)}
}

func (_c *Route_Handle_Call) Run(run func(ctx context.Context, event *kafka.Event)) *Route_Handle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*kafka.Event))
	})
	return _c
}

func (_c *Route_Handle_Call) Return(_a0 error) *Route_Handle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Route_Handle_Call) RunAndReturn(run func(context.Context, *kafka.Event) error) *Route_Handle_Call {
	_c.Call.Return(run)
	return _c
}

// Identifier provides a mock function with given fields:
func (_m *Route) Identifier() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Identifier")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Route_Identifier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Identifier'
type Route_Identifier_Call struct {
	*mock.Call
}

// Identifier is a helper method to define mock.On call
func (_e *Route_Expecter) Identifier() *Route_Identifier_Call {
	return &Route_Identifier_Call{Call: _e.mock.On("Identifier")}
}

func (_c *Route_Identifier_Call) Run(run func()) *Route_Identifier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Route_Identifier_Call) Return(_a0 string) *Route_Identifier_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Route_Identifier_Call) RunAndReturn(run func() string) *Route_Identifier_Call {
	_c.Call.Return(run)
	return _c
}

// Match provides a mock function with given fields: ctx, event
func (_m *Route) Match(ctx context.Context, event *kafka.Event) (bool, error) {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for Match")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *kafka.Event) (bool, error)); ok {
		return rf(ctx, event)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *kafka.Event) bool); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *kafka.Event) error); ok {
		r1 = rf(ctx, event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Route_Match_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Match'
type Route_Match_Call struct {
	*mock.Call
}

// Match is a helper method to define mock.On call
//   - ctx context.Context
//   - event *kafka.Event
func (_e *Route_Expecter) Match(ctx interface{}, event interface{}) *Route_Match_Call {
	return &Route_Match_Call{Call: _e.mock.On("Match", ctx, event)}
}

func (_c *Route_Match_Call) Run(run func(ctx context.Context, event *kafka.Event)) *Route_Match_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*kafka.Event))
	})
	return _c
}

func (_c *Route_Match_Call) Return(_a0 bool, _a1 error) *Route_Match_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Route_Match_Call) RunAndReturn(run func(context.Context, *kafka.Event) (bool, error)) *Route_Match_Call {
	_c.Call.Return(run)
	return _c
}

// NewRoute creates a new instance of Route. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRoute(t interface {
	mock.TestingT
	Cleanup(func())
}) *Route {
	mock := &Route{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

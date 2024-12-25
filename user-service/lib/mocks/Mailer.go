// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mail "github.com/harmonify/movie-reservation-system/user-service/lib/mail"
	mock "github.com/stretchr/testify/mock"
)

// Mailer is an autogenerated mock type for the Mailer type
type Mailer struct {
	mock.Mock
}

type Mailer_Expecter struct {
	mock *mock.Mock
}

func (_m *Mailer) EXPECT() *Mailer_Expecter {
	return &Mailer_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: ctx, message
func (_m *Mailer) Send(ctx context.Context, message mail.Message) (string, error) {
	ret := _m.Called(ctx, message)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, mail.Message) (string, error)); ok {
		return rf(ctx, message)
	}
	if rf, ok := ret.Get(0).(func(context.Context, mail.Message) string); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, mail.Message) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Mailer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type Mailer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - ctx context.Context
//   - message mail.Message
func (_e *Mailer_Expecter) Send(ctx interface{}, message interface{}) *Mailer_Send_Call {
	return &Mailer_Send_Call{Call: _e.mock.On("Send", ctx, message)}
}

func (_c *Mailer_Send_Call) Run(run func(ctx context.Context, message mail.Message)) *Mailer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(mail.Message))
	})
	return _c
}

func (_c *Mailer_Send_Call) Return(id string, err error) *Mailer_Send_Call {
	_c.Call.Return(id, err)
	return _c
}

func (_c *Mailer_Send_Call) RunAndReturn(run func(context.Context, mail.Message) (string, error)) *Mailer_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewMailer creates a new instance of Mailer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMailer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mailer {
	mock := &Mailer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
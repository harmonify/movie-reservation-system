// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	mock "github.com/stretchr/testify/mock"
)

// NotificationServiceServer is an autogenerated mock type for the NotificationServiceServer type
type NotificationServiceServer struct {
	mock.Mock
}

type NotificationServiceServer_Expecter struct {
	mock *mock.Mock
}

func (_m *NotificationServiceServer) EXPECT() *NotificationServiceServer_Expecter {
	return &NotificationServiceServer_Expecter{mock: &_m.Mock}
}

// BulkSendSms provides a mock function with given fields: _a0, _a1
func (_m *NotificationServiceServer) BulkSendSms(_a0 context.Context, _a1 *notification_proto.BulkSendSmsRequest) (*notification_proto.BulkSendSmsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for BulkSendSms")
	}

	var r0 *notification_proto.BulkSendSmsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *notification_proto.BulkSendSmsRequest) (*notification_proto.BulkSendSmsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *notification_proto.BulkSendSmsRequest) *notification_proto.BulkSendSmsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*notification_proto.BulkSendSmsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *notification_proto.BulkSendSmsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NotificationServiceServer_BulkSendSms_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BulkSendSms'
type NotificationServiceServer_BulkSendSms_Call struct {
	*mock.Call
}

// BulkSendSms is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *notification_proto.BulkSendSmsRequest
func (_e *NotificationServiceServer_Expecter) BulkSendSms(_a0 interface{}, _a1 interface{}) *NotificationServiceServer_BulkSendSms_Call {
	return &NotificationServiceServer_BulkSendSms_Call{Call: _e.mock.On("BulkSendSms", _a0, _a1)}
}

func (_c *NotificationServiceServer_BulkSendSms_Call) Run(run func(_a0 context.Context, _a1 *notification_proto.BulkSendSmsRequest)) *NotificationServiceServer_BulkSendSms_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*notification_proto.BulkSendSmsRequest))
	})
	return _c
}

func (_c *NotificationServiceServer_BulkSendSms_Call) Return(_a0 *notification_proto.BulkSendSmsResponse, _a1 error) *NotificationServiceServer_BulkSendSms_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NotificationServiceServer_BulkSendSms_Call) RunAndReturn(run func(context.Context, *notification_proto.BulkSendSmsRequest) (*notification_proto.BulkSendSmsResponse, error)) *NotificationServiceServer_BulkSendSms_Call {
	_c.Call.Return(run)
	return _c
}

// SendEmail provides a mock function with given fields: _a0, _a1
func (_m *NotificationServiceServer) SendEmail(_a0 context.Context, _a1 *notification_proto.SendEmailRequest) (*notification_proto.SendEmailResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SendEmail")
	}

	var r0 *notification_proto.SendEmailResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *notification_proto.SendEmailRequest) (*notification_proto.SendEmailResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *notification_proto.SendEmailRequest) *notification_proto.SendEmailResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*notification_proto.SendEmailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *notification_proto.SendEmailRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NotificationServiceServer_SendEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendEmail'
type NotificationServiceServer_SendEmail_Call struct {
	*mock.Call
}

// SendEmail is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *notification_proto.SendEmailRequest
func (_e *NotificationServiceServer_Expecter) SendEmail(_a0 interface{}, _a1 interface{}) *NotificationServiceServer_SendEmail_Call {
	return &NotificationServiceServer_SendEmail_Call{Call: _e.mock.On("SendEmail", _a0, _a1)}
}

func (_c *NotificationServiceServer_SendEmail_Call) Run(run func(_a0 context.Context, _a1 *notification_proto.SendEmailRequest)) *NotificationServiceServer_SendEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*notification_proto.SendEmailRequest))
	})
	return _c
}

func (_c *NotificationServiceServer_SendEmail_Call) Return(_a0 *notification_proto.SendEmailResponse, _a1 error) *NotificationServiceServer_SendEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NotificationServiceServer_SendEmail_Call) RunAndReturn(run func(context.Context, *notification_proto.SendEmailRequest) (*notification_proto.SendEmailResponse, error)) *NotificationServiceServer_SendEmail_Call {
	_c.Call.Return(run)
	return _c
}

// SendSms provides a mock function with given fields: _a0, _a1
func (_m *NotificationServiceServer) SendSms(_a0 context.Context, _a1 *notification_proto.SendSmsRequest) (*notification_proto.SendSmsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SendSms")
	}

	var r0 *notification_proto.SendSmsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *notification_proto.SendSmsRequest) (*notification_proto.SendSmsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *notification_proto.SendSmsRequest) *notification_proto.SendSmsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*notification_proto.SendSmsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *notification_proto.SendSmsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NotificationServiceServer_SendSms_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendSms'
type NotificationServiceServer_SendSms_Call struct {
	*mock.Call
}

// SendSms is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *notification_proto.SendSmsRequest
func (_e *NotificationServiceServer_Expecter) SendSms(_a0 interface{}, _a1 interface{}) *NotificationServiceServer_SendSms_Call {
	return &NotificationServiceServer_SendSms_Call{Call: _e.mock.On("SendSms", _a0, _a1)}
}

func (_c *NotificationServiceServer_SendSms_Call) Run(run func(_a0 context.Context, _a1 *notification_proto.SendSmsRequest)) *NotificationServiceServer_SendSms_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*notification_proto.SendSmsRequest))
	})
	return _c
}

func (_c *NotificationServiceServer_SendSms_Call) Return(_a0 *notification_proto.SendSmsResponse, _a1 error) *NotificationServiceServer_SendSms_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NotificationServiceServer_SendSms_Call) RunAndReturn(run func(context.Context, *notification_proto.SendSmsRequest) (*notification_proto.SendSmsResponse, error)) *NotificationServiceServer_SendSms_Call {
	_c.Call.Return(run)
	return _c
}

// mustEmbedUnimplementedNotificationServiceServer provides a mock function with given fields:
func (_m *NotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	_m.Called()
}

// NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedNotificationServiceServer'
type NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedNotificationServiceServer is a helper method to define mock.On call
func (_e *NotificationServiceServer_Expecter) mustEmbedUnimplementedNotificationServiceServer() *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call {
	return &NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call{Call: _e.mock.On("mustEmbedUnimplementedNotificationServiceServer")}
}

func (_c *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call) Run(run func()) *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call) Return() *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call) RunAndReturn(run func()) *NotificationServiceServer_mustEmbedUnimplementedNotificationServiceServer_Call {
	_c.Call.Return(run)
	return _c
}

// NewNotificationServiceServer creates a new instance of NotificationServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotificationServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotificationServiceServer {
	mock := &NotificationServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

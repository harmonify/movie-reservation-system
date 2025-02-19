// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	logger "github.com/harmonify/movie-reservation-system/pkg/logger"
	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"

	zapcore "go.uber.org/zap/zapcore"
)

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

type Logger_Expecter struct {
	mock *mock.Mock
}

func (_m *Logger) EXPECT() *Logger_Expecter {
	return &Logger_Expecter{mock: &_m.Mock}
}

// Debug provides a mock function with given fields: msg, fields
func (_m *Logger) Debug(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type Logger_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Debug(msg interface{}, fields ...interface{}) *Logger_Debug_Call {
	return &Logger_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Debug_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Debug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Debug_Call) Return() *Logger_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Debug_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: msg, fields
func (_m *Logger) Error(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type Logger_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Error(msg interface{}, fields ...interface{}) *Logger_Error_Call {
	return &Logger_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Error_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Error_Call) Return() *Logger_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Error_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Error_Call {
	_c.Call.Return(run)
	return _c
}

// GetZapLogger provides a mock function with given fields:
func (_m *Logger) GetZapLogger() *zap.Logger {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetZapLogger")
	}

	var r0 *zap.Logger
	if rf, ok := ret.Get(0).(func() *zap.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*zap.Logger)
		}
	}

	return r0
}

// Logger_GetZapLogger_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetZapLogger'
type Logger_GetZapLogger_Call struct {
	*mock.Call
}

// GetZapLogger is a helper method to define mock.On call
func (_e *Logger_Expecter) GetZapLogger() *Logger_GetZapLogger_Call {
	return &Logger_GetZapLogger_Call{Call: _e.mock.On("GetZapLogger")}
}

func (_c *Logger_GetZapLogger_Call) Run(run func()) *Logger_GetZapLogger_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Logger_GetZapLogger_Call) Return(_a0 *zap.Logger) *Logger_GetZapLogger_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Logger_GetZapLogger_Call) RunAndReturn(run func() *zap.Logger) *Logger_GetZapLogger_Call {
	_c.Call.Return(run)
	return _c
}

// Info provides a mock function with given fields: msg, fields
func (_m *Logger) Info(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type Logger_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Info(msg interface{}, fields ...interface{}) *Logger_Info_Call {
	return &Logger_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Info_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Info_Call) Return() *Logger_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Info_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Info_Call {
	_c.Call.Return(run)
	return _c
}

// Level provides a mock function with given fields:
func (_m *Logger) Level() zapcore.Level {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Level")
	}

	var r0 zapcore.Level
	if rf, ok := ret.Get(0).(func() zapcore.Level); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(zapcore.Level)
	}

	return r0
}

// Logger_Level_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Level'
type Logger_Level_Call struct {
	*mock.Call
}

// Level is a helper method to define mock.On call
func (_e *Logger_Expecter) Level() *Logger_Level_Call {
	return &Logger_Level_Call{Call: _e.mock.On("Level")}
}

func (_c *Logger_Level_Call) Run(run func()) *Logger_Level_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Logger_Level_Call) Return(_a0 zapcore.Level) *Logger_Level_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Logger_Level_Call) RunAndReturn(run func() zapcore.Level) *Logger_Level_Call {
	_c.Call.Return(run)
	return _c
}

// Log provides a mock function with given fields: debugLevel, msg, fields
func (_m *Logger) Log(debugLevel zapcore.Level, msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, debugLevel, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Log_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Log'
type Logger_Log_Call struct {
	*mock.Call
}

// Log is a helper method to define mock.On call
//   - debugLevel zapcore.Level
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Log(debugLevel interface{}, msg interface{}, fields ...interface{}) *Logger_Log_Call {
	return &Logger_Log_Call{Call: _e.mock.On("Log",
		append([]interface{}{debugLevel, msg}, fields...)...)}
}

func (_c *Logger_Log_Call) Run(run func(debugLevel zapcore.Level, msg string, fields ...zapcore.Field)) *Logger_Log_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(zapcore.Level), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Log_Call) Return() *Logger_Log_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Log_Call) RunAndReturn(run func(zapcore.Level, string, ...zapcore.Field)) *Logger_Log_Call {
	_c.Call.Return(run)
	return _c
}

// Warn provides a mock function with given fields: msg, fields
func (_m *Logger) Warn(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type Logger_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Warn(msg interface{}, fields ...interface{}) *Logger_Warn_Call {
	return &Logger_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Warn_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Warn_Call) Return() *Logger_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Warn_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Warn_Call {
	_c.Call.Return(run)
	return _c
}

// With provides a mock function with given fields: fields
func (_m *Logger) With(fields ...zapcore.Field) logger.Logger {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for With")
	}

	var r0 logger.Logger
	if rf, ok := ret.Get(0).(func(...zapcore.Field) logger.Logger); ok {
		r0 = rf(fields...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(logger.Logger)
		}
	}

	return r0
}

// Logger_With_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'With'
type Logger_With_Call struct {
	*mock.Call
}

// With is a helper method to define mock.On call
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) With(fields ...interface{}) *Logger_With_Call {
	return &Logger_With_Call{Call: _e.mock.On("With",
		append([]interface{}{}, fields...)...)}
}

func (_c *Logger_With_Call) Run(run func(fields ...zapcore.Field)) *Logger_With_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Logger_With_Call) Return(_a0 logger.Logger) *Logger_With_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Logger_With_Call) RunAndReturn(run func(...zapcore.Field) logger.Logger) *Logger_With_Call {
	_c.Call.Return(run)
	return _c
}

// WithCtx provides a mock function with given fields: ctx
func (_m *Logger) WithCtx(ctx context.Context) logger.Logger {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for WithCtx")
	}

	var r0 logger.Logger
	if rf, ok := ret.Get(0).(func(context.Context) logger.Logger); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(logger.Logger)
		}
	}

	return r0
}

// Logger_WithCtx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithCtx'
type Logger_WithCtx_Call struct {
	*mock.Call
}

// WithCtx is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Logger_Expecter) WithCtx(ctx interface{}) *Logger_WithCtx_Call {
	return &Logger_WithCtx_Call{Call: _e.mock.On("WithCtx", ctx)}
}

func (_c *Logger_WithCtx_Call) Run(run func(ctx context.Context)) *Logger_WithCtx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Logger_WithCtx_Call) Return(_a0 logger.Logger) *Logger_WithCtx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Logger_WithCtx_Call) RunAndReturn(run func(context.Context) logger.Logger) *Logger_WithCtx_Call {
	_c.Call.Return(run)
	return _c
}

// NewLogger creates a new instance of Logger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *Logger {
	mock := &Logger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

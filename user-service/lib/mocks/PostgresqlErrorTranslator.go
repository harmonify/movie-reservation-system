// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PostgresqlErrorTranslator is an autogenerated mock type for the PostgresqlErrorTranslator type
type PostgresqlErrorTranslator struct {
	mock.Mock
}

type PostgresqlErrorTranslator_Expecter struct {
	mock *mock.Mock
}

func (_m *PostgresqlErrorTranslator) EXPECT() *PostgresqlErrorTranslator_Expecter {
	return &PostgresqlErrorTranslator_Expecter{mock: &_m.Mock}
}

// Translate provides a mock function with given fields: err
func (_m *PostgresqlErrorTranslator) Translate(err error) error {
	ret := _m.Called(err)

	if len(ret) == 0 {
		panic("no return value specified for Translate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(error) error); ok {
		r0 = rf(err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostgresqlErrorTranslator_Translate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Translate'
type PostgresqlErrorTranslator_Translate_Call struct {
	*mock.Call
}

// Translate is a helper method to define mock.On call
//   - err error
func (_e *PostgresqlErrorTranslator_Expecter) Translate(err interface{}) *PostgresqlErrorTranslator_Translate_Call {
	return &PostgresqlErrorTranslator_Translate_Call{Call: _e.mock.On("Translate", err)}
}

func (_c *PostgresqlErrorTranslator_Translate_Call) Run(run func(err error)) *PostgresqlErrorTranslator_Translate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error))
	})
	return _c
}

func (_c *PostgresqlErrorTranslator_Translate_Call) Return(_a0 error) *PostgresqlErrorTranslator_Translate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PostgresqlErrorTranslator_Translate_Call) RunAndReturn(run func(error) error) *PostgresqlErrorTranslator_Translate_Call {
	_c.Call.Return(run)
	return _c
}

// NewPostgresqlErrorTranslator creates a new instance of PostgresqlErrorTranslator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostgresqlErrorTranslator(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostgresqlErrorTranslator {
	mock := &PostgresqlErrorTranslator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
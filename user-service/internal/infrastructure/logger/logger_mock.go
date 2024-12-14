package logger

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MockLogger struct {
	*zap.Logger
	mock.Mock
}

func NewMockLogger(patch *gomonkey.Patches) Logger {
	mockLogger := new(MockLogger)
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("WithCtx", context.Background()).Return(mockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything)
	patch.ApplyMethod(mockLogger.Logger, "Error", func(*zap.Logger, string, ...zapcore.Field) {})
	return &LoggerImpl{
		Logger: mockLogger.Logger,
	}
}

func (m *MockLogger) WithCtx(ctx context.Context) *MockLogger {
	args := m.Called(ctx)
	return args.Get(0).(*MockLogger)
}

func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) With(fields ...zap.Field) *MockLogger {
	args := m.Called(fields)
	return args.Get(0).(*MockLogger)
}

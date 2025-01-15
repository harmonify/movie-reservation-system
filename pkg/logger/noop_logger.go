package logger

import (
	"context"

	"go.uber.org/zap"
)

func NewNopLogger() Logger {
	logger := zap.NewNop()
	return &NopLoggerImpl{logger, nil}
}

func (l *NopLoggerImpl) GetZapLogger() *zap.Logger {
	return l.Logger
}

func (l *NopLoggerImpl) WithCtx(ctx context.Context) Logger {
	return l
}

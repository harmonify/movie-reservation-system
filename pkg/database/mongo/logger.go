package mongo

import (
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type mongoClientLoggerImpl struct {
	logger logger.Logger
}

func (m *mongoClientLoggerImpl) Error(err error, message string, keysAndValues ...interface{}) {
	fields := []zap.Field{
		zap.Error(err),
	}

	for i := 0; i < len(keysAndValues); i += 2 {
		fields = append(fields, zap.Any(keysAndValues[i].(string), keysAndValues[i+1]))
	}

	m.logger.Error(message, fields...)
}

func (m *mongoClientLoggerImpl) Info(level int, message string, keysAndValues ...interface{}) {
	fields := []zap.Field{}

	for i := 0; i < len(keysAndValues); i += 2 {
		fields = append(fields, zap.Any(keysAndValues[i].(string), keysAndValues[i+1]))
	}

	if level >= 2 {
		m.logger.Debug(message, fields...)
	} else {
		m.logger.Info(message, fields...)
	}
}

func NewMongoClientLogger(logger logger.Logger) options.LogSink {
	return &mongoClientLoggerImpl{
		logger: logger,
	}
}

package ratelimiter

import (
	"fmt"

	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"go.uber.org/zap"
)

type RateLimiterLogger struct {
	logger logger.Logger
}

func NewRateLimiterLogger(logger logger.Logger) *RateLimiterLogger {
	return &RateLimiterLogger{
		logger: logger,
	}
}

func (r *RateLimiterLogger) Log(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	r.logger.Info(msg, zap.Any("original_message", v))
}

package failsafe_object_logger

import (
	"reflect"

	"github.com/failsafe-go/failsafe-go"
	"go.uber.org/zap/zapcore"
)

type LoggableExecutionEvent[T any] struct {
	// The event to log.
	*failsafe.ExecutionEvent[T]
	// Whether to reflect the result or not in the log.
	UseReflect bool
}

func NewLoggableAnyExecutionEvent[T any](event failsafe.ExecutionEvent[T], useReflect bool) *LoggableExecutionEvent[T] {
	return &LoggableExecutionEvent[T]{
		ExecutionEvent: &event,
		UseReflect:     useReflect,
	}
}

func (e *LoggableExecutionEvent[T]) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("attempt.info.attempts", e.Attempts())
	enc.AddInt("attempt.info.executions", e.Executions())
	enc.AddInt("attempt.info.retries", e.Retries())
	enc.AddInt("attempt.info.hedges", e.Hedges())
	enc.AddTime("attempt.info.start_time", e.StartTime())
	enc.AddDuration("attempt.info.elapsed_time", e.ElapsedTime())

	if e.UseReflect {
		lastResult := reflect.ValueOf(e.LastResult())
		if lastResult.Kind() == reflect.Pointer || lastResult.Kind() == reflect.Interface {
			lastResult = lastResult.Elem()
		}
		if lastResult.Kind() != reflect.Invalid && !lastResult.IsZero() {
			enc.AddReflected("attempt.last_result", lastResult.Interface())
		}
	}

	if lastError := e.LastError(); lastError != nil {
		enc.AddString("attempt.last_error", lastError.Error())
	}

	enc.AddBool("attempt.is_first_attempt", e.IsFirstAttempt())
	enc.AddBool("attempt.is_retry", e.IsRetry())
	enc.AddBool("attempt.is_hedge", e.IsHedge())
	enc.AddTime("attempt.attempt_start_time", e.AttemptStartTime())
	enc.AddDuration("attempt.elapsed_attempt_time", e.ElapsedAttemptTime())

	return nil
}

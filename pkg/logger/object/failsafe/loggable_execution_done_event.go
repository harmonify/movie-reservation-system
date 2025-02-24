package failsafe_object_logger

import (
	"reflect"

	"github.com/failsafe-go/failsafe-go"
	"go.uber.org/zap/zapcore"
)

type LoggableExecutionDoneEvent[T any] struct {
	// The event to log.
	*failsafe.ExecutionDoneEvent[T]
	// Whether to reflect the result or not in the log.
	UseReflect bool
}

func NewLoggableAnyExecutionDoneEvent[T any](event failsafe.ExecutionDoneEvent[T], useReflect bool) *LoggableExecutionDoneEvent[T] {
	return &LoggableExecutionDoneEvent[T]{
		ExecutionDoneEvent: &event,
		UseReflect:         useReflect,
	}
}

func (e *LoggableExecutionDoneEvent[T]) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("info.attempts", e.Attempts())
	enc.AddInt("info.executions", e.Executions())
	enc.AddInt("info.retries", e.Retries())
	enc.AddInt("info.hedges", e.Hedges())
	enc.AddTime("info.start_time", e.StartTime())
	enc.AddDuration("info.elapsed_time", e.ElapsedTime())

	if e.UseReflect {
		result := reflect.ValueOf(e.Result)
		if result.Kind() == reflect.Pointer || result.Kind() == reflect.Interface {
			result = result.Elem()
		}
		if result.Kind() != reflect.Invalid && !result.IsZero() {
			enc.AddReflected("attempt.last_result", result.Interface())
		}
	}

	if e.Error != nil {
		enc.AddString("error", e.Error.Error())
	}

	return nil
}

package circuitbreaker_object_logger

import (
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"go.uber.org/zap/zapcore"
)

type LoggableStateChangedEvent struct {
	*circuitbreaker.StateChangedEvent
}

func NewLoggableStateChangedEvent(event circuitbreaker.StateChangedEvent) *LoggableStateChangedEvent {
	return &LoggableStateChangedEvent{&event}
}

func (e *LoggableStateChangedEvent) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("current_state", e.NewState.String())
	enc.AddString("old_state", e.OldState.String())
	enc.AddUint("executions", e.Metrics().Executions())
	enc.AddUint("failures", e.Metrics().Failures())
	enc.AddUint("failure_rate", e.Metrics().FailureRate())
	enc.AddUint("successes", e.Metrics().Successes())
	enc.AddUint("success_rate", e.Metrics().SuccessRate())
	return nil
}

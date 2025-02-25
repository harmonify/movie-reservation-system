package grpc_failsafe

import (
	"context"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

// NewUnaryClientInterceptorWithExecutor returns a grpc.UnaryClientInterceptor that wraps the invoker with a failsafe.Executor.
// R is the response type.
// NewUnaryClientInterceptorWithExecutorContext is a variant of the original function located at
// https://github.com/failsafe-go/failsafe-go/blob/bcad1475c8008421e9e60630031a2d27980ecba1/failsafegrpc/client.go#L22.
// This modified function provides executor context to its executions.
// The executor context is useful to correlate all execution units under the same executor using OpenTelemetry trace.
func NewUnaryClientInterceptorWithExecutorContext[R any](executor failsafe.Executor[R], tracer tracer.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return executor.WithContext(ctx).RunWithExecution(func(exec failsafe.Execution[R]) error {
			ctx, span := tracer.Start(ctx, method)
			defer span.End()

			attrs := []attribute.KeyValue{
				attribute.String("grpc.method", method),
				attribute.String("grpc.service", cc.Target()),
				attribute.Bool("failsafe.is_first_attempt", exec.IsFirstAttempt()),
				attribute.Bool("failsafe.is_retry", exec.IsRetry()),
				attribute.Bool("failsafe.is_hedge", exec.IsHedge()),
				attribute.String("failsafe.attempt_start_time", exec.AttemptStartTime().Format(time.RFC3339Nano)),
				attribute.String("failsafe.elapsed_attempt_time", exec.ElapsedAttemptTime().String()),
				attribute.Int("failsafe.info.attempts", exec.Attempts()),
				attribute.Int("failsafe.info.executions", exec.Executions()),
				attribute.Int("failsafe.info.retries", exec.Retries()),
				attribute.Int("failsafe.info.hedges", exec.Hedges()),
				attribute.String("failsafe.info.start_time", exec.StartTime().Format(time.RFC3339Nano)),
				attribute.String("failsafe.info.elapsed_time", exec.ElapsedTime().String()),
			}

			if err := exec.LastError(); err != nil {
				attrs = append(attrs, attribute.String("failsafe.last_error", err.Error()))
			}

			span.SetAttributes(attrs...)

			mergedCtx, cancel := MergeContexts(ctx, exec.Context())
			defer cancel(nil)

			return invoker(mergedCtx, method, req, reply, cc, opts...)
		})
	}
}

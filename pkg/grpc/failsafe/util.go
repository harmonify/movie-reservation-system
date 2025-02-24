package grpc_failsafe

import (
	"context"
)

func noop(_ error) {}

// MergeContexts returns a context that is canceled when either ctx1 or ctx2 are Done.
func MergeContexts(ctx1, ctx2 context.Context) (context.Context, context.CancelCauseFunc) {
	bgContext := context.Background()
	if ctx1 == bgContext {
		return ctx2, noop
	}
	if ctx2 == bgContext {
		return ctx1, noop
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		select {
		case <-ctx1.Done():
			cancel(ctx1.Err())
		case <-ctx2.Done():
			cancel(ctx2.Err())
		}
	}()
	return ctx, cancel
}

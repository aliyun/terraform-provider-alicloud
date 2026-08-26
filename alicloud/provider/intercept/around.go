package intercept

import "context"

// DiagBridge adapts a diagnostics collection D to the plain error the interceptor
// contract speaks. Supplied by the adapter, the two stacks having incompatible
// diagnostics types.
type DiagBridge[D any] struct {
	// The first error-severity entry of d, or nil.
	ToError func(d D) error
	// d with its error-severity entries replaced by one describing err, or dropped
	// if err is nil. Warnings must survive in their original order.
	WithError func(d D, err error) D
}

// Around runs the chain around a diagnostics-returning inner function and folds
// the chain's error back in. Only a *change* is folded back: hand the error back
// unaltered and inner's diagnostics are returned verbatim, wording intact.
func Around[D any](ctx context.Context, chain []Interceptor, call Call, inner func() D, bridge DiagBridge[D]) D {
	if len(chain) == 0 {
		return inner()
	}
	var (
		out      D
		innerErr error
	)
	outErr := Execute(ctx, chain, call, func() error {
		out = inner()
		innerErr = bridge.ToError(out)
		return innerErr
	})
	// Identity, not errors.Is: a hook that wraps the error means to reword it.
	if outErr == innerErr {
		return out
	}
	return bridge.WithError(out, outErr)
}

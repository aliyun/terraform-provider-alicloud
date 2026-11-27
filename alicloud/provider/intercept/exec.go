package intercept

import (
	"context"
	"log"
	"slices"
)

// Execute runs the chain around inner and restores any error an After hook
// dropped, so a failed operation is never reported as success.
func Execute(ctx context.Context, chain []Interceptor, call Call, inner func() error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var err error
	for _, o := range chain {
		if err = ctx.Err(); err != nil {
			break
		}
		if err = o.Before(ctx, call); err != nil {
			break
		}
	}
	if err == nil {
		err = ctx.Err()
	}
	if err == nil {
		err = inner()
	}
	for _, o := range slices.Backward(chain) {
		prev := err
		if err = o.After(ctx, call, err); err == nil && prev != nil {
			log.Printf("[ERROR] intercept: %T dropped %s %s error, restoring it: %s", o, call.Op, call.Name, prev)
			err = prev
		}
	}
	return err
}

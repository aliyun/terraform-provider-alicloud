package intercept

import (
	"context"
	"slices"
)

// Before hooks run in forward order and a non-nil return aborts the operation,
// inner never executing. After hooks always run, in reverse order, and may rewrite
// the error.
func Execute(ctx context.Context, chain []Interceptor, call Call, inner func() error) error {
	if len(chain) == 0 {
		return inner()
	}
	var err error
	aborted := false
	for _, o := range chain {
		if e := o.Before(ctx, call); e != nil {
			err = e
			aborted = true
			break
		}
	}
	if !aborted {
		err = inner()
	}
	for _, o := range slices.Backward(chain) {
		err = o.After(ctx, call, err)
	}
	return err
}

package intercept

import (
	"context"
	"errors"
	"testing"
)

// The plain error form of SDK v2 CRUD reaches Execute directly, with no bridge
// in between.
func TestExecuteAfterCannotSwallowError(t *testing.T) {
	inner := errors.New("the API said no")
	it := &funcInterceptor{after: func(error) error { return nil }}
	got := Execute(context.Background(), []Interceptor{it}, testCall, func() error { return inner })
	if got != inner {
		t.Fatalf("got %v, want the inner error restored", got)
	}
}

func TestExecuteAfterRewriteIsHonoured(t *testing.T) {
	rewritten := errors.New("normalised")
	it := &funcInterceptor{after: func(error) error { return rewritten }}
	got := Execute(context.Background(), []Interceptor{it}, testCall,
		func() error { return errors.New("the API said no") })
	if got != rewritten {
		t.Fatalf("got %v, want %v", got, rewritten)
	}
}

// The restore happens per hook, so a hook further out still sees the failure the
// one before it tried to drop.
func TestExecuteRestoredErrorReachesOuterHook(t *testing.T) {
	inner := errors.New("the API said no")
	var outerSaw error
	outer := &funcInterceptor{after: func(err error) error { outerSaw = err; return err }}
	dropper := &funcInterceptor{after: func(error) error { return nil }}

	Execute(context.Background(), []Interceptor{outer, dropper}, testCall, func() error { return inner })
	if outerSaw != inner {
		t.Fatalf("outer hook saw %v, want the restored error", outerSaw)
	}
}

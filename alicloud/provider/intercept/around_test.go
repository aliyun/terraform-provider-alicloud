package intercept

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type fakeDiags []string

func fakeToError(d fakeDiags) error {
	var msgs []string
	for _, e := range d {
		if strings.HasPrefix(e, "E:") {
			msgs = append(msgs, strings.TrimPrefix(e, "E:"))
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	return errors.New(strings.Join(msgs, "; "))
}

func fakeWithError(d fakeDiags, err error) fakeDiags {
	if err == nil {
		return d
	}
	return append(d, "E:"+err.Error())
}

var fakeBridge = DiagBridge[fakeDiags]{ToError: fakeToError, WithError: fakeWithError}

type funcInterceptor struct {
	before func() error
	after  func(err error) error
}

func (f *funcInterceptor) Before(ctx context.Context, call Call) error {
	if f.before == nil {
		return nil
	}
	return f.before()
}

func (f *funcInterceptor) After(ctx context.Context, call Call, err error) error {
	if f.after == nil {
		return err
	}
	return f.after(err)
}

var testCall = Call{Name: "alicloud_test", Op: OpRead}

func around(t *testing.T, it Interceptor, inner fakeDiags) fakeDiags {
	t.Helper()
	return Around(context.Background(), []Interceptor{it}, testCall,
		func() fakeDiags { return inner }, fakeBridge)
}

func TestAroundPassesDiagnosticsThroughVerbatim(t *testing.T) {
	inner := fakeDiags{"W:careful", "E:the API said no", "E:and again"}
	got := around(t, &funcInterceptor{}, inner)
	if fmt.Sprint(got) != fmt.Sprint(inner) {
		t.Fatalf("diagnostics were rewritten: got %v, want %v", got, inner)
	}
}

// Execute restores the error it returned nil for, which makes it identical to
// inner's, so the diagnostics come back verbatim rather than through the bridge.
func TestAroundAfterCannotSwallowError(t *testing.T) {
	it := &funcInterceptor{after: func(error) error { return nil }}
	inner := fakeDiags{"W:careful", "E:the API said no"}
	got := around(t, it, inner)
	if fmt.Sprint(got) != fmt.Sprint(inner) {
		t.Fatalf("got %v, want %v", got, inner)
	}
}

func TestAroundAfterRewritesError(t *testing.T) {
	it := &funcInterceptor{after: func(error) error { return errors.New("normalised") }}
	got := around(t, it, fakeDiags{"W:careful", "E:the API said no"})
	want := fakeDiags{"W:careful", "E:the API said no", "E:normalised"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAroundBeforeAbort(t *testing.T) {
	innerRan := false
	it := &funcInterceptor{before: func() error { return errors.New("abort") }}
	got := Around(context.Background(), []Interceptor{it}, testCall,
		func() fakeDiags { innerRan = true; return fakeDiags{"W:careful"} }, fakeBridge)
	if innerRan {
		t.Fatal("inner must be skipped after a Before abort")
	}
	want := fakeDiags{"E:abort"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAroundAfterRaisesOnSuccess(t *testing.T) {
	it := &funcInterceptor{after: func(err error) error {
		if err != nil {
			t.Fatalf("expected no inner error, got %v", err)
		}
		return errors.New("post-condition failed")
	}}
	got := around(t, it, nil)
	want := fakeDiags{"E:post-condition failed"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAroundEmptyChainIsIdentity(t *testing.T) {
	inner := fakeDiags{"E:untouched"}
	got := Around(context.Background(), nil, testCall,
		func() fakeDiags { return inner }, fakeBridge)
	if fmt.Sprint(got) != fmt.Sprint(inner) {
		t.Fatalf("got %v, want %v", got, inner)
	}
}

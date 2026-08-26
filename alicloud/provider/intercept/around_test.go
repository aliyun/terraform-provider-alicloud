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
	for _, e := range d {
		if strings.HasPrefix(e, "E:") {
			return errors.New(strings.TrimPrefix(e, "E:"))
		}
	}
	return nil
}

func fakeWithError(d fakeDiags, err error) fakeDiags {
	out := make(fakeDiags, 0, len(d)+1)
	for _, e := range d {
		if !strings.HasPrefix(e, "E:") {
			out = append(out, e)
		}
	}
	if err != nil {
		out = append(out, "E:"+err.Error())
	}
	return out
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

func TestAroundAfterSwallowsError(t *testing.T) {
	it := &funcInterceptor{after: func(error) error { return nil }}
	got := around(t, it, fakeDiags{"W:careful", "E:the API said no"})
	want := fakeDiags{"W:careful"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAroundAfterRewritesError(t *testing.T) {
	it := &funcInterceptor{after: func(error) error { return errors.New("normalised") }}
	got := around(t, it, fakeDiags{"W:careful", "E:the API said no"})
	want := fakeDiags{"W:careful", "E:normalised"}
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

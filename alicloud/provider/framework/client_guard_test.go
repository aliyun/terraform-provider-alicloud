package framework

import (
	"context"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

type carrier struct {
	fwadapt.DataSourceBase
}

type notCarrier struct{}

type observer struct{}

func (observer) Before(context.Context, intercept.Call) error { return nil }

func (observer) After(_ context.Context, _ intercept.Call, err error) error { return err }

func TestAppendUnlessClientCapable(t *testing.T) {
	if got := appendUnlessClientCapable(nil, "data source", "alicloud_unit_test", "DataSourceBase", &carrier{}); got != nil {
		t.Fatalf("a pointer to an embedder must pass the guard, got %v", got)
	}

	problems := appendUnlessClientCapable(nil, "data source", "alicloud_unit_test", "DataSourceBase", carrier{})
	if len(problems) != 1 {
		t.Fatalf("a value must fail the guard, got %d problems: %v", len(problems), problems)
	}

	problems = appendUnlessClientCapable(problems, "resource", "alicloud_unit_test_two", "ResourceBase", &notCarrier{})
	if len(problems) != 2 {
		t.Fatalf("got %d problems, want 2 (the guard accumulates rather than stopping at the first)", len(problems))
	}

	for _, want := range []string{"alicloud_unit_test_two", "fwadapt.ResourceBase", "return a pointer"} {
		if !strings.Contains(problems[1], want) {
			t.Errorf("problem %q does not mention %q", problems[1], want)
		}
	}
}

func TestAppendIfInterceptorsDropped(t *testing.T) {
	if got := appendIfInterceptorsDropped(nil, "action", "alicloud_unit_test", nil); got != nil {
		t.Fatalf("an empty chain must pass the guard, got %v", got)
	}
	if got := appendIfInterceptorsDropped(nil, "action", "alicloud_unit_test", []intercept.Interceptor{}); got != nil {
		t.Fatalf("a non-nil empty chain must pass the guard, got %v", got)
	}

	chain := []intercept.Interceptor{observer{}, observer{}}
	problems := appendIfInterceptorsDropped(nil, "action", "alicloud_unit_test", chain)
	if len(problems) != 1 {
		t.Fatalf("a non-empty chain must fail the guard, got %d problems: %v", len(problems), problems)
	}
	for _, want := range []string{"alicloud_unit_test", "never run", "wrapper"} {
		if !strings.Contains(problems[0], want) {
			t.Errorf("problem %q does not mention %q", problems[0], want)
		}
	}

	problems = appendUnlessClientCapable(problems, "action", "alicloud_unit_test_two", "ActionBase", &notCarrier{})
	if len(problems) != 2 {
		t.Fatalf("got %d problems, want 2 (the two guards accumulate into one list)", len(problems))
	}
}

func TestMustRegisterSilentWhenClean(t *testing.T) {
	mustRegister(nil)
	mustRegister([]string{})
}

func TestMustRegisterPanicsOnProblems(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("mustRegister did not panic on a non-empty problem list")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panicked with %T, want string", r)
		}
		if !strings.HasPrefix(msg, "provider registration error:") {
			t.Errorf("panic message %q does not carry the registration-error prefix", msg)
		}
		for _, want := range []string{"first problem", "second problem"} {
			if !strings.Contains(msg, want) {
				t.Errorf("panic message %q does not mention %q", msg, want)
			}
		}
	}()

	mustRegister([]string{"first problem", "second problem"})
}

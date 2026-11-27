package intercept

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
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

func TestAroundCanceledEntry(t *testing.T) {
	for _, withChain := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		var chain []Interceptor
		if withChain {
			chain = []Interceptor{&funcInterceptor{
				before: func() error { t.Fatal("预取消不应进入 Before"); return nil },
				after:  func(err error) error { t.Fatal("未进入的链不应执行 After"); return err },
			}}
		}
		got := Around(ctx, chain, testCall, func() fakeDiags {
			t.Fatal("预取消不应调用 inner")
			return nil
		}, fakeBridge)
		if !slices.Equal(got, fakeDiags{"E:context canceled"}) {
			t.Fatalf("取消诊断丢失：%v", got)
		}
	}
}

func TestAroundCancellationDuringBefore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var stages []string
	chain := []Interceptor{
		&funcInterceptor{before: func() error {
			stages = append(stages, "before-1")
			cancel()
			return nil
		}, after: func(err error) error {
			stages = append(stages, "after-1")
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("After 未收到取消：%v", err)
			}
			return nil
		}},
		&funcInterceptor{before: func() error {
			t.Fatal("取消后不应继续 Before")
			return nil
		}, after: func(err error) error {
			stages = append(stages, "after-2")
			return err
		}},
	}
	got := Around(ctx, chain, testCall, func() fakeDiags {
		t.Fatal("Before 中取消不应启动 CRUD")
		return nil
	}, fakeBridge)
	if !slices.Equal(stages, []string{"before-1", "after-2", "after-1"}) || !slices.Equal(got, fakeDiags{"E:context canceled"}) {
		t.Fatalf("取消后的收尾语义错误：stages=%v diags=%v", stages, got)
	}
}

func TestAroundCancellationDoesNotRewriteCompletedInner(t *testing.T) {
	for _, chain := range [][]Interceptor{nil, {&funcInterceptor{}}} {
		ctx, cancel := context.WithCancel(context.Background())
		got := Around(ctx, chain, testCall, func() fakeDiags {
			cancel()
			return fakeDiags{"W:completed"}
		}, fakeBridge)
		cancel()
		if !slices.Equal(got, fakeDiags{"W:completed"}) {
			t.Fatalf("完成后的取消不应重写 CRUD 结果：%v", got)
		}
	}
}

func TestBootstrapContextUsesDeclaredService(t *testing.T) {
	bootstrap := BootstrapContext("DeclaredService")
	for _, op := range []Op{OpCreate, OpRead, OpUpdate, OpDelete} {
		t.Run(string(op), func(t *testing.T) {
			var output bytes.Buffer
			ctx := bootstrap(tflogtest.RootLogger(context.Background(), &output), Call{
				Name: "alicloud_different_prefix", Op: op,
			})
			tflog.Info(ctx, "inner")
			entries, err := tflogtest.MultilineJSONDecode(&output)
			if err != nil || len(entries) != 1 {
				t.Fatalf("日志条目 = %v，错误 = %v", entries, err)
			}
			if entries[0]["alicloud_service"] != "DeclaredService" || entries[0]["alicloud_operation"] != string(op) {
				t.Fatalf("声明的 service 或操作字段丢失：%v", entries[0])
			}
		})
	}
}

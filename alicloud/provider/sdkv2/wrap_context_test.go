package sdkv2

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type contextRecorder struct {
	before, after context.Context
	beforeErr     error
}

func (r *contextRecorder) Before(ctx context.Context, _ intercept.Call) error {
	r.before = ctx
	return r.beforeErr
}

func (r *contextRecorder) After(ctx context.Context, _ intercept.Call, err error) error {
	r.after = ctx
	return err
}

func TestWrapResourceContextBootstrap(t *testing.T) {
	for _, withoutTimeout := range []bool{false, true} {
		for _, withChain := range []bool{false, true} {
			name := "Context"
			if withoutTimeout {
				name = "WithoutTimeout"
			}
			if withChain {
				name += "/chain"
			}
			t.Run(name, func(t *testing.T) {
				type key struct{}
				parent, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
				defer cancel()
				meta := &struct{}{}
				var seen context.Context
				inner := func(ctx context.Context, _ *schema.ResourceData, gotMeta interface{}) diag.Diagnostics {
					seen = ctx
					if gotMeta != meta {
						t.Fatal("meta 被替换")
					}
					return diag.Diagnostics{{Severity: diag.Warning, Summary: "original"}}
				}
				r := &schema.Resource{}
				if withoutTimeout {
					r.CreateWithoutTimeout, r.ReadWithoutTimeout = inner, inner
					r.UpdateWithoutTimeout, r.DeleteWithoutTimeout = inner, inner
				} else {
					r.CreateContext, r.ReadContext = inner, inner
					r.UpdateContext, r.DeleteContext = inner, inner
				}
				rec := &contextRecorder{}
				var chain []intercept.Interceptor
				if withChain {
					chain = []intercept.Interceptor{rec}
				}
				calls := 0
				bootstrap := func(ctx context.Context, call intercept.Call) context.Context {
					calls++
					if ctx != parent || call.Name != "alicloud_test" || call.Meta != meta {
						t.Fatalf("bootstrap 输入错误：%+v", call)
					}
					return context.WithValue(ctx, key{}, call.Op)
				}
				got := WrapResource("alicloud_test", r, chain, bootstrap)
				if got == r || got.Create != nil || got.Read != nil || got.Update != nil || got.Delete != nil {
					t.Fatal("必须包装 context-aware 字段，不能改变 CRUD 形态")
				}
				ops := []struct {
					op  intercept.Op
					run func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
				}{
					{intercept.OpCreate, got.CreateContext}, {intercept.OpRead, got.ReadContext},
					{intercept.OpUpdate, got.UpdateContext}, {intercept.OpDelete, got.DeleteContext},
				}
				if withoutTimeout {
					ops[0].run, ops[1].run = got.CreateWithoutTimeout, got.ReadWithoutTimeout
					ops[2].run, ops[3].run = got.UpdateWithoutTimeout, got.DeleteWithoutTimeout
				}
				for _, op := range ops {
					diags := op.run(parent, nil, meta)
					if len(diags) != 1 || diags[0].Severity != diag.Warning || diags[0].Summary != "original" {
						t.Fatalf("原始诊断未保留：%v", diags)
					}
					if seen == nil || seen.Value(key{}) != op.op || seen.Done() != parent.Done() {
						t.Fatalf("%s 未收到派生 ctx", op.op)
					}
					wantDeadline, _ := parent.Deadline()
					if deadline, ok := seen.Deadline(); !ok || !deadline.Equal(wantDeadline) {
						t.Fatal("原操作 deadline 被改变")
					}
					if withChain && (rec.before != seen || rec.after != seen) {
						t.Fatal("Before、inner、After 必须收到同一个派生 ctx")
					}
				}
				if calls != 4 {
					t.Fatalf("bootstrap 调用次数 = %d，期望每次操作一次", calls)
				}
				cancel()
				if seen.Err() != context.Canceled {
					t.Fatal("父取消未到达资源")
				}
			})
		}
	}
}

func TestWrapResourceBootstrapDoesNotPromotePlainCRUD(t *testing.T) {
	r := &schema.Resource{Create: func(*schema.ResourceData, interface{}) error { return nil }}
	bootstrap := func(context.Context, intercept.Call) context.Context {
		t.Fatal("plain CRUD 没有请求 ctx，不应调用 bootstrap")
		return nil
	}
	got := WrapResource("alicloud_test", r, nil, bootstrap)
	if got != r || got.CreateContext != nil || got.CreateWithoutTimeout != nil {
		t.Fatal("纯 plain 资源应保持原对象及分派形态")
	}
	if err := got.Create(nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestWrapDataSourceContextBootstrap(t *testing.T) {
	for _, withoutTimeout := range []bool{false, true} {
		type key struct{}
		calls := 0
		inner := func(ctx context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
			if ctx.Value(key{}) != intercept.OpRead {
				t.Fatal("数据源未收到 bootstrap 返回的 ctx")
			}
			return nil
		}
		ds := &schema.Resource{}
		if withoutTimeout {
			ds.ReadWithoutTimeout = inner
		} else {
			ds.ReadContext = inner
		}
		got := WrapDataSource("alicloud_test", ds, nil, func(ctx context.Context, call intercept.Call) context.Context {
			calls++
			return context.WithValue(ctx, key{}, call.Op)
		})
		run := got.ReadContext
		if withoutTimeout {
			run = got.ReadWithoutTimeout
		}
		if diags := run(context.Background(), nil, nil); diags.HasError() || calls != 1 {
			t.Fatalf("数据源调用失败：diags=%v calls=%d", diags, calls)
		}
	}
}

func TestWrapCanceledContext(t *testing.T) {
	for _, withoutTimeout := range []bool{false, true} {
		for _, withChain := range []bool{false, true} {
			for _, expired := range []bool{false, true} {
				ctx, cancel := context.WithCancel(context.Background())
				want := context.Canceled
				if expired {
					cancel()
					ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
					want = context.DeadlineExceeded
				}
				cancel()
				inner := func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
					t.Fatal("取消后不应调用 CRUD")
					return nil
				}
				r := &schema.Resource{CreateContext: inner, ReadContext: inner, UpdateContext: inner, DeleteContext: inner}
				if withoutTimeout {
					r = &schema.Resource{CreateWithoutTimeout: inner, ReadWithoutTimeout: inner, UpdateWithoutTimeout: inner, DeleteWithoutTimeout: inner}
				}
				rec := &contextRecorder{}
				var chain []intercept.Interceptor
				if withChain {
					chain = []intercept.Interceptor{rec}
				}
				wrapped := WrapResource("alicloud_test", r, chain, intercept.BootstrapContext("Test"))
				ds := WrapDataSource("alicloud_test", r, chain, intercept.BootstrapContext("Test"))
				ops := []func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics{wrapped.CreateContext, wrapped.ReadContext, wrapped.UpdateContext, wrapped.DeleteContext, ds.ReadContext}
				if withoutTimeout {
					ops = []func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics{wrapped.CreateWithoutTimeout, wrapped.ReadWithoutTimeout, wrapped.UpdateWithoutTimeout, wrapped.DeleteWithoutTimeout, ds.ReadWithoutTimeout}
				}
				for _, run := range ops {
					if diags := run(ctx, nil, nil); len(diags) != 1 || diags[0].Summary != want.Error() {
						t.Fatalf("取消诊断错误：%v", diags)
					}
				}
				if rec.before != nil || rec.after != nil {
					t.Fatal("预取消不应进入拦截链")
				}
			}
		}
	}
}

func TestWrapResourceBootstrapBeforeAbort(t *testing.T) {
	type key struct{}
	rec := &contextRecorder{beforeErr: errors.New("abort")}
	r := &schema.Resource{ReadContext: func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
		t.Fatal("Before 失败后不应调用 inner")
		return nil
	}}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec}, func(ctx context.Context, _ intercept.Call) context.Context {
		return context.WithValue(ctx, key{}, "prepared")
	})
	if diags := got.ReadContext(context.Background(), nil, nil); len(diags) != 1 || diags[0].Summary != "abort" {
		t.Fatalf("中止诊断改变：%v", diags)
	}
	if rec.before == nil || rec.before.Value(key{}) != "prepared" || rec.after != rec.before {
		t.Fatal("中止路径没有共享派生 ctx")
	}
}

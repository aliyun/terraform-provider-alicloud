package fwadapt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type contextResource struct {
	stubResource
	ResourceBase
	seen context.Context
}

func (r *contextResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	r.seen = ctx
	resp.Diagnostics.AddWarning("original", "detail")
}

func (r *contextResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	r.seen = ctx
	resp.Diagnostics.AddWarning("original", "detail")
}

func (r *contextResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.seen = ctx
	resp.Diagnostics.AddWarning("original", "detail")
}

func (r *contextResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.seen = ctx
	resp.Diagnostics.AddWarning("original", "detail")
}

type contextDataSource struct {
	stubDataSource
	DataSourceBase
	seen context.Context
}

func (ds *contextDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	ds.seen = ctx
	resp.Diagnostics.AddWarning("original", "detail")
}

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
	for _, withChain := range []bool{false, true} {
		name := "empty"
		if withChain {
			name = "chain"
		}
		t.Run(name, func(t *testing.T) {
			type key struct{}
			parent, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
			defer cancel()
			meta := &connectivity.AliyunClient{}
			inner := &contextResource{}
			rec := &contextRecorder{}
			var chain []intercept.Interceptor
			if withChain {
				chain = []intercept.Interceptor{rec}
			}
			calls := 0
			got := WrapResource("alicloud_test", inner, chain, func(ctx context.Context, call intercept.Call) context.Context {
				calls++
				if ctx != parent || call.Name != "alicloud_test" || call.Meta != meta {
					t.Fatalf("bootstrap 输入错误：%+v", call)
				}
				return context.WithValue(ctx, key{}, call.Op)
			})
			var configured resource.ConfigureResponse
			got.(resource.ResourceWithConfigure).Configure(context.Background(), resource.ConfigureRequest{ProviderData: meta}, &configured)
			if configured.Diagnostics.HasError() || inner.Meta() != meta || calls != 0 {
				t.Fatalf("Configure 不应补操作 ctx 或改变 meta：%v", configured.Diagnostics)
			}
			ops := []struct {
				op  intercept.Op
				run func() diag.Diagnostics
			}{
				{intercept.OpCreate, func() diag.Diagnostics {
					var resp resource.CreateResponse
					got.Create(parent, resource.CreateRequest{}, &resp)
					return resp.Diagnostics
				}},
				{intercept.OpRead, func() diag.Diagnostics {
					var resp resource.ReadResponse
					got.Read(parent, resource.ReadRequest{}, &resp)
					return resp.Diagnostics
				}},
				{intercept.OpUpdate, func() diag.Diagnostics {
					var resp resource.UpdateResponse
					got.Update(parent, resource.UpdateRequest{}, &resp)
					return resp.Diagnostics
				}},
				{intercept.OpDelete, func() diag.Diagnostics {
					var resp resource.DeleteResponse
					got.Delete(parent, resource.DeleteRequest{}, &resp)
					return resp.Diagnostics
				}},
			}
			for _, op := range ops {
				diags := op.run()
				if len(diags) != 1 || diags[0].Severity() != diag.SeverityWarning || diags[0].Summary() != "original" || diags[0].Detail() != "detail" {
					t.Fatalf("原始诊断丢失：%v", diags)
				}
				if inner.seen == nil || inner.seen.Value(key{}) != op.op || inner.seen.Done() != parent.Done() {
					t.Fatalf("%s 未收到派生 ctx", op.op)
				}
				deadline, ok := inner.seen.Deadline()
				want, _ := parent.Deadline()
				if !ok || !deadline.Equal(want) {
					t.Fatal("父 deadline 丢失")
				}
				if withChain && (rec.before != inner.seen || rec.after != inner.seen) {
					t.Fatal("Before、inner、After 的 ctx 不一致")
				}
			}
			if calls != 4 {
				t.Fatalf("bootstrap 调用次数 = %d，期望 4", calls)
			}
			cancel()
			if inner.seen.Err() != context.Canceled {
				t.Fatal("父取消未传播")
			}
		})
	}
}

func TestWrapDataSourceContextBootstrap(t *testing.T) {
	type key struct{}
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	inner := &contextDataSource{}
	calls := 0
	got := WrapDataSource("alicloud_test", inner, nil, func(ctx context.Context, call intercept.Call) context.Context {
		calls++
		if call.Name != "alicloud_test" || call.Meta != nil {
			t.Fatalf("数据源 bootstrap 输入错误：%+v", call)
		}
		return context.WithValue(ctx, key{}, call.Op)
	})
	var configured datasource.ConfigureResponse
	got.(datasource.DataSourceWithConfigure).Configure(parent, datasource.ConfigureRequest{}, &configured)
	if configured.Diagnostics.HasError() || calls != 0 {
		t.Fatalf("nil ProviderData 离线配置失败：%v", configured.Diagnostics)
	}
	var resp datasource.ReadResponse
	got.Read(parent, datasource.ReadRequest{}, &resp)
	if calls != 1 || inner.seen == nil || inner.seen.Value(key{}) != intercept.OpRead || len(resp.Diagnostics) != 1 || resp.Diagnostics[0].Summary() != "original" {
		t.Fatalf("空链数据源 ctx 或诊断丢失：calls=%d diags=%v", calls, resp.Diagnostics)
	}
	cancel()
	if inner.seen.Err() != context.Canceled {
		t.Fatal("父取消未到达数据源")
	}
}

func TestWrapCanceledContext(t *testing.T) {
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
			inner, source := &contextResource{}, &contextDataSource{}
			rec := &contextRecorder{}
			var chain []intercept.Interceptor
			if withChain {
				chain = []intercept.Interceptor{rec}
			}
			r := WrapResource("alicloud_test", inner, chain, intercept.BootstrapContext("Test"))
			ds := WrapDataSource("alicloud_test", source, chain, intercept.BootstrapContext("Test"))
			var create resource.CreateResponse
			var read resource.ReadResponse
			var update resource.UpdateResponse
			var remove resource.DeleteResponse
			var data datasource.ReadResponse
			r.Create(ctx, resource.CreateRequest{}, &create)
			r.Read(ctx, resource.ReadRequest{}, &read)
			r.Update(ctx, resource.UpdateRequest{}, &update)
			r.Delete(ctx, resource.DeleteRequest{}, &remove)
			ds.Read(ctx, datasource.ReadRequest{}, &data)
			for _, diagnostics := range []diag.Diagnostics{create.Diagnostics, read.Diagnostics, update.Diagnostics, remove.Diagnostics, data.Diagnostics} {
				if len(diagnostics) != 1 || !diagnostics.HasError() || diagnostics[0].Summary() != want.Error() {
					t.Fatalf("取消诊断错误：%v", diagnostics)
				}
			}
			if inner.seen != nil || source.seen != nil || rec.before != nil || rec.after != nil {
				t.Fatal("预取消仍进入了拦截链或 CRUD")
			}
		}
	}
}

func TestRunInterceptedPreservesExistingDiagnostics(t *testing.T) {
	for _, scenario := range []string{"success", "canceled", "before-abort"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			rec := &contextRecorder{}
			if scenario == "canceled" {
				cancel()
			}
			if scenario == "before-abort" {
				rec.beforeErr = errors.New("abort")
			}
			var diags diag.Diagnostics
			diags.AddWarning("existing", "detail")
			called := false
			runIntercepted(ctx, "alicloud_test", intercept.OpRead, []intercept.Interceptor{rec}, nil, &diags, nil, func(context.Context) {
				called = true
				diags.AddWarning("inner", "detail")
			})
			if len(diags) != 2 || diags[0].Summary() != "existing" || diags[0].Severity() != diag.SeverityWarning {
				t.Fatalf("已有诊断丢失或重复：%v", diags)
			}
			if called != (scenario == "success") || diags.HasError() != (scenario != "success") {
				t.Fatalf("调用或诊断结果错误：called=%v diags=%v", called, diags)
			}
		})
	}
}

func TestWrapResourceBootstrapBeforeAbort(t *testing.T) {
	type key struct{}
	inner := &contextResource{}
	rec := &contextRecorder{beforeErr: errors.New("abort")}
	calls := 0
	got := WrapResource("alicloud_test", inner, []intercept.Interceptor{rec}, func(ctx context.Context, _ intercept.Call) context.Context {
		calls++
		return context.WithValue(ctx, key{}, "prepared")
	})
	var resp resource.ReadResponse
	got.Read(context.Background(), resource.ReadRequest{}, &resp)
	if inner.seen != nil || calls != 1 || len(resp.Diagnostics) != 1 || resp.Diagnostics[0].Summary() != "abort" {
		t.Fatalf("Before 中止行为改变：calls=%d diags=%v", calls, resp.Diagnostics)
	}
	if rec.before == nil || rec.before.Value(key{}) != "prepared" || rec.after != rec.before {
		t.Fatal("中止路径没有共享派生 ctx")
	}
}

func TestWrapResourceBootstrapPreservesCapabilities(t *testing.T) {
	for _, inner := range []resource.Resource{&stubResource{}, &importStateStub{}, &identityStub{}, &importStateAndIdentityStub{}} {
		got := WrapResource("alicloud_test", inner, nil, func(context.Context, intercept.Call) context.Context {
			t.Fatal("非 CRUD 生命周期不应执行 bootstrap")
			return nil
		})
		ctx := context.Background()
		_, wantImport := inner.(resource.ResourceWithImportState)
		withImport, hasImport := got.(resource.ResourceWithImportState)
		_, wantIdentity := inner.(resource.ResourceWithIdentity)
		withIdentity, hasIdentity := got.(resource.ResourceWithIdentity)
		if hasImport != wantImport || hasIdentity != wantIdentity {
			t.Fatalf("%T 的可选能力被改变", inner)
		}
		got.Schema(ctx, resource.SchemaRequest{}, &resource.SchemaResponse{})
		got.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{}, &resource.ConfigureResponse{})
		got.(resource.ResourceWithValidateConfig).ValidateConfig(ctx, resource.ValidateConfigRequest{}, &resource.ValidateConfigResponse{})
		got.(resource.ResourceWithModifyPlan).ModifyPlan(ctx, resource.ModifyPlanRequest{}, &resource.ModifyPlanResponse{})
		if hasImport {
			withImport.ImportState(ctx, resource.ImportStateRequest{}, &resource.ImportStateResponse{})
		}
		if hasIdentity {
			withIdentity.IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &resource.IdentitySchemaResponse{})
		}
		if innerStub(inner).importStateCalled != wantImport || innerStub(inner).identitySchemaCalled != wantIdentity {
			t.Fatalf("%T 的可选能力未转发", inner)
		}
	}
}

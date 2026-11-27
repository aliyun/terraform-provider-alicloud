package provider_test

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type wireIncomingContextKey struct{}

type wireContextRecorder struct {
	t            *testing.T
	output       bytes.Buffer
	meta         *connectivity.AliyunClient
	wantName     string
	wantDeadline bool
	deadline     time.Time
	prepared     context.Context
	stages       []string
	cancel       context.CancelFunc
}

func (r *wireContextRecorder) bootstrap(ctx context.Context, call intercept.Call) context.Context {
	r.t.Helper()
	r.stages = append(r.stages, "bootstrap")
	if ctx.Value(wireIncomingContextKey{}) != r {
		r.t.Fatal("协议入口的 incoming value 未到达 bootstrap")
	}
	deadline, ok := ctx.Deadline()
	if ok != r.wantDeadline {
		r.t.Fatalf("SDK 原有 timeout 分派改变：有 deadline=%t，期望 %t", ok, r.wantDeadline)
	}
	r.deadline = deadline
	r.prepared = client.RequestContext(ctx, "DeclaredWireService", string(call.Op))
	r.checkCall(call)
	return r.prepared
}

func (r *wireContextRecorder) checkCall(call intercept.Call) {
	r.t.Helper()
	if call.Name != r.wantName || call.Op != intercept.OpRead || call.Meta != r.meta {
		r.t.Fatalf("操作信息或 meta 丢失：%+v", call)
	}
}

func (r *wireContextRecorder) record(ctx context.Context, call intercept.Call, stage string) {
	r.t.Helper()
	r.checkCall(call)
	r.stages = append(r.stages, stage)
	if ctx != r.prepared || ctx.Value(wireIncomingContextKey{}) != r {
		r.t.Fatalf("%s 未使用本次 bootstrap 返回的 ctx", stage)
	}
	deadline, ok := ctx.Deadline()
	if ok != r.wantDeadline || !deadline.Equal(r.deadline) {
		r.t.Fatalf("%s 改变了 SDK 已有 deadline", stage)
	}
}

func (r *wireContextRecorder) Before(ctx context.Context, call intercept.Call) error {
	r.record(ctx, call, "before")
	return nil
}

func (r *wireContextRecorder) After(ctx context.Context, call intercept.Call, err error) error {
	r.record(ctx, call, "after")
	if r.cancel != nil && (ctx.Err() != context.Canceled || err == nil) {
		r.t.Fatal("取消状态或 inner 错误未到达 After")
	}
	return err
}

func (r *wireContextRecorder) inner(ctx context.Context, call intercept.Call) error {
	r.record(ctx, call, "inner")
	tflog.Info(ctx, "request-context-inner")
	if ctx.Err() != nil {
		r.t.Fatalf("取消前 ctx 已失效：%v", ctx.Err())
	}
	if r.cancel != nil {
		r.cancel()
		if ctx.Err() != context.Canceled {
			r.t.Fatal("操作执行期间的父取消未传播到 inner")
		}
		return ctx.Err()
	}
	return nil
}

func (r *wireContextRecorder) assertResult(withChain bool, diags []*tfprotov5.Diagnostic) {
	r.t.Helper()
	want := []string{"bootstrap", "inner"}
	if withChain {
		want = []string{"bootstrap", "before", "inner", "after"}
	}
	if !slices.Equal(r.stages, want) {
		r.t.Fatalf("调用顺序或次数错误：got=%v want=%v", r.stages, want)
	}
	if r.cancel == nil {
		requireNoDiagErrors(r.t, "Read", diags)
	} else {
		found := false
		for _, d := range diags {
			if d.Severity == tfprotov5.DiagnosticSeverityError && d.Summary == context.Canceled.Error() {
				found = true
			}
		}
		if !found {
			r.t.Fatalf("取消诊断未穿过协议层：%v", diags)
		}
	}
	entries, err := tflogtest.MultilineJSONDecode(&r.output)
	if err != nil {
		r.t.Fatal(err)
	}
	count := 0
	for _, entry := range entries {
		if entry["@message"] != "request-context-inner" {
			continue
		}
		count++
		if entry["alicloud_service"] != "DeclaredWireService" || entry["alicloud_operation"] != "Read" || entry["incoming"] != "kept" {
			r.t.Fatalf("协议调用中的请求日志字段丢失：%v", entry)
		}
	}
	if count != 1 {
		r.t.Fatalf("inner 日志数量=%d，期望 1", count)
	}
}

func newWireRequestContext(t *testing.T, name string, hasDeadline, cancelInCall bool) (*wireContextRecorder, context.Context) {
	t.Helper()
	rec := &wireContextRecorder{t: t, meta: &connectivity.AliyunClient{}, wantName: name, wantDeadline: hasDeadline}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if cancelInCall {
		rec.cancel = cancel
	}
	ctx = context.WithValue(ctx, wireIncomingContextKey{}, rec)
	ctx = tflogtest.RootLogger(ctx, &rec.output)
	return rec, tflog.SetField(ctx, "incoming", "kept")
}

func TestUnitProviderWireRequestContextSDKv2(t *testing.T) {
	for _, mode := range []string{"Context", "WithoutTimeout"} {
		for _, dataSource := range []bool{false, true} {
			for _, withChain := range []bool{false, true} {
				for _, cancelInCall := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/dataSource=%t/chain=%t/cancel=%t", mode, dataSource, withChain, cancelInCall), func(t *testing.T) {
						const name = "alicloud_wire_context"
						rec, ctx := newWireRequestContext(t, name, mode == "Context", cancelInCall)
						inner := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
							if err := rec.inner(ctx, intercept.Call{Name: name, Op: intercept.OpRead, Meta: meta}); err != nil {
								return diag.FromErr(err)
							}
							d.SetId("unit-test")
							return nil
						}
						r := &schema.Resource{Schema: map[string]*schema.Schema{
							"name": {Type: schema.TypeString, Required: true},
						}}
						if mode == "Context" {
							r.ReadContext = inner
						} else {
							r.ReadWithoutTimeout = inner
							r.ReadContext = func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
								t.Error("WithoutTimeout 存在时不应分派到 Context")
								return diag.Errorf("错误分派")
							}
						}
						var chain []intercept.Interceptor
						if withChain {
							chain = []intercept.Interceptor{rec}
						}
						p := &schema.Provider{}
						if dataSource {
							p.DataSourcesMap = map[string]*schema.Resource{name: sdkv2.WrapDataSource(name, r, chain, rec.bootstrap)}
						} else {
							p.ResourcesMap = map[string]*schema.Resource{name: sdkv2.WrapResource(name, r, chain, rec.bootstrap)}
						}
						p.SetMeta(rec.meta)
						diags, state := readWireContext(t, p.GRPCProvider(), ctx, name, dataSource)
						rec.assertResult(withChain, diags)
						if !cancelInCall && stringOf(t, attrsOf(t, maximalFrameworkType, state)["id"]) != "unit-test" {
							t.Fatal("Read 的 state 丢失")
						}
					})
				}
			}
		}
	}
}

func TestUnitProviderWirePlainReadDoesNotBootstrap(t *testing.T) {
	for _, dataSource := range []bool{false, true} {
		t.Run(fmt.Sprintf("dataSource=%t", dataSource), func(t *testing.T) {
			const name = "alicloud_wire_plain"
			calls := 0
			meta := &connectivity.AliyunClient{}
			unreachable := func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				t.Error("plain Read 存在时不应执行 context-aware 入口")
				return diag.Errorf("错误分派")
			}
			r := &schema.Resource{
				Schema: map[string]*schema.Schema{"name": {Type: schema.TypeString, Required: true}},
				Read: func(d *schema.ResourceData, gotMeta interface{}) error {
					calls++
					if gotMeta != meta {
						t.Fatal("plain Read 的 meta 改变")
					}
					d.SetId("unit-test")
					return nil
				},
				ReadContext: unreachable, ReadWithoutTimeout: unreachable,
			}
			bootstrap := func(ctx context.Context, _ intercept.Call) context.Context {
				t.Error("plain Read 不应调用 bootstrap")
				return ctx
			}
			chain := &chainRecorder{}
			p := &schema.Provider{}
			if dataSource {
				p.DataSourcesMap = map[string]*schema.Resource{name: sdkv2.WrapDataSource(name, r, []intercept.Interceptor{chain}, bootstrap)}
			} else {
				p.ResourcesMap = map[string]*schema.Resource{name: sdkv2.WrapResource(name, r, []intercept.Interceptor{chain}, bootstrap)}
			}
			p.SetMeta(meta)
			diags, state := readWireContext(t, p.GRPCProvider(), context.Background(), name, dataSource)
			requireNoDiagErrors(t, "plain Read", diags)
			if calls != 1 || !slices.Equal(chain.beforeOps(), []string{"Read"}) || !slices.Equal(chain.afterOps(), []string{"Read"}) {
				t.Fatalf("plain Read 或拦截器分派改变：calls=%d before=%v after=%v", calls, chain.beforeOps(), chain.afterOps())
			}
			if stringOf(t, attrsOf(t, maximalFrameworkType, state)["id"]) != "unit-test" {
				t.Fatal("plain Read 的 state 丢失")
			}
		})
	}
}

type wireContextResource struct {
	*maximalFrameworkResource
	recorder *wireContextRecorder
}

func (r *wireContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if err := r.recorder.inner(ctx, intercept.Call{Name: maximalFrameworkName, Op: intercept.OpRead, Meta: r.Meta()}); err != nil {
		resp.Diagnostics.AddError(err.Error(), "父请求已取消")
		return
	}
	r.maximalFrameworkResource.Read(ctx, req, resp)
}

type wireContextDataSource struct {
	*maximalFrameworkDataSource
	recorder *wireContextRecorder
}

func (d *wireContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if err := d.recorder.inner(ctx, intercept.Call{Name: maximalFrameworkDataSourceName, Op: intercept.OpRead, Meta: d.Meta()}); err != nil {
		resp.Diagnostics.AddError(err.Error(), "父请求已取消")
		return
	}
	d.maximalFrameworkDataSource.Read(ctx, req, resp)
}

type wireContextProvider struct {
	maximalFrameworkTestProvider
	recorder *wireContextRecorder
}

func (p *wireContextProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{func() resource.Resource {
		r := &wireContextResource{newMaximalFrameworkResource(p.log), p.recorder}
		return fwadapt.WrapResource(maximalFrameworkName, r, p.chain, p.recorder.bootstrap)
	}}
}

func (p *wireContextProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{func() datasource.DataSource {
		d := &wireContextDataSource{newMaximalFrameworkDataSource(p.log), p.recorder}
		return fwadapt.WrapDataSource(maximalFrameworkDataSourceName, d, p.chain, p.recorder.bootstrap)
	}}
}

func TestUnitProviderWireRequestContextFramework(t *testing.T) {
	for _, dataSource := range []bool{false, true} {
		for _, withChain := range []bool{false, true} {
			for _, cancelInCall := range []bool{false, true} {
				t.Run(fmt.Sprintf("dataSource=%t/chain=%t/cancel=%t", dataSource, withChain, cancelInCall), func(t *testing.T) {
					name := maximalFrameworkName
					if dataSource {
						name = maximalFrameworkDataSourceName
					}
					rec, ctx := newWireRequestContext(t, name, false, cancelInCall)
					p := &wireContextProvider{
						maximalFrameworkTestProvider: maximalFrameworkTestProvider{log: &hookLog{}, client: rec.meta},
						recorder:                     rec,
					}
					if withChain {
						p.chain = []intercept.Interceptor{rec}
					}
					server := providerserver.NewProtocol5(p)()
					configured, err := server.ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{
						Config: dynamic(t, emptyObjectType, tftypes.NewValue(emptyObjectType, map[string]tftypes.Value{})),
					})
					if err != nil {
						t.Fatal(err)
					}
					requireNoDiagErrors(t, "ConfigureProvider", configured.Diagnostics)
					if len(rec.stages) != 0 {
						t.Fatal("Configure 不应执行操作 bootstrap")
					}
					diags, state := readWireContext(t, server, ctx, name, dataSource)
					rec.assertResult(withChain, diags)
					if !cancelInCall && stringOf(t, attrsOf(t, maximalFrameworkType, state)["id"]) != "unit-test" {
						t.Fatal("Framework Read 的 state 丢失")
					}
				})
			}
		}
	}
}

func readWireContext(t *testing.T, server tfprotov5.ProviderServer, ctx context.Context, name string, dataSource bool) ([]*tfprotov5.Diagnostic, *tfprotov5.DynamicValue) {
	t.Helper()
	attrs := map[string]tftypes.Value{"name": tftypes.NewValue(tftypes.String, "fixture")}
	if dataSource {
		resp, err := server.ReadDataSource(ctx, &tfprotov5.ReadDataSourceRequest{
			TypeName: name, Config: dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, attrs)),
		})
		if err != nil {
			t.Fatal(err)
		}
		return resp.Diagnostics, resp.State
	}
	attrs["id"] = tftypes.NewValue(tftypes.String, "unit-test")
	resp, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName: name, CurrentState: dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, attrs)),
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Diagnostics, resp.NewState
}

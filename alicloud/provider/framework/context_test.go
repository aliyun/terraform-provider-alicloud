package framework

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	providerregistry "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProviderConfigureContext(t *testing.T) {
	meta, err := connectivity.NewClient(&connectivity.Config{
		RegionId: "cn-hangzhou", Region: connectivity.Region("cn-hangzhou"),
		AccountType: "Domestic", AccessKey: "test-ak", SecretKey: "test-sk", TerraformTraceId: "provider-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	primary := &sdkschema.Provider{}
	primary.SetMeta(meta)
	p := NewProvider(primary)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var canceled provider.ConfigureResponse
	p.Configure(ctx, provider.ConfigureRequest{}, &canceled)
	if !canceled.Diagnostics.HasError() || canceled.ResourceData != nil || canceled.DataSourceData != nil {
		t.Fatalf("预取消 Configure 仍发布 meta：%v", canceled.Diagnostics)
	}
	var output bytes.Buffer
	incoming := tflogtest.RootLogger(context.Background(), &output)
	incoming = tflog.SetField(incoming, "tf_req_id", "framework-configure")
	var resp provider.ConfigureResponse
	p.Configure(incoming, provider.ConfigureRequest{}, &resp)
	if resp.Diagnostics.HasError() || resp.ResourceData != meta || resp.DataSourceData != meta || resp.EphemeralResourceData != meta || resp.ListResourceData != meta || resp.ActionData != meta {
		t.Fatalf("共享 meta 分发失败：%v", resp.Diagnostics)
	}
	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil || len(entries) == 0 {
		t.Fatalf("缺少 Framework Configure 日志：%v", err)
	}
	entry := entries[len(entries)-1]
	if entry["alicloud_trace_id"] != "provider-1" || entry["alicloud_region"] != "cn-hangzhou" || entry["alicloud_operation"] != "Configure" || entry["tf_req_id"] != "framework-configure" {
		t.Fatalf("Framework Configure 公共上下文缺失：%v", entry)
	}
}

func TestCollectPreservesServiceDeclaration(t *testing.T) {
	var want []string
	for _, sp := range providerregistry.ServicePackages() {
		for _, d := range sp.Resources {
			want = append(want, sp.Name+":"+d.TypeName)
		}
	}
	factories := collect(
		func(sp conns.ServicePackage) []conns.Resource { return sp.Resources },
		func(service string, d conns.Resource) func() string {
			return func() string { return service + ":" + d.TypeName }
		})
	var got []string
	for _, factory := range factories {
		got = append(got, factory())
	}
	if len(want) == 0 || !slices.Equal(got, want) {
		t.Fatalf("service 声明或工厂闭包串值：got=%v want=%v", got, want)
	}
}

func TestNewProviderRegistryRequestContext(t *testing.T) {
	packages := providerregistry.ServicePackages()
	if len(packages) < 2 {
		t.Fatal("registry 至少需要两个 service entry")
	}
	saved := [2]conns.ServicePackage{packages[0], packages[1]}
	t.Cleanup(func() {
		packages[0], packages[1] = saved[0], saved[1]
	})
	var seen context.Context
	meta := &connectivity.AliyunClient{}
	observe := func(ctx context.Context, gotMeta *connectivity.AliyunClient, fixture string) {
		if gotMeta != meta {
			t.Fatal("Configure 注入的 meta 未原样到达操作回调")
		}
		seen = ctx
		tflog.Info(ctx, "callback", map[string]interface{}{"fixture": fixture})
	}
	specs := []struct {
		service string
		name    string
	}{
		{"DeclaredAlpha", "alicloud_request_context_fixture_one"},
		{"DeclaredBeta", "alicloud_request_context_fixture_two"},
	}
	for i, spec := range specs {
		if len(intercept.ChainOf(spec.name, nil)) != 0 {
			t.Fatalf("%s 必须使用空拦截器链", spec.name)
		}
		packages[i] = conns.ServicePackage{
			Name: spec.service,
			Resources: []conns.Resource{{
				TypeName: spec.name,
				Factory: func() resource.Resource {
					return &registryContextResource{name: spec.name, observe: observe}
				},
			}},
			DataSources: []conns.DataSource{{
				TypeName: spec.name,
				Factory: func() datasource.DataSource {
					return &registryContextDataSource{name: spec.name, observe: observe}
				},
			}},
		}
	}
	p := NewProvider(nil)
	resourceFactories := p.Resources(context.Background())
	dataSourceFactories := p.DataSources(context.Background())
	if len(resourceFactories) < 2 || len(dataSourceFactories) < 2 {
		t.Fatal("fixture 未经真实 provider 装配")
	}
	for _, index := range []int{1, 0, 1} {
		spec := specs[index]
		r := resourceFactories[index]()
		ds := dataSourceFactories[index]()
		rc, ok := r.(resource.ResourceWithConfigure)
		if !ok {
			t.Fatalf("%T 不支持 Configure", r)
		}
		var resourceConfigured resource.ConfigureResponse
		rc.Configure(context.Background(), resource.ConfigureRequest{ProviderData: meta}, &resourceConfigured)
		if resourceConfigured.Diagnostics.HasError() {
			t.Fatal(resourceConfigured.Diagnostics)
		}
		dc, ok := ds.(datasource.DataSourceWithConfigure)
		if !ok {
			t.Fatalf("%T 不支持 Configure", ds)
		}
		var dataSourceConfigured datasource.ConfigureResponse
		dc.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: meta}, &dataSourceConfigured)
		if dataSourceConfigured.Diagnostics.HasError() {
			t.Fatal(dataSourceConfigured.Diagnostics)
		}
		ops := []struct {
			op, kind string
			run      func(context.Context) fwdiag.Diagnostics
		}{
			{"Create", "resource", func(ctx context.Context) fwdiag.Diagnostics {
				var resp resource.CreateResponse
				r.Create(ctx, resource.CreateRequest{}, &resp)
				return resp.Diagnostics
			}},
			{"Read", "resource", func(ctx context.Context) fwdiag.Diagnostics {
				var resp resource.ReadResponse
				r.Read(ctx, resource.ReadRequest{}, &resp)
				return resp.Diagnostics
			}},
			{"Update", "resource", func(ctx context.Context) fwdiag.Diagnostics {
				var resp resource.UpdateResponse
				r.Update(ctx, resource.UpdateRequest{}, &resp)
				return resp.Diagnostics
			}},
			{"Delete", "resource", func(ctx context.Context) fwdiag.Diagnostics {
				var resp resource.DeleteResponse
				r.Delete(ctx, resource.DeleteRequest{}, &resp)
				return resp.Diagnostics
			}},
			{"Read", "datasource", func(ctx context.Context) fwdiag.Diagnostics {
				var resp datasource.ReadResponse
				ds.Read(ctx, datasource.ReadRequest{}, &resp)
				return resp.Diagnostics
			}},
		}
		for _, op := range ops {
			assertAssembledRequestContext(t, spec.service, spec.name+"/"+op.kind, op.op, &seen,
				func(ctx context.Context) {
					if diags := op.run(ctx); diags.HasError() {
						t.Fatal(diags)
					}
				})
		}
	}
}

type registryContextResource struct {
	fwadapt.ResourceBase
	name    string
	observe func(context.Context, *connectivity.AliyunClient, string)
}

func (r *registryContextResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.name
}

func (*registryContextResource) Schema(context.Context, resource.SchemaRequest, *resource.SchemaResponse) {
}

func (r *registryContextResource) Create(ctx context.Context, _ resource.CreateRequest, _ *resource.CreateResponse) {
	r.observe(ctx, r.Meta(), r.name+"/resource")
}

func (r *registryContextResource) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	r.observe(ctx, r.Meta(), r.name+"/resource")
}

func (r *registryContextResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	r.observe(ctx, r.Meta(), r.name+"/resource")
}

func (r *registryContextResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	r.observe(ctx, r.Meta(), r.name+"/resource")
}

type registryContextDataSource struct {
	fwadapt.DataSourceBase
	name    string
	observe func(context.Context, *connectivity.AliyunClient, string)
}

func (ds *registryContextDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = ds.name
}

func (*registryContextDataSource) Schema(context.Context, datasource.SchemaRequest, *datasource.SchemaResponse) {
}

func (ds *registryContextDataSource) Read(ctx context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
	ds.observe(ctx, ds.Meta(), ds.name+"/datasource")
}

func assertAssembledRequestContext(t *testing.T, service, fixture, operation string, seen *context.Context, run func(context.Context)) {
	t.Helper()
	type key struct{}
	token := &struct{ value int }{value: 1}
	deadline := time.Now().Add(time.Hour)
	parent, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	parent = context.WithValue(parent, key{}, token)
	var output bytes.Buffer
	incoming := tflogtest.RootLogger(parent, &output)
	incoming = tflog.SetField(incoming, "incoming", "kept")
	incoming = tflog.SetField(incoming, "alicloud_service", "stale-service")
	incoming = tflog.SetField(incoming, "alicloud_operation", "stale-operation")
	*seen = nil
	run(incoming)
	got := *seen
	if got == nil {
		t.Fatalf("%s/%s 未到达操作回调", fixture, operation)
	}
	if got.Value(key{}) != token || got.Done() != incoming.Done() {
		t.Fatalf("%s/%s 丢失本次 incoming value 或取消通道", fixture, operation)
	}
	if actual, ok := got.Deadline(); !ok || !actual.Equal(deadline) {
		t.Fatalf("%s/%s 改变了 incoming deadline", fixture, operation)
	}
	if err := got.Err(); err != nil {
		t.Fatalf("父取消前回调 ctx 已失效：%v", err)
	}
	cancel()
	if !errors.Is(got.Err(), context.Canceled) {
		t.Fatalf("%s/%s 未传播父取消：%v", fixture, operation, got.Err())
	}
	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil || len(entries) != 1 {
		t.Fatalf("%s/%s 日志条目=%v，错误=%v", fixture, operation, entries, err)
	}
	entry := entries[0]
	if entry["fixture"] != fixture || entry["alicloud_service"] != service || entry["alicloud_operation"] != operation || entry["incoming"] != "kept" {
		t.Fatalf("%s/%s 日志上下文错误：%v", fixture, operation, entry)
	}
}

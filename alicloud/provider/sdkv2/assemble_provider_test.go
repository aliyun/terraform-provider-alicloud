package sdkv2

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	providerregistry "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAssembleProviderIdentityOnEmptyRegistries(t *testing.T) {
	res := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { return nil },
	}
	ds := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error { return nil },
	}
	p := &schema.Provider{
		ResourcesMap:   map[string]*schema.Resource{"alicloud_test": res},
		DataSourcesMap: map[string]*schema.Resource{"alicloud_test": ds},
	}

	diags := assembleProvider(p)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if p.ResourcesMap["alicloud_test"] != res {
		t.Fatal("resource pointer changed under empty registries")
	}
	if p.DataSourcesMap["alicloud_test"] != ds {
		t.Fatal("data source pointer changed under empty registries")
	}
}

func TestAssembleProviderBootstrapsLegacyContext(t *testing.T) {
	var output bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &output)
	inner := func(ctx context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
		tflog.Info(ctx, "inner")
		return nil
	}
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_context_test": {CreateContext: inner},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_context_test": {ReadContext: inner},
		},
	}
	if diags := assembleProvider(p); diags.HasError() {
		t.Fatal(diags)
	}
	p.ResourcesMap["alicloud_context_test"].CreateContext(ctx, nil, nil)
	p.DataSourcesMap["alicloud_context_test"].ReadContext(ctx, nil, nil)
	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil || len(entries) != 2 {
		t.Fatalf("日志条目 = %v，错误 = %v", entries, err)
	}
	for i, op := range []string{"Create", "Read"} {
		if entries[i]["alicloud_operation"] != op {
			t.Fatalf("存量 map 未接入操作 ctx：%v", entries[i])
		}
		if _, ok := entries[i]["alicloud_service"]; ok {
			t.Fatalf("不能从存量 Terraform 名称猜 service：%v", entries[i])
		}
	}
}

func TestAssembleProviderRegistryRequestContext(t *testing.T) {
	for _, mode := range []string{"Context", "WithoutTimeout"} {
		t.Run(mode, func(t *testing.T) {
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
			callback := func(fixture string) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
				return func(ctx context.Context, _ *schema.ResourceData, gotMeta interface{}) diag.Diagnostics {
					if gotMeta != meta {
						t.Fatal("操作回调的 meta 被改变")
					}
					seen = ctx
					tflog.Info(ctx, "callback", map[string]interface{}{"fixture": fixture})
					return nil
				}
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
					SDKResources: []conns.SDKResource{{
						TypeName: spec.name,
						Factory: func() *schema.Resource {
							inner := callback(spec.name + "/resource")
							r := &schema.Resource{}
							if mode == "WithoutTimeout" {
								r.CreateWithoutTimeout, r.ReadWithoutTimeout = inner, inner
								r.UpdateWithoutTimeout, r.DeleteWithoutTimeout = inner, inner
							} else {
								r.CreateContext, r.ReadContext = inner, inner
								r.UpdateContext, r.DeleteContext = inner, inner
							}
							return r
						},
					}},
					SDKDataSources: []conns.SDKDataSource{{
						TypeName: spec.name,
						Factory: func() *schema.Resource {
							inner := callback(spec.name + "/datasource")
							ds := &schema.Resource{}
							if mode == "WithoutTimeout" {
								ds.ReadWithoutTimeout = inner
							} else {
								ds.ReadContext = inner
							}
							return ds
						},
					}},
				}
			}
			p := &schema.Provider{
				ResourcesMap:   map[string]*schema.Resource{},
				DataSourcesMap: map[string]*schema.Resource{},
			}
			if diags := assembleProvider(p); diags.HasError() {
				t.Fatal(diags)
			}
			for _, index := range []int{1, 0, 1} {
				spec := specs[index]
				r, ds := p.ResourcesMap[spec.name], p.DataSourcesMap[spec.name]
				if r == nil || ds == nil {
					t.Fatalf("%s 未从 registry 装配", spec.name)
				}
				ops := []struct {
					op, kind string
					run      func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
				}{
					{"Create", "resource", r.CreateContext},
					{"Read", "resource", r.ReadContext},
					{"Update", "resource", r.UpdateContext},
					{"Delete", "resource", r.DeleteContext},
					{"Read", "datasource", ds.ReadContext},
				}
				if mode == "WithoutTimeout" {
					ops[0].run, ops[1].run = r.CreateWithoutTimeout, r.ReadWithoutTimeout
					ops[2].run, ops[3].run = r.UpdateWithoutTimeout, r.DeleteWithoutTimeout
					ops[4].run = ds.ReadWithoutTimeout
				}
				for _, op := range ops {
					if op.run == nil {
						t.Fatalf("%s/%s/%s 回调缺失", spec.name, op.kind, op.op)
					}
					assertAssembledRequestContext(t, spec.service, spec.name+"/"+op.kind, op.op, &seen,
						func(ctx context.Context) {
							if diags := op.run(ctx, nil, meta); diags.HasError() {
								t.Fatal(diags)
							}
						})
				}
			}
		})
	}
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
	if got.Err() != context.Canceled {
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

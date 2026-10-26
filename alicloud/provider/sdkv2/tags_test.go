package sdkv2

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func tagFilterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func tagFilterResource(read func(*schema.ResourceData, interface{}) error) *schema.Resource {
	return &schema.Resource{Read: read, Schema: tagFilterSchema()}
}

func tagFilterClient(keys, prefixes []string) *connectivity.AliyunClient {
	return &connectivity.AliyunClient{
		IgnoreTags: features.IgnoreTags{Keys: keys, KeyPrefixes: prefixes},
	}
}

func tagFilterData(t *testing.T, tags map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, tagFilterSchema(), map[string]interface{}{"tags": tags})
}

func tagFilterGet(t *testing.T, d *schema.ResourceData) map[string]interface{} {
	t.Helper()
	raw, ok := d.Get("tags").(map[string]interface{})
	if !ok {
		t.Fatalf("tags is %T, want map[string]interface{}", d.Get("tags"))
	}
	return raw
}

func tagFilterAssert(t *testing.T, d *schema.ResourceData, want map[string]string) {
	t.Helper()
	got := tagFilterGet(t, d)
	if len(got) != len(want) {
		t.Fatalf("tags = %v, want %v", got, want)
	}
	for key, value := range want {
		if fmt.Sprint(got[key]) != value {
			t.Fatalf("tags[%q] = %v, want %q (all: %v)", key, got[key], value, got)
		}
	}
}

func TestUnitTagFilterRunsWithoutInterceptorChain(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	if wrapped == r {
		t.Fatal("a resource declaring top level tags must be wrapped even without interceptors")
	}
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"team": "core"})
}

func TestUnitTagFilterKeepsBaselineWithoutConfiguration(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "acs:created": "yes", "aliyun:biz": "x"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, &connectivity.AliyunClient{}); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"env": "prod", "acs:created": "yes", "aliyun:biz": "x"})
}

func TestUnitTagFilterKeepsUnmatchedTags(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "Env": "upper", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"Env": "upper", "team": "core"})
}

func TestUnitTagFilterKeepsTagsOnReadError(t *testing.T) {
	inner := errors.New("boom")
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		if err := d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"}); err != nil {
			return err
		}
		return inner
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); !errors.Is(err, inner) {
		t.Fatalf("err = %v, want %v", err, inner)
	}
	tagFilterAssert(t, d, map[string]string{"env": "prod", "team": "core"})
}

func TestUnitTagFilterUnionsSystemTagPrefixes(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{
			"env":                   "prod",
			"acs:created":           "yes",
			"aliyun:biz":            "x",
			"http://example.com/k":  "v",
			"https://example.com/k": "v",
			"sae.do.not.delete":     "keep",
			"team":                  "core",
		})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"team": "core"})
}

func TestUnitTagFilterIsIdempotent(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	client := tagFilterClient(nil, []string{"tf-"})
	d := tagFilterData(t, nil)
	for i := 0; i < 3; i++ {
		if err := wrapped.Read(d, client); err != nil {
			t.Fatalf("read %d failed: %v", i, err)
		}
		tagFilterAssert(t, d, map[string]string{"env": "prod", "team": "core"})
	}
}

func TestUnitTagFilterKeepsInterceptorChain(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if want := []string{"before:rec", "after:rec"}; fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Fatalf("interceptor ops = %v, want %v", rec.ops, want)
	}
	tagFilterAssert(t, d, map[string]string{"team": "core"})
}

func TestUnitTagFilterSupportsContextRead(t *testing.T) {
	r := &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"}))
		},
		Schema: tagFilterSchema(),
	}
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if diags := wrapped.ReadContext(context.Background(), d, tagFilterClient(nil, []string{"env"})); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	tagFilterAssert(t, d, map[string]string{"team": "core"})
}

func TestUnitTagFilterContextReadErrorKeepsState(t *testing.T) {
	r := &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			if err := d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"}); err != nil {
				return diag.FromErr(err)
			}
			return diag.Errorf("boom")
		},
		Schema: tagFilterSchema(),
	}
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if diags := wrapped.ReadContext(context.Background(), d, tagFilterClient([]string{"env"}, nil)); !diags.HasError() {
		t.Fatal("expected the read error to reach the caller")
	}
	tagFilterAssert(t, d, map[string]string{"env": "prod", "team": "core"})
}

func TestUnitTagFilterLeavesDataSourcesAlone(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	ds := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
		},
		Schema: tagFilterSchema(),
	}
	wrapped := WrapDataSource("alicloud_test", ds, []intercept.Interceptor{rec})
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"env": "prod", "team": "core"})
}

func TestUnitTagFilterIgnoresNonMapTagSchemas(t *testing.T) {
	r := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return d.Set("tags", []interface{}{map[string]interface{}{"key": "env", "value": "prod"}})
		},
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"key":   {Type: schema.TypeString, Optional: true},
					"value": {Type: schema.TypeString, Optional: true},
				}},
			},
		},
	}
	if got := WrapResource("alicloud_test", r, nil); got != r {
		t.Fatal("a resource without a top level tags map must stay identity-wrapped")
	}
	wrapped := WrapResource("alicloud_test", r, []intercept.Interceptor{&recordInterceptor{name: "rec"}})
	d := schema.TestResourceDataRaw(t, r.Schema, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if got := fmt.Sprint(d.Get("tags")); got != fmt.Sprint([]interface{}{map[string]interface{}{"key": "env", "value": "prod"}}) {
		t.Fatalf("nested tags were rewritten: %v", got)
	}
}

func TestUnitTagFilterIsolatesProviderInstances(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "tf-owner": "a", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)

	first := tagFilterClient([]string{"env"}, nil)
	second := tagFilterClient(nil, []string{"tf-"})
	dFirst := tagFilterData(t, nil)
	dSecond := tagFilterData(t, nil)
	if err := wrapped.Read(dFirst, first); err != nil {
		t.Fatalf("first read failed: %v", err)
	}
	if err := wrapped.Read(dSecond, second); err != nil {
		t.Fatalf("second read failed: %v", err)
	}
	tagFilterAssert(t, dFirst, map[string]string{"tf-owner": "a", "team": "core"})
	tagFilterAssert(t, dSecond, map[string]string{"env": "prod", "team": "core"})
}

func TestUnitTagFilterToleratesMissingMeta(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	for _, meta := range []interface{}{nil, struct{}{}, (*connectivity.AliyunClient)(nil)} {
		d := tagFilterData(t, nil)
		if err := wrapped.Read(d, meta); err != nil {
			t.Fatalf("read with meta %T failed: %v", meta, err)
		}
		tagFilterAssert(t, d, map[string]string{"env": "prod"})
	}
}

func TestUnitTagFilterSkipsEmptyTags(t *testing.T) {
	r := tagFilterResource(func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Read(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if got := tagFilterGet(t, d); len(got) != 0 {
		t.Fatalf("tags = %v, want none", got)
	}
}

func tagFilterMutationResource(op string, inner func(*schema.ResourceData, interface{}) error) *schema.Resource {
	r := &schema.Resource{Schema: tagFilterSchema()}
	switch op {
	case "create":
		r.Create = inner
	case "update":
		r.Update = inner
	}
	return r
}

func TestUnitTagFilterFiltersCreateState(t *testing.T) {
	r := tagFilterMutationResource("create", func(d *schema.ResourceData, meta interface{}) error {
		d.SetId("id")
		// A generated Create usually finishes by calling its own Read, which bypasses the wrapped
		// Read; the ignored key comes back here and must not reach the state Apply returns.
		return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Create(d, tagFilterClient([]string{"env"}, nil)); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"team": "core"})
}

func TestUnitTagFilterFiltersUpdateState(t *testing.T) {
	r := tagFilterMutationResource("update", func(d *schema.ResourceData, meta interface{}) error {
		return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
	})
	wrapped := WrapResource("alicloud_test", r, nil)
	d := tagFilterData(t, nil)
	if err := wrapped.Update(d, tagFilterClient(nil, []string{"env"})); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	tagFilterAssert(t, d, map[string]string{"team": "core"})
}

func TestUnitTagFilterKeepsCreateAndUpdateErrors(t *testing.T) {
	inner := errors.New("boom")
	for _, op := range []string{"create", "update"} {
		r := tagFilterMutationResource(op, func(d *schema.ResourceData, meta interface{}) error {
			if err := d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"}); err != nil {
				return err
			}
			return inner
		})
		wrapped := WrapResource("alicloud_test", r, nil)
		d := tagFilterData(t, nil)
		var err error
		if op == "create" {
			err = wrapped.Create(d, tagFilterClient([]string{"env"}, nil))
		} else {
			err = wrapped.Update(d, tagFilterClient([]string{"env"}, nil))
		}
		if !errors.Is(err, inner) {
			t.Fatalf("%s err = %v, want %v", op, err, inner)
		}
		tagFilterAssert(t, d, map[string]string{"env": "prod", "team": "core"})
	}
}

func TestUnitTagFilterKeepsCreateAndUpdateInterceptorChain(t *testing.T) {
	for _, op := range []string{"create", "update"} {
		rec := &recordInterceptor{name: "rec"}
		r := tagFilterMutationResource(op, func(d *schema.ResourceData, meta interface{}) error {
			if op == "create" {
				d.SetId("id")
			}
			return d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"})
		})
		wrapped := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
		d := tagFilterData(t, nil)
		var err error
		if op == "create" {
			err = wrapped.Create(d, tagFilterClient([]string{"env"}, nil))
		} else {
			err = wrapped.Update(d, tagFilterClient([]string{"env"}, nil))
		}
		if err != nil {
			t.Fatalf("%s failed: %v", op, err)
		}
		if want := []string{"before:rec", "after:rec"}; fmt.Sprint(rec.ops) != fmt.Sprint(want) {
			t.Fatalf("%s interceptor ops = %v, want %v", op, rec.ops, want)
		}
		tagFilterAssert(t, d, map[string]string{"team": "core"})
	}
}

func TestUnitTagFilterSupportsContextCreateAndUpdate(t *testing.T) {
	for _, op := range []string{"create", "update"} {
		inner := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			d.SetId("id")
			return diag.FromErr(d.Set("tags", map[string]interface{}{"env": "prod", "team": "core"}))
		}
		r := &schema.Resource{Schema: tagFilterSchema()}
		if op == "create" {
			r.CreateContext = inner
		} else {
			r.UpdateContext = inner
		}
		wrapped := WrapResource("alicloud_test", r, nil)
		d := tagFilterData(t, nil)
		var diags diag.Diagnostics
		if op == "create" {
			diags = wrapped.CreateContext(context.Background(), d, tagFilterClient([]string{"env"}, nil))
		} else {
			diags = wrapped.UpdateContext(context.Background(), d, tagFilterClient([]string{"env"}, nil))
		}
		if diags.HasError() {
			t.Fatalf("%s failed: %v", op, diags)
		}
		tagFilterAssert(t, d, map[string]string{"team": "core"})
	}
}

func TestUnitTagFilterLeavesCreateAndUpdateAloneWithoutTagsMap(t *testing.T) {
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { return nil },
		Update: func(d *schema.ResourceData, meta interface{}) error { return nil },
		Schema: map[string]*schema.Schema{"name": {Type: schema.TypeString, Optional: true}},
	}
	if got := WrapResource("alicloud_test", r, nil); got != r {
		t.Fatal("a resource without a top level tags map must stay identity-wrapped")
	}
}

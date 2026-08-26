package alicloud

import (
	"testing"

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

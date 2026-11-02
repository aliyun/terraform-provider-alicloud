package acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitLegacyHelperProviderCompatibility(t *testing.T) {
	if Provider.ResourcesMap["alicloud_fcv3_function"] == nil || Provider.ResourcesMap["alicloud_vpc"] == nil {
		t.Fatal("legacy provider selection was narrowed")
	}
	check := ResourceCheckInit("alicloud_vpc.default", nil, nil, "vpc_name")
	if check == nil {
		t.Fatal("legacy resource checks lost their provider")
	}
	d := (&DataSourceAttr{ResourceId: "data.alicloud_vpcs.default"}).helper()
	p, err := d.ProviderFactories["alicloud"]()
	if err != nil || p != Provider {
		t.Fatal("legacy datasource checks lost their provider")
	}
	var _ *schema.Provider = p
}

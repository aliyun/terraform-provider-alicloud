package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpnCgwsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnCgwsDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_customer_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.name", "tf-testAccVpnCgwName_DataResource"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.ip_address", "40.104.22.228"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.description", "tf-testAccVpnCgwDesc_Create"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_customer_gateways.foo", "gateways.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_customer_gateways.foo", "gateways.0.create_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnCgwsDataCfg = `
resource "alicloud_vpn_customer_gateway" "foo" {
	name = "tf-testAccVpnCgwName_DataResource"
	ip_address = "40.104.22.228"
	description = "tf-testAccVpnCgwDesc_Create"
}

data "alicloud_vpn_customer_gateways" "foo" {
	name_regex = "tf-testAcc*"
	ids = ["${alicloud_vpn_customer_gateway.foo.id}"]
}
`

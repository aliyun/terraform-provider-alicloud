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
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_customer_gateways.customer_gateways"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.customer_gateways", "customer_gateways.0.name", "testAccVpnCgwName_Create"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.customer_gateways", "customer_gateways.0.ip_address", "40.104.22.228"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.customer_gateways", "customer_gateways.0.description", "testAccVpnCgwDesc_Create"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnCgwsDataCfg = `
resource "alicloud_vpn_customer_gateway" "customer_gateways" {
  name = "testAccVpnCgwName_Create"
  ip_address = "40.104.22.228"
  description = "testAccVpnCgwDesc_Create"
}

data "alicloud_vpn_customer_gateways" "customer_gateways" {
	name_regex = "${alicloud_vpn_customer_gateway.customer_gateways.name}"
	output_file = "/tmp/cgws"
}
`

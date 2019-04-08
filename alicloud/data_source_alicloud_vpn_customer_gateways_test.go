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
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "names.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_customer_gateways.foo", "gateways.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.name", "tf-testAccVpnCgwName_DataResource"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.ip_address", "40.104.22.228"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.description", "tf-testAccVpnCgwDesc_Create"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_customer_gateways.foo", "gateways.0.create_time"),
				),
			},
		},
	})
}

func TestAccAlicloudVpnCgwsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnCgwsDataEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_customer_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_customer_gateways.foo", "names.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.ip_address"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_customer_gateways.foo", "gateways.0.create_time"),
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

const testAccCheckAlicloudVpnCgwsDataEmpty = `
data "alicloud_vpn_customer_gateways" "foo" {
	name_regex = "tf-testAcc-fake"
}
`

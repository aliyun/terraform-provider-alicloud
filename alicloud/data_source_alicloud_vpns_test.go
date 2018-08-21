package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpnsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnsDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpns.vpn_gateways"),
					resource.TestCheckResourceAttr("data.alicloud_vpns.vpn_gateways", "vpn_gateways.0.name", "testAccVpnConfig_create"),
					resource.TestCheckResourceAttr("data.alicloud_vpns.vpn_gateways", "vpn_gateways.0.spec", "10M"),
					resource.TestCheckResourceAttr("data.alicloud_vpns.vpn_gateways", "vpn_gateways.0.description", "test_create_description"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnsDataCfg = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "testAccVpcConfig"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchConfig"
}

resource "alicloud_vpn" "foo" {
        name = "testAccVpnConfig_create"
        vpc_id = "${alicloud_vpc.foo.id}"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_create_description"
}

data "alicloud_vpns" "vpn_gateways" {
	vpc_id = "${alicloud_vpn.foo.vpc_id}"
	output_file = "/tmp/vpns"
}
`

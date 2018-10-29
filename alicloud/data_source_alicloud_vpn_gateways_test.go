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
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_gateways.vpn_gateways"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.name", "tf-testAccVpnDatasource"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.specification", "10M"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.description", "test_create_description"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.enable_ssl", "enable"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.enable_ipsec", "enable"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.business_status", "Normal"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.ssl_connections", "5"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnsDataCfg = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf-testAccVpcConfig"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "tf-testAccVswitchConfig"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "tf-testAccVpnDatasource"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

data "alicloud_vpn_gateways" "vpn_gateways" {
	vpc_id = "${alicloud_vpc.foo.id}"
	ids = ["${alicloud_vpn_gateway.foo.id}"]
	status = "Active"
	business_status = "Normal"
	name_regex = "tf-testAcc*"
}
`

package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpcsDataSource_cidr_block(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceCidrBlockConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "tf-testAccVpcsdatasourceNameRegex"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.description"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.create_time"),
				),
			},
		},
	})
}

func TestAccAlicloudVpcsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.create_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpcsDataSourceCidrBlockConfig = `
variable "name" {
  default = "tf-testAccVpcsdatasourceNameRegex"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "testAccVpcsdatasource*"
  cidr_block = "${alicloud_vpc.foo.cidr_block}"
}
`

const testAccCheckAlicloudVpcsDataSourceEmpty = `
data "alicloud_vpcs" "vpc" {
  name_regex = "^tf-fake-name"
}
`

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
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "testAccVpcsdatasourceNameRegex"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpcsDataSourceCidrBlockConfig = `
variable "name" {
  default = "testAccVpcsdatasourceNameRegex"
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

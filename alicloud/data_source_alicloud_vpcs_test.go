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
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.region_id", "cn-beijing"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpcsDataSourceCidrBlockConfig = `
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  cidr_block = "${alicloud_vpc.foo.cidr_block}"
}
`

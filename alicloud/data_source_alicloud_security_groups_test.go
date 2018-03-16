package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSecurityGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.name", "webaccess"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.description", "test security group"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupsDataSourceConfig = `
data "alicloud_vpcs" "main" {
}

resource "alicloud_security_group" "test" {
  name        = "webaccess"
  description = "test security group"
  vpc_id      = "${data.alicloud_vpcs.main.vpcs.0.id}"
}

data "alicloud_security_groups" "web" {
    name_regex = "^web"
    vpc_id     = "${alicloud_security_group.test.vpc_id}"
}
`

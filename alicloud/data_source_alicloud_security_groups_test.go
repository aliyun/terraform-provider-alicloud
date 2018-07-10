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
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.name", "testAccCheckAlicloudSecurityGroupsDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.description", "test security group"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupsDataSourceConfig = `
variable "name" {
	default = "testAccCheckAlicloudSecurityGroupsDataSourceConfig"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

data "alicloud_security_groups" "web" {
    name_regex = "(SecurityGroupsDataSourceConfig)$"
    vpc_id     = "${alicloud_security_group.test.vpc_id}"
}
`

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
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.name", "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.description", "test security group"),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupsDataSource_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupDataSourceTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.tags.from", "datasource"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.tags.usage1", "test"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupsDataSourceConfig = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
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

const testAccCheckAlicloudSecurityGroupDataSourceTags = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
}
resource "alicloud_vpc" "tf_vpc_foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.tf_vpc_foo.id}"
  tags {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
		usage3 = "test"
		usage4 = "test"
		usage5 = "test"
		usage6 = "test"

  }
}

data "alicloud_security_groups" "web" {
    name_regex = "(SecurityGroupsDataSourceConfig)$"
    vpc_id     = "${alicloud_security_group.test.vpc_id}"
	tags {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
		usage3 = "test"
		usage4 = "test"
		usage5 = "test"
		usage6 = "test"

	}
}
`

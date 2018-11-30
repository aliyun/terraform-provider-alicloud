package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbAclsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbAclsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_acls.slb_acls"),
					resource.TestCheckResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_acls.slb_acls", "acls.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.name", "tf-testAccSlbAclDataSourceBisic"),
					resource.TestCheckResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.ip_version", "ipv4"),
					resource.TestCheckResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.entry_list.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.related_listeners.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbAclsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbAclsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_acls.slb_acls"),
					resource.TestCheckResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.ip_version"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.entry_list.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_acls.slb_acls", "acls.0.related_listeners.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbAclsDataSourceBasic = `
variable "name" {
	default = "tf-testAccSlbAclDataSourceBisic"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "foo" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}

data "alicloud_slb_acls" "slb_acls" {
  ids = ["${alicloud_slb_acl.foo.id}"]
  name_regex = "${var.name}"
}
`

const testAccCheckAlicloudSlbAclsDataSourceEmpty = `
data "alicloud_slb_acls" "slb_acls" {
  name_regex = "tf-testacc-fake-name"
}
`

package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloud_AccessGroup_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_accessgroups.ag"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessgroups.ag", "accessgroups.0.rule_count"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessgroups.ag", "accessgroups.0.accessgroup_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessgroups.ag", "accessgroups.0.description"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessgroups.ag", "accessgroups.0.accessgroup_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessgroups.ag", "accessgroups.0.mounttarget_count"),
				),
			},
		},
	})
}

func TestAccAlicloud_AccessGroup_DataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_accessgroups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_accessgroups.ag", "accessgroups.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_nas_accessgroups.ag", "accessgroups.0.rule_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_nas_accessgroups.ag", "accessgroups.0.accessgroup_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_nas_accessgroups.ag", "accessgroups.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_nas_accessgroups.ag", "accessgroups.0.accessgroup_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_nas_accessgroups.ag", "accessgroups.0.mounttarget_count"),
				),
			},
		},
	})
}

const testAccCheckAlicloudAccessGroupsDataSource = `
variable "name" {
  default = "tf-testAccAccessGroupsdatasourceNameRegex"
}
resource "alicloud_nas_accessgroup" "foo" {
  accessgroup_name = "${var.name}"
  accessgroup_type = "Classic"
  description = "test_wang"
}
data "alicloud_nas_accessgroups" "ag" {
  accessgroup_name = "${alicloud_nas_accessgroup.foo.accessgroup_name}"
  accessgroup_type = "${alicloud_nas_accessgroup.foo.accessgroup_type}"
  description = "${alicloud_nas_accessgroup.foo.description}"
}
`

const testAccCheckAlicloudAccessGroupsDataSourceEmpty = `
data "alicloud_nas_accessgroups" "ag" {
 name_regex = "^tf-fake-name"
}
`

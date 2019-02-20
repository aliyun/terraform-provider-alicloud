package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloud_AccessRule_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_accessrules.ar"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessrules.ar", "accessrules.0.sourcecidr_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessrules.ar", "accessrules.0.priority"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessrules.ar", "accessrules.0.accessrule_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessrules.ar", "accessrules.0.user_access"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_accessrules.ar", "accessrules.0.rw_access"),
				),
			},
		},
	})
}

func TestAccAlicloud_AccessRule_DataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceEmpty,
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


const testAccCheckAlicloudAccessRulesDataSource = `
data "alicloud_nas_accessrules" "ar" {
  accessgroup_name = "classic"
}
`

const testAccCheckAlicloudAccessRulesDataSourceEmpty = `
data "alicloud_nas_accessgroups" "ag" {
 name_regex = "^tf-fake-name"
}
`

//const testAccCheckAlicloudAccessRulesDataSource = `
//variable "name" {
//  default = "tf-testAccAccessRulesdatasourceNameRegex"
//}
//resource "alicloud_nas_accessgroup" "foo" {
//  accessgroup_name = "${var.name}"
//  accessgroup_type = "Classic"
//  description = "test_wang"
//}
//resource "alicloud_nas_accessrule" "foo" {
//		accessgroup_name = "${alicloud_nas_accessgroup.foo.accessgroup_name}"
//		sourcecidr_ip = "168.1.1.0/16"
//		rwaccess_type = "RDWR"
//		useraccess_type = "no_squash"
//}
//data "alicloud_nas_accessrules" "ar" {
//  accessgroup_name = "${alicloud_nas_accessgroup.foo.accessgroup_name}"
//}
//`
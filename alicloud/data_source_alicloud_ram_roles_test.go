package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamRolesDataSource_for_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceForPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.0.name", "testAccCheckAlicloudRamRolesDataSourceForPolicyConfig"),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolesDataSource_for_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceForAllConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolesDataSource_role_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamRolesDataSourceForPolicyConfig = `
variable "name" {
  default = "testAccCheckAlicloudRamRolesDataSourceForPolicyConfig"
}
resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects",
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"]
    }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_role" "role" {
  name = "${var.name}"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}

data "alicloud_ram_roles" "role" {
  policy_name = "${alicloud_ram_role_policy_attachment.attach.policy_name}"
  policy_type = "Custom"
}`

const testAccCheckAlicloudRamRolesDataSourceForAllConfig = `
data "alicloud_ram_roles" "role" {
}`

const testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig = `
resource "alicloud_ram_role" "role" {
  name = "testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}
data "alicloud_ram_roles" "role" {
  name_regex = "${alicloud_ram_role.role.name}"
}`

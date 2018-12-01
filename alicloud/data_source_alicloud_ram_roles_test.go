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
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.0.name", "tf-testAccCheckAlicloudRamRolesDataSourceForPolicyConfig"),
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

func TestAccAlicloudRamRolesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.arn"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.assume_role_policy_document"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.document"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.create_date"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.update_date"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamRolesDataSourceForPolicyConfig = `
variable "name" {
  default = "tf-testAccCheckAlicloudRamRolesDataSourceForPolicyConfig"
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
  name = "tf-testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}
data "alicloud_ram_roles" "role" {
  name_regex = "${alicloud_ram_role.role.name}"
}`

const testAccCheckAlicloudRamRolesDataSourceEmpty = `
data "alicloud_ram_roles" "role" {
	name_regex = "^tf-testacc-fake-name"
}`

package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
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
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.0.name", "testrole"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.0.arn", "acs:ram::1307087942598154:role/testrole"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.0.id", "345148520161269882"),
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
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "3"),
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
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "2"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamRolesDataSourceForPolicyConfig = `
data "alicloud_ram_roles" "role" {
  policy_name = "AliyunACSDefaultAccess"
  policy_type = "Custom"
}`

const testAccCheckAlicloudRamRolesDataSourceForAllConfig = `
data "alicloud_ram_roles" "role" {
}`

const testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig = `
data "alicloud_ram_roles" "role" {
  name_regex = "^test"
}`

package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamUsersDataSource_for_group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamUsersDataSourceForGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.0.name", "yu"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUsersDataSource_for_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamUsersDataSourceForPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.0.name", "user3"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.0.id", "279233601570976655"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUsersDataSource_for_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamUsersDataSourceForAllConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "5"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUsersDataSource_user_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamGroupsDataSourceUserNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "2"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamUsersDataSourceForGroupConfig = `
data "alicloud_ram_users" "user" {
  type = "group"
  group_name = "group4"
}`

const testAccCheckAlicloudRamUsersDataSourceForPolicyConfig = `
data "alicloud_ram_users" "user" {
  type = "policy"
  policy_name = "AliyunACSDefaultAccess"
  policy_type = "Custom"
}`

const testAccCheckAlicloudRamUsersDataSourceForAllConfig = `
data "alicloud_ram_users" "user" {
}`

const testAccCheckAlicloudRamGroupsDataSourceUserNameRegexConfig = `
data "alicloud_ram_users" "user" {
  name_regex = "^yu"
}`

package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudRamGroupsDataSource_for_user(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamGroupsDataSourceForUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.name", "group3"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.comments", "33"),
				),
			},
		},
	})
}

func TestAccAlicloudRamGroupsDataSource_for_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.name", "group1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.comments", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamGroupsDataSource_for_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamGroupsDataSourceForAllConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "5"),
				),
			},
		},
	})
}

func TestAccAlicloudRamGroupsDataSource_group_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamGroupsDataSourceGroupNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "4"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamGroupsDataSourceForUserConfig = `
data "alicloud_ram_groups" "group" {
  user_name = "user1"
}`

const testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig = `
data "alicloud_ram_groups" "group" {
  policy_name = "AliyunMobileTestingFullAccess"
  policy_type = "System"
}`

const testAccCheckAlicloudRamGroupsDataSourceForAllConfig = `
data "alicloud_ram_groups" "group" {
}`

const testAccCheckAlicloudRamGroupsDataSourceGroupNameRegexConfig = `
data "alicloud_ram_groups" "group" {
  name_regex = "^group"
}`

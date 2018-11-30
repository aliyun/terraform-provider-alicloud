package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.name", "tf-testAccCheckAlicloudRamGroupsDataSourceForUserConfig"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.comments", "group comments"),
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
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.name", "tf-testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.0.comments", "group comments"),
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
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamGroupsDataSource_Empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamGroupsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_groups.group", "groups.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_groups.group", "groups.0.comments"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamGroupsDataSourceForUserConfig = `
variable "name" {
  default = "tf-testAccCheckAlicloudRamGroupsDataSourceForUserConfig"
}
resource "alicloud_ram_user" "user" {
  name = "${var.name}"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}

resource "alicloud_ram_group" "group" {
  name = "${var.name}"
  comments = "group comments"
  force=true
}
resource "alicloud_ram_group_membership" "membership" {
  group_name = "${alicloud_ram_group.group.name}"
  user_names = ["${alicloud_ram_user.user.name}"]
}

data "alicloud_ram_groups" "group" {
  user_name = "${alicloud_ram_group_membership.membership.user_names[0]}"
}`

const testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig = `
variable "name" {
  default = "tf-testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig"
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

resource "alicloud_ram_group" "group" {
  name = "${var.name}"
  comments = "group comments"
  force=true
}

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  group_name = "${alicloud_ram_group.group.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}
data "alicloud_ram_groups" "group" {
  policy_name = "${alicloud_ram_group_policy_attachment.attach.policy_name}"
  policy_type = "Custom"
}`

const testAccCheckAlicloudRamGroupsDataSourceForAllConfig = `
data "alicloud_ram_groups" "group" {
}`

const testAccCheckAlicloudRamGroupsDataSourceGroupNameRegexConfig = `
resource "alicloud_ram_group" "group" {
  name = "tf-testAccCheckAlicloudRamGroupsDataSourceGroupNameRegexConfig"
  comments = "group comments"
  force=true
}
data "alicloud_ram_groups" "group" {
  name_regex = "${alicloud_ram_group.group.name}"
}`

const testAccCheckAlicloudRamGroupsDataSourceEmpty = `
data "alicloud_ram_groups" "group" {
	name_regex = "^tf-testacc-fake-name"
}`

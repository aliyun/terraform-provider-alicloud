package alicloud

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
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
				Config: testAccCheckAlicloudRamGroupsDataSourceForUserConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_ram_groups.group", "groups.0.name",
						regexp.MustCompile("^tf-testAccRamGroupsDataSourceForUser-*")),
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
				Config: testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_ram_groups.group", "groups.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_ram_groups.group", "groups.0.name",
						regexp.MustCompile("^tf-testAccRamGroupsDataSourceForPolicy-*")),
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
				Config: testAccCheckAlicloudRamGroupsDataSourceGroupNameRegexConfig(acctest.RandIntRange(1000000, 99999999)),
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

func testAccCheckAlicloudRamGroupsDataSourceForUserConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamGroupsDataSourceForUser-%d"
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
	}`, rand)
}

func testAccCheckAlicloudRamGroupsDataSourceForPolicyConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamGroupsDataSourceForPolicy-%d"
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
	}`, rand)
}

const testAccCheckAlicloudRamGroupsDataSourceForAllConfig = `
data "alicloud_ram_groups" "group" {
}`

func testAccCheckAlicloudRamGroupsDataSourceGroupNameRegexConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "group" {
	  name = "tf-testAccRamGroupsDataSourceGroupNameRegex-%d"
	  comments = "group comments"
	  force=true
	}
	data "alicloud_ram_groups" "group" {
	  name_regex = "${alicloud_ram_group.group.name}"
	}`, rand)
}

const testAccCheckAlicloudRamGroupsDataSourceEmpty = `
data "alicloud_ram_groups" "group" {
	name_regex = "^tf-testacc-fake-name"
}`

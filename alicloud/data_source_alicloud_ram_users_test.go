package alicloud

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
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
				Config: testAccCheckAlicloudRamUsersDataSourceForGroupConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_ram_users.user", "users.0.name",
						regexp.MustCompile("tf-testAccRamUsersDataSourceForGroup-*")),
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
				Config: testAccCheckAlicloudRamUsersDataSourceForPolicyConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_ram_users.user", "users.0.name",
						regexp.MustCompile("^tf-testAccRamUsersDataSourceForPolicy-*")),
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
				Config: testAccCheckAlicloudRamGroupsDataSourceUserNameRegexConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUsersDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamUsersDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_users.user"),
					resource.TestCheckResourceAttr("data.alicloud_ram_users.user", "users.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_users.user", "users.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_users.user", "users.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_users.user", "users.0.create_date"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_users.user", "users.0.last_login_date"),
				),
			},
		},
	})
}

func testAccCheckAlicloudRamUsersDataSourceForGroupConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamUsersDataSourceForGroup-%d"
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

	data "alicloud_ram_users" "user" {
	  group_name = "${alicloud_ram_group_membership.membership.group_name}"
	}`, rand)
}

func testAccCheckAlicloudRamUsersDataSourceForPolicyConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamUsersDataSourceForPolicy-%d"
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

	resource "alicloud_ram_user" "user" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_user_policy_attachment" "attach" {
	  policy_name = "${alicloud_ram_policy.policy.name}"
	  user_name = "${alicloud_ram_user.user.name}"
	  policy_type = "${alicloud_ram_policy.policy.type}"
	}

	data "alicloud_ram_users" "user" {
	  policy_name = "${alicloud_ram_user_policy_attachment.attach.policy_name}"
	  policy_type = "Custom"
	}`, rand)
}

const testAccCheckAlicloudRamUsersDataSourceForAllConfig = `
data "alicloud_ram_users" "user" {
}`

func testAccCheckAlicloudRamGroupsDataSourceUserNameRegexConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamGroupsDataSourceUserNameRegex-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}
	data "alicloud_ram_users" "user" {
	  name_regex = "${alicloud_ram_user.user.name}"
	}`, rand)
}

const testAccCheckAlicloudRamUsersDataSourceEmpty = `
data "alicloud_ram_users" "user" {
	name_regex = "tf-testacc-fake-name"
}`

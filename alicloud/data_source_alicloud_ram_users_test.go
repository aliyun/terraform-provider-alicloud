package alicloud

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRAMUsersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	groupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name": `"${alicloud_ram_group_membership.default.group_name}"`,
		}),
	}

	policyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_user.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_user.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ram_user.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ram_user.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group_membership.default.group_name}"`,
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"ids":         `["${alicloud_ram_user.default.id}"]`,
			"name_regex":  `"${alicloud_ram_user.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group_membership.default.group_name}"`,
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"ids":         `["${alicloud_ram_user.default.id}"]`,
			"name_regex":  `"${alicloud_ram_user.default.name}_fake"`,
		}),
	}

	ramUsersCheckInfo.dataSourceTestCheck(t, rand, groupConf, policyConf, nameRegexConf, idsConf, allConf)

}
func testAccCheckAlicloudRamUsersDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamUsersDataSource-%d"
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}
	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}"]
	}
	resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}
	resource "alicloud_ram_user" "default" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_user_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  user_name = "${alicloud_ram_user.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}

	data "alicloud_ram_users" "default" {
	  %s
	}`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRamUsersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "1",
		"names.#":             "1",
		"users.#":             "1",
		"users.0.id":          CHECKSET,
		"users.0.name":        fmt.Sprintf("tf-testAcc%sRamUsersDataSource-%d", defaultRegionToTest, rand),
		"users.0.create_date": CHECKSET,
	}
}

var fakeRamUsersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"users.#": "0",
	}
}

var ramUsersCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_users.default",
	existMapFunc: existRamUsersMapFunc,
	fakeMapFunc:  fakeRamUsersMapFunc,
}

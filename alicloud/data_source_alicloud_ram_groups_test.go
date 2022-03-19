package alicloud

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRAMGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	userConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"user_name": `"${element(split(",",join(",",alicloud_ram_group_membership.default.user_names)), 0)}"`,
		}),
	}

	policyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_group_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_group.default.name}_fake"`,
		}),
	}

	userAndNameRegexconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"user_name":  `"${element(split(",",join(",",alicloud_ram_group_membership.default.user_names)), 0)}"`,
			"name_regex": `"${alicloud_ram_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"user_name":  `"${element(split(",",join(",",alicloud_ram_group_membership.default.user_names)), 0)}"`,
			"name_regex": `"${alicloud_ram_group.default.name}_fake"`,
		}),
	}

	userAndPolicyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"user_name":   `"${element(split(",",join(",",alicloud_ram_group_membership.default.user_names)), 0)}"`,
			"policy_name": `"${alicloud_ram_group_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	policyAndNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_group_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"name_regex":  `"${alicloud_ram_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_group_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"name_regex":  `"${alicloud_ram_group.default.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"user_name":   `"${element(split(",",join(",",alicloud_ram_group_membership.default.user_names)), 0)}"`,
			"policy_name": `"${alicloud_ram_group_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"name_regex":  `"${alicloud_ram_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamGroupsDataSourceConfig(rand, map[string]string{
			"user_name":   `"${element(split(",",join(",",alicloud_ram_group_membership.default.user_names)), 0)}"`,
			"policy_name": `"${alicloud_ram_group_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"name_regex":  `"${alicloud_ram_group.default.name}_fake"`,
		}),
	}

	RamGroupsCheckInfo.dataSourceTestCheck(t, rand, userConf, policyConf, nameRegexConf, userAndNameRegexconf, userAndPolicyConf, policyAndNameRegexConf, allConf)
}

func testAccCheckAlicloudRamGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamGroupsDataSource-%d"
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
	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}
	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}"]
	}
	resource "alicloud_ram_group_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  group_name = "${alicloud_ram_group.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}
	data "alicloud_ram_groups" "default" {
	  %s
	}`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRamGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":           "1",
		"groups.#":          "1",
		"groups.0.name":     fmt.Sprintf("tf-testAcc%sRamGroupsDataSource-%d", defaultRegionToTest, rand),
		"groups.0.comments": "group comments",
	}
}

var fakeRamGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":  "0",
		"groups.#": "0",
	}
}

var RamGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_groups.default",
	existMapFunc: existRamGroupsMapFunc,
	fakeMapFunc:  fakeRamGroupsMapFunc,
}

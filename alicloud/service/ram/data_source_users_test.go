package ram_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRamUsersDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 99999999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ram_user.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ram_user.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_user.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_user.default.name}_fake"`,
		}),
	}

	groupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name": `"${alicloud_ram_group_membership.default.group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name": `"${alicloud_ram_group.fake.name}"`,
		}),
	}

	policyConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_policy.fake.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	emptyGroupFullPolicyConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group.fake.name}"`,
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group.fake.name}"`,
			"policy_name": `"${alicloud_ram_policy.fake.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	emptyPolicyFullGroupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group_membership.default.group_name}"`,
			"policy_name": `"${alicloud_ram_policy.fake.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group.fake.name}"`,
			"policy_name": `"${alicloud_ram_policy.fake.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	groupAndPolicyConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group_membership.default.group_name}"`,
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group.fake.name}"`,
			"policy_name": `"${alicloud_ram_policy.fake.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group_membership.default.group_name}"`,
			"policy_name": `"${alicloud_ram_user_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"ids":         `["${alicloud_ram_user.default.id}"]`,
			"name_regex":  `"${alicloud_ram_user.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudRamUsersDataSourceConfig(rand, map[string]string{
			"group_name":  `"${alicloud_ram_group.fake.name}"`,
			"policy_name": `"${alicloud_ram_policy.fake.policy_name}"`,
			"policy_type": `"Custom"`,
			"ids":         `["${alicloud_ram_user.default.id}"]`,
			"name_regex":  `"${alicloud_ram_user.default.name}_fake"`,
		}),
	}

	ramUsersCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, groupConf, policyConf, emptyGroupFullPolicyConf, emptyPolicyFullGroupConf, groupAndPolicyConf, allConf)
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
  		name     = "${var.name}"
  		comments = "group comments"
  		force    = true
	}

	resource "alicloud_ram_group" "fake" {
  		name     = "${var.name}-fake"
  		comments = "fake group comments"
  		force    = true
	}

	resource "alicloud_ram_user" "default" {
  		name         = "${var.name}"
  		display_name = "displayname"
  		mobile       = "86-18888888888"
  		email        = "hello.uuu@aaa.com"
  		comments     = "yoyoyo"
	}

	resource "alicloud_ram_group_membership" "default" {
  		group_name = "${alicloud_ram_group.default.name}"
  		user_names = ["${alicloud_ram_user.default.name}"]
	}

	resource "alicloud_ram_policy" "default" {
  		name        = "${var.name}"
  		document    = <<EOF
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
  		force       = true
	}

	resource "alicloud_ram_policy" "fake" {
  		name        = "${var.name}-fake"
  		document    = <<EOF
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
  		description = "this is a fake policy test"
  		force       = true
	}

	resource "alicloud_ram_user_policy_attachment" "default" {
  		policy_name = "${alicloud_ram_policy.default.name}"
  		user_name   = "${alicloud_ram_user.default.name}"
  		policy_type = "${alicloud_ram_policy.default.type}"
	}

	data "alicloud_ram_users" "default" {
		%s
	}
`, acctest.DefaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRamUsersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "1",
		"names.#":             "1",
		"users.#":             "1",
		"users.0.id":          acctest.CHECKSET,
		"users.0.name":        fmt.Sprintf("tf-testAcc%sRamUsersDataSource-%d", acctest.DefaultRegionToTest, rand),
		"users.0.create_date": acctest.CHECKSET,
	}
}

var fakeRamUsersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"users.#": "0",
	}
}

var ramUsersCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ram_users.default",
	ExistMapFunc: existRamUsersMapFunc,
	FakeMapFunc:  fakeRamUsersMapFunc,
}

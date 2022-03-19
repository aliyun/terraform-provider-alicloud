package alicloud

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRAMPoliciesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)

	groupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"group_name":     `"${alicloud_ram_group_policy_attachment.default.group_name}"`,
			"enable_details": "true",
		}),
	}

	roleConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"role_name":      `"${alicloud_ram_role_policy_attachment.default.role_name}"`,
			"type":           `"Custom"`,
			"enable_details": "true",
		}),
	}

	userConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"user_name":      `"${alicloud_ram_user_policy_attachment.default.user_name}"`,
			"enable_details": "true",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_ram_policy.default.name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_ram_policy.default.name}_fake"`,
			"enable_details": "true",
		}),
	}

	policyType := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_ram_policy.default.name}"`,
			"type":           `"Custom"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_ram_policy.default.name}_fake"`,
			"type":           `"Custom"`,
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"group_name":     `"${alicloud_ram_group_policy_attachment.default.group_name}"`,
			"role_name":      `"${alicloud_ram_role_policy_attachment.default.role_name}"`,
			"user_name":      `"${alicloud_ram_user_policy_attachment.default.user_name}"`,
			"name_regex":     `"${alicloud_ram_policy.default.name}"`,
			"type":           `"Custom"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudRamPoliciesDataSourceConfig(rand, map[string]string{
			"group_name":     `"${alicloud_ram_group_policy_attachment.default.group_name}"`,
			"role_name":      `"${alicloud_ram_role_policy_attachment.default.role_name}"`,
			"user_name":      `"${alicloud_ram_user_policy_attachment.default.user_name}"`,
			"name_regex":     `"${alicloud_ram_policy.default.name}_fake"`,
			"type":           `"Custom"`,
			"enable_details": "true",
		}),
	}

	ramPoliciesCheckInfo.dataSourceTestCheck(t, rand, groupConf, roleConf, userConf, nameRegexConf, policyType, allConf)
}

func testAccCheckAlicloudRamPoliciesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamPoliciessDataSource-%d"
	}
	resource "alicloud_ram_policy" "default" {
	  policy_name = "${var.name}"
	  policy_document = <<EOF
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

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}
	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
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
	resource "alicloud_ram_role_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  role_name = "${alicloud_ram_role.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}
	resource "alicloud_ram_group_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  group_name = "${alicloud_ram_group.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}
	data "alicloud_ram_policies" "default" {
	  %s
	}`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRamPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"names.#":                     "1",
		"policies.#":                  "1",
		"policies.0.name":             fmt.Sprintf("tf-testAcc%sRamPoliciessDataSource-%d", defaultRegionToTest, rand),
		"policies.0.policy_name":      fmt.Sprintf("tf-testAcc%sRamPoliciessDataSource-%d", defaultRegionToTest, rand),
		"policies.0.type":             CHECKSET,
		"policies.0.description":      "this is a policy test",
		"policies.0.default_version":  CHECKSET,
		"policies.0.create_date":      CHECKSET,
		"policies.0.update_date":      CHECKSET,
		"policies.0.attachment_count": CHECKSET,
		"policies.0.document":         CHECKSET,
		"policies.0.policy_document":  CHECKSET,
	}
}

var fakeRamPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"policies.#": "0",
	}
}

var ramPoliciesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_policies.default",
	existMapFunc: existRamPoliciesMapFunc,
	fakeMapFunc:  fakeRamPoliciesMapFunc,
}

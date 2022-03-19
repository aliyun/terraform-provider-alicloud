package alicloud

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRAMRolesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)

	policyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_role.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ram_role.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ram_role.default.role_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ram_role.default.role_id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"ids":         `["${alicloud_ram_role.default.role_id}"]`,
			"name_regex":  `"${alicloud_ram_role.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
			"ids":         `["${alicloud_ram_role.default.role_id}_fake"]`,
			"name_regex":  `"${alicloud_ram_role.default.name}_fake"`,
		}),
	}

	ramRolesCheckInfo.dataSourceTestCheck(t, rand, policyConf, nameRegexConf, idsConf, allConf)

}

func testAccCheckAlicloudRamRolesDataSourceConfig(rand int, attrMap map[string]string) string {
	region := defaultRegionToTest
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamRolesDataSourceForPolicy-%d"
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

	resource "alicloud_ram_role_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  role_name = "${alicloud_ram_role.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}
	data "alicloud_ram_roles" "default" {
	%s
	}`, region, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRamRolesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "1",
		"names.#":             "1",
		"roles.#":             "1",
		"roles.0.id":          CHECKSET,
		"roles.0.name":        fmt.Sprintf("tf-testAcc%sRamRolesDataSourceForPolicy-%d", defaultRegionToTest, rand),
		"roles.0.arn":         CHECKSET,
		"roles.0.description": "this is a test",
		"roles.0.document":    CHECKSET,
	}
}

var fakeRamRolesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"roles.#": "0",
	}
}

var ramRolesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_roles.default",
	existMapFunc: existRamRolesMapFunc,
	fakeMapFunc:  fakeRamRolesMapFunc,
}

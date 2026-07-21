package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerPolicyAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	policyNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_resource_manager_policy_attachment.this.policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand, map[string]string{
			"policy_name": `"fake"`,
		}),
	}

	principalNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand, map[string]string{
			"principal_name": `"${alicloud_resource_manager_policy_attachment.this.principal_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand, map[string]string{
			"principal_name": `"fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand, map[string]string{
			"policy_name":       `"${alicloud_resource_manager_policy_attachment.this.policy_name}"`,
			"policy_type":       `"Custom"`,
			"principal_name":    `"${alicloud_resource_manager_policy_attachment.this.principal_name}"`,
			"principal_type":    `"IMSUser"`,
			"resource_group_id": `"${alicloud_resource_manager_policy_attachment.this.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand, map[string]string{
			"policy_name":       `"fake"`,
			"policy_type":       `"Custom"`,
			"principal_name":    `"${alicloud_resource_manager_policy_attachment.this.principal_name}"`,
			"principal_type":    `"IMSUser"`,
			"resource_group_id": `"${alicloud_resource_manager_policy_attachment.this.resource_group_id}"`,
		}),
	}

	var existResourceManagerPolicyAttachmentsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"attachments.#":                   "1",
			"attachments.0.id":                CHECKSET,
			"attachments.0.attach_date":       CHECKSET,
			"attachments.0.description":       CHECKSET,
			"attachments.0.policy_name":       fmt.Sprintf("tf-testaccrdpolicy%d", rand),
			"attachments.0.policy_type":       "Custom",
			"attachments.0.principal_name":    CHECKSET,
			"attachments.0.principal_type":    "IMSUser",
			"attachments.0.resource_group_id": CHECKSET,
		}
	}

	var fakeResourceManagerPolicyAttachmentsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"attachments.#": "0",
		}
	}

	var policyAttachmentsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_policy_attachments.this",
		existMapFunc: existResourceManagerPolicyAttachmentsRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerPolicyAttachmentsRecordsMapFunc,
	}

	policyAttachmentsRecordsCheckInfo.dataSourceTestCheck(t, rand, policyNameConf, principalNameConf, allConf)

}

func testAccCheckAlicloudResourceManagerPolicyAttachmentsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testaccrdpolicy%d"
}

resource "alicloud_ram_user" "this" {
  name = "${var.name}"
}

locals {
  principal_name = format("%%s@%%s.onaliyun.com", alicloud_ram_user.this.name, data.alicloud_account.this.id)
}

data "alicloud_account" "this" {}

data "alicloud_resource_manager_resource_groups" "this" {
  name_regex = "default"
}

resource "alicloud_resource_manager_policy" "this" {
  policy_name     = "${var.name}"
  description 	  = "policy_attachment"
  policy_document = <<EOF
        {
            "Statement": [{
                "Action": ["oss:*"],
                "Effect": "Allow",
                "Resource": ["acs:oss:*:*:*"]
            }],
            "Version": "1"
        }
    EOF
}

resource "alicloud_resource_manager_policy_attachment" "this" {
  policy_name = "${alicloud_resource_manager_policy.this.policy_name}"
  policy_type = "Custom"
  principal_name = "${local.principal_name}"
  principal_type = "IMSUser"
  resource_group_id = data.alicloud_resource_manager_resource_groups.this.groups.0.id
}

data "alicloud_resource_manager_policy_attachments" "this"{
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}

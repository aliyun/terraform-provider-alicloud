// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRamRolePolicyAttachmentDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	RoleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolePolicyAttachmentSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ram_role_policy_attachment.default.id}"]`,
			"role_name": `"${alicloud_ram_role.defaultzEJqM7.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamRolePolicyAttachmentSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ram_role_policy_attachment.default.id}_fake"]`,
			"role_name": `"${alicloud_ram_role.defaultzEJqM7.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolePolicyAttachmentSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ram_role_policy_attachment.default.id}"]`,
			"role_name": `"${alicloud_ram_role.defaultzEJqM7.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRamRolePolicyAttachmentSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ram_role_policy_attachment.default.id}_fake"]`,
			"role_name": `"${alicloud_ram_role.defaultzEJqM7.id}"`,
		}),
	}

	RamRolePolicyAttachmentCheckInfo.dataSourceTestCheck(t, rand, RoleNameConf, allConf)
}

var existRamRolePolicyAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#":             "1",
		"attachments.0.policy_type": CHECKSET,
		"attachments.0.attach_date": CHECKSET,
		"attachments.0.policy_name": CHECKSET,
	}
}

var fakeRamRolePolicyAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#": "0",
	}
}

var RamRolePolicyAttachmentCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_role_policy_attachments.default",
	existMapFunc: existRamRolePolicyAttachmentMapFunc,
	fakeMapFunc:  fakeRamRolePolicyAttachmentMapFunc,
}

func testAccCheckAlicloudRamRolePolicyAttachmentSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRamRolePolicyAttachment%d"
}
resource "alicloud_ram_policy" "defaultGNmT8H" {
  policy_name     = "zaijiuRPTestPolicy"
  policy_document = "{\"Statement\": [{\"Effect\": \"Allow\",\"Action\": \"ecs:Describe*\",\"Resource\": \"acs:ecs:cn-qingdao:*:instance/*\"}],\"Version\": \"1\"}"
}

resource "alicloud_ram_role" "defaultzEJqM7" {
  name = "zaijiuRPTestPolicy"
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}



resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.defaultzEJqM7.id
  policy_name = alicloud_ram_policy.defaultGNmT8H.id
  policy_type = "Custom"
}

data "alicloud_ram_role_policy_attachments" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

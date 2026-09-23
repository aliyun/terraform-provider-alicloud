package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerControlPolicyAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerControlPolicyAttachmentsDataSourceName(rand, map[string]string{
			"target_id":   `"${alicloud_resource_manager_control_policy_attachment.default.target_id}"`,
			"policy_type": `"Custom"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerControlPolicyAttachmentsDataSourceName(rand, map[string]string{
			"target_id":   `"${alicloud_resource_manager_control_policy_attachment.default.target_id}"`,
			"policy_type": `"fake"`,
		}),
	}
	var existAlicloudResourceManagerControlPolicyAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"attachments.#":             "1",
			"attachments.0.policy_id":   CHECKSET,
			"attachments.0.attach_date": CHECKSET,
			"attachments.0.description": fmt.Sprintf("tf-testAccControlPolicyAttachment-%d", rand),
			"attachments.0.policy_name": fmt.Sprintf("tf-testAccControlPolicyAttachment-%d", rand),
			"attachments.0.policy_type": "Custom",
			"attachments.0.id":          CHECKSET,
		}
	}
	var fakeAlicloudResourceManagerControlPolicyAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudResourceManagerControlPolicyAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_control_policy_attachments.default",
		existMapFunc: existAlicloudResourceManagerControlPolicyAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudResourceManagerControlPolicyAttachmentsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	alicloudResourceManagerControlPolicyAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}
func testAccCheckAlicloudResourceManagerControlPolicyAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccControlPolicyAttachment-%[1]d"
}

resource "alicloud_resource_manager_folder" "example" {
    folder_name = "tf-testAcc-%[1]d"
}

resource "alicloud_resource_manager_control_policy" "example" {
	control_policy_name = var.name
	description = var.name
	effect_scope = "RAM"
	policy_document = "{\"Version\":\"1\",\"Statement\":[{\"Effect\":\"Deny\",\"Action\":[\"ram:UpdateRole\",\"ram:DeleteRole\",\"ram:AttachPolicyToRole\",\"ram:DetachPolicyFromRole\"],\"Resource\":\"acs:ram:*:*:role/ResourceDirectoryAccountAccessRole\"}]}"
}

resource "alicloud_resource_manager_control_policy_attachment" "default" {
	policy_id = alicloud_resource_manager_control_policy.example.id
	target_id = alicloud_resource_manager_folder.example.id
}

data "alicloud_resource_manager_control_policy_attachments" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

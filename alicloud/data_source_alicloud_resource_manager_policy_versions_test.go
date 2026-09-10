package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerPolicyVersionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPolicyVersionsSourceConfig(rand, map[string]string{
			"ids":         fmt.Sprintf(`["tf-testAccPolicy-%d:v1"]`, rand),
			"policy_name": `"${alicloud_resource_manager_policy.default.policy_name}"`,
			"policy_type": `"${alicloud_resource_manager_policy.default.policy_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPolicyVersionsSourceConfig(rand, map[string]string{
			"ids":         `["fake"]`,
			"policy_name": `"${alicloud_resource_manager_policy.default.policy_name}"`,
			"policy_type": `"${alicloud_resource_manager_policy.default.policy_type}"`,
		}),
	}

	var existResourceManagerPolicyVersionsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"versions.#":                    "1",
			"versions.0.policy_document":    CHECKSET,
			"versions.0.is_default_version": "true",
			"versions.0.version_id":         "v1",
		}
	}

	var fakeResourceManagerPolicyVersionsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}

	var policyVersionsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_policy_versions.default",
		existMapFunc: existResourceManagerPolicyVersionsRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerPolicyVersionsRecordsMapFunc,
	}

	policyVersionsRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudResourceManagerPolicyVersionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_resource_manager_policy" "default"{
	policy_name = "tf-testAccPolicy-%d"
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

data "alicloud_resource_manager_policy_versions" "default"{
	enable_details = true
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}

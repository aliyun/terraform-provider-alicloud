package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerPoliciesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_resource_manager_policies.example"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_policy.example.policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_policy.example.policy_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_policy.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_policy.example.id}_fake"]`,
		}),
	}

	policyTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}"`,
			"policy_type": `"System"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}"`,
			"ids":         `["${alicloud_resource_manager_policy.example.id}"]`,
			"policy_type": `"Custom"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}_fake"`,
			"ids":         `["${alicloud_resource_manager_policy.example.id}"]`,
			"policy_type": `"Custom"`,
		}),
	}

	var existResourceManagerPoliciesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#":                  "1",
			"names.#":                     "1",
			"ids.#":                       "1",
			"policies.0.attachment_count": CHECKSET,
			"policies.0.default_version":  "v1",
			"policies.0.description":      "policy_test",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      fmt.Sprintf("tf-testAccPolicy-%d", rand),
			"policies.0.policy_type":      "Custom",
			"policies.0.update_date":      CHECKSET,
		}
	}

	var fakeResourceManagerPoliciesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var policiesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existResourceManagerPoliciesRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerPoliciesRecordsMapFunc,
	}

	policiesRecordsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, policyTypeConf, allConf)

}

func testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_resource_manager_policy" "example"{
	policy_name = "tf-testAccPolicy-%d"
	description = "policy_test"
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

data "alicloud_resource_manager_policies" "example"{
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}

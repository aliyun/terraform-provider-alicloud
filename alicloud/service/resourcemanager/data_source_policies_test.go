package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_resource_manager_policies.example"

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_policy.example.policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_policy.example.policy_name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_policy.example.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_policy.example.id}_fake"]`,
		}),
	}

	policyTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}"`,
			"policy_type": `"System"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_resource_manager_policy.example.policy_name}"`,
			"ids":         `["${alicloud_resource_manager_policy.example.id}"]`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerPoliciesSourceConfig(rand, map[string]string{
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
			"policies.0.attachment_count": acctest.CHECKSET,
			"policies.0.default_version":  "v1",
			"policies.0.description":      "policy_test",
			"policies.0.id":               acctest.CHECKSET,
			"policies.0.policy_name":      fmt.Sprintf("tf-testAccPolicy-%d", rand),
			"policies.0.policy_type":      "Custom",
			"policies.0.update_date":      acctest.CHECKSET,
		}
	}

	var fakeResourceManagerPoliciesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var policiesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existResourceManagerPoliciesRecordsMapFunc,
		FakeMapFunc:  fakeResourceManagerPoliciesRecordsMapFunc,
	}

	policiesRecordsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, policyTypeConf, allConf)

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

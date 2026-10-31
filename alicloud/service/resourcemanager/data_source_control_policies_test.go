package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerControlPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_control_policy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_control_policy.default.id}_fake"]`,
		}),
	}
	policyTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_resource_manager_control_policy.default.id}"]`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_resource_manager_control_policy.default.id}"]`,
			"policy_type": `"System"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_control_policy.default.control_policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_control_policy.default.control_policy_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_resource_manager_control_policy.default.id}"]`,
			"name_regex":  `"${alicloud_resource_manager_control_policy.default.control_policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_resource_manager_control_policy.default.id}_fake"]`,
			"name_regex":  `"${alicloud_resource_manager_control_policy.default.control_policy_name}_fake"`,
			"policy_type": `"System"`,
		}),
	}
	var existAlicloudResourceManagerControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"policies.#":                     "1",
			"policies.0.control_policy_name": fmt.Sprintf("tf-testAccControlPolicy-%d", rand),
			"policies.0.description":         fmt.Sprintf("tf-testAccControlPolicy-%d", rand),
			"policies.0.effect_scope":        `RAM`,
			"policies.0.policy_document":     `{"Version":"1","Statement":[{"Effect":"Deny","Action":["ram:UpdateRole","ram:DeleteRole","ram:AttachPolicyToRole","ram:DetachPolicyFromRole"],"Resource":"acs:ram:*:*:role/ResourceDirectoryAccountAccessRole"}]}`,
		}
	}
	var fakeAlicloudResourceManagerControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudResourceManagerControlPoliciesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_control_policies.default",
		ExistMapFunc: existAlicloudResourceManagerControlPoliciesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudResourceManagerControlPoliciesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}
	alicloudResourceManagerControlPoliciesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, policyTypeConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudResourceManagerControlPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccControlPolicy-%d"
}

resource "alicloud_resource_manager_control_policy" "default" {
	control_policy_name = var.name
	description = var.name
	effect_scope = "RAM"
	policy_document = "{\"Version\":\"1\",\"Statement\":[{\"Effect\":\"Deny\",\"Action\":[\"ram:UpdateRole\",\"ram:DeleteRole\",\"ram:AttachPolicyToRole\",\"ram:DetachPolicyFromRole\"],\"Resource\":\"acs:ram:*:*:role/ResourceDirectoryAccountAccessRole\"}]}"
}

data "alicloud_resource_manager_control_policies" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

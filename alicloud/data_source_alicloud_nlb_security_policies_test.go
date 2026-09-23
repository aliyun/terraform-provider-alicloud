package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNlbSecurityPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_security_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_security_policy.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_security_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_nlb_security_policy.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_security_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_nlb_security_policy.default.resource_group_id}_fake"`,
		}),
	}
	securityPolicyNamesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_nlb_security_policy.default.id}"]`,
			"security_policy_names": `["${alicloud_nlb_security_policy.default.security_policy_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_nlb_security_policy.default.id}"]`,
			"security_policy_names": `["${alicloud_nlb_security_policy.default.security_policy_name}_fake"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_security_policy.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_security_policy.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nlb_security_policy.default.security_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nlb_security_policy.default.security_policy_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nlb_security_policy.default.id}"]`,
			"status": `"${alicloud_nlb_security_policy.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nlb_security_policy.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_nlb_security_policy.default.id}"]`,
			"name_regex":            `"${alicloud_nlb_security_policy.default.security_policy_name}"`,
			"resource_group_id":     `"${alicloud_nlb_security_policy.default.resource_group_id}"`,
			"security_policy_names": `["${alicloud_nlb_security_policy.default.security_policy_name}"]`,
			"status":                `"${alicloud_nlb_security_policy.default.status}"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_nlb_security_policy.default.id}_fake"]`,
			"name_regex":            `"${alicloud_nlb_security_policy.default.security_policy_name}_fake"`,
			"resource_group_id":     `"${alicloud_nlb_security_policy.default.resource_group_id}_fake"`,
			"security_policy_names": `["${alicloud_nlb_security_policy.default.security_policy_name}_fake"]`,
			"status":                `"Configuring"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	var existAlicloudNlbSecurityPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"policies.#":                      "1",
			"policies.0.ciphers.#":            "2",
			"policies.0.resource_group_id":    CHECKSET,
			"policies.0.security_policy_name": fmt.Sprintf("tf-testAccSecurityPolicy-%d", rand),
			"policies.0.tags.%":               "2",
			"policies.0.tags.Created":         "TF",
			"policies.0.tags.For":             "Acceptance-test",
			"policies.0.tls_versions.#":       "3",
			"policies.0.status":               CHECKSET,
			"policies.0.id":                   CHECKSET,
		}
	}
	var fakeAlicloudNlbSecurityPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudNlbSecurityPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nlb_security_policies.default",
		existMapFunc: existAlicloudNlbSecurityPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNlbSecurityPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudNlbSecurityPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, securityPolicyNamesConf, tagsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudNlbSecurityPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSecurityPolicy-%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_nlb_security_policy" "default" {
	resource_group_id   =    data.alicloud_resource_manager_resource_groups.default.ids.0
	security_policy_name = var.name
	ciphers  =      ["ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES128-SHA"]
	tls_versions  =    ["TLSv1.0", "TLSv1.1", "TLSv1.2"]
	tags = {
		Created = "TF"
		For = "Acceptance-test"
	}
}

data "alicloud_nlb_security_policies" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

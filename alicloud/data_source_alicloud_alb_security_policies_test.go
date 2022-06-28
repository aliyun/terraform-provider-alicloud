package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBSecurityPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_security_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_security_policy.default.id}_fake"]`,
		}),
	}
	policyIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_ids": `["${alicloud_alb_security_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_ids": `["${alicloud_alb_security_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
		}),
	}
	policyNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
		}),
	}
	groupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_security_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_security_policy.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_security_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_security_policy.default.resource_group_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_security_policy.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_security_policy.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_alb_security_policy.default.id}"]`,
			"name_regex":           `"${alicloud_alb_security_policy.default.security_policy_name}"`,
			"security_policy_ids":  `["${alicloud_alb_security_policy.default.id}"]`,
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
			"resource_group_id":    `"${alicloud_alb_security_policy.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_alb_security_policy.default.id}_fake"]`,
			"name_regex":           `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
			"security_policy_ids":  `["${alicloud_alb_security_policy.default.id}_fake"]`,
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
			"resource_group_id":    `"${alicloud_alb_security_policy.default.resource_group_id}_fake"`,
		}),
	}

	var existDataAlicloudAlbSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"policies.#":                      "1",
			"policies.0.security_policy_name": fmt.Sprintf("tf-testAccSecurityPolicy-%d", rand),
			"policies.0.tls_versions.#":       "1",
			"policies.0.ciphers.#":            "2",
			"policies.0.resource_group_id":    CHECKSET,
			"policies.0.id":                   CHECKSET,
			"policies.0.security_policy_id":   CHECKSET,
			"policies.0.status":               "Available",
		}
	}
	var fakeDataAlicloudAlbSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudAlbSecurityPolicyCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_security_policies.default",
		existMapFunc: existDataAlicloudAlbSecurityPoliciesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbSecurityPoliciesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbSecurityPolicyCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, policyIdsConf, nameRegexConf, policyNameConf, groupIdConf, statusConf, allConf)
}
func testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSecurityPolicy-%d"
}

resource "alicloud_alb_security_policy" "default" {
	security_policy_name = var.name
	tls_versions = ["TLSv1.0"]
	ciphers = ["ECDHE-ECDSA-AES128-SHA","AES256-SHA"]
}

data "alicloud_alb_security_policies" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

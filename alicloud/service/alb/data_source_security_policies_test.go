package alb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudALBSecurityPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_security_policy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_security_policy.default.id}_fake"]`,
		}),
	}
	policyIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_ids": `["${alicloud_alb_security_policy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_ids": `["${alicloud_alb_security_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
		}),
	}
	policyNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
		}),
	}
	groupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_security_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_security_policy.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_security_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_security_policy.default.resource_group_id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_security_policy.default.id}"]`,
			"status": `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_security_policy.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_alb_security_policy.default.id}"]`,
			"name_regex":           `"${alicloud_alb_security_policy.default.security_policy_name}"`,
			"security_policy_ids":  `["${alicloud_alb_security_policy.default.id}"]`,
			"security_policy_name": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
			"resource_group_id":    `"${alicloud_alb_security_policy.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
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
			"policies.0.resource_group_id":    acctest.CHECKSET,
			"policies.0.id":                   acctest.CHECKSET,
			"policies.0.security_policy_id":   acctest.CHECKSET,
			"policies.0.status":               "Available",
		}
	}
	var fakeDataAlicloudAlbSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudAlbSecurityPolicyCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alb_security_policies.default",
		ExistMapFunc: existDataAlicloudAlbSecurityPoliciesSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudAlbSecurityPoliciesSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbSecurityPolicyCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, policyIdsConf, nameRegexConf, policyNameConf, groupIdConf, statusConf, allConf)
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

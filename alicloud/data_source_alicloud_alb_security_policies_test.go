package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlbSecurityPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_security_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_security_policy.default.id}_fake"]`,
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
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alb_security_policy.default.id}"]`,
			"name_regex": `"${alicloud_alb_security_policy.default.security_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbSecurityPolicieDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alb_security_policy.default.id}_fake"]`,
			"name_regex": `"${alicloud_alb_security_policy.default.security_policy_name}_fake"`,
		}),
	}

	var existDataAlicloudAlbSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"policies.#":                      "1",
			"policies.0.security_policy_name": fmt.Sprintf("tf-testAccSecurityPolicy-%d", rand),
			"policies.0.tls_versions.#":       "1",
			"policies.0.ciphers.#":            "2",
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
	alicloudAlbSecurityPolicyCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
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

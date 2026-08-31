package alicloud

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlicloudSLBTlsCipherPoliciesDataSource_basic(t *testing.T) {

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}_fake"`,
		}),
	}

	policyNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"tls_cipher_policy_name": `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"tls_cipher_policy_name": `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids": `["${alicloud_slb_tls_cipher_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids": `["${alicloud_slb_tls_cipher_policy.default.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids":    `["${alicloud_slb_tls_cipher_policy.default.id}"]`,
			"status": `"normal"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids":    `["${alicloud_slb_tls_cipher_policy.default.id}_fake"]`,
			"status": `"configuring"`,
		}),
	}

	includeListenerConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids":              `["${alicloud_slb_tls_cipher_policy.default.id}"]`,
			"include_listener": `true`,
			"status":           `"normal"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids":              `["${alicloud_slb_tls_cipher_policy.default.id}_fake"]`,
			"include_listener": `false`,
			"status":           `"configuring"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids":                    `["${alicloud_slb_tls_cipher_policy.default.id}"]`,
			"name_regex":             `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}"`,
			"tls_cipher_policy_name": `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}"`,
			"status":                 `"normal"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(map[string]string{
			"ids":                    `["${alicloud_slb_tls_cipher_policy.default.id}_fake"]`,
			"name_regex":             `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}"`,
			"tls_cipher_policy_name": `"${alicloud_slb_tls_cipher_policy.default.tls_cipher_policy_name}_fake"`,
			"status":                 `"configuring"`,
		}),
	}

	var existSLBTlsCipherPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"policies.#":                        "1",
			"policies.0.id":                     CHECKSET,
			"policies.0.tls_cipher_policy_name": "Tf-testAccSlbTlsBasic",
		}
	}

	var fakeSLBTlsCipherPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"policies.#": "0",
		}
	}

	var slbTlsCipherPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_tls_cipher_policies.default",
		existMapFunc: existSLBTlsCipherPoliciesMapFunc,
		fakeMapFunc:  fakeSLBTlsCipherPoliciesMapFunc,
	}

	slbTlsCipherPoliciesCheckInfo.dataSourceTestCheck(t, -1, nameRegexConf, policyNameConf, idsConf, statusConf, includeListenerConf, allConf)
}

func testAccCheckAlicloudSlbTlsCipherPoliciesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "Tf-testAccSlbTlsBasic"
}

resource "alicloud_slb_tls_cipher_policy" "default" {
  tls_cipher_policy_name = var.name
  tls_versions           = ["TLSv1.2"]
  ciphers                = ["AES256-SHA256", "AES128-GCM-SHA256"]
}

data "alicloud_slb_tls_cipher_policies" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

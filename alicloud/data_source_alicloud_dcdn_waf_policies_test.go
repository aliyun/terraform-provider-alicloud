package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDcdnWafPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dcdn_waf_policy.default.policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dcdn_waf_policy.default.policy_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dcdn_waf_policy.default.id}"]`,
			"status": `"on"`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dcdn_waf_policy.default.id}"]`,
			"status": `"off"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dcdn_waf_policy.default.id}"]`,
			"status":     `"on"`,
			"name_regex": `"${alicloud_dcdn_waf_policy.default.policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dcdn_waf_policy.default.id}_fake"]`,
			"status":     `"off"`,
			"name_regex": `"${alicloud_dcdn_waf_policy.default.policy_name}_fake"`,
		}),
	}
	var existAlicloudDcdnWafPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"policies.#":                    "1",
			"policies.0.defense_scene":      "waf_group",
			"policies.0.policy_name":        fmt.Sprintf("tf_testAccWafPolicy_%d", rand),
			"policies.0.policy_type":        "custom",
			"policies.0.status":             "on",
			"policies.0.id":                 CHECKSET,
			"policies.0.dcdn_waf_policy_id": CHECKSET,
			"policies.0.rule_count":         CHECKSET,
			"policies.0.gmt_modified":       CHECKSET,
			"policies.0.domain_count":       CHECKSET,
		}
	}
	var fakeAlicloudDcdnWafPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDcdnWafPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dcdn_waf_policies.default",
		existMapFunc: existAlicloudDcdnWafPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDcdnWafPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDcdnWafPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudDcdnWafPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testAccWafPolicy_%d"
}

resource "alicloud_dcdn_waf_policy" "default" {
	defense_scene = "waf_group"
	policy_name = var.name
	policy_type = "custom"
	status = "on"
}

data "alicloud_dcdn_waf_policies" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

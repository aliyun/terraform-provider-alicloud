package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASLifecyclePoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_lifecycle_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_lifecycle_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_lifecycle_policy.default.lifecycle_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_lifecycle_policy.default.lifecycle_policy_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nas_lifecycle_policy.default.id}"]`,
			"name_regex": `"${alicloud_nas_lifecycle_policy.default.lifecycle_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nas_lifecycle_policy.default.id}_fake"]`,
			"name_regex": `"${alicloud_nas_lifecycle_policy.default.lifecycle_policy_name}_fake"`,
		}),
	}
	var existAlicloudNasLifecyclePoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"policies.#":                       "1",
			"policies.0.create_time":           CHECKSET,
			"policies.0.id":                    CHECKSET,
			"policies.0.file_system_id":        CHECKSET,
			"policies.0.lifecycle_policy_name": fmt.Sprintf("tf-testAccLifecyclePolicy-%d", rand),
			"policies.0.lifecycle_rule_name":   "DEFAULT_ATIME_14",
			"policies.0.paths.#":               "1",
			"policies.0.storage_type":          "InfrequentAccess",
		}
	}
	var fakeAlicloudNasLifecyclePoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudNasLifecyclePoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nas_lifecycle_policies.default",
		existMapFunc: existAlicloudNasLifecyclePoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNasLifecyclePoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudNasLifecyclePoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudNasLifecyclePoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccLifecyclePolicy-%d"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_lifecycle_policy" "default" {
  file_system_id        = "${alicloud_nas_file_system.default.id}"
  lifecycle_policy_name = var.name
  lifecycle_rule_name   = "DEFAULT_ATIME_14"
  paths                 = ["/"]
  storage_type          = "InfrequentAccess"
}

data "alicloud_nas_lifecycle_policies" "default" {
  file_system_id = alicloud_nas_file_system.default.id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

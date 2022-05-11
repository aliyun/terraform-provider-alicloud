package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASAutoSnapshotPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_auto_snapshot_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nas_auto_snapshot_policy.default.id}_fake"]`,
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}_fake"`,
			"status":     `"Creating"`,
		}),
	}
	var existAlicloudNasAutoSnapshotPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"policies.#":                           "1",
			"policies.0.auto_snapshot_policy_name": fmt.Sprintf("tf-testAccAutoSnapshotPolicy-%d", rand),
			"policies.0.create_time":               CHECKSET,
			"policies.0.auto_snapshot_policy_id":   CHECKSET,
			"policies.0.id":                        CHECKSET,
			"policies.0.file_system_nums":          CHECKSET,
			"policies.0.repeat_weekdays.#":         "3",
			"policies.0.retention_days":            "30",
			"policies.0.time_points.#":             "3",
			"policies.0.status":                    "Available",
		}
	}
	var fakeAlicloudNasAutoSnapshotPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudNasAutoSnapshotPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nas_auto_snapshot_policies.default",
		existMapFunc: existAlicloudNasAutoSnapshotPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNasAutoSnapshotPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudNasAutoSnapshotPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccAutoSnapshotPolicy-%d"
}

resource "alicloud_nas_auto_snapshot_policy" "default" {
  auto_snapshot_policy_name = var.name
  repeat_weekdays           = ["3","4","5"]
  retention_days            = 30
  time_points               = ["1","2","3"]
}

data "alicloud_nas_auto_snapshot_policies" "default" {	
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

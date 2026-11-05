package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASAutoSnapshotPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_auto_snapshot_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
			"status": `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nas_auto_snapshot_policy.default.id}"]`,
			"name_regex": `"${alicloud_nas_auto_snapshot_policy.default.auto_snapshot_policy_name}"`,
			"status":     `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudNasAutoSnapshotPoliciesDataSourceName(rand, map[string]string{
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
			"policies.0.create_time":               acctest.CHECKSET,
			"policies.0.auto_snapshot_policy_id":   acctest.CHECKSET,
			"policies.0.id":                        acctest.CHECKSET,
			"policies.0.file_system_nums":          acctest.CHECKSET,
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
	var alicloudNasAutoSnapshotPoliciesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nas_auto_snapshot_policies.default",
		ExistMapFunc: existAlicloudNasAutoSnapshotPoliciesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudNasAutoSnapshotPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudNasAutoSnapshotPoliciesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
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

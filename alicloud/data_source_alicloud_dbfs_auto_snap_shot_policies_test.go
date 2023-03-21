package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDbfsAutoSnapShotPolicyDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.DBFSSystemSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dbfs_auto_snap_shot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dbfs_auto_snap_shot_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dbfs_auto_snap_shot_policy.default.id}"]`,
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dbfs_auto_snap_shot_policy.default.id}_fake"]`,
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}_fake"`,
		}),
	}

	DbfsAutoSnapShotPolicyCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existDbfsAutoSnapShotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"names.#":                      "1",
		"auto_snap_shot_policies.#":    "1",
		"auto_snap_shot_policies.0.id": CHECKSET,
		"auto_snap_shot_policies.0.applied_dbfs_number": CHECKSET,
		"auto_snap_shot_policies.0.create_time":         CHECKSET,
		"auto_snap_shot_policies.0.last_modified":       CHECKSET,
		"auto_snap_shot_policies.0.policy_id":           CHECKSET,
		"auto_snap_shot_policies.0.policy_name":         CHECKSET,
		"auto_snap_shot_policies.0.repeat_weekdays.#":   "1",
		"auto_snap_shot_policies.0.retention_days":      CHECKSET,
		"auto_snap_shot_policies.0.status":              CHECKSET,
		"auto_snap_shot_policies.0.status_detail":       "",
		"auto_snap_shot_policies.0.time_points.#":       "1",
	}
}

var fakeDbfsAutoSnapShotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "0",
		"names.#":                   "0",
		"auto_snap_shot_policies.#": "0",
	}
}

var DbfsAutoSnapShotPolicyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dbfs_auto_snap_shot_policies.default",
	existMapFunc: existDbfsAutoSnapShotPolicyMapFunc,
	fakeMapFunc:  fakeDbfsAutoSnapShotPolicyMapFunc,
}

func testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDbfsAutoSnapShotPolicy%d"
}

resource "alicloud_dbfs_auto_snap_shot_policy" "default" {
  time_points = ["01"]
  policy_name    = var.name
  retention_days = 1
  repeat_weekdays = ["2"]
}

data "alicloud_dbfs_auto_snap_shot_policies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

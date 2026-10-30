package dbfs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDbfsAutoSnapShotPolicyDataSource(t *testing.T) {
	t.Skip("This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/")
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.DBFSSystemSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dbfs_auto_snap_shot_policy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dbfs_auto_snap_shot_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dbfs_auto_snap_shot_policy.default.id}"]`,
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsAutoSnapShotPolicySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dbfs_auto_snap_shot_policy.default.id}_fake"]`,
			"name_regex": `"${alicloud_dbfs_auto_snap_shot_policy.default.policy_name}_fake"`,
		}),
	}

	DbfsAutoSnapShotPolicyCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existDbfsAutoSnapShotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"names.#":                      "1",
		"auto_snap_shot_policies.#":    "1",
		"auto_snap_shot_policies.0.id": acctest.CHECKSET,
		"auto_snap_shot_policies.0.applied_dbfs_number": acctest.CHECKSET,
		"auto_snap_shot_policies.0.create_time":         acctest.CHECKSET,
		"auto_snap_shot_policies.0.last_modified":       acctest.CHECKSET,
		"auto_snap_shot_policies.0.policy_id":           acctest.CHECKSET,
		"auto_snap_shot_policies.0.policy_name":         acctest.CHECKSET,
		"auto_snap_shot_policies.0.repeat_weekdays.#":   "1",
		"auto_snap_shot_policies.0.retention_days":      acctest.CHECKSET,
		"auto_snap_shot_policies.0.status":              acctest.CHECKSET,
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

var DbfsAutoSnapShotPolicyCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_dbfs_auto_snap_shot_policies.default",
	ExistMapFunc: existDbfsAutoSnapShotPolicyMapFunc,
	FakeMapFunc:  fakeDbfsAutoSnapShotPolicyMapFunc,
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

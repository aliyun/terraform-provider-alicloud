// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcsAutoSnapshotPolicyDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsAutoSnapshotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_auto_snapshot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsAutoSnapshotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_auto_snapshot_policy.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsAutoSnapshotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_auto_snapshot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsAutoSnapshotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_auto_snapshot_policy.default.id}_fake"]`,
		}),
	}

	EcsAutoSnapshotPolicyCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existEcsAutoSnapshotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#":                                 "1",
		"policies.0.status":                          CHECKSET,
		"policies.0.record_total":                    CHECKSET,
		"policies.0.time_points.#":                   CHECKSET,
		"policies.0.volume_nums":                     CHECKSET,
		"policies.0.create_time":                     CHECKSET,
		"policies.0.auto_snapshot_policy_id":         CHECKSET,
		"policies.0.retention_days":                  "-1",
		"policies.0.repeat_weekdays.#":               CHECKSET,
		"policies.0.disk_nums":                       CHECKSET,
		"policies.0.copied_snapshots_retention_days": CHECKSET,
		"policies.0.target_copy_regions.#":           CHECKSET,
		"policies.0.association_type":                "AssociatedWithDisk",
		"policies.0.enable_cross_region_copy":        CHECKSET,
		"policies.0.target_tags.#":                   CHECKSET,
		"policies.0.region_id":                       CHECKSET,
		"policies.0.auto_snapshot_policy_name":       CHECKSET,
		"policies.0.tags.%":                          CHECKSET,
	}
}

var fakeEcsAutoSnapshotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#": "0",
	}
}

var EcsAutoSnapshotPolicyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ecs_auto_snapshot_policies.default",
	existMapFunc: existEcsAutoSnapshotPolicyMapFunc,
	fakeMapFunc:  fakeEcsAutoSnapshotPolicyMapFunc,
}

func testAccCheckAlicloudEcsAutoSnapshotPolicySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEcsAutoSnapshotPolicy%d"
}


resource "alicloud_ecs_auto_snapshot_policy" "default" {
	name            = var.name
	repeat_weekdays = ["1"]
	retention_days  = -1
	time_points     = ["1"]
}

data "alicloud_ecs_auto_snapshot_policies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

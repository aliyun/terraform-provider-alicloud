package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSAutoSnapshotPoliciesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_auto_snapshot_policies.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsAutoSnapshotPoliciesTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsAutoSnapshotPoliciesDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_auto_snapshot_policy.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Normal",
			"ids":    []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "Expire",
			"ids":    []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
			"ids": []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
			"ids": []string{"${alicloud_ecs_auto_snapshot_policy.default.id}"},
		}),
	}
	var existEcsAutoSnapshotPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"ids.0":                              CHECKSET,
			"names.#":                            "1",
			"names.0":                            name,
			"policies.#":                         "1",
			"policies.0.id":                      CHECKSET,
			"policies.0.auto_snapshot_policy_id": CHECKSET,
			"policies.0.copied_snapshots_retention_days": "-1",
			"policies.0.disk_nums":                       "0",
			"policies.0.enable_cross_region_copy":        "false",
			"policies.0.name":                            name,
			"policies.0.repeat_weekdays.#":               "1",
			"policies.0.retention_days":                  "-1",
			"policies.0.status":                          "Normal",
			"policies.0.volume_nums":                     "0",
			"policies.0.time_points.#":                   "1",
		}
	}

	var fakeEcsAutoSnapshotPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"policies.#": "0",
		}
	}

	var EcsAutoSnapshotPoliciesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsAutoSnapshotPoliciesMapFunc,
		fakeMapFunc:  fakeEcsAutoSnapshotPoliciesMapFunc,
	}

	EcsAutoSnapshotPoliciesInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, statusConf, tagsConf)
}

func dataSourceEcsAutoSnapshotPoliciesDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ecs_auto_snapshot_policy" "default" {
		name              = "%s"
		repeat_weekdays   = ["1"]
		retention_days    =  -1
		time_points       = ["1"]
		tags 	 = {
			Created = "TF"
			For 	= "acceptance test"
		}
	}`, name)
}

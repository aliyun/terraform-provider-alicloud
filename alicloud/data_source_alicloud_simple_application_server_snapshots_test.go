package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerSnapshotsDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_snapshots.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.SWASSupportRegions)
	name := fmt.Sprintf("tf-testacc-swas_snapshots-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerSnapshotsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_snapshot.default.snapshot_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_snapshot.default.snapshot_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_snapshot.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_snapshot.default.id}-fake"},
		}),
	}
	diskIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id": "${alicloud_simple_application_server_snapshot.default.disk_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id": "${alicloud_simple_application_server_snapshot.default.disk_id}-fake",
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}-fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"status": "Accomplished",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"status": "Failed",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_simple_application_server_snapshot.default.snapshot_name}",
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id":     "${alicloud_simple_application_server_snapshot.default.disk_id}",
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}",
			"status":      "Accomplished",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_simple_application_server_snapshot.default.snapshot_name}-fake",
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id":     "${alicloud_simple_application_server_snapshot.default.disk_id}-fake",
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}-fake",
			"status":      "Failed",
		}),
	}
	var existSimpleApplicationServerSnapshotMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"snapshots.#":                  "1",
			"snapshots.0.snapshot_name":    fmt.Sprintf("tf-testacc-swas_snapshots-%d", rand),
			"snapshots.0.status":           "Accomplished",
			"snapshots.0.create_time":      CHECKSET,
			"snapshots.0.disk_id":          CHECKSET,
			"snapshots.0.progress":         CHECKSET,
			"snapshots.0.remark":           "",
			"snapshots.0.id":               CHECKSET,
			"snapshots.0.snapshot_id":      CHECKSET,
			"snapshots.0.source_disk_type": CHECKSET,
		}
	}

	var fakeSimpleApplicationServerSnapshotMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"snapshots.#": "0",
		}
	}

	var SimpleApplicationServerSnapshotCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSimpleApplicationServerSnapshotMapFunc,
		fakeMapFunc:  fakeSimpleApplicationServerSnapshotMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}

	SimpleApplicationServerSnapshotCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, diskIdConf, instanceIdConf, statusConf, allConf)
}

func dataSourceSimpleApplicationServerSnapshotsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_simple_application_server_images" "default" {
	platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
	platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

data "alicloud_simple_application_server_disks" "default" {
  instance_id = alicloud_simple_application_server_instance.default.id
}

resource "alicloud_simple_application_server_snapshot" "default" {
  snapshot_name = var.name
  disk_id       = data.alicloud_simple_application_server_disks.default.ids.0
}`, name)
}

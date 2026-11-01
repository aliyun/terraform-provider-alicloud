package swasopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerSnapshotsDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_snapshots.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.SWASSupportRegions)
	name := fmt.Sprintf("tf-testacc-swas_snapshots-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerSnapshotsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_snapshot.default.snapshot_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_snapshot.default.snapshot_name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_snapshot.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_snapshot.default.id}-fake"},
		}),
	}
	diskIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id": "${alicloud_simple_application_server_snapshot.default.disk_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id": "${alicloud_simple_application_server_snapshot.default.disk_id}-fake",
		}),
	}
	instanceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}-fake",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"status": "Accomplished",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"status": "Failed",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_simple_application_server_snapshot.default.snapshot_name}",
			"ids":         []string{"${alicloud_simple_application_server_snapshot.default.id}"},
			"disk_id":     "${alicloud_simple_application_server_snapshot.default.disk_id}",
			"instance_id": "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}",
			"status":      "Accomplished",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                        acctest.CHECKSET,
			"snapshots.#":                  "1",
			"snapshots.0.snapshot_name":    fmt.Sprintf("tf-testacc-swas_snapshots-%d", rand),
			"snapshots.0.status":           "accomplished",
			"snapshots.0.create_time":      acctest.CHECKSET,
			"snapshots.0.disk_id":          acctest.CHECKSET,
			"snapshots.0.progress":         acctest.CHECKSET,
			"snapshots.0.remark":           "",
			"snapshots.0.id":               acctest.CHECKSET,
			"snapshots.0.snapshot_id":      acctest.CHECKSET,
			"snapshots.0.source_disk_type": acctest.CHECKSET,
		}
	}

	var fakeSimpleApplicationServerSnapshotMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"snapshots.#": "0",
		}
	}

	var SimpleApplicationServerSnapshotCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existSimpleApplicationServerSnapshotMapFunc,
		FakeMapFunc:  fakeSimpleApplicationServerSnapshotMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}

	SimpleApplicationServerSnapshotCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, diskIdConf, instanceIdConf, statusConf, allConf)
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

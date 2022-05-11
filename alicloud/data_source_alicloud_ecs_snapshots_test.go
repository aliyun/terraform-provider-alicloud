package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSSnapshotsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAcc%d", rand)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_snapshot.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_snapshot.default.id}_fake"]`,
		}),
	}
	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":      `["${alicloud_ecs_snapshot.default.id}"]`,
			"category": `"${alicloud_ecs_snapshot.default.category}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":      `["${alicloud_ecs_snapshot.default.id}"]`,
			"category": `"flash"`,
		}),
	}
	usageConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":   `["${alicloud_ecs_snapshot.default.id}"]`,
			"usage": `"none"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":   `["${alicloud_ecs_snapshot.default.id}"]`,
			"usage": `"image"`,
		}),
	}

	snapshotNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_name": `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_name": `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	sourceDiskTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":              `["${alicloud_ecs_snapshot.default.id}"]`,
			"source_disk_type": `"Data"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":              `["${alicloud_ecs_snapshot.default.id}"]`,
			"source_disk_type": `"System"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_snapshot.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_snapshot.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":    `["${alicloud_ecs_snapshot.default.id}"]`,
			"status": `"accomplished"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":    `["${alicloud_ecs_snapshot.default.id}"]`,
			"status": `"failed"`,
		}),
	}
	snapTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_type": `"user"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_type": `"auto"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"category":         `"${alicloud_ecs_snapshot.default.category}"`,
			"ids":              `["${alicloud_ecs_snapshot.default.id}"]`,
			"name_regex":       `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
			"snapshot_name":    `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
			"source_disk_type": `"Data"`,
			"snapshot_type":    `"user"`,
			"status":           `"accomplished"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotsDataSourceName(name, map[string]string{
			"category":         `"${alicloud_ecs_snapshot.default.category}"`,
			"ids":              `["${alicloud_ecs_snapshot.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
			"snapshot_name":    `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
			"source_disk_type": `"System"`,
			"status":           `"failed"`,
			"snapshot_type":    `"auto"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	var existAlicloudEcsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"snapshots.#":                     "1",
			"snapshots.0.category":            `standard`,
			"snapshots.0.description":         `Test For Terraform`,
			"snapshots.0.disk_id":             CHECKSET,
			"snapshots.0.retention_days":      `20`,
			"snapshots.0.snapshot_name":       CHECKSET,
			"snapshots.0.usage":               CHECKSET,
			"snapshots.0.source_storage_type": CHECKSET,
			"snapshots.0.status":              "accomplished",
			"snapshots.0.source_disk_id":      CHECKSET,
			"snapshots.0.snapshot_type":       CHECKSET,
			"snapshots.0.snapshot_sn":         CHECKSET,
			"snapshots.0.product_code":        "",
			"snapshots.0.progress":            CHECKSET,
			"snapshots.0.tags.%":              "2",
			"snapshots.0.tags.Created":        "TF",
			"snapshots.0.tags.For":            "Acceptance-test",
		}
	}
	var fakeAlicloudEcsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcsSnapshotsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_snapshots.default",
		existMapFunc: existAlicloudEcsSnapshotsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsSnapshotsDataSourceNameMapFunc,
	}
	alicloudEcsSnapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConf, categoryConf, usageConfig, snapshotNameConf, sourceDiskTypeConf, tagsConf, nameRegexConf, statusConf, snapTypeConf, allConf)
}
func testAccCheckAlicloudEcsSnapshotsDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
			default = "%s"
		}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_disk" "default" {
  count = "2"
  name = "${var.name}"
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category          = "cloud_efficiency"
  size              = "20"
}

data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_disk_attachment" "default" {
  count = "2"
  disk_id     = "${element(alicloud_disk.default.*.id,count.index)}"
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot" "default" {
	category = "standard"
	description = "Test For Terraform"
	disk_id = alicloud_disk_attachment.default.0.disk_id
	retention_days = "20"
	snapshot_name = var.name
	tags 				 = {
		Created = "TF"
		For 	= "Acceptance-test"
	}
}

data "alicloud_ecs_snapshots" "default" {	
	%s
}
`, name, strings.Join(pairs, " \n "))
	return config
}

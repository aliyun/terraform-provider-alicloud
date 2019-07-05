package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudSnapshotsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccSnapshotDataSourceBasic%d", rand)
	resourceId := "data.alicloud_snapshots.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSnapshotsConfigDependence)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_snapshot.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_snapshot.default.id}_fake"},
		}),
	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_instance.default.id}_fake",
		}),
	}

	diskIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${alicloud_disk_attachment.default.disk_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${alicloud_disk_attachment.default.disk_id}_fake",
		}),
	}

	encryptedConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alicloud_snapshot.default.id}"},
			"encrypted": "false",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alicloud_snapshot.default.id}"},
			"encrypted": "true",
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "_fake",
		}),
	}

	statusConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_snapshot.default.id}"},
			"status": "accomplished",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_snapshot.default.id}"},
			"status": "failed",
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_snapshot.default.id}"},
			"type": "user",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_snapshot.default.id}"},
			"type": "auto",
		}),
	}

	sourceDiskTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_snapshot.default.id}"},
			"source_disk_type": "Data",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_snapshot.default.id}"},
			"source_disk_type": "System",
		}),
	}

	usageConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":   []string{"${alicloud_snapshot.default.id}"},
			"usage": "none",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":   []string{"${alicloud_snapshot.default.id}"},
			"usage": "image",
		}),
	}

	tagsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_snapshot.default.id}"},
			"tags": map[string]interface{}{
				"version": "1.0",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_snapshot.default.id}"},
			"tags": map[string]interface{}{
				"version": "1.0_fake",
			},
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_snapshot.default.id}"},
			"instance_id":      "${alicloud_instance.default.id}",
			"disk_id":          "${alicloud_disk_attachment.default.disk_id}",
			"encrypted":        "false",
			"name_regex":       name,
			"status":           "accomplished",
			"type":             "user",
			"source_disk_type": "Data",
			"usage":            "none",
			"tags": map[string]interface{}{
				"version": "1.0",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_snapshot.default.id}"},
			"instance_id":      "${alicloud_instance.default.id}",
			"disk_id":          "${alicloud_disk_attachment.default.disk_id}",
			"encrypted":        "false",
			"name_regex":       name,
			"status":           "accomplished",
			"type":             "user",
			"source_disk_type": "Data",
			"usage":            "none",
			"tags": map[string]interface{}{
				"version": "1.0_fake",
			},
		}),
	}

	var existSnapshotsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"snapshots.#":                  "1",
			"snapshots.0.id":               CHECKSET,
			"snapshots.0.name":             name,
			"snapshots.0.description":      name,
			"snapshots.0.progress":         CHECKSET,
			"snapshots.0.source_disk_id":   CHECKSET,
			"snapshots.0.source_disk_size": "20",
			"snapshots.0.source_disk_type": CHECKSET,
			"snapshots.0.product_code":     "",
			"snapshots.0.retention_days":   CHECKSET,
			"snapshots.0.remain_time":      CHECKSET,
			"snapshots.0.creation_time":    CHECKSET,
			"snapshots.0.status":           "accomplished",
			"snapshots.0.usage":            "none",
		}
	}

	var fakeSnapshotsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"snapshots.#": "0",
		}
	}

	var snapshotsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSnapshotsMapFunc,
		fakeMapFunc:  fakeSnapshotsMapFunc,
	}

	snapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConfig, instanceIdConfig, diskIdConfig, encryptedConfig, nameRegexConfig,
		statusConfig, typeConfig, sourceDiskTypeConfig, usageConfig, tagsConfig, allConfig)
}

func dataSourceSnapshotsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_instance_types" "default" {
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "default" {
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
  name        = "${var.name}"
  description = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_disk" "default" {
  availability_zone = "${alicloud_instance.default.availability_zone}"
  category          = "cloud_efficiency"
  size              = "20"
}

data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  instance_name   = "${var.name}"
  image_id        = "${data.alicloud_images.default.images.0.id}"
  instance_type   = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.default.id}"]
  vswitch_id      = "${alicloud_vswitch.default.id}"
}

resource "alicloud_disk_attachment" "default" {
  disk_id     = "${alicloud_disk.default.id}"
  instance_id = "${alicloud_instance.default.id}"
}

resource "alicloud_snapshot" "default" {
  disk_id = "${alicloud_disk_attachment.default.disk_id}"
  name = "${var.name}"
  description = "${var.name}"
  tags = {
    version = "1.0"
  }
}
`, name)
}

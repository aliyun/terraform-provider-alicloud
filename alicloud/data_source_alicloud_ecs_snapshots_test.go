package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudECSSnapshotsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAcc%d", rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_snapshot.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_snapshot.default.id}_fake"]`,
		}),
	}

	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":      `["${alicloud_ecs_snapshot.default.id}"]`,
			"category": `"${alicloud_ecs_snapshot.default.category}"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":      `["${alicloud_ecs_snapshot.default.id}"]`,
			"category": `"flash"`,
		}),
	}

	usageConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":   `["${alicloud_ecs_snapshot.default.id}"]`,
			"usage": `"none"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":   `["${alicloud_ecs_snapshot.default.id}"]`,
			"usage": `"image"`,
		}),
	}

	snapshotNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_name": `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_name": `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
		}),
	}

	sourceDiskTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":              `["${alicloud_ecs_snapshot.default.id}"]`,
			"source_disk_type": `"Data"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":              `["${alicloud_ecs_snapshot.default.id}"]`,
			"source_disk_type": `"System"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"tags": `{
				"Created" = "TF"
				"For" = "Snapshot"
		}`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Snapshot-fake"
			}`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":    `["${alicloud_ecs_snapshot.default.id}"]`,
			"status": `"accomplished"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":    `["${alicloud_ecs_snapshot.default.id}"]`,
			"status": `"failed"`,
		}),
	}

	snapTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_type": `"user"`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"ids":           `["${alicloud_ecs_snapshot.default.id}"]`,
			"snapshot_type": `"auto"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"category":         `"${alicloud_ecs_snapshot.default.category}"`,
			"ids":              `["${alicloud_ecs_snapshot.default.id}"]`,
			"name_regex":       `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
			"snapshot_name":    `"${alicloud_ecs_snapshot.default.snapshot_name}"`,
			"source_disk_type": `"Data"`,
			"snapshot_type":    `"user"`,
			"status":           `"accomplished"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Snapshot"
		}`,
		}),
		fakeConfig: testAccCheckAliCloudEcsSnapshotsDataSourceName(name, map[string]string{
			"category":         `"${alicloud_ecs_snapshot.default.category}"`,
			"ids":              `["${alicloud_ecs_snapshot.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
			"snapshot_name":    `"${alicloud_ecs_snapshot.default.snapshot_name}_fake"`,
			"source_disk_type": `"System"`,
			"status":           `"failed"`,
			"snapshot_type":    `"auto"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Snapshot-fake"
			}`,
		}),
	}

	var existAliCloudEcsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"snapshots.#":                     "1",
			"snapshots.0.category":            CHECKSET,
			"snapshots.0.description":         CHECKSET,
			"snapshots.0.disk_id":             CHECKSET,
			"snapshots.0.retention_days":      CHECKSET,
			"snapshots.0.snapshot_name":       CHECKSET,
			"snapshots.0.usage":               CHECKSET,
			"snapshots.0.source_storage_type": CHECKSET,
			"snapshots.0.status":              CHECKSET,
			"snapshots.0.source_disk_id":      CHECKSET,
			"snapshots.0.snapshot_type":       CHECKSET,
			"snapshots.0.snapshot_sn":         CHECKSET,
			"snapshots.0.product_code":        "",
			"snapshots.0.progress":            CHECKSET,
			"snapshots.0.tags.%":              "2",
			"snapshots.0.tags.Created":        "TF",
			"snapshots.0.tags.For":            "Snapshot",
		}
	}
	var fakeAliCloudEcsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"snapshots.#": "0",
		}
	}
	var alicloudEcsSnapshotsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_snapshots.default",
		existMapFunc: existAliCloudEcsSnapshotsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudEcsSnapshotsDataSourceNameMapFunc,
	}
	alicloudEcsSnapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConf, categoryConf, usageConfig, snapshotNameConf, sourceDiskTypeConf, tagsConf, nameRegexConf, statusConf, snapTypeConf, allConf)
}

func testAccCheckAliCloudEcsSnapshotsDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
  		available_resource_creation = "VSwitch"
	}
	
	data "alicloud_images" "default" {
  		most_recent = true
  		owners      = "system"
	}
	
	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		image_id             = data.alicloud_images.default.images.0.id
  		system_disk_category = "cloud_essd"
	}
	
	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_vswitch" "default" {
		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
	
	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
		data_disks {
			category = "cloud_essd"
			size     = 20
  		}
	}
	
	resource "alicloud_ecs_disk" "default" {
  		disk_name = var.name
  		zone_id   = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		category  = "cloud_essd"
  		size      = 500
	}
	
	resource "alicloud_ecs_disk_attachment" "default" {
  		disk_id     = alicloud_ecs_disk.default.id
  		instance_id = alicloud_instance.default.id
	}

	resource "alicloud_ecs_snapshot" "default" {
  		disk_id        = alicloud_ecs_disk_attachment.default.disk_id
  		category       = "standard"
  		retention_days = 20
  		snapshot_name  = var.name
  		description    = var.name
  		tags = {
    		Created = "TF"
    		For     = "Snapshot"
  		}
	}

	data "alicloud_ecs_snapshots" "default" {	
		%s
	}
`, name, strings.Join(pairs, " \n "))
	return config
}

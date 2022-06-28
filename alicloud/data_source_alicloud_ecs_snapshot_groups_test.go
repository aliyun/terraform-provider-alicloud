package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSSnapshotGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_snapshot_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_snapshot_group.default.id}_fake"]`,
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"instance_id": `"${alicloud_ecs_snapshot_group.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_ecs_snapshot_group.default.id}_fake"]`,
			"instance_id": `"${alicloud_ecs_snapshot_group.default.instance_id}"`,
		}),
	}
	snapshotGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"snapshot_group_name": `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"snapshot_group_name": `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"status": `"accomplished"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"status": `"failed"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_ecs_snapshot_group.default.id}"]`,
			"instance_id":         `"${alicloud_ecs_snapshot_group.default.instance_id}"`,
			"name_regex":          `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}"`,
			"snapshot_group_name": `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}"`,
			"status":              `"accomplished"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_ecs_snapshot_group.default.id}_fake"]`,
			"instance_id":         `"${alicloud_ecs_snapshot_group.default.instance_id}"`,
			"name_regex":          `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}_fake"`,
			"snapshot_group_name": `"${alicloud_ecs_snapshot_group.default.snapshot_group_name}_fake"`,
			"status":              `"failed"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	var existAlicloudEcsSnapshotGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"groups.#":                     "1",
			"groups.0.id":                  CHECKSET,
			"groups.0.status":              "accomplished",
			"groups.0.description":         fmt.Sprintf("tf-testAccSnapshotGroup-%d", rand),
			"groups.0.instance_id":         CHECKSET,
			"groups.0.resource_group_id":   CHECKSET,
			"groups.0.snapshot_group_name": fmt.Sprintf("tf-testAccSnapshotGroup-%d", rand),
			"groups.0.tags.%":              "2",
			"groups.0.tags.Created":        "TF",
			"groups.0.tags.For":            "Acceptance-test",
		}
	}
	var fakeAlicloudEcsSnapshotGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcsSnapshotGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_snapshot_groups.default",
		existMapFunc: existAlicloudEcsSnapshotGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsSnapshotGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsSnapshotGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, instanceIdConf, snapshotGroupNameConf, tagsConf, nameRegexConf, statusConf, allConf)
}

func testAccCheckAlicloudEcsSnapshotGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSnapshotGroup-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" default {
  available_resource_creation = "Instance"
  available_disk_category     = "cloud_essd"
}

data "alicloud_instance_types" "default" {
  availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count       = 2
  memory_size          = 4
  system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name        = "${var.name}"
  description = "New security group"
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_disk" "default" {
  count     = "2"
  disk_name = "${var.name}"
  zone_id   = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category  = "cloud_essd"
  size      = "20"
}

data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name     = "${var.name}"
  host_name         = "tf-testAcc"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_disk_attachment" "default" {
  count       = "2"
  disk_id     = "${element(alicloud_disk.default.*.id,count.index)}"
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot_group" "default" {
	description = var.name
	disk_id = [alicloud_disk_attachment.default.0.disk_id,alicloud_disk_attachment.default.1.disk_id]
	snapshot_group_name = var.name
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	instance_id = alicloud_disk_attachment.default.0.instance_id
	instant_access = true
	instant_access_retention_days = 1
	tags = {
		Created = "TF"
		For 	= "Acceptance-test"
	}
}

data "alicloud_ecs_snapshot_groups" "default" {	
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

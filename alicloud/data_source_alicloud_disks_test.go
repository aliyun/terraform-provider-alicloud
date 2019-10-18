package alicloud

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"
)

func TestAccAlicloudDisksDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_disk.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_disk.default.id}_fake" ]`,
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"type":       `"data"`,
		}),
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"type":       `"system"`,
		}),
	}

	categoryConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"category":   `"cloud_efficiency"`,
		}),
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"category":   `"cloud"`,
		}),
	}

	encryptedConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"encrypted":  `"off"`,
		}),
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"encrypted":  `"on"`,
		}),
	}

	tagsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"tags":       `"${alicloud_disk.default.tags}"`,
		}),
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_disk.default.name}"`,
			"tags": `{
                           Name = "TerraformTest_fake"
                           Name1 = "TerraformTest_fake"
                        }`,
		}),
	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${alicloud_disk_attachment.default.instance_id}"`,
			"name_regex":  `"${alicloud_disk.default.name}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${alicloud_disk_attachment.default.instance_id}_fake"`,
			"name_regex":  `"${alicloud_disk.default.name}"`,
		}),
	}

	resourceGroupIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id":       `"${alicloud_disk_attachment.default.instance_id}"`,
			"resource_group_id": `"${var.resource_group_id}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id":       `"${alicloud_disk_attachment.default.instance_id}_fake"`,
			"resource_group_id": `"${var.resource_group_id}_fake"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":               `[ "${alicloud_disk.default.id}" ]`,
			"name_regex":        `"${alicloud_disk.default.name}"`,
			"type":              `"data"`,
			"category":          `"cloud_efficiency"`,
			"encrypted":         `"off"`,
			"tags":              `"${alicloud_disk.default.tags}"`,
			"instance_id":       `"${alicloud_disk_attachment.default.instance_id}"`,
			"resource_group_id": `"${var.resource_group_id}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":               `[ "${alicloud_disk.default.id}" ]`,
			"name_regex":        `"${alicloud_disk.default.name}"`,
			"type":              `"data"`,
			"category":          `"cloud_efficiency"`,
			"encrypted":         `"on"`,
			"tags":              `"${alicloud_disk.default.tags}"`,
			"instance_id":       `"${alicloud_disk_attachment.default.instance_id}"`,
			"resource_group_id": `"${var.resource_group_id}"`,
		}),
	}

	disksCheckInfo.dataSourceTestCheck(t, rand, idsConfig, nameRegexConfig, typeConfig, categoryConfig, encryptedConfig,
		tagsConfig, instanceIdConfig, resourceGroupIdConfig, allConfig)
}

func testAccCheckAlicloudDisksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "resource_group_id" {
	default = "%s"
}

variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_disk" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	description = "${var.name}_description"
	size = "20"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
	resource_group_id = "${var.resource_group_id}"
}

data "alicloud_disks" "default" {
	%s
}
	`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), rand, strings.Join(pairs, "\n	"))
	return config
}

func testAccCheckAlicloudDisksDataSourceConfigWithCommon(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "resource_group_id" {
	default = "%s"
}

variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

resource "alicloud_disk" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	description = "${var.name}_description"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
	size = "20"
	resource_group_id = "${var.resource_group_id}"
}

resource "alicloud_instance" "default" {
	vswitch_id = "${alicloud_vswitch.default.id}"
	private_ip = "172.16.0.10"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.default.id}"]
}

resource "alicloud_disk_attachment" "default" {
	disk_id = "${alicloud_disk.default.id}"
	instance_id = "${alicloud_instance.default.id}"
}

data "alicloud_disks" "default" {
	%s
}
`, EcsInstanceCommonTestCase, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), rand, strings.Join(pairs, "\n	"))
	return config
}

var existDisksMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"disks.#":                   "1",
		"disks.0.id":                CHECKSET,
		"disks.0.name":              fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand),
		"disks.0.description":       fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand),
		"disks.0.region_id":         CHECKSET,
		"disks.0.availability_zone": CHECKSET,
		"disks.0.status":            "Available",
		"disks.0.type":              "data",
		"disks.0.category":          "cloud_efficiency",
		"disks.0.encrypted":         "off",
		"disks.0.size":              "20",
		"disks.0.image_id":          "",
		"disks.0.snapshot_id":       "",
		"disks.0.instance_id":       "",
		"disks.0.creation_time":     CHECKSET,
		"disks.0.attached_time":     "",
		"disks.0.detached_time":     "",
		"disks.0.tags.%":            "2",
		"disks.0.resource_group_id": CHECKSET,
	}
}

var fakeDisksMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"disks.#": "0",
	}
}

var disksCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_disks.default",
	existMapFunc: existDisksMapFunc,
	fakeMapFunc:  fakeDisksMapFunc,
}

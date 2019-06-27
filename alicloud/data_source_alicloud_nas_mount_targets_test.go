package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudNasMountTargetDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	fileSystemIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
		}),
	}
	accessGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":    `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":    `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"type":           `"${alicloud_nas_access_group.default.type}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"type":           `"${alicloud_nas_access_group.default.type}_fake"`,
		}),
	}
	mountTargetDomainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"mount_target_domain": `"${alicloud_nas_mount_target.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"mount_target_domain": `"${alicloud_nas_mount_target.default.id}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${alicloud_vpc.default.id}_fake"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vswitch_id":     `"${alicloud_nas_mount_target.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vswitch_id":     `"${alicloud_nas_mount_target.default.vswitch_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alicloud_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alicloud_nas_mount_target.default.vswitch_id}"`,
			"type":                `"${alicloud_nas_access_group.default.type}"`,
			"vpc_id":              `"${alicloud_vpc.default.id}"`,
			"mount_target_domain": `"${alicloud_nas_mount_target.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alicloud_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alicloud_nas_mount_target.default.vswitch_id}_fake"`,
			"type":                `"${alicloud_nas_access_group.default.type}_fake"`,
			"vpc_id":              `"${alicloud_vpc.default.id}_fake"`,
			"mount_target_domain": `"${alicloud_nas_mount_target.default.id}_fake"`,
		}),
	}
	mountTargetCheckInfo.dataSourceTestCheck(t, rand, fileSystemIdConf, accessGroupNameConf, typeConf, mountTargetDomainConf, vpcIdConf, vswitchIdConf, allConf)
}

func testAccCheckAlicloudMountTargetDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
			default = "tf-testAccVswitch"
}
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_zones" "default" {
			available_resource_creation= "VSwitch"
}
resource "alicloud_vpc" "default" {
			name = "${var.name}"
			cidr_block = "172.16.0.0/12"
}
resource "alicloud_vswitch" "default" {
			vpc_id = "${alicloud_vpc.default.id}"
			cidr_block = "172.16.0.0/24"
			availability_zone = "${data.alicloud_zones.default.zones.0.id}"
			name = "${var.name}-1"
}
variable "storage_type" {
  default = "Capacity"
}
data "alicloud_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alicloud_nas_file_system" "default" {
  description = "${var.description}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
resource "alicloud_nas_access_group" "default" {
			name = "tf-testAccNasConfig-%d"
			type = "Vpc"
			description = "tf-testAccNasConfig"
}
resource "alicloud_nas_mount_target" "default" {
			file_system_id = "${alicloud_nas_file_system.default.id}"
			access_group_name = "${alicloud_nas_access_group.default.id}"
			vswitch_id = "${alicloud_vswitch.default.id}"               
}
data "alicloud_nas_mount_targets" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"targets.0.type":                "Vpc",
		"targets.0.vpc_id":              CHECKSET,
		"targets.0.mount_target_domain": CHECKSET,
		"targets.0.vswitch_id":          CHECKSET,
		"targets.0.access_group_name":   fmt.Sprintf("tf-testAccNasConfig-%d", rand),
		"ids.#":                         "1",
		"ids.0":                         CHECKSET,
	}
}

var fakeMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"targets.#": "0",
		"ids.#":     "0",
	}
}

var mountTargetCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_mount_targets.default",
	existMapFunc: existMountTargetMapCheck,
	fakeMapFunc:  fakeMountTargetMapCheck,
}

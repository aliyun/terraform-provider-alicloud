package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASMountTargetDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	fileSystemIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
		}),
	}
	accessGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":    `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":    `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}_fake"`,
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
	netWorkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"network_type":   `"${alicloud_nas_access_group.default.access_group_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"network_type":   `"${alicloud_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	mountTargetDomainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"mount_target_domain": `split(":",alicloud_nas_mount_target.default.id)[1]`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"mount_target_domain": `"fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${data.alicloud_vpcs.default.vpcs.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${data.alicloud_vpcs.default.vpcs.0.id}_fake"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vswitch_id":     `"${alicloud_nas_mount_target.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"vswitch_id":     `"fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"ids":            `[split(":",alicloud_nas_mount_target.default.id)[1]]`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"ids":            `["fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"status":         `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"status":         `"Inactive"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alicloud_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alicloud_nas_mount_target.default.vswitch_id}"`,
			"type":                `"${alicloud_nas_access_group.default.type}"`,
			"network_type":        `"${alicloud_nas_access_group.default.access_group_type}"`,
			"vpc_id":              `"${data.alicloud_vpcs.default.vpcs.0.id}"`,
			"mount_target_domain": `split(":",alicloud_nas_mount_target.default.id)[1]`,
			"status":              `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alicloud_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alicloud_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alicloud_nas_mount_target.default.vswitch_id}_fake"`,
			"type":                `"${alicloud_nas_access_group.default.type}_fake"`,
			"network_type":        `"${alicloud_nas_access_group.default.access_group_type}_fake}"`,
			"vpc_id":              `"${data.alicloud_vpcs.default.vpcs.0.id}_fake"`,
			"mount_target_domain": `"fake"`,
			"status":              `"Inactive"`,
		}),
	}
	mountTargetCheckInfo.dataSourceTestCheck(t, rand, fileSystemIdConf, accessGroupNameConf, typeConf, netWorkTypeConf, mountTargetDomainConf, vpcIdConf, vswitchIdConf, idsConf, statusConf, allConf)
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
data "alicloud_vpcs" "default" {
			name_regex = "default-NODELETING"
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
			access_group_name = "tf-testAccNasConfig-%d"
			access_group_type = "Vpc"
			description = "tf-testAccNasConfig"
}
resource "alicloud_nas_mount_target" "default" {
			file_system_id = "${alicloud_nas_file_system.default.id}"
			access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
}
data "alicloud_nas_mount_targets" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"targets.0.type":                "Vpc",
		"targets.0.network_type":        "Vpc",
		"targets.0.status":              "Active",
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

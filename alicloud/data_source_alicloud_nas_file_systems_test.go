package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASFileSystem_DataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	storageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alicloud_nas_file_system.default.storage_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	protocolTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"protocol_type":     `"${alicloud_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_file_system.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alicloud_nas_file_system.default.storage_type}"`,
			"protocol_type":     `"${alicloud_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
			"ids":               `["${alicloud_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}_fake"`,
			"ids":               `["${alicloud_nas_file_system.default.id}_fake"]`,
		}),
	}

	fileSystemCheckInfo.dataSourceTestCheck(t, rand, storageTypeConf, protocolTypeConf,
		descriptionConf, idsConf, allConf)
}

func testAccCheckAlicloudFileSystemDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
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
data "alicloud_nas_file_systems" "default" {
	%s
}`, strings.Join(pairs, "\n  "))
	return config
}

var existFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.0.id":            CHECKSET,
		"systems.0.region_id":     CHECKSET,
		"systems.0.description":   "tf-testAccCheckAlicloudFileSystemsDataSource",
		"systems.0.protocol_type": CHECKSET,
		"systems.0.storage_type":  "Capacity",
		"systems.0.metered_size":  CHECKSET,
		"systems.0.create_time":   CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"descriptions.#":          "1",
		"descriptions.0":          "tf-testAccCheckAlicloudFileSystemsDataSource",
	}
}

var fakeFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.#":      "0",
		"ids.#":          "0",
		"descriptions.#": "0",
	}
}

var fileSystemCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_file_systems.default",
	existMapFunc: existFileSystemMapCheck,
	fakeMapFunc:  fakeFileSystemMapCheck,
}

package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASFileSystem_DataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	storageTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alicloud_nas_file_system.default.storage_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	protocolTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"protocol_type":     `"${alicloud_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_file_system.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_file_system.default.id}_fake"]`,
		}),
	}
	fileSystemTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"file_system_type":  `"standard"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"file_system_type":  `"extreme"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alicloud_nas_file_system.default.storage_type}"`,
			"protocol_type":     `"${alicloud_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
			"ids":               `["${alicloud_nas_file_system.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}_fake"`,
			"ids":               `["${alicloud_nas_file_system.default.id}_fake"]`,
		}),
	}

	fileSystemCheckInfo.DataSourceTestCheck(t, rand, storageTypeConf, protocolTypeConf,
		descriptionConf, idsConf, fileSystemTypeConf, allConf)
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
		"systems.0.id":            acctest.CHECKSET,
		"systems.0.region_id":     acctest.CHECKSET,
		"systems.0.description":   "tf-testAccCheckAlicloudFileSystemsDataSource",
		"systems.0.protocol_type": acctest.CHECKSET,
		"systems.0.storage_type":  "Capacity",
		"systems.0.metered_size":  acctest.CHECKSET,
		"systems.0.create_time":   acctest.CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   acctest.CHECKSET,
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

var fileSystemCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_nas_file_systems.default",
	ExistMapFunc: existFileSystemMapCheck,
	FakeMapFunc:  fakeFileSystemMapCheck,
}

package dfs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDfsFileSystemsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsFileSystemsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_file_system.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDfsFileSystemsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_file_system.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dfs_file_system.default.file_system_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDfsFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dfs_file_system.default.file_system_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsFileSystemsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dfs_file_system.default.id}"]`,
			"name_regex": `"${alicloud_dfs_file_system.default.file_system_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDfsFileSystemsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dfs_file_system.default.id}_fake"]`,
			"name_regex": `"${alicloud_dfs_file_system.default.file_system_name}_fake"`,
		}),
	}
	var existAlicloudDfsFileSystemsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"systems.#":                  "1",
			"systems.0.description":      fmt.Sprintf("tf-testAccFileSystem-%d", rand),
			"systems.0.file_system_name": fmt.Sprintf("tf-testAccFileSystem-%d", rand),
			"systems.0.protocol_type":    "HDFS",
			"systems.0.storage_type":     "STANDARD",
			"systems.0.throughput_mode":  "Standard",
			"systems.0.space_capacity":   "1024",
		}
	}
	var fakeAlicloudDfsFileSystemsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDfsFileSystemsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dfs_file_systems.default",
		ExistMapFunc: existAlicloudDfsFileSystemsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDfsFileSystemsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudDfsFileSystemsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudDfsFileSystemsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccFileSystem-%d"
}

data "alicloud_dfs_zones" "default" {}

resource "alicloud_dfs_file_system" "default" {
  storage_type     = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
  zone_id          = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type    = "HDFS"
  description      = var.name
  file_system_name = var.name
  throughput_mode  = "Standard"
  space_capacity   = "1024"
}

data "alicloud_dfs_file_systems" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

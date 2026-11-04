package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASFilesetsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_fileset.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_fileset.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}"]`,
			"status": `"CREATED"`,
		}),
		FakeConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}_fake"]`,
			"status": `"CREATING"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}"]`,
			"status": `"CREATED"`,
		}),
		FakeConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}_fake"]`,
			"status": `"CREATING"`,
		}),
	}
	var existAlicloudNasFilesetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"filesets.#":                  "1",
			"filesets.0.create_time":      acctest.CHECKSET,
			"filesets.0.description":      fmt.Sprintf("tf-testAccFileset-%d", rand),
			"filesets.0.file_system_id":   acctest.CHECKSET,
			"filesets.0.file_system_path": "/tf-testAcc-Path/",
			"filesets.0.id":               acctest.CHECKSET,
			"filesets.0.fileset_id":       acctest.CHECKSET,
			"filesets.0.status":           acctest.CHECKSET,
			"filesets.0.update_time":      acctest.CHECKSET,
		}
	}
	var fakeAlicloudNasFilesetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"filesets.#": "0",
		}
	}
	var alicloudNasFilesetsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nas_filesets.default",
		ExistMapFunc: existAlicloudNasFilesetsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudNasFilesetsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudNasFilesetsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudNasFilesetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccFileset-%d"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_nas_fileset" "default" {
  file_system_id   = alicloud_nas_file_system.default.id
  description      = var.name
  file_system_path = "/tf-testAcc-Path/"
}

data "alicloud_nas_filesets" "default" {	
	file_system_id = alicloud_nas_file_system.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

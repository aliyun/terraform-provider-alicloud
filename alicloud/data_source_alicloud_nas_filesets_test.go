package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASFilesetsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_fileset.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_fileset.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}"]`,
			"status": `"CREATED"`,
		}),
		fakeConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}_fake"]`,
			"status": `"CREATING"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}"]`,
			"status": `"CREATED"`,
		}),
		fakeConfig: testAccCheckAlicloudNasFilesetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_fileset.default.id}_fake"]`,
			"status": `"CREATING"`,
		}),
	}
	var existAlicloudNasFilesetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"filesets.#":                  "1",
			"filesets.0.create_time":      CHECKSET,
			"filesets.0.description":      fmt.Sprintf("tf-testAccFileset-%d", rand),
			"filesets.0.file_system_id":   CHECKSET,
			"filesets.0.file_system_path": "/tf-testAcc-Path/",
			"filesets.0.id":               CHECKSET,
			"filesets.0.fileset_id":       CHECKSET,
			"filesets.0.status":           CHECKSET,
			"filesets.0.update_time":      CHECKSET,
		}
	}
	var fakeAlicloudNasFilesetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"filesets.#": "0",
		}
	}
	var alicloudNasFilesetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nas_filesets.default",
		existMapFunc: existAlicloudNasFilesetsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNasFilesetsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudNasFilesetsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
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
  name_regex = "default-NODELETING"
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

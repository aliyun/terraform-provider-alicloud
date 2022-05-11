package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASDataFlowsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_data_flow.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_data_flow.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_data_flow.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_data_flow.default.id}"]`,
			"status": `"Stopped"`,
		}),
	}
	var existAlicloudNasDataFlowsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"flows.#":                      "1",
			"flows.0.create_time":          CHECKSET,
			"flows.0.id":                   CHECKSET,
			"flows.0.data_flow_id":         CHECKSET,
			"flows.0.error_message":        "",
			"flows.0.description":          fmt.Sprintf("tf-testaccdataflow-%d", rand),
			"flows.0.file_system_id":       CHECKSET,
			"flows.0.file_system_path":     CHECKSET,
			"flows.0.fset_description":     CHECKSET,
			"flows.0.fset_id":              CHECKSET,
			"flows.0.source_security_type": `SSL`,
			"flows.0.source_storage":       CHECKSET,
			"flows.0.status":               "Running",
			"flows.0.throughput":           "600",
		}
	}
	var fakeAlicloudNasDataFlowsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudNasDataFlowsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nas_data_flows.default",
		existMapFunc: existAlicloudNasDataFlowsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNasDataFlowsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudNasDataFlowsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
}
func testAccCheckAlicloudNasDataFlowsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccdataflow-%d"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "private"
  tags   = {
    cpfs-dataflow = "true"
  }
}

data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = local.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_nas_mount_target" "default" {
  file_system_id = "${alicloud_nas_file_system.default.id}"
  vswitch_id     = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_nas_fileset" "default" {
  depends_on       = ["alicloud_nas_mount_target.default"]
  file_system_id   = alicloud_nas_file_system.default.id
  description      = var.name
  file_system_path = "/tf-testAcc-Path/"
}

resource "alicloud_nas_data_flow" "default" {
  fset_id              = alicloud_nas_fileset.default.fileset_id
  description          = var.name
  file_system_id       = alicloud_nas_file_system.default.id
  source_security_type = "SSL"
  source_storage       = join("", ["oss://", alicloud_oss_bucket.default.bucket])
  throughput           = 600
}

data "alicloud_nas_data_flows" "default" {	
  file_system_id = alicloud_nas_file_system.default.id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

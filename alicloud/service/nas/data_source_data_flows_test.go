package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASDataFlowsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_data_flow.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_data_flow.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_data_flow.default.id}"]`,
			"status": `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudNasDataFlowsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_data_flow.default.id}"]`,
			"status": `"Stopped"`,
		}),
	}
	var existAlicloudNasDataFlowsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"flows.#":                      "1",
			"flows.0.create_time":          acctest.CHECKSET,
			"flows.0.id":                   acctest.CHECKSET,
			"flows.0.data_flow_id":         acctest.CHECKSET,
			"flows.0.error_message":        "",
			"flows.0.description":          fmt.Sprintf("tf-testaccdataflow-%d", rand),
			"flows.0.file_system_id":       acctest.CHECKSET,
			"flows.0.file_system_path":     acctest.CHECKSET,
			"flows.0.fset_description":     acctest.CHECKSET,
			"flows.0.fset_id":              acctest.CHECKSET,
			"flows.0.source_security_type": `SSL`,
			"flows.0.source_storage":       acctest.CHECKSET,
			"flows.0.status":               "Running",
			"flows.0.throughput":           "600",
		}
	}
	var fakeAlicloudNasDataFlowsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudNasDataFlowsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nas_data_flows.default",
		ExistMapFunc: existAlicloudNasDataFlowsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudNasDataFlowsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudNasDataFlowsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
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
  name_regex = "^default-NODELETING$"
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

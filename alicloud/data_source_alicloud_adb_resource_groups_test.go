package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAdbResourceGroupDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.ADBResourceGroupSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}

	GroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"group_name":    `"${var.name}"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"group_name":    `"${var.name}_fake"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}
	DbClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"group_name":    `"${var.name}"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"group_name":    `"${var.name}_fake"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}

	AdbResourceGroupCheckInfo.dataSourceTestCheck(t, rand, idsConf, GroupNameConf, DbClusterIdConf, allConf)
}

var existAdbResourceGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":               "1",
		"groups.0.id":            CHECKSET,
		"groups.0.create_time":   CHECKSET,
		"groups.0.db_cluster_id": CHECKSET,
		"groups.0.group_name":    CHECKSET,
		"groups.0.group_type":    CHECKSET,
		"groups.0.node_num":      CHECKSET,
		"groups.0.user":          CHECKSET,
	}
}

var fakeAdbResourceGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AdbResourceGroupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_adb_resource_groups.default",
	existMapFunc: existAdbResourceGroupMapFunc,
	fakeMapFunc:  fakeAdbResourceGroupMapFunc,
}

func testAccCheckAlicloudAdbResourceGroupSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "TF_TESTACCADBRG%d"
}

variable "creation" {
  default = "ADB"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_adb_zones" "zones_ids" {}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_adb_zones.zones_ids.zones.0.id
  vswitch_name = var.name
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "32Core128GB"
  payment_type        = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.vswitch.id
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  security_ips      = ["10.168.1.12", "10.168.1.11"]
}

resource "alicloud_adb_resource_group" "default" {
  group_name    = "${var.name}"
  group_type    = "batch"
  db_cluster_id = "${alicloud_adb_db_cluster.default.id}"
}

data "alicloud_adb_resource_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

// Package adb_test.
package adb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudAdbResourceGroupDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ADBResourceGroupSupportRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}

	GroupNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"group_name":    `"${var.name}"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"group_name":    `"${var.name}_fake"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}
	DbClusterIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}"]`,
			"group_name":    `"${var.name}"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbResourceGroupSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_adb_resource_group.default.id}_fake"]`,
			"group_name":    `"${var.name}_fake"`,
			"db_cluster_id": `"${alicloud_adb_db_cluster.default.id}"`,
		}),
	}

	AdbResourceGroupCheckInfo.DataSourceTestCheck(t, rand, idsConf, GroupNameConf, DbClusterIdConf, allConf)
}

var existAdbResourceGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":               "1",
		"groups.0.id":            acctest.CHECKSET,
		"groups.0.create_time":   acctest.CHECKSET,
		"groups.0.db_cluster_id": acctest.CHECKSET,
		"groups.0.group_name":    acctest.CHECKSET,
		"groups.0.group_type":    acctest.CHECKSET,
		"groups.0.node_num":      acctest.CHECKSET,
		"groups.0.user":          acctest.CHECKSET,
	}
}

var fakeAdbResourceGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AdbResourceGroupCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_adb_resource_groups.default",
	ExistMapFunc: existAdbResourceGroupMapFunc,
	FakeMapFunc:  fakeAdbResourceGroupMapFunc,
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

data "alicloud_adb_zones" "default" {}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "32Core128GB"
  payment_type        = "PayAsYouGo"
  vswitch_id          = data.alicloud_vswitches.default.ids.0
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

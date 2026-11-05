package mse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudMSEZnodesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_znode.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_znode.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_znode.default.path}"`,
		}),
		FakeConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_znode.default.path}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mse_znode.default.id}"]`,
			"name_regex": `"${alicloud_mse_znode.default.path}"`,
		}),
		FakeConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mse_znode.default.id}_fake"]`,
			"name_regex": `"${alicloud_mse_znode.default.path}_fake"`,
		}),
	}
	var existAlicloudMseZnodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"znodes.#":            "1",
			"znodes.0.cluster_id": acctest.CHECKSET,
			// znodes.0.data is intentionally not asserted: the backing
			// ListZnodeChildren API only returns child metadata (path/name/dir),
			// not the znode payload, so this field is always empty in the data
			// source's list view.
			"znodes.0.path":       fmt.Sprintf("/tf-testAccZnode-%d", rand),
			"znodes.0.znode_name": acctest.CHECKSET,
			"znodes.0.dir":        acctest.CHECKSET,
			"znodes.0.id":         acctest.CHECKSET,
		}
	}
	var fakeAlicloudMseZnodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudMseZnodesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mse_znodes.default",
		ExistMapFunc: existAlicloudMseZnodesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMseZnodesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudMseZnodesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudMseZnodesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccZnode-%d"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_mse_cluster" "default" {
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type          = "ZooKeeper"
  cluster_version       = "ZooKeeper_3_8_0"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = alicloud_vswitch.default.id
  pub_network_flow      = "1"
  cluster_alias_name    = var.name
  mse_version = "mse_dev"
}
resource "alicloud_mse_znode" "default" {
  cluster_id = alicloud_mse_cluster.default.cluster_id
  data       = "/${var.name}"
  path       = "/${var.name}"
}
data "alicloud_mse_znodes" "default" {
  cluster_id = alicloud_mse_cluster.default.cluster_id
  path = "/"
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

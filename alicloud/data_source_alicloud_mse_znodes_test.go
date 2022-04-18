package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMSEZnodesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_znode.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_znode.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_znode.default.path}"`,
		}),
		fakeConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_znode.default.path}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mse_znode.default.id}"]`,
			"name_regex": `"${alicloud_mse_znode.default.path}"`,
		}),
		fakeConfig: testAccCheckAlicloudMseZnodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mse_znode.default.id}_fake"]`,
			"name_regex": `"${alicloud_mse_znode.default.path}_fake"`,
		}),
	}
	var existAlicloudMseZnodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"znodes.#":            "1",
			"znodes.0.cluster_id": CHECKSET,
			"znodes.0.data":       fmt.Sprintf("/tf-testAccZnode-%d", rand),
			"znodes.0.path":       fmt.Sprintf("/tf-testAccZnode-%d", rand),
			"znodes.0.znode_name": CHECKSET,
			"znodes.0.dir":        CHECKSET,
			"znodes.0.id":         CHECKSET,
		}
	}
	var fakeAlicloudMseZnodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudMseZnodesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mse_znodes.default",
		existMapFunc: existAlicloudMseZnodesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMseZnodesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMseZnodesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudMseZnodesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "/tf-testAccZnode-%d"
}
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
resource "alicloud_mse_cluster" "default" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type          = "ZooKeeper"
  cluster_version       = "ZooKeeper_3_5_5"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = var.name
}
resource "alicloud_mse_znode" "default" {
  cluster_id = alicloud_mse_cluster.default.cluster_id
  data       = var.name
  path       = var.name
}
data "alicloud_mse_znodes" "default" {
  cluster_id = alicloud_mse_cluster.default.cluster_id
  path = "/"
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

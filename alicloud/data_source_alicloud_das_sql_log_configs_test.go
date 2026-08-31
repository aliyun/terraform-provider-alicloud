package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDasSqlLogConfigDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDasSqlLogConfigDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_das_sql_log_config.default.instance_id}"`,
		}),
	}

	DasSqlLogConfigCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, func() {
		testAccPreCheck(t)
	}, instanceIdConf)
}

var existDasSqlLogConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#":               "1",
		"configs.0.instance_id":   CHECKSET,
		"configs.0.retention":     "30",
		"configs.0.hot_retention": "7",
		"ids.#":                   "1",
	}
}

var fakeDasSqlLogConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#": "1",
		"ids.#":     "1",
	}
}

var DasSqlLogConfigCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_das_sql_log_configs.default",
	existMapFunc: existDasSqlLogConfigMapFunc,
	fakeMapFunc:  fakeDasSqlLogConfigMapFunc,
}

func testAccCheckAliCloudDasSqlLogConfigDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDasSqlLogConfig%d"
}

data "alicloud_polardb_zones" "default" {}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  zone_id    = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  zone_id      = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
  vpc_id       = alicloud_vpc.default.id
  vswitch_name = var.name
  cidr_block   = "172.16.0.0/24"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_das_sql_log_config" "default" {
  instance_id    = alicloud_polardb_cluster.default.id
  enable         = true
  request_enable = true
  retention      = 30
  hot_retention  = 7
}

data "alicloud_das_sql_log_configs" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

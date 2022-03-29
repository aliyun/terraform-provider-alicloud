package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudClickHouseAccountDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_click_house_account.default.id}"]`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_click_house_account.default.id}_fake"]`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"name_regex":    `"${alicloud_click_house_account.default.account_name}"`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"name_regex":    `"${alicloud_click_house_account.default.account_name}_fake"`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"status":        `"${alicloud_click_house_account.default.status}"`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
			"ids":           `["${alicloud_click_house_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"status":        `"Deleting"`,
			"ids":           `["${alicloud_click_house_account.default.id}"]`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_click_house_account.default.id}"]`,
			"name_regex":    `"${alicloud_click_house_account.default.account_name}"`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
			"status":        `"${alicloud_click_house_account.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudClickHouseAccountDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_click_house_account.default.id}_fake"]`,
			"name_regex":    `"${alicloud_click_house_account.default.account_name}_fake"`,
			"db_cluster_id": `"${alicloud_click_house_db_cluster.default.id}"`,
			"status":        `"Deleting"`,
		}),
	}
	var existAlicloudClickHouseAccountDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"accounts.#":              "1",
			"accounts.0.account_name": fmt.Sprintf("tf_testacc%d", rand),
		}
	}
	var fakeAlicloudClickHouseAccountDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"accounts.#": "0",
		}
	}
	var alicloudClickHouseAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_click_house_accounts.default",
		existMapFunc: existAlicloudClickHouseAccountDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudClickHouseAccountDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
	}
	alicloudClickHouseAccountBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudClickHouseAccountDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf_testacc%d"
}
variable "pwd" {
  default = "Tf-test%d"
}

data "alicloud_click_house_regions" "default" {
  current = true
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = 1
  payment_type            = "PayAsYouGo"
  db_node_storage         = "100"
  storage_type            = "cloud_essd"
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}

resource "alicloud_click_house_account" "default" {
  db_cluster_id    = alicloud_click_house_db_cluster.default.id
  account_name     = var.name
  account_password = var.pwd
}
data "alicloud_click_house_accounts" "default" {	
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}

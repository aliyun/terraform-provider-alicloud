package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGPDBAccountsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_account.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_account.default.account_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_account.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_account.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_gpdb_account.default.id}"]`,
			"name_regex": `"${alicloud_gpdb_account.default.account_name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_gpdb_account.default.id}_fake"]`,
			"name_regex": `"${alicloud_gpdb_account.default.account_name}_fake"`,
			"status":     `"Creating"`,
		}),
	}
	var existAlicloudGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"accounts.#":                     "1",
			"accounts.0.id":                  CHECKSET,
			"accounts.0.account_name":        fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.account_description": fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.db_instance_id":      CHECKSET,
			"accounts.0.status":              "Active",
		}
	}
	var fakeAlicloudGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudGpdbAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_gpdb_accounts.default",
		existMapFunc: existAlicloudGpdbAccountsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGpdbAccountsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.GpdbElasticInstanceSupportRegions)
	}
	alicloudGpdbAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudGpdbAccountsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tftestacc%d"
}
data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.zones.3.id
  vswitch_name = var.name
}

resource "alicloud_gpdb_elastic_instance" "default" {
  engine                   = "gpdb"
  engine_version           = "6.0"
  seg_storage_type         = "cloud_essd"
  seg_node_num             = 4
  storage_size             = 50
  instance_spec            = "2C16G"
  db_instance_description  = var.name
  instance_network_type    = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id               = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}

resource "alicloud_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = alicloud_gpdb_elastic_instance.default.id
  account_password    = "TFTest123"
  account_description = var.name
}

data "alicloud_gpdb_accounts" "default" {	
	db_instance_id = alicloud_gpdb_elastic_instance.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

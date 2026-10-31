package gpdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudGPDBAccountsDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandIntRange(1000, 9999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_account.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_account.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_account.default.account_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_account.default.account_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_account.default.id}"]`,
			"status": `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_account.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_gpdb_account.default.id}"]`,
			"name_regex": `"${alicloud_gpdb_account.default.account_name}"`,
			"status":     `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
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
			"accounts.0.id":                  acctest.CHECKSET,
			"accounts.0.account_name":        fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.account_description": fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.db_instance_id":      acctest.CHECKSET,
			"accounts.0.status":              "Active",
		}
	}
	var fakeAlicloudGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudGpdbAccountsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_gpdb_accounts.default",
		ExistMapFunc: existAlicloudGpdbAccountsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudGpdbAccountsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	}
	alicloudGpdbAccountsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
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

data "alicloud_gpdb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = alicloud_gpdb_instance.default.id
  account_password    = "TFTest123"
  account_description = var.name
}

data "alicloud_gpdb_accounts" "default" {	
	db_instance_id = alicloud_gpdb_account.default.db_instance_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

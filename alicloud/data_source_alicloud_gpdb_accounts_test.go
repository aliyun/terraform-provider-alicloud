package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGpdbAccountsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	checkoutSupportedRegions(t, true, connectivity.GPDBAccountSupportRegions)
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
			"status": `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_account.default.id}"]`,
			"status": `"0"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_gpdb_account.default.id}"]`,
			"name_regex": `"${alicloud_gpdb_account.default.account_name}"`,
			"status":     `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_gpdb_account.default.id}_fake"]`,
			"name_regex": `"${alicloud_gpdb_account.default.account_name}_fake"`,
			"status":     `"0"`,
		}),
	}
	var existAlicloudGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"accounts.#":                     "1",
			"accounts.0.account_name":        fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.account_description": fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.db_instance_id":      CHECKSET,
			"accounts.0.status":              "1",
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

resource "alicloud_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = "gp-bp16dlw51zj999483"
  account_password    = "TFTest123"
  account_description = var.name
}

data "alicloud_gpdb_accounts" "default" {	
	db_instance_id = "gp-bp16dlw51zj999483" 
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

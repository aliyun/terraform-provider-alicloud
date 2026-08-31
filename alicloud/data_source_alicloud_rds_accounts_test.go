package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsAccountsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}_fake"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_rds_account.default.account_name}"`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_rds_account.default.account_name}_fake"`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"status":         `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}_fake"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"status":         `"Unavailable"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"name_regex":     `"${alicloud_rds_account.default.account_name}"`,
			"status":         `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}_fake"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"name_regex":     `"${alicloud_rds_account.default.account_name}_fake"`,
			"status":         `"Unavailable"`,
		}),
	}
	var existAlicloudRdsAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"accounts.#":                       "1",
			"accounts.0.account_description":   ``,
			"accounts.0.account_name":          "tftestnormal000",
			"accounts.0.account_type":          `Normal`,
			"accounts.0.database_privileges.#": `0`,
			"accounts.0.priv_exceeded":         "0",
			"accounts.0.status":                `Available`,
		}
	}
	var fakeAlicloudRdsAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudRdsAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_accounts.default",
		existMapFunc: existAlicloudRdsAccountsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsAccountsDataSourceNameMapFunc,
	}
	alicloudRdsAccountsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudRdsAccountsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alicloud_db_instances" "default" { 
}
resource "alicloud_rds_account" "default" {
  db_instance_id = data.alicloud_db_instances.default.ids.0
  account_name        = "tftestnormal000"
  account_password    = "Test12345"
}
data "alicloud_rds_accounts" "default" {	
	%s
}`, strings.Join(pairs, "\n"))
}

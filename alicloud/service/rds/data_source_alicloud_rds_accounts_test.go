package rds_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsAccountsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}_fake"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_rds_account.default.account_name}"`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_rds_account.default.account_name}_fake"`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"status":         `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}_fake"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"status":         `"Unavailable"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_account.default.account_name}"]`,
			"db_instance_id": `"${alicloud_rds_account.default.db_instance_id}"`,
			"name_regex":     `"${alicloud_rds_account.default.account_name}"`,
			"status":         `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudRdsAccountsDataSourceName(rand, map[string]string{
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
	var alicloudRdsAccountsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_rds_accounts.default",
		ExistMapFunc: existAlicloudRdsAccountsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudRdsAccountsDataSourceNameMapFunc,
	}
	alicloudRdsAccountsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
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

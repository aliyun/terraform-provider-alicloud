package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcDedicatedHostAccountsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_account.default.id}_fake"]`,
		}),
	}
	dedicatedHostIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_cddc_dedicated_host_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_cddc_dedicated_host_account.default.account_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}"`,
			"name_regex":        `"${alicloud_cddc_dedicated_host_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}_fake"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}_fake"`,
			"name_regex":        `"${alicloud_cddc_dedicated_host_account.default.account_name}_fake"`,
		}),
	}
	var existAlicloudCddcDedicatedHostAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"accounts.#":                   "1",
			"accounts.0.id":                CHECKSET,
			"accounts.0.account_name":      name,
			"accounts.0.dedicated_host_id": CHECKSET,
		}
	}
	var fakeAlicloudCddcDedicatedHostAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCddcDedicatedHostAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cddc_dedicated_host_accounts.default",
		existMapFunc: existAlicloudCddcDedicatedHostAccountsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCddcDedicatedHostAccountsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
	}

	alicloudCddcDedicatedHostAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dedicatedHostIdConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "%s"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_cddc_dedicated_host_groups.default.groups.0.vpc_id
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "^NO-DELETING"
  engine     = "MySQL"
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = data.alicloud_cddc_dedicated_host_groups.default.ids.0
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  payment_type            = "Subscription"
}

resource "alicloud_cddc_dedicated_host_account" "default" {
  account_name      = var.name
  account_password  = "Test1234+!"
  dedicated_host_id = alicloud_cddc_dedicated_host.default.dedicated_host_id
  account_type      = "Normal"
}

data "alicloud_cddc_dedicated_host_accounts" "default" {	
	%s
}
`, name, strings.Join(pairs, " \n "))
	return config
}

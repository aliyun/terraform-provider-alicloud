package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcDedicatedHostAccountsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.CddcSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_account.default.id}_fake"]`,
		}),
	}
	dedicatedHostIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cddc_dedicated_host_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cddc_dedicated_host_account.default.account_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cddc_dedicated_host_account.default.id}"]`,
			"dedicated_host_id": `"${alicloud_cddc_dedicated_host_account.default.dedicated_host_id}"`,
			"name_regex":        `"${alicloud_cddc_dedicated_host_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand, map[string]string{
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
			"accounts.0.account_name":      fmt.Sprintf("tftestacc%d", rand),
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
		testAccPreCheckWithTime(t, []int{1})
	}
	alicloudCddcDedicatedHostAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dedicatedHostIdConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCddcDedicatedHostAccountsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tftestacc%d"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type        = "mssql"
  zone_id        = data.alicloud_cddc_zones.default.ids.0
  storage_type   = "cloud_essd"
  image_category = "WindowsWithMssqlStdLicense"

}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "default-NODELETING"
  engine     = "mssql"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  count                     = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
  engine                    = "SQLServer"
  vpc_id                    = data.alicloud_vpcs.default.ids.0
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = var.name
  open_permission           = true
}

data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.groups[0].vpc_id : data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_vswitch" "default" {
  count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = data.alicloud_vpcs.default.vpcs[0].cidr_block
  zone_id    = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : alicloud_cddc_dedicated_host_group.default[0].id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : alicloud_vswitch.default[0].id
  payment_type            = "Subscription"
  image_category          = "WindowsWithMssqlStdLicense"
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
`, rand, strings.Join(pairs, " \n "))
	return config
}

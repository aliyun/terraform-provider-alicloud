package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDSimpleOfficeSitesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_simple_office_site.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_simple_office_site.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_simple_office_site.default.id}"]`,
			"name_regex": `"${alicloud_ecd_simple_office_site.default.office_site_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_simple_office_site.default.id}"]`,
			"name_regex": `"${alicloud_ecd_simple_office_site.default.office_site_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_simple_office_site.default.id}"]`,
			"name_regex": `"${alicloud_ecd_simple_office_site.default.office_site_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_simple_office_site.default.id}_fake"]`,
			"name_regex": `"${alicloud_ecd_simple_office_site.default.office_site_name}_fake"`,
		}),
	}
	var existAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"sites.#":                             "1",
			"sites.0.cidr_block":                  "172.16.0.0/12",
			"sites.0.desktop_access_type":         "Internet",
			"sites.0.status":                      "REGISTERED",
			"sites.0.office_site_id":              CHECKSET,
			"sites.0.enable_cross_desktop_access": "false",
			// todo: need to check the `bandwidth` and `enable_internet_access` after fixing the issue occurred in ap-southeast-1
			//"sites.0.bandwidth": "10",
		}
	}
	var fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"sites.#": "0",
		}
	}
	var alicloudEventBridgeEventBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_simple_office_sites.default",
		existMapFunc: existAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}

	alicloudEventBridgeEventBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudSimpleOfficeSitesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccSimpleOfficeSite-%d"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
  #  bandwidth           = 10
  # enable_internet_access= true
  enable_internet_access = false
}

data "alicloud_ecd_simple_office_sites" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

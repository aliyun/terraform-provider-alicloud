package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCHavipsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_havip.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_havip.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_havip.default.havip_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_havip.default.havip_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_havip.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_havip.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_havip.default.id}"]`,
			"name_regex": `"${alicloud_havip.default.havip_name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudHavipsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_havip.default.id}_fake"]`,
			"name_regex": `"${alicloud_havip.default.havip_name}_fake"`,
			"status":     `"Pending"`,
		}),
	}
	var existAlicloudHavipsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"havips.#":                            "1",
			"havips.0.description":                fmt.Sprintf("tf-testacchavip-%d", rand),
			"havips.0.havip_name":                 fmt.Sprintf("tf-testacchavip-%d", rand),
			"havips.0.ip_address":                 CHECKSET,
			"havips.0.vswitch_id":                 CHECKSET,
			"havips.0.associated_eip_addresses.#": CHECKSET,
			"havips.0.associated_instances.#":     CHECKSET,
			"havips.0.id":                         CHECKSET,
			"havips.0.havip_id":                   CHECKSET,
			"havips.0.master_instance_id":         "",
			"havips.0.status":                     "Available",
			"havips.0.vpc_id":                     CHECKSET,
		}
	}
	var fakeAlicloudHavipsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudHavipsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_havips.default",
		existMapFunc: existAlicloudHavipsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudHavipsDataSourceNameMapFunc,
	}
	alicloudHavipsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudHavipsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacchavip-%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  zone_id = data.alicloud_zones.default.zones.0.id
  cidr_block = "192.168.0.0/21"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_havip" "default" {
  havip_name = var.name
  vswitch_id = alicloud_vswitch.default.id
  description = var.name
}

data "alicloud_havips" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPvtzEndpointsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_endpoint.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_endpoint.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_pvtz_endpoint.default.id}"]`,
			"status": `"SUCCESS"`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_pvtz_endpoint.default.id}"]`,
			"status": `"INIT"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_pvtz_endpoint.default.id}"]`,
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}"`,
			"status":     `"SUCCESS"`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_pvtz_endpoint.default.id}_fake"]`,
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}_fake"`,
			"status":     `"INIT"`,
		}),
	}
	var existAlicloudPvtzEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"endpoints.#":                         "1",
			"endpoints.0.endpoint_name":           fmt.Sprintf("tf-testacc%d", rand),
			"endpoints.0.id":                      CHECKSET,
			"endpoints.0.status":                  CHECKSET,
			"endpoints.0.security_group_id":       CHECKSET,
			"endpoints.0.create_time":             CHECKSET,
			"endpoints.0.vpc_id":                  CHECKSET,
			"endpoints.0.vpc_region_id":           CHECKSET,
			"endpoints.0.vpc_name":                CHECKSET,
			"endpoints.0.ip_configs.#":            "2",
			"endpoints.0.ip_configs.0.zone_id":    CHECKSET,
			"endpoints.0.ip_configs.0.cidr_block": CHECKSET,
			"endpoints.0.ip_configs.0.vswitch_id": CHECKSET,
			"endpoints.0.ip_configs.0.ip":         CHECKSET,
			"endpoints.0.ip_configs.1.zone_id":    CHECKSET,
			"endpoints.0.ip_configs.1.cidr_block": CHECKSET,
			"endpoints.0.ip_configs.1.vswitch_id": CHECKSET,
			"endpoints.0.ip_configs.1.ip":         CHECKSET,
		}
	}
	var fakeAlicloudPvtzEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudPvtzEndpointsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pvtz_endpoints.default",
		existMapFunc: existAlicloudPvtzEndpointsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudPvtzEndpointsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudPvtzEndpointsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudPvtzEndpointsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {  
   default = "tf-testacc%d"
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

data "alicloud_vpcs" "default" {
   name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
   vpc_id = data.alicloud_vpcs.default.ids.0
   zone_id      = data.alicloud_pvtz_resolver_zones.default.zones.0.zone_id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_pvtz_resolver_zones.default.zones.0.zone_id
  vswitch_name      = var.name
}

data "alicloud_vswitches" "default1" {
   vpc_id = data.alicloud_vpcs.default.ids.0
   zone_id      = data.alicloud_pvtz_resolver_zones.default.zones.1.zone_id
}

resource "alicloud_vswitch" "vswitch1" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_pvtz_resolver_zones.default.zones.1.zone_id
  vswitch_name      = var.name
}

locals {
  zone_id_1 =  data.alicloud_pvtz_resolver_zones.default.zones.0.zone_id
  zone_id_2 =  data.alicloud_pvtz_resolver_zones.default.zones.1.zone_id
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  vswitch_id_1 = length(data.alicloud_vswitches.default1.ids) > 0 ? data.alicloud_vswitches.default1.ids[0] : concat(alicloud_vswitch.vswitch1.*.id, [""])[0]

}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name   = var.name
}

resource "alicloud_pvtz_endpoint" "default" {
  endpoint_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vpc_id            = data.alicloud_vpcs.default.ids.0
  vpc_region_id     = "%s"
  ip_configs {
    zone_id    = local.zone_id_1
    cidr_block = data.alicloud_vswitches.default.vswitches[0].cidr_block
    vswitch_id = local.vswitch_id
  }
 ip_configs {
   zone_id    = local.zone_id_2
    cidr_block = data.alicloud_vswitches.default1.vswitches[0].cidr_block
    vswitch_id = local.vswitch_id_1
  }
}

data "alicloud_pvtz_endpoints" "default" { 
   %s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}

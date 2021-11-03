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

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  count      = 2
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = var.name
}

resource "alicloud_pvtz_endpoint" "default" {
  endpoint_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vpc_region_id     = "%s"
  ip_configs {
    zone_id    = alicloud_vswitch.default[0].zone_id
    cidr_block = alicloud_vswitch.default[0].cidr_block
    vswitch_id = alicloud_vswitch.default[0].id
  }
  ip_configs {
    zone_id    = alicloud_vswitch.default[1].zone_id
    cidr_block = alicloud_vswitch.default[1].cidr_block
    vswitch_id = alicloud_vswitch.default[1].id
  }
}

data "alicloud_pvtz_endpoints" "default" {	
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}

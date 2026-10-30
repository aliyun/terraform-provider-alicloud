package privatezone_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPvtzEndpointsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 99999999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_pvtz_endpoint.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"ids": `["${alicloud_pvtz_endpoint.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"ids":    `["${alicloud_pvtz_endpoint.default.id}"]`,
			"status": `"SUCCESS"`,
		}),
		FakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"ids":    `["${alicloud_pvtz_endpoint.default.id}"]`,
			"status": `"INIT"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
			"ids":        `["${alicloud_pvtz_endpoint.default.id}"]`,
			"name_regex": `"${alicloud_pvtz_endpoint.default.endpoint_name}"`,
			"status":     `"SUCCESS"`,
		}),
		FakeConfig: testAccCheckAlicloudPvtzEndpointsDataSourceName(name, map[string]string{
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
			"endpoints.0.endpoint_name":           name,
			"endpoints.0.id":                      acctest.CHECKSET,
			"endpoints.0.status":                  acctest.CHECKSET,
			"endpoints.0.security_group_id":       acctest.CHECKSET,
			"endpoints.0.create_time":             acctest.CHECKSET,
			"endpoints.0.vpc_id":                  acctest.CHECKSET,
			"endpoints.0.vpc_region_id":           acctest.CHECKSET,
			"endpoints.0.vpc_name":                acctest.CHECKSET,
			"endpoints.0.ip_configs.#":            "2",
			"endpoints.0.ip_configs.0.zone_id":    acctest.CHECKSET,
			"endpoints.0.ip_configs.0.cidr_block": acctest.CHECKSET,
			"endpoints.0.ip_configs.0.vswitch_id": acctest.CHECKSET,
			"endpoints.0.ip_configs.0.ip":         acctest.CHECKSET,
			"endpoints.0.ip_configs.1.zone_id":    acctest.CHECKSET,
			"endpoints.0.ip_configs.1.cidr_block": acctest.CHECKSET,
			"endpoints.0.ip_configs.1.vswitch_id": acctest.CHECKSET,
			"endpoints.0.ip_configs.1.ip":         acctest.CHECKSET,
		}
	}
	var fakeAlicloudPvtzEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudPvtzEndpointsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_pvtz_endpoints.default",
		ExistMapFunc: existAlicloudPvtzEndpointsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudPvtzEndpointsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudPvtzEndpointsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudPvtzEndpointsDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {  
   default = "%s"
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
   vpc_id = data.alicloud_vpcs.default.ids.0
   zone_id      = data.alicloud_pvtz_resolver_zones.default.zones.0.zone_id
}

data "alicloud_vswitches" "default1" {
   vpc_id = data.alicloud_vpcs.default.ids.0
   zone_id      = data.alicloud_pvtz_resolver_zones.default.zones.1.zone_id
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
    zone_id    = data.alicloud_vswitches.default.vswitches[0].zone_id
    cidr_block = data.alicloud_vswitches.default.vswitches[0].cidr_block
    vswitch_id = data.alicloud_vswitches.default.vswitches[0].vswitch_id
  }
 ip_configs {
   zone_id    = data.alicloud_vswitches.default1.vswitches[0].zone_id
    cidr_block = data.alicloud_vswitches.default1.vswitches[0].cidr_block
    vswitch_id = data.alicloud_vswitches.default1.vswitches[0].vswitch_id
  }
}

data "alicloud_pvtz_endpoints" "default" { 
   %s
}
`, name, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}

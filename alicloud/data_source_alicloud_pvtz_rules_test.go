package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPvtzRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_pvtz_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_pvtz_rule.default.id}_fake"]`,
		}),
	}
	endpointIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_pvtz_rule.default.id}"]`,
			"endpoint_id": `"${alicloud_pvtz_rule.default.endpoint_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_pvtz_rule.default.id}"]`,
			"endpoint_id": `"${alicloud_pvtz_rule.default.endpoint_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_pvtz_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_pvtz_rule.default.rule_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"endpoint_id": `"${alicloud_pvtz_rule.default.endpoint_id}"`,
			"ids":         `["${alicloud_pvtz_rule.default.id}"]`,
			"name_regex":  `"${alicloud_pvtz_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudPvtzRulesDataSourceName(name, map[string]string{
			"endpoint_id": `"${alicloud_pvtz_rule.default.endpoint_id}_fake"`,
			"ids":         `["${alicloud_pvtz_rule.default.id}_fake"]`,
			"name_regex":  `"${alicloud_pvtz_rule.default.rule_name}_fake"`,
		}),
	}
	var existAlicloudPvtzRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"rules.#":                    "1",
			"rules.0.endpoint_id":        CHECKSET,
			"rules.0.endpoint_name":      CHECKSET,
			"rules.0.create_time":        CHECKSET,
			"rules.0.id":                 CHECKSET,
			"rules.0.rule_id":            CHECKSET,
			"rules.0.rule_name":          name,
			"rules.0.type":               "OUTBOUND",
			"rules.0.zone_name":          name,
			"rules.0.forward_ips.#":      "1",
			"rules.0.forward_ips.0.ip":   "114.114.114.114",
			"rules.0.forward_ips.0.port": "8080",
		}
	}
	var fakeAlicloudPvtzRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudPvtzRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pvtz_rules.default",
		existMapFunc: existAlicloudPvtzRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudPvtzRulesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudPvtzRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, endpointIdConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudPvtzRulesDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
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
  vpc_region_id     = var.region
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

resource "alicloud_pvtz_rule" "default" {
  endpoint_id = alicloud_pvtz_endpoint.default.id
  rule_name   = var.name
  type        = "OUTBOUND"
  zone_name   = var.name
  forward_ips {
    ip   = "114.114.114.114"
    port = 8080
  }
}

data "alicloud_pvtz_rules" "default" {	
	%s
}
`, name, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}

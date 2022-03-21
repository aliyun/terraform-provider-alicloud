package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBRulesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_rule.default.id}_fake"]`,
		}),
	}

	listenerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"listener_ids": `["${alicloud_alb_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"listener_ids": `["${alicloud_alb_listener.default.id}_fake"]`,
		}),
	}

	loadBalancerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}

	ruleIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"rule_ids": `["${alicloud_alb_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"rule_ids": `["${alicloud_alb_rule.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_rule.default.rule_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_rule.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_rule.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_rule.default.id}"]`,
			"rule_ids":          `["${alicloud_alb_rule.default.id}"]`,
			"name_regex":        `"${alicloud_alb_rule.default.rule_name}"`,
			"status":            `"Available"`,
			"listener_ids":      `["${alicloud_alb_listener.default.id}"]`,
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbRuleDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_rule.default.id}_fake"]`,
			"rule_ids":          `["${alicloud_alb_rule.default.id}_fake"]`,
			"status":            `"Configuring"`,
			"name_regex":        `"${alicloud_alb_rule.default.rule_name}_fake"`,
			"listener_ids":      `["${alicloud_alb_listener.default.id}_fake"]`,
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}

	var existDataAlicloudAlbRulesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"rules.#":                                      "1",
			"rules.0.rule_name":                            fmt.Sprintf("tf-testAccAlbRule%d", rand),
			"rules.0.status":                               "Available",
			"rules.0.listener_id":                          CHECKSET,
			"rules.0.priority":                             "555",
			"rules.0.rule_actions.#":                       "2",
			"rules.0.rule_actions.0.order":                 CHECKSET,
			"rules.0.rule_actions.0.type":                  CHECKSET,
			"rules.0.rule_conditions.#":                    "1",
			"rules.0.rule_conditions.0.type":               "SourceIp",
			"rules.0.rule_conditions.0.source_ip_config.#": "1",
			"rules.0.rule_conditions.0.source_ip_config.0.values.#": "1",
			"rules.0.rule_conditions.0.source_ip_config.0.values.0": "192.168.0.0/24",
		}
	}
	var fakeDataAlicloudAlbRulesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"rules.#": "0",
		}
	}
	var alicloudAlbRuleCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_rules.default",
		existMapFunc: existDataAlicloudAlbRulesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbRulesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbRuleCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, ruleIdsConf, idsConf, nameRegexConf, listenerIdsConf, statusConf, loadBalancerIdsConf, allConf)
}
func testAccCheckAlicloudAlbRuleDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAlbRule%d"
}

data "alicloud_alb_zones" "default"{}

data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count             = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id =  data.alicloud_alb_zones.default.zones.0.id
  vswitch_name              = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count             = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name              = var.name
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id =              data.alicloud_vpcs.default.ids.0
  address_type =        "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name =    var.name
  load_balancer_edition = "Standard"
  load_balancer_billing_config {
    pay_type = 	"PayAsYouGo"
  }
  zone_mappings{
		vswitch_id =  length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
		zone_id =  data.alicloud_alb_zones.default.zones.0.id
	}
  zone_mappings{
		vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
		zone_id =   data.alicloud_alb_zones.default.zones.1.id
	}
}

resource "alicloud_alb_server_group" "default" {
	count = 3
	protocol = "HTTP"
	vpc_id = data.alicloud_vpcs.default.vpcs.0.id
	server_group_name = var.name
	health_check_config {
       health_check_enabled = "false"
	}
	sticky_session_config {
       sticky_session_enabled = "false"
	}
}

resource "alicloud_alb_listener" "default" {
	load_balancer_id = alicloud_alb_load_balancer.default.id
	listener_protocol =  "HTTP"
	listener_port = 8080
	listener_description = var.name
	default_actions{
		type = "ForwardGroup"
		forward_group_config{
			server_group_tuples{
				server_group_id = alicloud_alb_server_group.default.0.id
			}
		}
	}
}

resource "alicloud_alb_rule" "default" {
  rule_name   = var.name
  listener_id = alicloud_alb_listener.default.id
  priority    = "555"
  rule_conditions {
    source_ip_config {
      values = ["192.168.0.0/24"]
    }
    type = "SourceIp"
  }
  rule_actions {
    traffic_mirror_config {
      target_type = "ForwardGroupMirror"
      mirror_group_config {
        server_group_tuples {
          server_group_id = alicloud_alb_server_group.default.2.id
        }
      }
    }
    order = 1
    type  = "TrafficMirror"
  }
  rule_actions {
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.0.id
        weight          = 1
      }
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.1.id
        weight          = 2
      }
    }
    order = 2
    type  = "ForwardGroup"
  }
}

data "alicloud_alb_rules" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

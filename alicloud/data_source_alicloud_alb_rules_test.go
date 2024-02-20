package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudALBRulesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_rule.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_rule.default.rule_name}_fake"`,
		}),
	}
	loadBalancerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}
	listenerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"listener_ids": `["${alicloud_alb_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"listener_ids": `["${alicloud_alb_listener.default.id}_fake"]`,
		}),
	}
	ruleIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"rule_ids": `["${alicloud_alb_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"rule_ids": `["${alicloud_alb_rule.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"status": `"Provisioning"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_rule.default.id}"]`,
			"name_regex":        `"${alicloud_alb_rule.default.rule_name}"`,
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}"]`,
			"listener_ids":      `["${alicloud_alb_listener.default.id}"]`,
			"rule_ids":          `["${alicloud_alb_rule.default.id}"]`,
			"status":            `"Available"`,
		}),
		fakeConfig: testAccCheckAliCloudAlbRuleDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_rule.default.id}_fake"]`,
			"name_regex":        `"${alicloud_alb_rule.default.rule_name}_fake"`,
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"listener_ids":      `["${alicloud_alb_listener.default.id}_fake"]`,
			"rule_ids":          `["${alicloud_alb_rule.default.id}_fake"]`,
			"status":            `"Provisioning"`,
		}),
	}

	var existDataAliCloudAlbRulesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"rules.#":                      "1",
			"rules.0.rule_name":            fmt.Sprintf("tf-testAccAlbRule%d", rand),
			"rules.0.status":               "Available",
			"rules.0.listener_id":          CHECKSET,
			"rules.0.priority":             "555",
			"rules.0.rule_actions.#":       "2",
			"rules.0.rule_actions.0.order": CHECKSET,
			"rules.0.rule_actions.0.type":  CHECKSET,
			"rules.0.rule_actions.0.traffic_limit_config.#":         "1",
			"rules.0.rule_actions.0.traffic_limit_config.0.qps":     "120",
			"rules.0.rule_actions.1.order":                          CHECKSET,
			"rules.0.rule_actions.1.type":                           CHECKSET,
			"rules.0.rule_actions.1.redirect_config.#":              "1",
			"rules.0.rule_actions.1.redirect_config.0.port":         "10",
			"rules.0.rule_conditions.#":                             "1",
			"rules.0.rule_conditions.0.type":                        "SourceIp",
			"rules.0.rule_conditions.0.source_ip_config.#":          "1",
			"rules.0.rule_conditions.0.source_ip_config.0.values.#": "1",
			"rules.0.rule_conditions.0.source_ip_config.0.values.0": "192.168.0.0/24",
		}
	}
	var fakeDataAliCloudAlbRulesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}
	var alicloudAlbRuleCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_rules.default",
		existMapFunc: existDataAliCloudAlbRulesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAliCloudAlbRulesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbRuleCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, ruleIdsConf, idsConf, nameRegexConf, listenerIdsConf, statusConf, loadBalancerIdsConf, allConf)
}

func testAccCheckAliCloudAlbRuleDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccAlbRule%d"
	}

	data "alicloud_alb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "vswitch_1" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  		zone_id      = data.alicloud_alb_zones.default.zones.0.id
  		vswitch_name = var.name
	}

	resource "alicloud_vswitch" "vswitch_2" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
  		zone_id      = data.alicloud_alb_zones.default.zones.1.id
  		vswitch_name = var.name
	}

	resource "alicloud_alb_load_balancer" "default" {
  		vpc_id                 = alicloud_vpc.default.id
  		address_type           = "Internet"
  		address_allocated_mode = "Fixed"
  		load_balancer_name     = var.name
  		load_balancer_edition  = "Standard"
  		load_balancer_billing_config {
    		pay_type = "PayAsYouGo"
  		}
  		zone_mappings {
    		vswitch_id = alicloud_vswitch.vswitch_1.id
    		zone_id    = data.alicloud_alb_zones.default.zones.0.id
  		}
  		zone_mappings {
    		vswitch_id = alicloud_vswitch.vswitch_2.id
    		zone_id    = data.alicloud_alb_zones.default.zones.1.id
  		}
	}

	resource "alicloud_alb_server_group" "default" {
  		protocol          = "HTTP"
  		vpc_id            = alicloud_vpc.default.id
  		server_group_name = var.name
  		health_check_config {
    		health_check_enabled = "false"
  		}
  		sticky_session_config {
    		sticky_session_enabled = "false"
  		}
	}

	resource "alicloud_alb_listener" "default" {
  		load_balancer_id     = alicloud_alb_load_balancer.default.id
  		listener_protocol    = "HTTP"
  		listener_port        = 8080
  		listener_description = var.name
  		default_actions {
    		type = "ForwardGroup"
    		forward_group_config {
      			server_group_tuples {
        			server_group_id = alicloud_alb_server_group.default.id
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
    		traffic_limit_config {
      			qps        = "120"
      			per_ip_qps = "120"
    		}
    		order = 1
    		type  = "TrafficLimit"
  		}
  		rule_actions {
    		redirect_config {
      			host      = "ww.ali.com"
      			http_code = "301"
      			path      = "/test"
      			port      = "10"
      			protocol  = "HTTP"
      			query     = "query"
    		}
    		order = 2
    		type  = "Redirect"
  		}
	}

	data "alicloud_alb_rules" "default" {
	%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}

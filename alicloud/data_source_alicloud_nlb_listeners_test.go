package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNlbListenersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_listener.default.id}_fake"]`,
		}),
	}
	loadBalancerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_nlb_listener.default.load_balancer_id}"]`,
			"ids":               `["${alicloud_nlb_listener.default.id}"]`,
		}),
		fakeConfig: "",
	}
	listenerProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"listener_protocol": `"${alicloud_nlb_listener.default.listener_protocol}"`,
			"ids":               `["${alicloud_nlb_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_listener.default.id}"]`,
			"listener_protocol": `"UDP"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_listener.default.id}"]`,
			"listener_protocol": `"${alicloud_nlb_listener.default.listener_protocol}"`,
			"load_balancer_ids": `["${alicloud_nlb_listener.default.load_balancer_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbListenersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_listener.default.id}_fake"]`,
			"listener_protocol": `"UDP"`,
			"load_balancer_ids": `["${alicloud_nlb_listener.default.load_balancer_id}"]`,
		}),
	}
	var existAlicloudNlbListenersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"listeners.#":                        "1",
			"listeners.0.id":                     CHECKSET,
			"listeners.0.alpn_enabled":           CHECKSET,
			"listeners.0.alpn_policy":            "",
			"listeners.0.ca_certificate_ids.#":   "0",
			"listeners.0.certificate_ids.#":      "0",
			"listeners.0.ca_enabled":             CHECKSET,
			"listeners.0.cps":                    "10000",
			"listeners.0.end_port":               "",
			"listeners.0.idle_timeout":           "900",
			"listeners.0.listener_description":   CHECKSET,
			"listeners.0.listener_id":            CHECKSET,
			"listeners.0.listener_port":          "80",
			"listeners.0.listener_protocol":      "TCP",
			"listeners.0.load_balancer_id":       CHECKSET,
			"listeners.0.mss":                    "0",
			"listeners.0.proxy_protocol_enabled": "true",
			"listeners.0.sec_sensor_enabled":     "true",
			"listeners.0.security_policy_id":     "",
			"listeners.0.server_group_id":        CHECKSET,
			"listeners.0.start_port":             "",
			"listeners.0.status":                 CHECKSET,
		}
	}
	var fakeAlicloudNlbListenersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var AlicloudNlbListenersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nlb_listeners.default",
		existMapFunc: existAlicloudNlbListenersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNlbListenersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	AlicloudNlbListenersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, loadBalancerIdsConf, listenerProtocolConf, allConf)
}
func testAccCheckAlicloudNlbListenersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccListener-%d"
}
resource "alicloud_nlb_server_group" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name = var.name
  server_group_type = "Instance"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
  health_check {
	health_check_url =           "/test/index.html"
	health_check_domain =       "tf-testAcc.com"
    health_check_enabled         = true
    health_check_type            = "TCP"
    health_check_connect_port    = 0
    healthy_threshold            = 2
    unhealthy_threshold          = 2
    health_check_connect_timeout = 5
    health_check_interval        = 10
    http_check_method            = "GET"
    health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
  }
  connection_drain           = true
  connection_drain_timeout   = 60
  preserve_client_ip_enabled = true
  tags = {
    Created = "TF"
  }
  address_ip_version = "Ipv4"
}

data "alicloud_nlb_zones" "default" {}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.0.id
}
data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.1.id
}
locals {
  zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_id_1 = data.alicloud_vswitches.default_1.ids[0]
  zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  vswitch_id_2 = data.alicloud_vswitches.default_2.ids[0]
}
resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  tags               = {
    Created = "tfTestAcc0"
    For     = "Tftestacc 0"
  }
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_mappings {
    vswitch_id = local.vswitch_id_1
    zone_id    = local.zone_id_1
  }
  zone_mappings {
    vswitch_id = local.vswitch_id_2
    zone_id    = local.zone_id_2
  }
}

resource "alicloud_nlb_listener" "default" {
  listener_protocol      = "TCP"
  listener_port          = "80"
  listener_description   = var.name
  load_balancer_id       = alicloud_nlb_load_balancer.default.id
  server_group_id        = alicloud_nlb_server_group.default.id
  idle_timeout           = "900"
  proxy_protocol_enabled = "true"
  sec_sensor_enabled     = "true"
  cps                    = "10000"
  mss                    = "0"
}

data "alicloud_nlb_listeners" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsLoadBalancerUDPListenerDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_load_balancer_udp_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_load_balancer_udp_listener.default.id}_fake"]`,
		}),
	}

	LoadBalancerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_load_balancer_udp_listener.default.id}"]`,
			"load_balancer_id": `"${alicloud_ens_load_balancer.defaultgNxO1j.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_load_balancer_udp_listener.default.id}_fake"]`,
			"load_balancer_id": `"${alicloud_ens_load_balancer.defaultgNxO1j.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_load_balancer_udp_listener.default.id}"]`,
			"load_balancer_id": `"${alicloud_ens_load_balancer.defaultgNxO1j.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_load_balancer_udp_listener.default.id}_fake"]`,
			"load_balancer_id": `"${alicloud_ens_load_balancer.defaultgNxO1j.id}_fake"`,
		}),
	}

	EnsLoadBalancerUDPListenerCheckInfo.dataSourceTestCheck(t, rand, idsConf, LoadBalancerIdConf, allConf)
}

var existEnsLoadBalancerUDPListenerMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"listeners.#":                  "1",
		"listeners.0.status":           CHECKSET,
		"listeners.0.listener_port":    CHECKSET,
		"listeners.0.description":      CHECKSET,
		"listeners.0.load_balancer_id": CHECKSET,
	}
}

var fakeEnsLoadBalancerUDPListenerMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"listeners.#": "0",
	}
}

var EnsLoadBalancerUDPListenerCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_load_balancer_udp_listeners.default",
	existMapFunc: existEnsLoadBalancerUDPListenerMapFunc,
	fakeMapFunc:  fakeEnsLoadBalancerUDPListenerMapFunc,
}

func testAccCheckAlicloudEnsLoadBalancerUDPListenerSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsLoadBalancerUDPListener%d"
}
variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "default8QXHtu" {
  network_name  = "测试用例-测试udp监听"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultN8wZgT" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "测试用例-测试udp监听"
  ens_region_id = alicloud_ens_network.default8QXHtu.ens_region_id
  network_id    = alicloud_ens_network.default8QXHtu.id
}

resource "alicloud_ens_load_balancer" "defaultgNxO1j" {
  load_balancer_name = "测试用例-测试udp监听"
  vswitch_id         = alicloud_ens_vswitch.defaultN8wZgT.id
  payment_type       = "PayAsYouGo"
  ens_region_id      = alicloud_ens_vswitch.defaultN8wZgT.ens_region_id
  network_id         = alicloud_ens_vswitch.defaultN8wZgT.network_id
  load_balancer_spec = "elb.s1.small"
}



resource "alicloud_ens_load_balancer_udp_listener" "default" {
  listener_port                = "53"
  health_check_interval        = "1"
  description                  = "test1"
  unhealthy_threshold          = "2"
  scheduler                    = "rr"
  health_check_connect_timeout = "1"
  load_balancer_id             = alicloud_ens_load_balancer.defaultgNxO1j.id
  backend_server_port          = "53"
  health_check_connect_port    = "53"
  health_check_req             = "hello"
  healthy_threshold            = "2"
  health_check_exp             = "rep"
  eip_transmit                 = "on"
  status                       = "Stopped"
  established_timeout          = "100"
}

data "alicloud_ens_load_balancer_udp_listeners" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

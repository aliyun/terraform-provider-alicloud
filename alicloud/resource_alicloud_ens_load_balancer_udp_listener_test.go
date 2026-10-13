// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Ens LoadBalancerUdpListener. >>> Resource test cases, automatically generated.
// Case 负载均衡实例-UDP监听_20240430 6636
func TestAccAliCloudEnsLoadBalancerUdpListener_basic6636(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer_udp_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsLoadBalancerUdpListenerMap6636)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancerUdpListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsLoadBalancerUdpListenerBasicDependence6636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_port":                "53",
					"health_check_interval":        "1",
					"description":                  "test1",
					"unhealthy_threshold":          "2",
					"scheduler":                    "rr",
					"health_check_connect_timeout": "1",
					"load_balancer_id":             "${alicloud_ens_load_balancer.defaultgNxO1j.id}",
					"backend_server_port":          "53",
					"health_check_connect_port":    "53",
					"health_check_req":             "hello",
					"healthy_threshold":            "2",
					"health_check_exp":             "rep",
					"eip_transmit":                 "on",
					"status":                       "Stopped",
					"established_timeout":          "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_port":                "53",
						"health_check_interval":        "1",
						"description":                  "test1",
						"unhealthy_threshold":          "2",
						"scheduler":                    "rr",
						"health_check_connect_timeout": "1",
						"load_balancer_id":             CHECKSET,
						"backend_server_port":          "53",
						"health_check_connect_port":    "53",
						"health_check_req":             "hello",
						"healthy_threshold":            "2",
						"health_check_exp":             "rep",
						"eip_transmit":                 "on",
						"status":                       "Stopped",
						"established_timeout":          "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":        "10",
					"description":                  "test2",
					"unhealthy_threshold":          "10",
					"scheduler":                    "sch",
					"health_check_connect_timeout": "300",
					"health_check_connect_port":    "1000",
					"health_check_req":             "hello210",
					"healthy_threshold":            "10",
					"health_check_exp":             "exp2",
					"eip_transmit":                 "off",
					"established_timeout":          "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":        "10",
						"description":                  "test2",
						"unhealthy_threshold":          "10",
						"scheduler":                    "sch",
						"health_check_connect_timeout": "300",
						"health_check_connect_port":    "1000",
						"health_check_req":             "hello210",
						"healthy_threshold":            "10",
						"health_check_exp":             "exp2",
						"eip_transmit":                 "off",
						"established_timeout":          "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsLoadBalancerUdpListenerMap6636 = map[string]string{}

func AlicloudEnsLoadBalancerUdpListenerBasicDependence6636(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
    default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "default8QXHtu" {
        network_name = "测试用例-测试udp监听"
        cidr_block = "10.0.0.0/8"
        ens_region_id = "${var.ens_region_id}"
}

resource "alicloud_ens_vswitch" "defaultN8wZgT" {
        cidr_block = "10.0.6.0/24"
        vswitch_name = "测试用例-测试udp监听"
        ens_region_id = "${alicloud_ens_network.default8QXHtu.ens_region_id}"
        network_id = "${alicloud_ens_network.default8QXHtu.id}"
}

resource "alicloud_ens_load_balancer" "defaultgNxO1j" {
        load_balancer_name = "测试用例-测试udp监听"
        vswitch_id = "${alicloud_ens_vswitch.defaultN8wZgT.id}"
        payment_type = "PayAsYouGo"
        ens_region_id = "${alicloud_ens_vswitch.defaultN8wZgT.ens_region_id}"
        network_id = "${alicloud_ens_vswitch.defaultN8wZgT.network_id}"
        load_balancer_spec = "elb.s1.small"
}


`, name)
}

// Test Ens LoadBalancerUdpListener. <<< Resource test cases, automatically generated.

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSLBRuleUpdate(t *testing.T) {
	var v *slb.DescribeRuleAttributeResponse
	resourceId := "alicloud_slb_rule.default"
	ra := resourceAttrInit(resourceId, ruleMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbRuleBasic")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbRuleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"frontend_port":    "${alicloud_slb_listener.default.frontend_port}",
					"domain":           "*.aliyun.com",
					"url":              "/image",
					"server_group_id":  "${alicloud_slb_server_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_sync": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_sync": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session":      "on",
					"sticky_session_type": "server",
					"cookie":              "23ffsa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session":      "on",
						"sticky_session_type": "server",
						"cookie":              "23ffsa",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_http_code": "http_2xx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_code": "http_2xx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_uri": "/test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "wrr",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_domain": "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_domain": "test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cookie": "23ffsa1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie": "23ffsa1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cookie_timeout":      "100",
					"sticky_session_type": "insert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_type": "insert",
						"cookie_timeout":      "100",
						"cookie":              "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cookie_timeout": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie_timeout": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check":              "off",
						"health_check_http_code":    "",
						"health_check_interval":     "0",
						"health_check_domain":       "",
						"health_check_uri":          "",
						"health_check_connect_port": "0",
						"health_check_timeout":      "0",
						"healthy_threshold":         "0",
						"unhealthy_threshold":       "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session":      "off",
						"sticky_session_type": "",
						"cookie_timeout":      "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbRuleBasic_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      "${var.name}",
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"frontend_port":             "${alicloud_slb_listener.default.frontend_port}",
					"domain":                    "*.aliyun.com",
					"url":                       "/image",
					"server_group_id":           "${alicloud_slb_server_group.default.id}",
					"cookie":                    "23ffsa",
					"cookie_timeout":            "100",
					"health_check_http_code":    "http_2xx",
					"health_check_interval":     "10",
					"health_check_uri":          "/test",
					"health_check_connect_port": "80",
					"health_check_timeout":      "10",
					"healthy_threshold":         "3",
					"unhealthy_threshold":       "3",
					"sticky_session":            "on",
					"sticky_session_type":       "server",
					"listener_sync":             "off",
					"scheduler":                 "rr",
					"health_check_domain":       "test",
					"health_check":              "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"frontend_port":             "22",
						"name":                      "tf-testAccSlbRuleBasic",
						"domain":                    "*.aliyun.com",
						"url":                       "/image",
						"server_group_id":           CHECKSET,
						"cookie":                    "23ffsa",
						"health_check_http_code":    "http_2xx",
						"health_check_interval":     "10",
						"health_check_domain":       "test",
						"health_check_uri":          "/test",
						"health_check_connect_port": "80",
						"health_check_timeout":      "10",
						"healthy_threshold":         "3",
						"unhealthy_threshold":       "3",
						"sticky_session":            "on",
						"sticky_session_type":       "server",
						"health_check":              "on",
						"listener_sync":             "off",
						"scheduler":                 "rr",
					}),
				),
			},
		},
	})
}

func resourceSlbRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
  availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_instance_types" "new" {
	eni_amount = 2
    availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  instance_name = "${var.name}"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  servers {
      server_ids = "${alicloud_instance.default.*.id}"
      port = 80
      weight = 100
    }
}
`, name)
}

var ruleMap = map[string]string{
	"load_balancer_id": CHECKSET,
	"frontend_port":    "22",
	"domain":           "*.aliyun.com",
	"url":              "/image",
	"server_group_id":  CHECKSET,
}

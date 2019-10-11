package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSlbRuleUpdate(t *testing.T) {
	var v *slb.DescribeRuleAttributeResponse
	resourceId := "alicloud_slb_rule.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbRuleUpdate%d", rand)
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
					"name":             "${var.name}",
					"load_balancer_id": "${alicloud_slb.default.id}",
					"frontend_port":    "${alicloud_slb_listener.default.frontend_port}",
					"domain":           "*.aliyun.com",
					"url":              "/image",
					"server_group_id":  "${alicloud_slb_server_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"frontend_port":    "22",
						"name":             "tf-testAccSlbRuleBasic",
						"domain":           "*.aliyun.com",
						"url":              "/image",
						"server_group_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      "${var.name}",
					"load_balancer_id":          "${alicloud_slb.default.id}",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_http_code": "http_3xx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_code": "http_3xx",
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
					"health_check_http_code":    "http_2xx",
					"health_check_interval":     "10",
					"health_check_uri":          "/test",
					"health_check_connect_port": "90",
					"health_check_timeout":      "30",
					"healthy_threshold":         "4",
					"unhealthy_threshold":       "5",
					"health_check":              "off",
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
					"sticky_session":      "off",
					"health_check_domain": "test",
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
						"sticky_session":      "off",
						"sticky_session_type": "",
						"cookie_timeout":      "0",
						"name":                "tf-testAccSlbRuleBasic_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":           "${var.name}",
					"sticky_session": "off",
					"listener_sync":  "on",
					"scheduler":      "wrr",
					"health_check":   "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_sync":  "on",
						"sticky_session": "",
						"health_check":   "",
						"scheduler":      "",
						"name":           "tf-testAccSlbRuleBasic",
					}),
				),
			},
		},
	})
}

func resourceSlbRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbRuleBasic"
}
data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  name = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
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
  vswitch_id = "${alicloud_vswitch.default.id}"
  instance_name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  servers {
      server_ids = "${alicloud_instance.default.*.id}"
      port = 80
      weight = 100
    }
}
`)
}

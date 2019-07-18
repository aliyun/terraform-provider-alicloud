package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbRuleBasic,
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
				Config: testAccSlbRuleBasic_server_group_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccSlbRuleBasic_health_check,
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
				Config: testAccSlbRuleBasic_update_health_check_http_code,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_code": "http_3xx",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_health_check_interval,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "20",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_health_check_uri,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/test1",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_health_check_connect_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "90",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_health_check_time_out,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "20",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_healthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "5",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_unhealthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "5",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_scheduler,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "wrr",
					}),
				),
			},

			{
				Config: testAccSlbRuleBasic_update_health_check_domain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_domain": "test1",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_cookie,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie": "23ffsa1",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_sticky_session_type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_type": "insert",
						"cookie_timeout":      "100",
						"cookie":              "",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_cookie_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie_timeout": "200",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_health_check,
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
				Config: testAccSlbRuleBasic_update_sticky_session,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session":      "off",
						"sticky_session_type": "",
						"cookie_timeout":      "0",
					}),
				),
			},
			{
				Config: testAccSlbRuleBasic_update_listener_sync,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_sync":  "on",
						"sticky_session": "",
						"health_check":   "",
						"scheduler":      "",
					}),
				),
			},
		},
	})
}

func testAccCheckSlbRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_rule" {
			continue
		}

		// Try to find the Slb server group
		if _, err := slbService.DescribeSlbRule(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("SLB Rule %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccSlbRuleBasic = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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

resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
}
`

const testAccSlbRuleBasic_server_group_id = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name = "${var.name}_test"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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

resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
}
`
const testAccSlbRuleBasic_health_check = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_2xx"
  health_check_interval = 10
  health_check_uri = "/test"
  health_check_connect_port = 80
  health_check_timeout = 10
  healthy_threshold = 3
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check_http_code = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 10
  health_check_uri = "/test"
  health_check_connect_port = 80
  health_check_timeout = 10
  healthy_threshold = 3
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check_interval = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test"
  health_check_connect_port = 80
  health_check_timeout = 10
  healthy_threshold = 3
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check_uri = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 80
  health_check_timeout = 10
  healthy_threshold = 3
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check_connect_port = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 10
  healthy_threshold = 3
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check_time_out = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 3
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_healthy_threshold = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 3
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_unhealthy_threshold = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "rr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_scheduler = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check_domain = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test1"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_cookie = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa1"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "server"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test1"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_sticky_session_type = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa1"
  cookie_timeout = 100
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "insert"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test1"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_cookie_timeout = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa1"
  cookie_timeout = 200
  health_check_http_code = "http_3xx"
  health_check_interval = 20
  health_check_uri = "/test1"
  health_check_connect_port = 90
  health_check_timeout = 20
  healthy_threshold = 5
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "insert"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test1"
  health_check = "on"
}`

const testAccSlbRuleBasic_update_health_check = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 200
  health_check_http_code = "http_2xx"
  health_check_interval = 10
  health_check_uri = "/test"
  health_check_connect_port = 90
  health_check_timeout = 30
  healthy_threshold = 4
  unhealthy_threshold = 5
  sticky_session = "on"
  sticky_session_type = "insert"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test"
  health_check = "off"
}`

const testAccSlbRuleBasic_update_sticky_session = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 200
  health_check_http_code = "http_2xx"
  health_check_interval = 10
  health_check_uri = "/test"
  health_check_connect_port = 90
  health_check_timeout = 30
  healthy_threshold = 4
  unhealthy_threshold = 5
  sticky_session = "off"
  sticky_session_type = "insert"
  listener_sync = "off"
  scheduler = "wrr"
  health_check_domain = "test"
  health_check = "off"
}`

const testAccSlbRuleBasic_update_listener_sync = `
variable "name" {
	default = "tf-testAccSlbRuleBasic"
}

data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
  cookie = "23ffsa"
  cookie_timeout = 200
  health_check_http_code = "http_2xx"
  health_check_interval = 10
  health_check_uri = "/test"
  health_check_connect_port = 90
  health_check_timeout = 30
  healthy_threshold = 4
  unhealthy_threshold = 5
  sticky_session = "off"
  sticky_session_type = "insert"
  listener_sync = "on"
  scheduler = "wrr"
  health_check_domain = "test"
  health_check = "off"
}`

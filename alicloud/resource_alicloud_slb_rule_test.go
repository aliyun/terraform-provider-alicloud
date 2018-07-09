package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSlbRule_basic(t *testing.T) {
	var rule slb.DescribeRuleAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_rule.rule",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbRuleExists("alicloud_slb_rule.rule", &rule),
					resource.TestCheckResourceAttr(
						"alicloud_slb_rule.rule", "name", "testAccSlbRuleBasic"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_rule.rule", "domain", "*.aliyun.com"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbRule_url(t *testing.T) {
	var rule slb.DescribeRuleAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_rule.rule",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbRuleUrl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbRuleExists("alicloud_slb_rule.rule", &rule),
					resource.TestCheckResourceAttr(
						"alicloud_slb_rule.rule", "name", "testAccSlbRuleUrl"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_rule.rule", "url", "/image"),
				),
			},
		},
	})
}

func testAccCheckSlbRuleExists(n string, rule *slb.DescribeRuleAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB Rule ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		r, err := client.DescribeLoadBalancerRuleAttribute(rs.Primary.ID)
		if err != nil {
			return err
		}

		*rule = *r

		return nil
	}
}

func testAccCheckSlbRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_rule" {
			continue
		}

		// Try to find the Slb server group
		if _, err := client.DescribeLoadBalancerRuleAttribute(rs.Primary.ID); err != nil {
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
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "image" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccSlbRuleBasic"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  depends_on = [
    "alicloud_vpc.main"]
}
resource "alicloud_security_group" "group" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  image_id = "${data.alicloud_images.image.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.group.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.main.id}"
  instance_name = "${var.name}"
}

resource "alicloud_slb" "instance" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  servers = [
    {
      server_ids = ["${alicloud_instance.instance.*.id}"]
      port = 80
      weight = 100
    }
  ]
}

resource "alicloud_slb_rule" "rule" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  frontend_port = "${alicloud_slb_listener.listener.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.group.id}"
}
`

const testAccSlbRuleUrl = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "image" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccSlbRuleUrl"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  depends_on = [
    "alicloud_vpc.main"]
}
resource "alicloud_security_group" "group" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  image_id = "${data.alicloud_images.image.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.group.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.main.id}"
  instance_name = "${var.name}"
}

resource "alicloud_slb" "instance" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  servers = [
    {
      server_ids = ["${alicloud_instance.instance.*.id}"]
      port = 80
      weight = 100
    }
  ]
}

resource "alicloud_slb_rule" "rule" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  frontend_port = "${alicloud_slb_listener.listener.frontend_port}"
  name = "${var.name}"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.group.id}"
}
`

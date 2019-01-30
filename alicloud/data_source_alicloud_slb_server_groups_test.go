package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbServerGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbServerGroupsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_server_groups.slb_server_groups"),
					resource.TestCheckResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.name", "tf-testAccslbservergroupsdatasourcebasic"),
					resource.TestCheckResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.servers.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.servers.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.servers.0.weight", "100"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbServerGroupsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbServerGroupsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_server_groups.slb_server_groups"),
					resource.TestCheckResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_server_groups.slb_server_groups", "slb_server_groups.0.servers.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbServerGroupsDataSourceBasic = `
variable "name" {
	default = "tf-testAccslbservergroupsdatasourcebasic"
}

data "alicloud_zones" "az" {
	"available_resource_creation" = "VSwitch"
}
data "alicloud_images" "images" {
  name_regex = "^ubuntu_16.*_64"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "instance_types" {
 	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}

resource "alicloud_vpc" "sample_vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "sample_vswitch" {
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.az.zones.0.id}"
}

resource "alicloud_security_group" "sample_security_group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.sample_vpc.id}"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_slb_listener" "sample_slb_listener" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "${var.name}"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for = {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}

resource "alicloud_instance" "sample_instance" {
  image_id = "${data.alicloud_images.images.images.0.id}"

  instance_type = "${data.alicloud_instance_types.instance_types.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.sample_security_group.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_slb_server_group" "sample_server_group" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.sample_instance.id}"]
      port = 80
      weight = 100
    }
  ]
}

resource "alicloud_slb_rule" "sample_rule" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  frontend_port = "${alicloud_slb_listener.sample_slb_listener.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.sample_server_group.id}"
}

data "alicloud_slb_server_groups" "slb_server_groups" {
  load_balancer_id = "${alicloud_slb_rule.sample_rule.load_balancer_id}"
  ids = ["${alicloud_slb_server_group.sample_server_group.id}"]
  name_regex = "${var.name}"
}
`
const testAccCheckAlicloudSlbServerGroupsDataSourceEmpty = `
variable "name" {
	default = "tf-testAccslbservergroupsdatasourcebasic"
}

data "alicloud_zones" "az" {
	"available_resource_creation" = "VSwitch"
}
data "alicloud_images" "images" {
  name_regex = "^ubuntu_16.*_64"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "instance_types" {
 	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}

resource "alicloud_vpc" "sample_vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "sample_vswitch" {
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.az.zones.0.id}"
}

resource "alicloud_security_group" "sample_security_group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.sample_vpc.id}"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

data "alicloud_slb_server_groups" "slb_server_groups" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  name_regex = "tf-testacc-fake*"
}
`

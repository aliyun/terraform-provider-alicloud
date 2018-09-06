package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbListenersDataSource_http(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbListenersDataSourceHttp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.backend_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.frontend_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.protocol", "http"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.status", "running"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.bandwidth", "10"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.scheduler", "wrr"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.sticky_session", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.sticky_session_type", "insert"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.cookie_timeout", "86400"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_uri", "/cons"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_port", "20"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.healthy_threshold", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.unhealthy_threshold", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_timeout", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_interval", "5"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_http_code", "http_2xx,http_3xx"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.gzip", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.x_forwarded_for", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.x_forwarded_for_slb_ip", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.x_forwarded_for_slb_id", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.x_forwarded_for_slb_proto", "off"),

					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners_with_filters"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners_with_filters", "slb_listeners.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListenersDataSource_tcp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbListenersDataSourceTcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.backend_port", "22"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.frontend_port", "22"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.status", "running"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.bandwidth", "10"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.scheduler", "wrr"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.persistence_timeout", "0"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.established_timeout", "900"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_type", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_port", "20"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_timeout", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.healthy_threshold", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.unhealthy_threshold", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_timeout", "0"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_interval", "5"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListenersDataSource_udp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbListenersDataSourceUdp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.backend_port", "10"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.frontend_port", "11"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.protocol", "udp"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.status", "running"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.bandwidth", "10"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.scheduler", "wrr"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_port", "20"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_timeout", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.healthy_threshold", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.unhealthy_threshold", "8"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_timeout", "0"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_interval", "5"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbListenersDataSourceHttp = `
variable "name" {
	default = "testAccCheckAlicloudSlbListenersDataSourceHttp"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
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

data "alicloud_slb_listeners" "slb_listeners" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
}

data "alicloud_slb_listeners" "slb_listeners_with_filters" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
  frontend_port = 80
  protocol = "http"
}
`

const testAccCheckAlicloudSlbListenersDataSourceTcp = `
variable "name" {
	default = "testAccCheckAlicloudSlbListenersDataSourceTcp"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_slb_listener" "sample_slb_listener" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "tcp"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_type = "tcp"
  bandwidth = 10
}

data "alicloud_slb_listeners" "slb_listeners" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
}
`

const testAccCheckAlicloudSlbListenersDataSourceUdp = `
variable "name" {
	default = "testAccCheckAlicloudSlbListenersDataSourceUdp"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_slb_listener" "sample_slb_listener" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  backend_port = 10
  frontend_port = 11
  protocol = "udp"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  bandwidth = 10
}

data "alicloud_slb_listeners" "slb_listeners" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
}
`

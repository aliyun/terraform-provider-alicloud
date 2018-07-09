package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSlbListener_http(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_listener.http",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbListenerHttp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbListenerExists("alicloud_slb_listener.http", 80),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "protocol", "http"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "backend_port", "80"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "health_check", "on"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_tcp(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_listener.tcp",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbListenerTcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbListenerExists("alicloud_slb_listener.tcp", 22),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.tcp", "backend_port", "22"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.tcp", "healthy_threshold", "8"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_udp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_listener.udp",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbListenerUdp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbListenerExists("alicloud_slb_listener.udp", 2001),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.udp", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.udp", "persistence_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.udp", "healthy_threshold", "8"),
				),
			},
		},
	})
}

func testAccCheckSlbListenerExists(n string, port int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB listener ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		parts := strings.Split(rs.Primary.ID, ":")
		loadBalancer, err := client.DescribeLoadBalancerAttribute(parts[0])
		if err != nil {
			return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
		}
		for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if portAndProtocol.ListenerPort == port {
				return nil
			}
		}

		return fmt.Errorf("The Listener %d not found.", port)
	}
}

func testAccCheckSlbListenerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_listener" {
			continue
		}

		// Try to find the Slb
		parts := strings.Split(rs.Primary.ID, ":")
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("Parsing SlbListener's id got an error: %#v", err)
		}
		loadBalancer, err := client.DescribeLoadBalancerAttribute(parts[0])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
		}
		for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if portAndProtocol.ListenerPort == port {
				return fmt.Errorf("SLB listener still exist")
			}
		}

	}

	return nil
}

const testAccSlbListenerHttp = `
resource "alicloud_slb" "instance" {
  name = "tf_test_slb_http"
  internet_charge_type = "paybytraffic"
  internet = true
}
resource "alicloud_slb_listener" "http" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
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
}
`

const testAccSlbListenerTcp = `
resource "alicloud_slb" "instance" {
  name = "tf_test_slb_tcp"
  internet_charge_type = "paybytraffic"
  internet = true
}
resource "alicloud_slb_listener" "tcp" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3600
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx"
  health_check_timeout = 8
  health_check_connect_port = 20
  health_check_uri = "/console"
}
`

const testAccSlbListenerUdp = `
resource "alicloud_slb" "instance" {
  name = "tf_test_slb_udp"
  internet_charge_type = "paybybandwidth"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "udp" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 4
  health_check_timeout = 8
  health_check_connect_port = 20
}
`

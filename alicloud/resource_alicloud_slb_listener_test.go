package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
						"alicloud_slb_listener.http", "acl_status", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "acl_type", string(AclTypeWhite)),
					resource.TestCheckResourceAttrSet(
						"alicloud_slb_listener.http", "acl_id"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "health_check", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "gzip", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "x_forwarded_for.0.retrive_client_ip", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "x_forwarded_for.0.retrive_slb_ip", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "x_forwarded_for.0.retrive_slb_id", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "x_forwarded_for.0.retrive_slb_proto", "false"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "idle_timeout", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.http", "request_timeout", "80"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_listener.https",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbListenerHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbListenerExists("alicloud_slb_listener.https", 80),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "protocol", "https"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "backend_port", "80"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "acl_status", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "acl_type", string(AclTypeWhite)),
					resource.TestCheckResourceAttrSet(
						"alicloud_slb_listener.https", "acl_id"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "health_check", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "gzip", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "x_forwarded_for.0.retrive_client_ip", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "x_forwarded_for.0.retrive_slb_ip", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "x_forwarded_for.0.retrive_slb_id", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "x_forwarded_for.0.retrive_slb_proto", "false"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "idle_timeout", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "request_timeout", "80"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "enable_http2", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.https", "tls_cipher_policy", "tls_cipher_policy_1_2"),
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
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.tcp", "acl_status", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.tcp", "acl_type", string(AclTypeWhite)),
					resource.TestCheckResourceAttrSet(
						"alicloud_slb_listener.tcp", "acl_id"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.tcp", "established_timeout", "600"),
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
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.udp", "acl_status", "on"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_listener.udp", "acl_type", string(AclTypeBlack)),
					resource.TestCheckResourceAttrSet(
						"alicloud_slb_listener.udp", "acl_id"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		slbService := SlbService{client}
		parts := strings.Split(rs.Primary.ID, ":")
		loadBalancer, err := slbService.DescribeLoadBalancerAttribute(parts[0])
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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

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
		loadBalancer, err := slbService.DescribeLoadBalancerAttribute(parts[0])
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
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
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
  x_forwarded_for = {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.acl.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
variable "name" {
  default = "tf-testAcc-http-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "acl" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}
`

const testAccSlbListenerTcp = `
resource "alicloud_slb" "instance" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
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
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.acl.id}"
  established_timeout = 600
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "acl" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}
`

const testAccSlbListenerUdp = `
resource "alicloud_slb" "instance" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
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
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.acl.id}"
}
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "acl" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}
`

const testAccSlbListenerHttps = `
resource "alicloud_slb" "instance" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s2.small"
}
resource "alicloud_slb_listener" "https" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
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
  x_forwarded_for = {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.acl.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.foo.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_2"
}
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "acl" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}
resource "alicloud_slb_server_certificate" "foo" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

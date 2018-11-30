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
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.idle_timeout", "30"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.request_timeout", "80"),
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners_with_filters"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners_with_filters", "slb_listeners.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListenersDataSource_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbListenersDataSourceHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.backend_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.frontend_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.protocol", "https"),
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
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.idle_timeout", "30"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.request_timeout", "80"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.enable_http2", "on"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.tls_cipher_policy", "tls_cipher_policy_1_2"),
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

func TestAccAlicloudSlbListenersDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbListenersDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.backend_port"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.frontend_port"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.protocol"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.bandwidth"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.scheduler"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_port"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_connect_timeout"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.healthy_threshold"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.unhealthy_threshold"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_timeout"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.health_check_interval"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbListenersDataSourceHttp = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceHttp"
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
  name = "${var.name}"
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
  request_timeout           = 80
  idle_timeout              = 30
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

const testAccCheckAlicloudSlbListenersDataSourceHttps = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceHttps"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s2.small"
}

resource "alicloud_slb_listener" "sample_slb_listener" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
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
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.acl.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.foo.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_2"
}

data "alicloud_slb_listeners" "slb_listeners" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
}

data "alicloud_slb_listeners" "slb_listeners_with_filters" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
  frontend_port = 80
  protocol = "https"
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

const testAccCheckAlicloudSlbListenersDataSourceTcp = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceTcp"
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
  name = "${var.name}"
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
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceUdp"
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
  name = "${var.name}"
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

const testAccCheckAlicloudSlbListenersDataSourceEmpty = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceUdp"
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
  name = "${var.name}"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

data "alicloud_slb_listeners" "slb_listeners" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
}
`

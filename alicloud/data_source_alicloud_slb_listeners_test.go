package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSLBListenersDataSource_http(t *testing.T) {
	rand := acctest.RandInt()
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_listener.default.load_balancer_id}"`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfig(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfig(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}-fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfig(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"80"`,
			"protocol":          `"http"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfig(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"81"`,
			"protocol":          `"http"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
	}

	var existSlbRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#":                           "1",
			"slb_listeners.0.backend_port":              "80",
			"slb_listeners.0.frontend_port":             "80",
			"slb_listeners.0.protocol":                  "http",
			"slb_listeners.0.status":                    "running",
			"slb_listeners.0.bandwidth":                 "10",
			"slb_listeners.0.scheduler":                 "wrr",
			"slb_listeners.0.sticky_session":            "on",
			"slb_listeners.0.sticky_session_type":       "insert",
			"slb_listeners.0.cookie_timeout":            "86400",
			"slb_listeners.0.health_check":              "on",
			"slb_listeners.0.health_check_uri":          "/cons",
			"slb_listeners.0.health_check_connect_port": "20",
			"slb_listeners.0.healthy_threshold":         "8",
			"slb_listeners.0.unhealthy_threshold":       "8",
			"slb_listeners.0.health_check_timeout":      "8",
			"slb_listeners.0.health_check_interval":     "5",
			"slb_listeners.0.health_check_http_code":    "http_2xx,http_3xx",
			"slb_listeners.0.gzip":                      "on",
			"slb_listeners.0.x_forwarded_for":           "on",
			"slb_listeners.0.x_forwarded_for_slb_ip":    "on",
			"slb_listeners.0.x_forwarded_for_slb_id":    "on",
			"slb_listeners.0.x_forwarded_for_slb_proto": "off",
			"slb_listeners.0.idle_timeout":              "30",
			"slb_listeners.0.request_timeout":           "80",
			"slb_listeners.0.description":               fmt.Sprintf("tf-testAccCheckAlicloudSlbListenersDataSourceHttp-%d", rand),
		}
	}

	var fakeSlbRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#": "0",
		}
	}

	var slbListenersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_listeners.default",
		existMapFunc: existSlbRecordsMapFunc,
		fakeMapFunc:  fakeSlbRecordsMapFunc,
	}

	slbListenersCheckInfo.dataSourceTestCheck(t, rand, basicConf, descriptionConf, allConf)
}

func TestAccAlicloudSLBListenersDataSource_https(t *testing.T) {
	rand := acctest.RandInt()
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigHttps(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_listener.default.load_balancer_id}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigHttps(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfigHttps(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}-fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigHttps(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"80"`,
			"protocol":          `"https"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfigHttps(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"81"`,
			"protocol":          `"https"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
	}

	var existSlbRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#":                           "1",
			"slb_listeners.0.backend_port":              "80",
			"slb_listeners.0.frontend_port":             "80",
			"slb_listeners.0.protocol":                  "https",
			"slb_listeners.0.status":                    "running",
			"slb_listeners.0.bandwidth":                 "10",
			"slb_listeners.0.scheduler":                 "wrr",
			"slb_listeners.0.sticky_session":            "on",
			"slb_listeners.0.sticky_session_type":       "insert",
			"slb_listeners.0.cookie_timeout":            "86400",
			"slb_listeners.0.health_check":              "on",
			"slb_listeners.0.health_check_uri":          "/cons",
			"slb_listeners.0.health_check_connect_port": "20",
			"slb_listeners.0.healthy_threshold":         "8",
			"slb_listeners.0.unhealthy_threshold":       "8",
			"slb_listeners.0.health_check_timeout":      "8",
			"slb_listeners.0.health_check_interval":     "5",
			"slb_listeners.0.health_check_http_code":    "http_2xx,http_3xx",
			"slb_listeners.0.gzip":                      "on",
			"slb_listeners.0.x_forwarded_for":           "on",
			"slb_listeners.0.x_forwarded_for_slb_ip":    "on",
			"slb_listeners.0.x_forwarded_for_slb_id":    "on",
			"slb_listeners.0.x_forwarded_for_slb_proto": "off",
			"slb_listeners.0.idle_timeout":              "30",
			"slb_listeners.0.request_timeout":           "80",
			"slb_listeners.0.description":               fmt.Sprintf("tf-testAccCheckAlicloudSlbListenersDataSourceHttps-%d", rand),
		}
	}

	var fakeSlbRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#": "0",
		}
	}

	var slbListenersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_listeners.default",
		existMapFunc: existSlbRecordsMapFunc,
		fakeMapFunc:  fakeSlbRecordsMapFunc,
	}

	slbListenersCheckInfo.dataSourceTestCheck(t, rand, basicConf, descriptionConf, allConf)
}

func TestAccAlicloudSLBListenersDataSource_tcp(t *testing.T) {
	rand := acctest.RandInt()
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigTcp(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_listener.default.load_balancer_id}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigTcp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfigTcp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}-fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigTcp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"22"`,
			"protocol":          `"tcp"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfigTcp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"21"`,
			"protocol":          `"tcp"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
	}

	var existSLBListenersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#":                           "1",
			"slb_listeners.0.backend_port":              "22",
			"slb_listeners.0.frontend_port":             "22",
			"slb_listeners.0.protocol":                  "tcp",
			"slb_listeners.0.status":                    "running",
			"slb_listeners.0.bandwidth":                 "10",
			"slb_listeners.0.scheduler":                 "wrr",
			"slb_listeners.0.persistence_timeout":       "0",
			"slb_listeners.0.established_timeout":       "900",
			"slb_listeners.0.health_check":              "on",
			"slb_listeners.0.health_check_type":         "tcp",
			"slb_listeners.0.health_check_connect_port": "20",
			"slb_listeners.0.healthy_threshold":         "8",
			"slb_listeners.0.unhealthy_threshold":       "8",
			"slb_listeners.0.health_check_timeout":      "0",
			"slb_listeners.0.health_check_interval":     "5",
			"slb_listeners.0.description":               fmt.Sprintf("tf-testAccCheckAlicloudSlbListenersDataSourceTcp-%d", rand),
		}
	}

	var fakeSLBListenersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#": "0",
		}
	}

	var slbListenersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_listeners.default",
		existMapFunc: existSLBListenersMapFunc,
		fakeMapFunc:  fakeSLBListenersMapFunc,
	}

	slbListenersCheckInfo.dataSourceTestCheck(t, rand, basicConf, descriptionConf, allConf)
}

func TestAccAlicloudSLBListenersDataSource_udp(t *testing.T) {
	rand := acctest.RandInt()
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigUdp(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_listener.default.load_balancer_id}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigUdp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfigUdp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"description_regex": `"${alicloud_slb_listener.default.description}-fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbListenersDataSourceConfigUdp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"11"`,
			"protocol":          `"udp"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbListenersDataSourceConfigUdp(rand, map[string]string{
			"load_balancer_id":  `"${alicloud_slb_listener.default.load_balancer_id}"`,
			"frontend_port":     `"21"`,
			"protocol":          `"udp"`,
			"description_regex": `"${alicloud_slb_listener.default.description}"`,
		}),
	}

	var existSLBListenersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#":                              "1",
			"slb_listeners.0.backend_port":                 "10",
			"slb_listeners.0.frontend_port":                "11",
			"slb_listeners.0.protocol":                     "udp",
			"slb_listeners.0.status":                       "running",
			"slb_listeners.0.bandwidth":                    "10",
			"slb_listeners.0.scheduler":                    "wrr",
			"slb_listeners.0.health_check":                 "on",
			"slb_listeners.0.health_check_connect_port":    "20",
			"slb_listeners.0.health_check_connect_timeout": "8",
			"slb_listeners.0.healthy_threshold":            "8",
			"slb_listeners.0.unhealthy_threshold":          "8",
			"slb_listeners.0.health_check_timeout":         "0",
			"slb_listeners.0.health_check_interval":        "5",
			"slb_listeners.0.description":                  fmt.Sprintf("tf-testAccCheckAlicloudSlbListenersDataSourceUdp-%d", rand),
		}
	}

	var fakeSLBListenersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_listeners.#": "0",
		}
	}

	var slbListenersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_listeners.default",
		existMapFunc: existSLBListenersMapFunc,
		fakeMapFunc:  fakeSLBListenersMapFunc,
	}

	slbListenersCheckInfo.dataSourceTestCheck(t, rand, basicConf, descriptionConf, allConf)
}

func testAccCheckAlicloudSlbListenersDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceHttp-%d"
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

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
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
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  request_timeout           = 80
  idle_timeout              = 30
  description = "${var.name}"
}

data "alicloud_slb_listeners" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudSlbListenersDataSourceConfigHttps(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceHttps-%d"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  internet_charge_type = "PayByTraffic"
  address_type = "internet"
  load_balancer_spec = "slb.s2.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
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
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  server_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_2"
  description = "${var.name}"
}

variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}

data "alicloud_slb_listeners" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudSlbListenersDataSourceConfigTcp(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceTcp-%d"
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

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
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
  description = "${var.name}"
}

data "alicloud_slb_listeners" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudSlbListenersDataSourceConfigUdp(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSlbListenersDataSourceUdp-%d"
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

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  backend_port = 10
  frontend_port = 11
  protocol = "udp"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  bandwidth = 10
  description = "${var.name}"
}

data "alicloud_slb_listeners" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

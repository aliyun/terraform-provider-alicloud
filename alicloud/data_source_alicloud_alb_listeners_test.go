package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBListenersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_listener.default.id}_fakeid"]`,
		}),
	}

	listenerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"listener_ids": `["${alicloud_alb_listener.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"listener_ids": `["${alicloud_alb_listener.default.id}_fake"]`,
		}),
	}

	loadBalancerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"listener_ids":      `["${alicloud_alb_listener.default.id}"]`,
			"load_balancer_ids": `["${alicloud_alb_listener.default.load_balancer_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"listener_ids":      `["${alicloud_alb_listener.default.id}"]`,
			"load_balancer_ids": `["${alicloud_alb_listener.default.load_balancer_id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_listener.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_listener.default.id}"]`,
			"status": `"Stopped"`,
		}),
	}

	listenerProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_listener.default.id}"]`,
			"listener_protocol": `"${alicloud_alb_listener.default.listener_protocol}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_listener.default.id}"]`,
			"listener_protocol": `"HTTP"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_listener.default.id}"]`,
			"listener_ids":      `["${alicloud_alb_listener.default.id}"]`,
			"load_balancer_ids": `["${alicloud_alb_listener.default.load_balancer_id}"]`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbListenerDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_listener.default.id}_fake"]`,
			"listener_ids":      `["${alicloud_alb_listener.default.id}_fake"]`,
			"load_balancer_ids": `["${alicloud_alb_listener.default.load_balancer_id}_fake"]`,
			"status":            `"Stopped"`,
		}),
	}

	var existDataAlicloudAlbListenersSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"listeners.#":                      "1",
			"listeners.0.status":               "Running",
			"listeners.0.listener_description": fmt.Sprintf("tf-testaccalblistener%d", rand),
		}
	}
	var fakeDataAlicloudAlbListenersSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"listeners.#": "0",
		}
	}
	var alicloudAlbListenerCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_listeners.default",
		existMapFunc: existDataAlicloudAlbListenersSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbListenersSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}

	alicloudAlbListenerCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, listenerIdsConf, loadBalancerIdsConf, statusConf, listenerProtocolConf, allConf)
}
func testAccCheckAlicloudAlbListenerDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "port" {	
	default = "%d"
}

variable "name" {	
	default = "tf-testaccalblistener%d"
}

data "alicloud_alb_zones" "default"{}

data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count             = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id =  data.alicloud_alb_zones.default.zones.0.id
  vswitch_name              = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count             = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name              = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_alb_load_balancer" "default_3" {
  vpc_id =                data.alicloud_vpcs.default.ids.0
  address_type =        "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name =    var.name
  load_balancer_edition = "Standard"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = 	"PayAsYouGo"
  }
  tags = {
		Created = "TF"
  }
  zone_mappings{
		vswitch_id =  length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
		zone_id =  data.alicloud_alb_zones.default.zones.0.id
	}
  zone_mappings{
		vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
		zone_id =   data.alicloud_alb_zones.default.zones.1.id
	}
  modification_protection_config{
	status = "NonProtection"
  }
  access_log_config{
  	log_project = alicloud_log_project.default.name
  	log_store =   alicloud_log_store.default.name
  }
}

resource "alicloud_alb_server_group" "default_4" {
	protocol = "HTTP"
	vpc_id = data.alicloud_vpcs.default.vpcs.0.id
	server_group_name = var.name
    resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	health_check_config {
       health_check_enabled = "false"
	}
	sticky_session_config {
       sticky_session_enabled = "false"
	}
	tags = {
		Created = "TF"
	}
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP
MA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0
ZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow
djELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE
ChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG
9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ
AoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB
coG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook
KOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw
HQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy
+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC
QkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN
MAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ
AJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT
cQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1
Ofi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd
DUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV
kg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM
ywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB
AoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd
6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP
hwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4
MdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz
71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm
Ev9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE
qygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8
9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM
zWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe
DrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_alb_listener" "default" {
	load_balancer_id = alicloud_alb_load_balancer.default_3.id
	listener_protocol =  "HTTPS"
	listener_port = var.port
	listener_description = var.name
	default_actions{
		type = "ForwardGroup"
		forward_group_config{
			server_group_tuples{
				server_group_id = alicloud_alb_server_group.default_4.id
			}
		}
	}
	certificates{
		certificate_id = join("",[alicloud_ssl_certificates_service_certificate.default.id,"-%s"])
	}
}

data "alicloud_alb_listeners" "default" {	
	%s
}
`, rand, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}

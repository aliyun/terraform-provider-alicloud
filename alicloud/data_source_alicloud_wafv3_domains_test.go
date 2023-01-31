package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudWafv3DomainDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_domain.default.id}_fake"]`,
		}),
	}
	DomainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_wafv3_domain.default.id}"]`,
			"domain": `"${alicloud_wafv3_domain.default.domain}"`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_wafv3_domain.default.id}_fake"]`,
			"domain": `"${alicloud_wafv3_domain.default.domain}_fake"`,
		}),
	}
	backendConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_wafv3_domain.default.id}"]`,
			"backend": `"1.1.1.1"`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_wafv3_domain.default.id}_fake"]`,
			"backend": `"1.1.1.2"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_wafv3_domain.default.id}"]`,
			"domain":  `"${alicloud_wafv3_domain.default.domain}"`,
			"backend": `"1.1.1.1"`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_wafv3_domain.default.id}_fake"]`,
			"domain":  `"${alicloud_wafv3_domain.default.domain}_fake"`,
			"backend": `"1.1.1.2"`,
		}),
	}

	Wafv3DomainCheckInfo.dataSourceTestCheck(t, rand, idsConf, DomainConf, backendConf, allConf)
}

var existWafv3DomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                   "1",
		"domains.#":                               "1",
		"domains.0.id":                            CHECKSET,
		"domains.0.listen.#":                      "1",
		"domains.0.listen.0.https_ports.#":        "1",
		"domains.0.listen.0.http_ports.#":         "1",
		"domains.0.listen.0.cert_id":              CHECKSET,
		"domains.0.listen.0.cipher_suite":         "99",
		"domains.0.listen.0.xff_header_mode":      "2",
		"domains.0.listen.0.protection_resource":  "share",
		"domains.0.listen.0.tls_version":          "tlsv1",
		"domains.0.listen.0.enable_tlsv3":         "true",
		"domains.0.listen.0.http2_enabled":        "true",
		"domains.0.listen.0.custom_ciphers.#":     "1",
		"domains.0.listen.0.focus_https":          "false",
		"domains.0.listen.0.exclusive_ip":         "false",
		"domains.0.listen.0.xff_headers.#":        "2",
		"domains.0.listen.0.ipv6_enabled":         "false",
		"domains.0.redirect.#":                    "1",
		"domains.0.redirect.0.backends.#":         "1",
		"domains.0.redirect.0.loadbalance":        "iphash",
		"domains.0.redirect.0.sni_enabled":        "true",
		"domains.0.redirect.0.sni_host":           "www.aliyun.com",
		"domains.0.redirect.0.focus_http_backend": "false",
		"domains.0.redirect.0.keepalive":          "true",
		"domains.0.redirect.0.retry":              "true",
		"domains.0.redirect.0.keepalive_requests": "80",
		"domains.0.redirect.0.keepalive_timeout":  "30",
		"domains.0.redirect.0.connect_timeout":    "30",
		"domains.0.redirect.0.read_timeout":       "30",
		"domains.0.redirect.0.write_timeout":      "30",
		"domains.0.redirect.0.request_headers.#":  "1",
	}
}

var fakeWafv3DomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"domains.#": "0",
	}
}

var Wafv3DomainCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_wafv3_domains.default",
	existMapFunc: existWafv3DomainMapFunc,
	fakeMapFunc:  fakeWafv3DomainMapFunc,
}

func testAccCheckAlicloudWafv3DomainSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tftest%d.tftest.top"
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

locals {
  certificate_id       = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "cn-hangzhou"])
}

data "alicloud_wafv3_instances" "default" {}

resource "alicloud_wafv3_domain" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    https_ports = [443]
    http_ports = [80]
    cert_id             = local.certificate_id
    cipher_suite        = 99
    xff_header_mode     = 2
    protection_resource = "share"
    tls_version         = "tlsv1"
    enable_tlsv3        = true
    http2_enabled       = true
    custom_ciphers = ["ECDHE-ECDSA-AES128-GCM-SHA256"]
    focus_https  = false
    ipv6_enabled = false
    exclusive_ip = false
	xff_headers = ["A", "B"]
  }
  redirect {
    backends = ["1.1.1.1"]
    loadbalance        = "iphash"
    sni_enabled        = true
    sni_host           = "www.aliyun.com"
    focus_http_backend = false
    keepalive          = true
    retry              = true
    keepalive_requests = 80
    keepalive_timeout  = 30
    connect_timeout    = 30
    read_timeout       = 30
    write_timeout      = 30
    request_headers {
      key   = "A"
      value = "B"
    }
  }
  domain      = var.name
  access_type = "share"
}

data "alicloud_wafv3_domains" "default" {
instance_id = data.alicloud_wafv3_instances.default.ids.0
enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

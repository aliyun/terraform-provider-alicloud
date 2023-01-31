package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudWafv3Domain_basic2308(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DomainMap2308)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftest%d.tftest.top", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DomainBasicDependence2308)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_wafv3_instances.default.ids.0}",
					"listen": []map[string]interface{}{
						{
							"https_ports":         []string{"443"},
							"http_ports":          []string{"80"},
							"cert_id":             "${local.certificate_id}",
							"cipher_suite":        "99",
							"xff_header_mode":     "2",
							"protection_resource": "share",
							"tls_version":         "tlsv1",
							"enable_tlsv3":        "true",
							"http2_enabled":       "true",
							"custom_ciphers":      []string{"ECDHE-ECDSA-AES128-GCM-SHA256"},
							"focus_https":         "false",
							"exclusive_ip":        "false",
							"xff_headers":         []string{"A", "B"},
							"ipv6_enabled":        "false",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends":           []string{"1.1.1.1"},
							"loadbalance":        "iphash",
							"sni_enabled":        "true",
							"sni_host":           "www.aliyun.com",
							"focus_http_backend": "false",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "1000",
							"keepalive_timeout":  "15",
							"connect_timeout":    "5",
							"read_timeout":       "30",
							"write_timeout":      "30",
							"request_headers": []map[string]interface{}{
								{
									"key":   "A",
									"value": "B",
								},
							},
						},
					},
					"domain":      "${var.name}",
					"access_type": "share",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":                   CHECKSET,
						"domain":                        name,
						"status":                        CHECKSET,
						"access_type":                   "share",
						"listen.#":                      "1",
						"listen.0.https_ports.#":        "1",
						"listen.0.http_ports.#":         "1",
						"listen.0.cert_id":              CHECKSET,
						"listen.0.cipher_suite":         "99",
						"listen.0.xff_header_mode":      "2",
						"listen.0.protection_resource":  "share",
						"listen.0.tls_version":          "tlsv1",
						"listen.0.enable_tlsv3":         "true",
						"listen.0.http2_enabled":        "true",
						"listen.0.custom_ciphers.#":     "1",
						"listen.0.focus_https":          "false",
						"listen.0.exclusive_ip":         "false",
						"listen.0.xff_headers.#":        "2",
						"listen.0.ipv6_enabled":         "false",
						"redirect.#":                    "1",
						"redirect.0.backends.#":         "1",
						"redirect.0.loadbalance":        "iphash",
						"redirect.0.sni_enabled":        "true",
						"redirect.0.sni_host":           "www.aliyun.com",
						"redirect.0.focus_http_backend": "false",
						"redirect.0.keepalive":          "true",
						"redirect.0.retry":              "true",
						"redirect.0.keepalive_requests": "1000",
						"redirect.0.keepalive_timeout":  "15",
						"redirect.0.connect_timeout":    "5",
						"redirect.0.read_timeout":       "30",
						"redirect.0.write_timeout":      "30",
						"redirect.0.request_headers.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listen": []map[string]interface{}{
						{
							"https_ports":         []string{"443"},
							"http_ports":          []string{"80", "88"},
							"cert_id":             "${local.certificate_id}",
							"cipher_suite":        "99",
							"xff_header_mode":     "2",
							"protection_resource": "share",
							"tls_version":         "tlsv1",
							"enable_tlsv3":        "true",
							"http2_enabled":       "true",
							"custom_ciphers":      []string{"ECDHE-ECDSA-AES128-GCM-SHA256"},
							"focus_https":         "false",
							"exclusive_ip":        "false",
							"xff_headers":         []string{"A", "B", "C"},
							"ipv6_enabled":        "false",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends":           []string{"1.1.1.1"},
							"loadbalance":        "iphash",
							"sni_enabled":        "true",
							"sni_host":           "www.aliyun.com",
							"focus_http_backend": "false",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "1000",
							"keepalive_timeout":  "15",
							"connect_timeout":    "5",
							"read_timeout":       "60",
							"write_timeout":      "60",
							"request_headers": []map[string]interface{}{
								{
									"key":   "A",
									"value": "B",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listen.#":                      "1",
						"listen.0.https_ports.#":        "1",
						"listen.0.http_ports.#":         "2",
						"listen.0.cert_id":              CHECKSET,
						"listen.0.cipher_suite":         "99",
						"listen.0.xff_header_mode":      "2",
						"listen.0.protection_resource":  "share",
						"listen.0.tls_version":          "tlsv1",
						"listen.0.enable_tlsv3":         "true",
						"listen.0.http2_enabled":        "true",
						"listen.0.custom_ciphers.#":     "1",
						"listen.0.focus_https":          "false",
						"listen.0.exclusive_ip":         "false",
						"listen.0.xff_headers.#":        "3",
						"listen.0.ipv6_enabled":         "false",
						"redirect.#":                    "1",
						"redirect.0.backends.#":         "1",
						"redirect.0.loadbalance":        "iphash",
						"redirect.0.sni_enabled":        "true",
						"redirect.0.sni_host":           "www.aliyun.com",
						"redirect.0.focus_http_backend": "false",
						"redirect.0.keepalive":          "true",
						"redirect.0.retry":              "true",
						"redirect.0.keepalive_requests": "1000",
						"redirect.0.keepalive_timeout":  "15",
						"redirect.0.connect_timeout":    "5",
						"redirect.0.read_timeout":       "60",
						"redirect.0.write_timeout":      "60",
						"redirect.0.request_headers.#":  "1",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_type"},
			},
		},
	})
}

var AlicloudWafv3DomainMap2308 = map[string]string{}

func AlicloudWafv3DomainBasicDependence2308(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
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

data "alicloud_wafv3_instances" "default" {}

locals {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "cn-hangzhou"])
}


`, name)
}

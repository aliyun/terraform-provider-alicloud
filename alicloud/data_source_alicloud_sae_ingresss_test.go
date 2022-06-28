package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAEIngressesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 100)
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeIngressesDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_sae_ingress.default.id}"]`,
			"namespace_id": `"${alicloud_sae_namespace.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeIngressesDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_sae_ingress.default.id}_fake"]`,
			"namespace_id": `"${alicloud_sae_namespace.default.id}"`,
		}),
	}

	var existAlicloudSaeIngressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ingresses.#":              "1",
			"ingresses.0.namespace_id": fmt.Sprintf("%s:tftestacc%d", defaultRegionToTest, rand),
		}
	}
	var fakeAlicloudSaeIngressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudSaeIngressesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_ingresses.default",
		existMapFunc: existAlicloudSaeIngressesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeIngressesDataSourceNameMapFunc,
	}
	alicloudSaeIngressesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}
func testAccCheckAlicloudSaeIngressesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {
  default = "tftestacc%d"
}

variable "namespace_id" {
  default = "%s:tftestacc%d"
}

resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/sae-demo-image/consumer:1.0"
  package_type    = "Image"
  vswitch_id      = data.alicloud_vswitches.default.vswitches.0.id
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Beijing"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = var.namespace_id
  namespace_name        = var.name
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_slb" "default" {
  load_balancer_name = var.name
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = data.alicloud_vswitches.default.vswitches.0.id
}

resource "alicloud_slb_server_certificate" "default" {
  name               = "slbservercertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}

resource "alicloud_sae_ingress" "default" {
  slb_id        = alicloud_slb.default.id
  namespace_id  = alicloud_sae_namespace.default.id
  listener_port = "80"
  rules {
    app_id         = alicloud_sae_application.default.id
    container_port = "443"
    domain         = "www.alicloud.com"
    app_name       = alicloud_sae_application.default.app_name
    path           = "/"
  }
  cert_id = alicloud_slb_server_certificate.default.id
  default_rule {
    app_id         = alicloud_sae_application.default.id
    container_port = "443"
  }
}

data "alicloud_sae_ingresses" "default" {	
	%s
}
`, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}

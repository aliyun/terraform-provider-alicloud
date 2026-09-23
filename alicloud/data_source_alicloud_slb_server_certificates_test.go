package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSLBServerCertificatesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_server_certificate.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_server_certificate.default.name}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_server_certificate.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_server_certificate.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_slb_server_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_slb_server_certificate.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_slb_server_certificate.default.id}"]`,
			"name_regex": `"${alicloud_slb_server_certificate.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_slb_server_certificate.default.id}_fake"]`,
			"name_regex": `"${alicloud_slb_server_certificate.default.name}"`,
		}),
	}

	var existSLBServerCertificatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#":                         "1",
			"ids.#":                                  "1",
			"names.#":                                "1",
			"certificates.0.id":                      CHECKSET,
			"certificates.0.name":                    fmt.Sprintf("tf-testAccSlbServerCertificatesDataSourceBasic-%d", rand),
			"certificates.0.fingerprint":             CHECKSET,
			"certificates.0.common_name":             CHECKSET,
			"certificates.0.expired_timestamp":       CHECKSET,
			"certificates.0.created_time":            CHECKSET,
			"certificates.0.created_timestamp":       CHECKSET,
			"certificates.0.is_alicloud_certificate": CHECKSET,
			"certificates.0.tags.%":                  "2",
			"certificates.0.tags.Created":            "TF",
			"certificates.0.tags.For":                "acceptance test",
		}
	}

	var fakeSLBServerCertificatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#": "0",
			"ids.#":          "0",
			"names.#":        "0",
		}
	}

	var slbServerCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_server_certificates.default",
		existMapFunc: existSLBServerCertificatesMapFunc,
		fakeMapFunc:  fakeSLBServerCertificatesMapFunc,
	}

	slbServerCertificatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, tagsConf, idsConf, allConf)

}

func testAccCheckAlicloudSlbServerCertificatesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbServerCertificatesDataSourceBasic-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  tags = {
      Created = "TF"
       For    = "acceptance test"
    }
}

data "alicloud_slb_server_certificates" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

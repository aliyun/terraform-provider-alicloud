package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSslCertificatesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ssl_certificates_service_certificate.default.certificate_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_certificates_service_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_certificates_service_certificate.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ssl_certificates_service_certificate.default.certificate_name}"`,
			"ids":        `["${alicloud_ssl_certificates_service_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake"`,
			"ids":        `["${alicloud_ssl_certificates_service_certificate.default.id}_fake"]`,
		}),
	}

	var existSslCertificateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#":                  "1",
			"ids.#":                           "1",
			"names.#":                         "1",
			"certificates.0.id":               CHECKSET,
			"certificates.0.certificate_name": fmt.Sprintf("tf-testAccSslCertificatesDataSourceBasic-%d", rand),
			"certificates.0.fingerprint":      CHECKSET,
			"certificates.0.start_date":       CHECKSET,
		}
	}

	var fakeSslCertificateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#": "0",
			"ids.#":          "0",
			"names.#":        "0",
		}
	}

	var slbServerCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ssl_certificates_service_certificates.default",
		existMapFunc: existSslCertificateMapFunc,
		fakeMapFunc:  fakeSslCertificateMapFunc,
	}

	slbServerCertificatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)

}

func testAccCheckAlicloudSslCertificatesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslCertificatesDataSourceBasic-%d"
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

data "alicloud_ssl_certificates_service_certificates" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

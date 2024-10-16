package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudSslCertificatesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ssl_certificates_service_certificates.default"
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAliCloudSslCertificatesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
		}),
	}

	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword": "${alicloud_ssl_certificates_service_certificate.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword": "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
			"keyword":    "${alicloud_ssl_certificates_service_certificate.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
			"keyword":    "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
		}),
	}

	var existAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"certificates.#":                  "1",
			"certificates.0.id":               CHECKSET,
			"certificates.0.cert_id":          CHECKSET,
			"certificates.0.certificate_name": CHECKSET,
			"certificates.0.fingerprint":      CHECKSET,
			"certificates.0.common":           CHECKSET,
			"certificates.0.sans":             CHECKSET,
			"certificates.0.org_name":         CHECKSET,
			"certificates.0.issuer":           CHECKSET,
			"certificates.0.country":          CHECKSET,
			"certificates.0.province":         CHECKSET,
			"certificates.0.city":             CHECKSET,
			"certificates.0.expired":          CHECKSET,
			"certificates.0.start_date":       CHECKSET,
			"certificates.0.end_date":         CHECKSET,
			"certificates.0.cert":             "",
			"certificates.0.key":              "",
			"certificates.0.name":             CHECKSET,
		}
	}

	var fakeAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"certificates.#": "0",
		}
	}

	var aliCloudSslCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ssl_certificates_service_certificates.default",
		existMapFunc: existAliCloudSslCertificatesDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudSslCertificatesDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudSslCertificatesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keyWordConf, allConf)
}

func TestAccAliCloudSslCertificatesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ssl_certificates_service_certificates.default"
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAliCloudSslCertificatesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
			"enable_details": "true",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
			"enable_details": "true",
		}),
	}

	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
			"enable_details": "true",
		}),
	}

	var existAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"certificates.#":                  "1",
			"certificates.0.id":               CHECKSET,
			"certificates.0.cert_id":          CHECKSET,
			"certificates.0.certificate_name": CHECKSET,
			"certificates.0.fingerprint":      CHECKSET,
			"certificates.0.common":           CHECKSET,
			"certificates.0.sans":             CHECKSET,
			"certificates.0.org_name":         CHECKSET,
			"certificates.0.issuer":           CHECKSET,
			"certificates.0.country":          CHECKSET,
			"certificates.0.province":         CHECKSET,
			"certificates.0.city":             CHECKSET,
			"certificates.0.expired":          CHECKSET,
			"certificates.0.start_date":       CHECKSET,
			"certificates.0.end_date":         CHECKSET,
			"certificates.0.cert":             CHECKSET,
			"certificates.0.key":              CHECKSET,
			"certificates.0.buy_in_aliyun":    CHECKSET,
			"certificates.0.name":             CHECKSET,
		}
	}

	var fakeAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"certificates.#": "0",
		}
	}

	var aliCloudSslCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ssl_certificates_service_certificates.default",
		existMapFunc: existAliCloudSslCertificatesDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudSslCertificatesDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudSslCertificatesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keyWordConf, allConf)
}

func dataSourceAliCloudSslCertificatesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
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
  		key              = <<EOF
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
`, name)
}

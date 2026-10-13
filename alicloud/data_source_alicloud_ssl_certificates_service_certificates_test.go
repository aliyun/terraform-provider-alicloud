package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudSslCertificatesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
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
MIIDUDCCAjigAwIBAgIJANNtoCzB6ORlMA0GCSqGSIb3DQEBCwUAMFQxCzAJBgNV
BAYTAkNOMRAwDgYDVQQIDAdCZWlKaW5nMRAwDgYDVQQHDAdCZWlKaW5nMQswCQYD
VQQKDAJBQTEUMBIGA1UEAwwLKi50ZnRlc3R0b3AwHhcNMjYwNDMwMDk1NTE2WhcN
MzYwNDI3MDk1NTE2WjBUMQswCQYDVQQGEwJDTjEQMA4GA1UECAwHQmVpSmluZzEQ
MA4GA1UEBwwHQmVpSmluZzELMAkGA1UECgwCQUExFDASBgNVBAMMCyoudGZ0ZXN0
dG9wMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArTqlarQgui1IdKGS
wGHImal5eC3WSwrp1BNaktp+OtDE1J+uLTGtzk5v/XSjRooKZvb+HiBZw4Kq2Im/
KZAnbdcQk53f5vSJGrQRgW7fPLBAoazomADG0sQijBOrLFB8ZbOfz9HmaHT5z6yG
NBTcbjiqQeBnoI1qLZU3FkC4OnjfVMn8ORAZH3rhehFUFd/4oLiWPeTezappCoiZ
VbOea6iM9aqaPIxiqDDNgeJ5Azdowifnj3WhHR9Jtul3WMyPll7SNNkJZWEbfYkp
0aimLlchpPd1HsX1sSuBvJtShnVPzZsJHE4vHPrZ7aW3e+Wo2eAgnZPFK2JQ5LM8
D0XjqQIDAQABoyUwIzAhBgNVHREEGjAYggsqLnRmdGVzdHRvcIIJdGZ0ZXN0dG9w
MA0GCSqGSIb3DQEBCwUAA4IBAQCCS9leIZhucIekuYjA++MXIxDwBBkXGtfYgsIM
kBDdbs7jAg0slxlVl3yVAwhPFR+Nz1lqCgEU5MbK+VsOPGDFNRqUSaVxdQU51719
dWY7o+oaiueXHMXEDqAEByjwVzxQogGGivVqcajuUcrJM8zgrJlUe8/9If6Hj4EV
TccrAdaxCIVaHRXLm49Mji4zwu32Mzce4DMyNO+0dI4OuR016UzpN1H8Gt8y/vAA
G3bcacQLw8uXg6CX6+UTePlDlRUrKzolCas6k3ReQNrthXDL6WTI8XWkr1ZEf9vF
5Vx6wpxA+4GBbaEmDEdrXgT2MuUBirV6Jgaok5XDMRfJP8ZU
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCtOqVqtCC6LUh0
oZLAYciZqXl4LdZLCunUE1qS2n460MTUn64tMa3OTm/9dKNGigpm9v4eIFnDgqrY
ib8pkCdt1xCTnd/m9IkatBGBbt88sEChrOiYAMbSxCKME6ssUHxls5/P0eZodPnP
rIY0FNxuOKpB4GegjWotlTcWQLg6eN9Uyfw5EBkfeuF6EVQV3/iguJY95N7NqmkK
iJlVs55rqIz1qpo8jGKoMM2B4nkDN2jCJ+ePdaEdH0m26XdYzI+WXtI02QllYRt9
iSnRqKYuVyGk93UexfWxK4G8m1KGdU/NmwkcTi8c+tntpbd75ajZ4CCdk8UrYlDk
szwPReOpAgMBAAECggEALH5QuHx1n3w8DQDnEZTMEnPOKEkVroaqEne6HtgR7gdr
VvWApkkJTGEnMKGBwn94NoQgNxq2E5p/SUjwGbvV9Md1kYTOWqLaiEpR6L3ShuL6
ZsxjnY8tgjswaGww3hv6J6EgXh+an5/T/jxg6AWnrIPPb0Y1N6664+RcltGewB2Z
hp9PugmBAvhNpoHbjVSm1KEqQONmIs9VzlSQv/GGDKzP0ILQNXUQ9lB8DjOxLIgA
Qh4GjxO181eS052n1/ypIPitSANbE5hAOngOmn0+kAImcAHwg7CYXun1ECCpoDs1
p+TWMbaTODbYbUMNMLhk2o3WiE8I0f4Gj5uK8loA1QKBgQDVqiU9J++H7FnOc6kG
jDt6ZrbYVXHskkFN6vjdCKBAqEUW2i9iiMJPfjKsC0BOYK8dzaMxs/Kg49OBysCe
IVdowFGxmaOQf8uh5CBlVr9bQ8yVfMuAMwne6XAnXpsrcATHEgCrzlYYwpnzYrwd
FeRVr/eX4Kjb5t/uHA41ows07wKBgQDPjXYgupLfFsqfELIXnyTVMrjdpZGLd4lm
GttCJBHztpm7aAPGxvKRZtWBAZhGhfvg5CZovCnZkE87C71dcT2/E3cMfZLVi8og
6L8HBeLnhs9fhLU/pPwCGr+gec4updNAnSYM92Rmu1CSh+D5IkoOmxe1bsTnZYdZ
vY5m66jg5wKBgAX2gqWLSMcNVJBLcAsrvLk0xqOQ9uX4SvRGu31HsVk8mPkDS9E/
KdUYG4frpHbLgfed9pD7iajt0shLhUakfZEB1QXU2Ni+iEtTV4gcfKiqYrpFSlSg
mATtlOC0ZjY3IRsBKJ2i24pDXBKLzd4t7zpo98r4TR4d+l/lzou2qihnAoGBAMwB
bsfyu2ReQaEhxvti6NBJ/92U3T9pqnrbEQKs85xKska7kWKzsD9tBZS9HZDMJ2zA
tGQj9zqrFsWh32rWYOvMVSrIzyfpNC8utn1Pst0B7BkmFwVMxIvJAOvabef6BHAx
RvJdxKkZmrO8rUwUARjkJuSl0RLk/6ocoLjXD/KnAoGAaw7y0oIvSxzwioZL2L2+
RGQ+TRGnWwGDNLljO8jiNXO+hr1Fxvkk4MuYMP0FUNBCRQ4YZi3bY/Nn3QiWDuJe
TCr1M75G/sSQV6EwJa4xzSzAs8WQvFQiwfAh3gJlOsl6eTA/1QbSw8KUIdHOob6U
000PZ6e8JW7OGMRlRPMHlkw=
-----END PRIVATE KEY-----
EOF
	}
`, name)
}

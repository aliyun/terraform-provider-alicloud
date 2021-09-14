package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudWAFCertificate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_waf_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudWAFCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Waf_openapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafcertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWAFCertificateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_name": name,
					"instance_id":      "${alicloud_waf_domain.domain.instance_id}",
					"domain":           "${alicloud_waf_domain.domain.domain_name}",
					"private_key":      `-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA5KeYEdE3moKDwwB4DV+yB44BNOqJF6KOS3wSF0IhjnWRTJ13\nUoCTN0CwIJTPzgu7hPXuvoRe2Mgn/5CclHaF3x1+FIGVskyC1hm5I29nuP21MDkN\njJ0e9TLzSLfKNmSiJKbUqeLa6l46U/5rUdI3a+qOADXNIV4chZlgYXng+HykXfDK\nXDdsuGE8h9Ue2+1WLM7KlWpbFlG8JCTTxIKhuKqORvdZRPE0u20tVqELtSjrglkR\n6yc1ZXZ3MIpS9L6QynysxnCJ/CpYKxlR+SkVR5Uozj3+J4tPtYy7kQR7AsWAIxnr\n5JCAW96yFjK+LSdtiJm5fvcmXPJt8Byui67sJwIDAQABAoIBAHkIBUSZG/ebUids\nHh/mIZRCJH1gEAUEtJFcMS9CgASqUqjdSQXQrRw4D4KPgpesOiwdCayq4JGbESED\n6P3k5uO89rngBd8FBNRTdV4+YAbZtIi8CmH94R8lQwWg01aLWhBQho0OWZZpHjLL\nnv8GaPOkilmuINCiwjIth4jRPLYKM8rQUL0Y9DSLW1dZfxUMj5yPHdpV81/es21k\ns3ZUETOsygAfv9ybhxFW/1h8vldH9bCRvVq4oXHXkpF9Wp9K6z4UjR60k7hnxzTw\nUgFoEheHz58LOJNzh+dad/2HCCEVutQXMIMomoMvSmQkODhOMYsiJC8L15l6/193\ngj6Ah1ECgYEA8AE/efv7Q8DGwVRP51a4whzkn2Mws/1cURxXLflYNdrOek3Vmdda\nR4k9N6xYKTjQ4xSV+yYHpCMBdSkFcKFOOMf+DQLB5NAb3C4E+K6ea8GrfIUNnFZM\nYKwfjhDEvNFky8FqLscQGr9HCzcx+lOrwtnMK/srvS8r9IuNfXcalaUCgYEA8+Sy\n4brUb9SNqAFY68lOsovuhQpCK8UnomgbVZG52wziA368HYfdmLFSzSAMaOYA4wog\nUvC1JtOqkP6tyAFdo8qWgrGtmBpVmPSvWA2ser0zqugvRyg4uY6vd3Hwj9VEGguB\nvotYxUhzetzyapTZazzRuE7aUB9dpbnvdDhoyNsCgYBgpbABGFajfwLKoYAXwZVf\nHbF2+cOIB5PgWdBFhC5gaX9SQvusGsuRGRPc0nsiBm4fs4999l+HWk1g96boJxzP\nwsFsTdr0oFVHwgRgfDjxXakH2LCVby8MkuOWGuyOuKelYXq34ZN7oeEjBBQNIRp+\nuaP9ZgTEBzXlITlV99ttvQKBgQCEJijeslKk/XB71a8OoxbnG3bz0ykjekINd4dI\nCBvCGurjpenbvmBNedc5meHffLCCVFxLVG1zSkEjKzuSIVtRglWnHwrGXY5/wCS5\n+z63iojSU9g6IAsMT8m3WJ9V7+JkklOMoQhKbQVVTx5yrZBY0K15xg/4VeZyA5tB\nR8dO/QKBgGqkJ1AB+qi9Tl10ic8IX2blqPt3FU6MkkVmmDl8vA5R35DBuZAD4VTs\nvsc0Y79mSqP9XL3KRAfA04tbGme30gJWz4NOoNsaaF2T2fv0gNNnVma87unFS6Y4\nFv64CkXzShjd16ov4eZetsIAZYn/bVn8zp61I6V50iT6AjpX1ptX\n-----END RSA PRIVATE KEY-----`,
					"certificate":      `-----BEGIN CERTIFICATE-----\nMIID8DCCAtigAwIBAgIQTvMBGm/PRXSj352aOU7GSjANBgkqhkiG9w0BAQsFADBe\nMQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl\nc3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe\nFw0yMTA5MDcxNDM5NTBaFw0yMjA5MDcxNDM5NTBaMC0xCzAJBgNVBAYTAkNOMR4w\nHAYDVQQDExV0Zi10ZXN0YWNjLndhZnFhMy5jb20wggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQDkp5gR0TeagoPDAHgNX7IHjgE06okXoo5LfBIXQiGOdZFM\nnXdSgJM3QLAglM/OC7uE9e6+hF7YyCf/kJyUdoXfHX4UgZWyTILWGbkjb2e4/bUw\nOQ2MnR71MvNIt8o2ZKIkptSp4trqXjpT/mtR0jdr6o4ANc0hXhyFmWBheeD4fKRd\n8MpcN2y4YTyH1R7b7VYszsqValsWUbwkJNPEgqG4qo5G91lE8TS7bS1WoQu1KOuC\nWRHrJzVldncwilL0vpDKfKzGcIn8KlgrGVH5KRVHlSjOPf4ni0+1jLuRBHsCxYAj\nGevkkIBb3rIWMr4tJ22Imbl+9yZc8m3wHK6LruwnAgMBAAGjgdowgdcwDgYDVR0P\nAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME\nGDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB\nBQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov\nL2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MCAGA1UdEQQZMBeCFXRmLXRl\nc3RhY2Mud2FmcWEzLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAPjbt2H1HmEc8DzyD\npi4IF1CvaNlYgKjPojYlt/gpj2n0MfntL8Ihly3e2fdSMEeVeTnFWFd34L4uZxMa\nxE/hx6VJWfNdnYW7FGCZr0rGj/KrtAox83H1dRrZ63hynpgCMIbc5lhA7wDe0R16\nP/1l3c50ZEmidicGhK/qmzsSQIVXC0kJf6hDQCxW6LVaDrmT8mvbhRh4ZNb2pgJ5\nQIWJHnlOmZkUVsR5cMBGzK2ModADjHXHmeoHHr3Tw7mPioE4Xh5EmMTXPLG22BKN\nRBFG9gSFri+3RxqdXwi1ZJajO3Nup5mcdGaHJbUbNUf16YKIq50PJlrVxzCZV31f\n7cOGfw==\n-----END CERTIFICATE-----`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_name": name,
						"instance_id":      CHECKSET,
						"domain":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
		},
	})
}

func TestAccAlicloudWAFCertificate_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_waf_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudWAFCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Waf_openapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafcertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWAFCertificateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_waf_domain.domain.instance_id}",
					"domain":             "${alicloud_waf_domain.domain.domain_name}",
					"certificate_id": "${alicloud_ssl_certificates_service_certificate.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"domain":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
		},
	})
}

var AlicloudWAFCertificateMap0 = map[string]string{}

func AlicloudWAFCertificateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_waf_instances" "default" {
  status          = "1"
  instance_source = "waf-cloud"
}

data "alicloud_waf_domains" "default" {
  instance_id = data.alicloud_waf_instances.default.ids.0
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = "tf-testaccSslCertificate"
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIID8DCCAtigAwIBAgIQTvMBGm/PRXSj352aOU7GSjANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yMTA5MDcxNDM5NTBaFw0yMjA5MDcxNDM5NTBaMC0xCzAJBgNVBAYTAkNOMR4w
HAYDVQQDExV0Zi10ZXN0YWNjLndhZnFhMy5jb20wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQDkp5gR0TeagoPDAHgNX7IHjgE06okXoo5LfBIXQiGOdZFM
nXdSgJM3QLAglM/OC7uE9e6+hF7YyCf/kJyUdoXfHX4UgZWyTILWGbkjb2e4/bUw
OQ2MnR71MvNIt8o2ZKIkptSp4trqXjpT/mtR0jdr6o4ANc0hXhyFmWBheeD4fKRd
8MpcN2y4YTyH1R7b7VYszsqValsWUbwkJNPEgqG4qo5G91lE8TS7bS1WoQu1KOuC
WRHrJzVldncwilL0vpDKfKzGcIn8KlgrGVH5KRVHlSjOPf4ni0+1jLuRBHsCxYAj
GevkkIBb3rIWMr4tJ22Imbl+9yZc8m3wHK6LruwnAgMBAAGjgdowgdcwDgYDVR0P
AQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME
GDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB
BQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov
L2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MCAGA1UdEQQZMBeCFXRmLXRl
c3RhY2Mud2FmcWEzLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAPjbt2H1HmEc8DzyD
pi4IF1CvaNlYgKjPojYlt/gpj2n0MfntL8Ihly3e2fdSMEeVeTnFWFd34L4uZxMa
xE/hx6VJWfNdnYW7FGCZr0rGj/KrtAox83H1dRrZ63hynpgCMIbc5lhA7wDe0R16
P/1l3c50ZEmidicGhK/qmzsSQIVXC0kJf6hDQCxW6LVaDrmT8mvbhRh4ZNb2pgJ5
QIWJHnlOmZkUVsR5cMBGzK2ModADjHXHmeoHHr3Tw7mPioE4Xh5EmMTXPLG22BKN
RBFG9gSFri+3RxqdXwi1ZJajO3Nup5mcdGaHJbUbNUf16YKIq50PJlrVxzCZV31f
7cOGfw==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA5KeYEdE3moKDwwB4DV+yB44BNOqJF6KOS3wSF0IhjnWRTJ13
UoCTN0CwIJTPzgu7hPXuvoRe2Mgn/5CclHaF3x1+FIGVskyC1hm5I29nuP21MDkN
jJ0e9TLzSLfKNmSiJKbUqeLa6l46U/5rUdI3a+qOADXNIV4chZlgYXng+HykXfDK
XDdsuGE8h9Ue2+1WLM7KlWpbFlG8JCTTxIKhuKqORvdZRPE0u20tVqELtSjrglkR
6yc1ZXZ3MIpS9L6QynysxnCJ/CpYKxlR+SkVR5Uozj3+J4tPtYy7kQR7AsWAIxnr
5JCAW96yFjK+LSdtiJm5fvcmXPJt8Byui67sJwIDAQABAoIBAHkIBUSZG/ebUids
Hh/mIZRCJH1gEAUEtJFcMS9CgASqUqjdSQXQrRw4D4KPgpesOiwdCayq4JGbESED
6P3k5uO89rngBd8FBNRTdV4+YAbZtIi8CmH94R8lQwWg01aLWhBQho0OWZZpHjLL
nv8GaPOkilmuINCiwjIth4jRPLYKM8rQUL0Y9DSLW1dZfxUMj5yPHdpV81/es21k
s3ZUETOsygAfv9ybhxFW/1h8vldH9bCRvVq4oXHXkpF9Wp9K6z4UjR60k7hnxzTw
UgFoEheHz58LOJNzh+dad/2HCCEVutQXMIMomoMvSmQkODhOMYsiJC8L15l6/193
gj6Ah1ECgYEA8AE/efv7Q8DGwVRP51a4whzkn2Mws/1cURxXLflYNdrOek3Vmdda
R4k9N6xYKTjQ4xSV+yYHpCMBdSkFcKFOOMf+DQLB5NAb3C4E+K6ea8GrfIUNnFZM
YKwfjhDEvNFky8FqLscQGr9HCzcx+lOrwtnMK/srvS8r9IuNfXcalaUCgYEA8+Sy
4brUb9SNqAFY68lOsovuhQpCK8UnomgbVZG52wziA368HYfdmLFSzSAMaOYA4wog
UvC1JtOqkP6tyAFdo8qWgrGtmBpVmPSvWA2ser0zqugvRyg4uY6vd3Hwj9VEGguB
votYxUhzetzyapTZazzRuE7aUB9dpbnvdDhoyNsCgYBgpbABGFajfwLKoYAXwZVf
HbF2+cOIB5PgWdBFhC5gaX9SQvusGsuRGRPc0nsiBm4fs4999l+HWk1g96boJxzP
wsFsTdr0oFVHwgRgfDjxXakH2LCVby8MkuOWGuyOuKelYXq34ZN7oeEjBBQNIRp+
uaP9ZgTEBzXlITlV99ttvQKBgQCEJijeslKk/XB71a8OoxbnG3bz0ykjekINd4dI
CBvCGurjpenbvmBNedc5meHffLCCVFxLVG1zSkEjKzuSIVtRglWnHwrGXY5/wCS5
+z63iojSU9g6IAsMT8m3WJ9V7+JkklOMoQhKbQVVTx5yrZBY0K15xg/4VeZyA5tB
R8dO/QKBgGqkJ1AB+qi9Tl10ic8IX2blqPt3FU6MkkVmmDl8vA5R35DBuZAD4VTs
vsc0Y79mSqP9XL3KRAfA04tbGme30gJWz4NOoNsaaF2T2fv0gNNnVma87unFS6Y4
Fv64CkXzShjd16ov4eZetsIAZYn/bVn8zp61I6V50iT6AjpX1ptX
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_waf_domain" "domain" {
  domain_name       = "tf-testacc.wafqa3.com"
  instance_id       = data.alicloud_waf_instances.default.ids.0
  is_access_product = "On"
  source_ips        = ["1.1.1.1"]
  cluster_type      = "PhysicalCluster"
  http2_port        = [443]
  http_port         = [80]
  https_port        = [443]
  http_to_user_ip   = "Off"
  https_redirect    = "Off"
  load_balancing    = "IpHash"
  log_headers {
    key   = "foo"
    value = "http"
  }
}
`, name)
}

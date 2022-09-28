package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
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
					"instance_id":    "${alicloud_waf_domain.domain.instance_id}",
					"domain":         "${alicloud_waf_domain.domain.domain_name}",
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

func TestUnitAlicloudWAFCertificate(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dId, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dIdCreate, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dIdCreate.MarkNewResource()
	dCreateError, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"certificate_name": "certificate_name",
		"instance_id":      "instance_id",
		"domain":           "domain",
		"private_key":      `-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA5KeYEdE3moKDwwB4DV+yB44BNOqJF6KOS3wSF0IhjnWRTJ13\nUoCTN0CwIJTPzgu7hPXuvoRe2Mgn/5CclHaF3x1+FIGVskyC1hm5I29nuP21MDkN\njJ0e9TLzSLfKNmSiJKbUqeLa6l46U/5rUdI3a+qOADXNIV4chZlgYXng+HykXfDK\nXDdsuGE8h9Ue2+1WLM7KlWpbFlG8JCTTxIKhuKqORvdZRPE0u20tVqELtSjrglkR\n6yc1ZXZ3MIpS9L6QynysxnCJ/CpYKxlR+SkVR5Uozj3+J4tPtYy7kQR7AsWAIxnr\n5JCAW96yFjK+LSdtiJm5fvcmXPJt8Byui67sJwIDAQABAoIBAHkIBUSZG/ebUids\nHh/mIZRCJH1gEAUEtJFcMS9CgASqUqjdSQXQrRw4D4KPgpesOiwdCayq4JGbESED\n6P3k5uO89rngBd8FBNRTdV4+YAbZtIi8CmH94R8lQwWg01aLWhBQho0OWZZpHjLL\nnv8GaPOkilmuINCiwjIth4jRPLYKM8rQUL0Y9DSLW1dZfxUMj5yPHdpV81/es21k\ns3ZUETOsygAfv9ybhxFW/1h8vldH9bCRvVq4oXHXkpF9Wp9K6z4UjR60k7hnxzTw\nUgFoEheHz58LOJNzh+dad/2HCCEVutQXMIMomoMvSmQkODhOMYsiJC8L15l6/193\ngj6Ah1ECgYEA8AE/efv7Q8DGwVRP51a4whzkn2Mws/1cURxXLflYNdrOek3Vmdda\nR4k9N6xYKTjQ4xSV+yYHpCMBdSkFcKFOOMf+DQLB5NAb3C4E+K6ea8GrfIUNnFZM\nYKwfjhDEvNFky8FqLscQGr9HCzcx+lOrwtnMK/srvS8r9IuNfXcalaUCgYEA8+Sy\n4brUb9SNqAFY68lOsovuhQpCK8UnomgbVZG52wziA368HYfdmLFSzSAMaOYA4wog\nUvC1JtOqkP6tyAFdo8qWgrGtmBpVmPSvWA2ser0zqugvRyg4uY6vd3Hwj9VEGguB\nvotYxUhzetzyapTZazzRuE7aUB9dpbnvdDhoyNsCgYBgpbABGFajfwLKoYAXwZVf\nHbF2+cOIB5PgWdBFhC5gaX9SQvusGsuRGRPc0nsiBm4fs4999l+HWk1g96boJxzP\nwsFsTdr0oFVHwgRgfDjxXakH2LCVby8MkuOWGuyOuKelYXq34ZN7oeEjBBQNIRp+\nuaP9ZgTEBzXlITlV99ttvQKBgQCEJijeslKk/XB71a8OoxbnG3bz0ykjekINd4dI\nCBvCGurjpenbvmBNedc5meHffLCCVFxLVG1zSkEjKzuSIVtRglWnHwrGXY5/wCS5\n+z63iojSU9g6IAsMT8m3WJ9V7+JkklOMoQhKbQVVTx5yrZBY0K15xg/4VeZyA5tB\nR8dO/QKBgGqkJ1AB+qi9Tl10ic8IX2blqPt3FU6MkkVmmDl8vA5R35DBuZAD4VTs\nvsc0Y79mSqP9XL3KRAfA04tbGme30gJWz4NOoNsaaF2T2fv0gNNnVma87unFS6Y4\nFv64CkXzShjd16ov4eZetsIAZYn/bVn8zp61I6V50iT6AjpX1ptX\n-----END RSA PRIVATE KEY-----`,
		"certificate":      `-----BEGIN CERTIFICATE-----\nMIID8DCCAtigAwIBAgIQTvMBGm/PRXSj352aOU7GSjANBgkqhkiG9w0BAQsFADBe\nMQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl\nc3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe\nFw0yMTA5MDcxNDM5NTBaFw0yMjA5MDcxNDM5NTBaMC0xCzAJBgNVBAYTAkNOMR4w\nHAYDVQQDExV0Zi10ZXN0YWNjLndhZnFhMy5jb20wggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQDkp5gR0TeagoPDAHgNX7IHjgE06okXoo5LfBIXQiGOdZFM\nnXdSgJM3QLAglM/OC7uE9e6+hF7YyCf/kJyUdoXfHX4UgZWyTILWGbkjb2e4/bUw\nOQ2MnR71MvNIt8o2ZKIkptSp4trqXjpT/mtR0jdr6o4ANc0hXhyFmWBheeD4fKRd\n8MpcN2y4YTyH1R7b7VYszsqValsWUbwkJNPEgqG4qo5G91lE8TS7bS1WoQu1KOuC\nWRHrJzVldncwilL0vpDKfKzGcIn8KlgrGVH5KRVHlSjOPf4ni0+1jLuRBHsCxYAj\nGevkkIBb3rIWMr4tJ22Imbl+9yZc8m3wHK6LruwnAgMBAAGjgdowgdcwDgYDVR0P\nAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME\nGDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB\nBQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov\nL2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MCAGA1UdEQQZMBeCFXRmLXRl\nc3RhY2Mud2FmcWEzLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAPjbt2H1HmEc8DzyD\npi4IF1CvaNlYgKjPojYlt/gpj2n0MfntL8Ihly3e2fdSMEeVeTnFWFd34L4uZxMa\nxE/hx6VJWfNdnYW7FGCZr0rGj/KrtAox83H1dRrZ63hynpgCMIbc5lhA7wDe0R16\nP/1l3c50ZEmidicGhK/qmzsSQIVXC0kJf6hDQCxW6LVaDrmT8mvbhRh4ZNb2pgJ5\nQIWJHnlOmZkUVsR5cMBGzK2ModADjHXHmeoHHr3Tw7mPioE4Xh5EmMTXPLG22BKN\nRBFG9gSFri+3RxqdXwi1ZJajO3Nup5mcdGaHJbUbNUf16YKIq50PJlrVxzCZV31f\n7cOGfw==\n-----END CERTIFICATE-----`,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for key, value := range map[string]interface{}{
		"certificate_name": "certificate_name",
		"instance_id":      "instance_id",
		"domain":           "domain",
	} {
		err := dCreateError.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for idKey, idValue := range map[string]interface{}{
		"instance_id":    "instance_id",
		"domain":         "domain",
		"certificate_id": "123456",
	} {
		err := dIdCreate.Set(idKey, idValue)
		assert.Nil(t, err)
		err = dId.Set(idKey, idValue)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Certificates": []interface{}{
			map[string]interface{}{
				"CertificateId":   "123456",
				"CertificateName": "certificate_name",
				"InstanceId":      "instance_id",
				"Domain":          "domain",
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_waf_certificate", "MockCertificateId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["CertificateId"] = "123456"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewWafClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudWafCertificateCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dCreateError, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("CreateCertificateByCertificateIdAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dId, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateCertificateByCertificateIdNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dIdCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "domain", ":", "123456"))
	//Delete
	t.Run("DeleteNormal", func(t *testing.T) {
		err := resourceAlicloudWafCertificateDelete(d, rawClient)
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeWafCertificateNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudWafCertificateRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeWafCertificateAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudWafCertificateRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

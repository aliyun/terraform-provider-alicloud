package alicloud

import (
	"fmt"
	"log"
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

func TestAccAliCloudGaListener_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence0)
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
					"accelerator_id": "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"port_ranges.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_protocol": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_protocol": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"idle_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "100",
							"to_port":   "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_ranges.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaListener_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence0)
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
					"accelerator_id":  "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"protocol":        "TCP",
					"proxy_protocol":  "true",
					"listener_type":   "Standard",
					"idle_timeout":    "10",
					"client_affinity": "SOURCE_IP",
					"name":            name,
					"description":     name,
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":  CHECKSET,
						"protocol":        "TCP",
						"proxy_protocol":  "true",
						"listener_type":   "Standard",
						"idle_timeout":    "10",
						"client_affinity": "SOURCE_IP",
						"name":            name,
						"description":     name,
						"port_ranges.#":   "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaListener_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence1)
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
					"accelerator_id": "${alicloud_ga_accelerator.default.id}",
					"protocol":       "HTTP",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"protocol":       "HTTP",
						"port_ranges.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTPS",
					"certificates": []map[string]string{
						{
							"id": "${local.certificate_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":       "HTTPS",
						"certificates.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_id": "tls_cipher_policy_1_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_id": "tls_cipher_policy_1_1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_version": "http3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_version": "http3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"idle_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_timeout": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificates": []map[string]interface{}{
						{
							"id": "${local.certificate_id_update}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificates.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "80",
							"to_port":   "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_ranges.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forwarded_for_config": []map[string]interface{}{
						{
							"forwarded_for_ga_id_enabled": "true",
							"forwarded_for_ga_ap_enabled": "true",
							"forwarded_for_proto_enabled": "true",
							"forwarded_for_port_enabled":  "true",
							"real_ip_enabled":             "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarded_for_config.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaListener_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence1)
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
					"accelerator_id":     "${alicloud_ga_accelerator.default.id}",
					"protocol":           "HTTPS",
					"security_policy_id": "tls_cipher_policy_1_1",
					"listener_type":      "Standard",
					"http_version":       "http3",
					"idle_timeout":       "60",
					"request_timeout":    "80",
					"client_affinity":    "SOURCE_IP",
					"name":               name,
					"description":        name,
					"certificates": []map[string]interface{}{
						{
							"id": "${local.certificate_id}",
						},
					},
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "80",
							"to_port":   "80",
						},
					},
					"forwarded_for_config": []map[string]interface{}{
						{
							"forwarded_for_ga_id_enabled": "true",
							"forwarded_for_ga_ap_enabled": "true",
							"forwarded_for_proto_enabled": "true",
							"forwarded_for_port_enabled":  "true",
							"real_ip_enabled":             "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":         CHECKSET,
						"protocol":               "HTTPS",
						"security_policy_id":     "tls_cipher_policy_1_1",
						"listener_type":          "Standard",
						"http_version":           "http3",
						"idle_timeout":           "60",
						"request_timeout":        "80",
						"client_affinity":        "SOURCE_IP",
						"name":                   name,
						"description":            name,
						"certificates.#":         "1",
						"port_ranges.#":          "1",
						"forwarded_for_config.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudGaListenerMap0 = map[string]string{
	"idle_timeout": CHECKSET,
	"status":       CHECKSET,
}

var AliCloudGaListenerMap1 = map[string]string{
	"idle_timeout":    CHECKSET,
	"request_timeout": CHECKSET,
	"status":          CHECKSET,
}

func AliCloudGaListenerBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status                 = "active"
  		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth      = 100
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		payment_type   = "PayAsYouGo"
  		billing_type   = "PayBy95"
  		ratio          = 30
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}
`, name)
}

func AliCloudGaListenerBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ga_accelerator" "default" {
  		bandwidth_billing_type = "CDT"
  		payment_type           = "PayAsYouGo"
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID1zCCAr+gAwIBAgIRAOrWWz1qmkcSg90JDHjuzFwwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjQxMTI2MDczNjA4WhcNMjkxMTI1MDczNjA4WjAgMQswCQYDVQQGEwJDTjER
MA8GA1UEAxMIdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDa7HDGbQ1Km0f4ZaFzYbjVN0q8KkvZ+oQUd4naGOZnlH5k0XFwmjg+TWf88YX3
5IF8c45/rXrTWucPLg7FeqR96Wq9HZEmzEhs6VG031V9Hqa32saRScCOAyhiW7Hj
OWf6BZveuxbZNbgQCR59QzX4CeAIC68xavIDAy3wcTAH9cIkD71BxEPJGGR7BIVH
9DcWXaMAnJqQfrkth0xHBjflZABHAI0wPYPfaw8fd9DRkMYOIkfjwrrcL5IvhI1u
D3wdHJQWA2vR8hjoU4dHiJLbUtQ+xV1UGVkF67CpQ6LDjSQdX7xlZ7WJMc/7dCJ9
a7tr0ZTwq4/3KSgcRvm62oGvAgMBAAGjgc0wgcowDgYDVR0PAQH/BAQDAgWgMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQogSYF0TQa
P8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGGFWh0dHA6
Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15c3NsLmNv
bS9teXNzbHRlc3Ryc2EuY3J0MBMGA1UdEQQMMAqCCHRlc3QuY29tMA0GCSqGSIb3
DQEBCwUAA4IBAQAxPOlK5WBA9kITzxYyjqe/YvWzfMlsmj0yvpyHrPeZf7HZTTFz
ebYkzrHL8ZLyOHBhag0nL7Poj6ek98NoXTuCYCi8LspdadapOeYQzLce3beu/frk
sqU0A6WLHG9Ol9yUDMCX7xvLoAY/LDrcOM3Z87C/u/ykB4wKfFN2XfR3EZx3PQqw
sV77LOnyQixB4FMHpHlKuDoUkSN9uvxwEPOeGnLZXm96hPsjPwk1bDM8qerNPpVI
CwJ6kNuZ2eLz2Umqu2Gh3l4aADdIwxRY1OOjjZNut8STosABKWVGIwQbbAdRPQze
qHZ05oVTjFy9L1DAzhQ5Zn3oUjLl5KW4tYBA
-----END CERTIFICATE-----
EOF
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2uxwxm0NSptH+GWhc2G41TdKvCpL2fqEFHeJ2hjmZ5R+ZNFx
cJo4Pk1n/PGF9+SBfHOOf61601rnDy4OxXqkfelqvR2RJsxIbOlRtN9VfR6mt9rG
kUnAjgMoYlux4zln+gWb3rsW2TW4EAkefUM1+AngCAuvMWryAwMt8HEwB/XCJA+9
QcRDyRhkewSFR/Q3Fl2jAJyakH65LYdMRwY35WQARwCNMD2D32sPH3fQ0ZDGDiJH
48K63C+SL4SNbg98HRyUFgNr0fIY6FOHR4iS21LUPsVdVBlZBeuwqUOiw40kHV+8
ZWe1iTHP+3QifWu7a9GU8KuP9ykoHEb5utqBrwIDAQABAoIBAQCErEfIKOymKybZ
pZXLnAxswt563FMtngGPecZEM1TmrvpOVROffwbY0wZTJ3fd/FBwwIM6Y0MNdYiU
DYCMM0AewmeahqGh1qmJv3hx2eswMXQt9driz8RvDADcYt+SagbWYbHNsKovJrwO
k8gzd5jsYeewWIxqsXpLUxDzJ1VJbIqoHgkrirRRPo0onpixPWeA0RbElSwjwIUw
y43cC4WF8N7wot3cTST8yeKM8ujtqpN22ZtKnbkHTd03vnwQTMeUMJeDQmSmY5aJ
yFr7yw/Z66+7Amh6pkWhzZSDHsjI4y/S3CCdpwFlMA7ID590umJB6HFxWsmVacSe
MSs2vIJZAoGBAOiecPH1HVDQqH6PcrN/X9E3pDKSyAj+nHsVDGIZsie9f5g/qA0A
tcJtQLS0CzrpMTLsAnsfdh2T7Lg6pYFz5jnOUyMjOImAEbCtgvqBxqgFea//OhdP
8s/RmxKIAenBsk7Wbwx8/KPhbZLUNe8OnILVHDfS6kLSa49Iu+4UvrpNAoGBAPDt
mky5MMHKdHwbqxPo9jYrz1m3gqqIvv+VihO4t/DE6t2Zg43ctfFm1BVEDSwPjYs/
YV69KfVrVRUnzMZVdtHZ/dBK784YTY0OujemoaIzMKFIL8tbJFldVv2IgB+IelTX
e675hVdHjNUqZhHwccd8X6d/8icohZw62SNHb/HrAoGBAN1HSt1/c6Gau42Y212Q
fw9ARLuvEQYtXaFfxmXTV7uh8axccXndAQmwb+r1kfE6PojYJQwGQ4+jVX1ynFnm
bEz0zfUQ3gk+gJV2mK+/n7/ZZYZb3WCrtqimFUOtiVRZ40pHhV91zcX+/QK9R4je
d1elbbBUvG9QRu0IHW0+4qfJAoGAOmlQvIM1l/ZOsXw/yO71KoMKnXTJYDERJYQK
2ucw6VXEn39FjtJQ5jsI9jLugp0usvDl2YNBNfgUw7FHi1pTGWOhjqtsYmov+x/z
8+QZUerZQnDu7X2mXWgs3AEJFxwOlJ09pllmg5ecRF4oKvdBjpzP0BtMCURgyFTY
Kh56vIsCgYBMbneMvFY6PCESKIAXj16BF4lqYVXFqHVoxyfxIuVlAy3TMNwxvpbS
yDETk05Ux9yNES0WyTb1SWVG1o1wXc0dnDXCwJqLC1tzJUNUSD1AYvktoNIFErcN
gs3ercrzBTX5ezORPj9ErRAPrSq+V3z1Lge5Gl+EqgDvAfnknww75w==
-----END RSA PRIVATE KEY-----
EOF
	}

resource "alicloud_ssl_certificates_service_certificate" "update" {
  		certificate_name = "${var.name}-update"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID1jCCAr6gAwIBAgIQGKYS2rt7QuCbV3mpxs2D9DANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNDExMjYwNzM3NTFaFw0yOTExMjUwNzM3NTFaMCAxCzAJBgNVBAYTAkNOMREw
DwYDVQQDEwh0ZXN0LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
ANqXuMEuRfqQ94tlv44PmYbbuJN6d1qaYu5DozPfOrqKMUD5TRkY1+hZtmC+36Ze
EuwQplYK6O+eaMXAloXL9Ofo4oz5Ny6fo6vjN32dcwD3iCYuQY6YpNQlnpl2jb7K
yh8CQYWbkGQ+U3Yg7K2ewp2HjWLBR0ODzGrcej0csbQ2WJtVzm5ptAbRfdLQADQ0
Q9ZmQ2RU4vmCqHGN7xZdnCEoWUMlvec++DsRB94URyAEsU+Z7hDzRR7723HSszry
Q+3aZfqlu4iq852lRQGUYJ8KoGyUWGlynnREB93KyGchG+x/lgADAYWlJh/19CgM
ElY4s1bqTbCUrltlgSA5qZMCAwEAAaOBzTCByjAOBgNVHQ8BAf8EBAMCBaAwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiBJgXRNBo/
wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYVaHR0cDov
L29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlzc2wuY29t
L215c3NsdGVzdHJzYS5jcnQwEwYDVR0RBAwwCoIIdGVzdC5jb20wDQYJKoZIhvcN
AQELBQADggEBABhk5x58ctbUkoM5z18bT8Ny0Ko9p0P6wn5XduK7JWD9QwjM5ZKr
kA39pHQU9D4sGhEhLR9SlWvSmrVQmSRn5tn03eHRXhhGv87IWmkTPHBYkoz8LP4L
ArYjAZpo9odmWpH6C+IkhqUw9nPg31na9wwVdUBCYxuIlL36PoII16FNsWwBnKMi
X81UCm+1UHp4qF3dT6s34ttEVNRoYw/u3rnwqVtnwTDs4svcLMaRyyNZrgV2RG5L
LC5tM9mrqvbKQIvQRxc47V1FV+t4jNun7se4St5nWEAavdLwmS1K/1QLb9UmYJOv
Nw8ocKgnHvrCoI59SQSO+oin+weDMchDK6U=
-----END CERTIFICATE-----
EOF
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2pe4wS5F+pD3i2W/jg+Zhtu4k3p3Wppi7kOjM986uooxQPlN
GRjX6Fm2YL7fpl4S7BCmVgro755oxcCWhcv05+jijPk3Lp+jq+M3fZ1zAPeIJi5B
jpik1CWemXaNvsrKHwJBhZuQZD5TdiDsrZ7CnYeNYsFHQ4PMatx6PRyxtDZYm1XO
bmm0BtF90tAANDRD1mZDZFTi+YKocY3vFl2cIShZQyW95z74OxEH3hRHIASxT5nu
EPNFHvvbcdKzOvJD7dpl+qW7iKrznaVFAZRgnwqgbJRYaXKedEQH3crIZyEb7H+W
AAMBhaUmH/X0KAwSVjizVupNsJSuW2WBIDmpkwIDAQABAoIBAQCK4uuYknYUBhfC
khtrf63kaaaUzbMX9g/1ozQGuUbvTu6MgdnioE5OavHd9mjTo+IR62JEORpXZSbc
vsjkqfopf2aye4X8MaIkjHGtdmSjsKLo32r31zSjNmPWzeSx3NcfbKeE5JqRlqgg
3jqC9eRhgsbqgDNvSkaPfxaLzbd67/+KSdituRNqMGvKCgyZAT63yLiO7ArdtEaY
Ij+BSECjABmhue+sBWtObmovI+MGJ7RetnBRaFh5/3I6rd0bY3dyhwab0A2rWuM7
T1usQSZ/Z8c4s1V2anQ8AgvcAe2bAfSSRoCUNwuPMtyj1LJVk2MaxaiLfpUNpKb7
r5P3fP2BAoGBAN18C/Tp9duDoAf+LYUtV4riXQ1CeCw/wxL6g9X0e+dnP7L4kLe4
m+/YSZUkv7IlRf56p4t9r0if2+w7u95zzXt3k7PLuFRigjSpxnZ1hrYKKctNY2Oj
urEUV+dkoekplFC0kSFtOFYaNwewaTV/fkWa0Apd/ZnccbViLB2zYOIZAoGBAPyo
ThxDE69gAOEiQm8B29bAi1lM98Dx69KSrOXP8Yf/mOUJyRnRphYVmYHKA/T5Gubt
Rn6o849Le4mJNWyyrgllg5QEMfShBBzndc5tL4ltsIrzQTKu/we/GlwxqRVccJ1B
Tn4+76gvMKpFvFEDtHQp9/XWy/FhYY3/VO7qTtaLAoGAUe/TKI7pIoVmTa6tvmgQ
y9OEYyRk+tG35CyDW0KwF+JtgVNNjnogTjGwvxkyRcBeTY+orgUYNIDXRmSu0tP6
f6O0I77I+Ybb7omkXyyJYo0N+yUtEK6AoYQKJRNohq6YLOcwDbKvNcNK+nA768u3
th5Yuo0dBa+07UpdUbuLqvkCgYEApLKN4Gx1S5AYYqmzhqs+hDoVXEwJANRytlx4
qoIn31BleYAsgFEipCjGXU2z0KAFwl0P5Ab8Zf99c0Vm9wlu258562XkrqO7i5/y
MnMIVtyTBbDWYlSi2IjhhRG2N79/hXMJ2M/r58WDQqucu27f1g15nt67KQkiz66O
zgMdC0sCgYAd3QLHQfHBxqlHBokcjdHxWoX2fkKwdQlKlKuk6Q+quyrY0dIF2dxr
/suURAMr4407dP4cjrG9LfWGGYfpcqt79/QDa7rbp9z6zdu6CU+RzqZyfgAtcd6r
1LeiSMDF5dMPJoxkrA9/aKMmp4UbYv/UTexUQ41tK/PFTG6fye44pA==
-----END RSA PRIVATE KEY-----
EOF
	}

	locals {
  		certificate_id        = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
  		certificate_id_update = join("-", [alicloud_ssl_certificates_service_certificate.update.id, "%s"])
	}
`, name, defaultRegionToTest, defaultRegionToTest)
}

func TestUnitAliCloudGaListener(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id": "CreateListenerValue",
		"description":    "CreateListenerValue",
		"name":           "CreateListenerValue",
		"port_ranges": []map[string]interface{}{
			{
				"from_port": 60,
				"to_port":   70,
			},
		},
		"proxy_protocol": true,
		"certificates": []map[string]interface{}{
			{
				"id": "CreateListenerValue",
			},
		},
		"client_affinity": "CreateListenerValue",
		"protocol":        "CreateListenerValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeListener
		"Description": "CreateListenerValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "CreateListenerValue",
			},
		},
		"ClientAffinity": "CreateListenerValue",
		"Name":           "CreateListenerValue",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 60,
				"ToPort":   70,
			},
		},
		"State":      "active",
		"Protocol":   "CreateListenerValue",
		"ListenerId": "CreateListenerValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateListener
		"ListenerId": "CreateListenerValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_listener", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeListener Response
		"ListenerId": "CreateListenerValue",
	}
	errorCodes := []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateListener" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaListenerCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateListener
	attributesDiff := map[string]interface{}{
		"accelerator_id": "UpdateListenerValue",
		"description":    "UpdateListenerValue",
		"name":           "UpdateListenerValue",
		"port_ranges": []map[string]interface{}{
			{
				"from_port": 70,
				"to_port":   80,
			},
		},
		"proxy_protocol": true,
		"certificates": []map[string]interface{}{
			{
				"id": "UpdateListenerValue",
			},
		},
		"client_affinity": "UpdateListenerValue",
		"protocol":        "UpdateListenerValue",
	}
	diff, err := newInstanceDiff("alicloud_ga_listener", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeListener Response
		"Description": "UpdateListenerValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "UpdateListenerValue",
			},
		},
		"ClientAffinity": "UpdateListenerValue",
		"Name":           "UpdateListenerValue",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 70,
				"ToPort":   80,
			},
		},
		"Protocol":   "UpdateListenerValue",
		"ListenerId": "UpdateListenerValue",
	}
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateListener" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaListenerUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeListener" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaListenerRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteListener" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DescribeListener" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaListenerDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

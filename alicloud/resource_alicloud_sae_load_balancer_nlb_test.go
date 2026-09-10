package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudSaeLoadBalancerNlb_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_load_balancer_nlb.default"
	ra := resourceAttrInit(resourceId, AlicloudSAELoadBalancerNlbMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApplicationNlb")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAELoadBalancerNlbBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alicloud_sae_application.default.id}",
					"nlb_id":       "${alicloud_nlb_load_balancer.default.id}",
					"address_type": "Internet",
					"listeners": []map[string]interface{}{
						{
							"protocol":    "TCP",
							"port":        "80",
							"target_port": "8080",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default.id}",
							"zone_id":    "${data.alicloud_nlb_zones.default.zones.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          CHECKSET,
						"address_type":    "Internet",
						"listeners.#":     "1",
						"zone_mappings.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alicloud_sae_application.default.id}",
					"nlb_id":       "${alicloud_nlb_load_balancer.default.id}",
					"address_type": "Internet",
					"listeners": []map[string]interface{}{
						{
							"protocol":    "TCPSSL",
							"port":        "443",
							"target_port": "9090",
							"cert_ids":    "${local.certificate_id}",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.secondary.id}",
							"zone_id":    "${data.alicloud_nlb_zones.default.zones.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          CHECKSET,
						"address_type":    "Internet",
						"listeners.#":     "1",
						"zone_mappings.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nlb_id", "address_type", "zone_mappings"},
			},
		},
	})
}

func AlicloudSAELoadBalancerNlbBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_nlb_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "secondary" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.2.0/24"
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_sae_namespace" "default" {
  namespace_id          = "%s:%s"
  namespace_name        = var.name
  namespace_description = var.name
}

resource "alicloud_sae_application" "default" {
  app_name          = var.name
  namespace_id      = alicloud_sae_namespace.default.namespace_id
  image_url         = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type      = "Image"
  cpu               = "500"
  memory            = "2048"
  replicas          = "2"
  vswitch_id        = alicloud_vswitch.default.id
  vpc_id            = alicloud_vpc.default.id
  security_group_id = alicloud_security_group.default.id
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name             = var.name
  address_type                   = "Internet"
  vpc_id                         = alicloud_vpc.default.id
  load_balancer_type             = "Network"
  bandwidth_package_id           = ""
  deletion_protection_enabled    = "false"
  modification_protection_status = "NonProtection"
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.secondary.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
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

locals {
  certificate_id = join("", [alicloud_ssl_certificates_service_certificate.default.id, "-%s"])
}
`, name, defaultRegionToTest, name, defaultRegionToTest)
}

var AlicloudSAELoadBalancerNlbMap0 = map[string]string{}

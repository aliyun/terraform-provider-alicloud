package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test SslCertificatesService Instance. >>> Resource test cases.
func TestAccAliCloudSslCertificatesServiceInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcasinstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_type":        "cas",
					"period":              "12",
					"pricing_cycle":       "2",
					"instance_name":       name,
					"domain":              "tf-test.certqa.cn",
					"key_algorithm":       "RSA_2048",
					"generate_csr_method": "online",
					"validation_method":   "DNS",
					"country_code":        "CN",
					"auto_reissue":        "disable",
					"company_id":          "${alicloud_ssl_certificates_service_company.default.id}",
					"contact_id_list":     []string{"${alicloud_ssl_certificates_service_contact.default.id}"},
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"parameter": []map[string]interface{}{
						{
							"code":  "fullSpec",
							"value": "ws.dv.f",
						},
						{
							"code":  "fullDomainCount",
							"value": "1",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_type":        "cas",
						"period":              "12",
						"pricing_cycle":       "2",
						"instance_name":       name,
						"domain":              "tf-test.certqa.cn",
						"key_algorithm":       "RSA_2048",
						"generate_csr_method": "online",
						"validation_method":   "DNS",
						"city":                CHECKSET,
						"province":            CHECKSET,
						"country_code":        "CN",
						"auto_reissue":        "disable",
						"company_id":          CHECKSET,
						"contact_id_list.#":   "1",
						"resource_group_id":   CHECKSET,
						"parameter.#":         "2",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Test",
						// instance-level facts returned by the order/read APIs
						"spec":                  CHECKSET,
						"brand":                 CHECKSET,
						"certificate_type":      CHECKSET,
						"full_domain_count":     CHECKSET,
						"wildcard_domain_count": CHECKSET,
						"order_start_time":      CHECKSET,
						"order_end_time":        CHECKSET,
						"upgrade_status":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name + "_update",
					"company_id":        "${alicloud_ssl_certificates_service_company.update.id}",
					"contact_id_list":   []string{"${alicloud_ssl_certificates_service_contact.default.id}", "${alicloud_ssl_certificates_service_contact.update.id}"},
					"domain":            "tf-test-update.certqa.cn",
					"key_algorithm":     "RSA_4096",
					"validation_method": "HTTP",
					"country_code":      "US",
					"auto_reissue":      "enable",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
						"Extra":   "Value",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"company_id":        CHECKSET,
						"contact_id_list.#": "2",
						"domain":            "tf-test-update.certqa.cn",
						"key_algorithm":     "RSA_4096",
						"validation_method": "HTTP",
						"city":              CHECKSET,
						"province":          CHECKSET,
						"country_code":      "US",
						"auto_reissue":      "enable",
						"tags.%":            "3",
						"tags.Created":      "TF-update",
						"tags.For":          "Test-update",
						"tags.Extra":        "Value",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"generate_csr_method": "upload",
					"csr":                 "${local.test_csr}",
					// Test steps accumulate their configuration, so the domain set in the previous
					// step would otherwise still be declared here. It is dropped rather than
					// restated: uploading a CSR hands authority over the domain to that CSR, and
					// declaring one alongside it either contradicts the CSR — which the server
					// rejects — or supplies the very answer the check below is meant to observe.
					"domain": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"generate_csr_method": "upload",
						"csr":                 CHECKSET,
						// Nothing in the configuration asks for this domain — the previous step set
						// tf-test-update.certqa.cn and this one names none at all. Reading back the
						// common name the test CSR was issued for is therefore evidence that the
						// server took it out of the CSR.
						"domain": "tf-test.certqa.cn",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
				// Per-attribute verification is off because the order parameters
				// (parameter/period/pricing_cycle/product_type) are consumed by BssOpenApi at
				// purchase time and never reported back by GetInstanceDetail, and company_id is
				// likewise absent from the read response until a certificate application has been
				// submitted. An imported instance cannot recover any of them, so the import is
				// exercised for the identity round trip only.
				ImportStateVerify: false,
			},
		},
	})
}

var AliCloudSslCertificatesServiceInstanceMap = map[string]string{
	"status":        CHECKSET,
	"instance_type": CHECKSET,
	"id":            CHECKSET,
}

func AliCloudSslCertificatesServiceInstanceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

locals {
  test_csr = <<CSR
-----BEGIN CERTIFICATE REQUEST-----
MIICsDCCAZgCAQAwazELMAkGA1UEBhMCQ04xETAPBgNVBAgMCFpoZWppYW5nMREw
DwYDVQQHDAhIYW5nemhvdTEMMAoGA1UECgwDcmRrMQwwCgYDVQQLDANyZGsxGjAY
BgNVBAMMEXRmLXRlc3QuY2VydHFhLmNuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAuNebmEk8hdwHlin+8QT4QqlcV9oSmK90lrH6jOKC9jQAGK2rfmjm
M3BeMUrWTfq7MaQeUOCLrAU08n56C3y4FIfxjsGqFCr9tY5VkC1vh9yeepwiiFAW
UhqMBpdeXLnAOX6JT9jDGckjd04RlOyfji1Iy05bJNBi1r/b9lCW3sWbloR5lW8I
pct+SH+gT/sOus+aBOnEEk/xHe6cDKV5wJnilzTT96fpwXgJ99h+pJlv/id8/tHy
l9Xni0CzuUK5pwYfuZwQlqSHgFNG/kBUqoW/zfn4Ko9SBfKoM/7B5PjyPt2vcXRW
TTtdFxrYM585JkGSoVlXxdLoQ9NWv5rmhwIDAQABoAAwDQYJKoZIhvcNAQELBQAD
ggEBAD9Z/A8Xm3gTOwLeDSsaP7HajLlTym16r7wpitNZ4Cf7dtSeSTcpADrHeXGr
kG3sjhNz6Ef/q+WsjSs52/Rz4xuhv2NWHIxXpEZCgmHl4hzL+LLdOQXsI44k1x58
9KyjpbJ8IH7Pe3IQyz8OlvX9cPpnCqyf1gTwNz0Ki6T/5BjTdaBjuRa6Zeu9A9jg
/naOLO8BScwD84Dgp/yLRHOfjo1tWmZoXcpKfM8ob3hTDBAFHec6uBUnwShlGqLY
J9jZJW/mwgFvgiSGlrK/EKtGXsr4qZGgi6b7aNt+Si+jF1Rt9rWrgCpKRsPd3cV9
5mWJaKIo6y3R8dPBmH6Ip8RCEQ0=
-----END CERTIFICATE REQUEST-----
CSR
}

resource "alicloud_ssl_certificates_service_contact" "default" {
  name   = var.name
  mobile = "1390000${substr(var.name, -4, 4)}"
  email  = "${var.name}@example.com"
}

resource "alicloud_ssl_certificates_service_contact" "update" {
  name   = "${var.name}-update"
  mobile = "1391111${substr(var.name, -4, 4)}"
  email  = "${var.name}-update@example.com"
}

resource "alicloud_ssl_certificates_service_company" "default" {
  company_name    = var.name
  company_address = "西安市"
  department      = "测试部门"
  city            = "西安"
  province        = "陕西"
  country_code    = "111122"
  post_code       = "11112233"
  company_type    = "1"
  company_code    = "12312311"
  company_phone   = "15101081174"
  company_email   = "test@example.com"
  lang            = "zh"
}

resource "alicloud_ssl_certificates_service_company" "update" {
  company_name    = "${var.name}-update"
  company_address = "上海市"
  department      = "测试部门"
  city            = "上海"
  province        = "上海"
  country_code    = "111122"
  post_code       = "11112233"
  company_type    = "1"
  company_code    = "12312312"
  company_phone   = "15101081175"
  company_email   = "test-update@example.com"
  lang            = "zh"
}
`, name)
}

// Test SslCertificatesService Instance. <<< Resource test cases.

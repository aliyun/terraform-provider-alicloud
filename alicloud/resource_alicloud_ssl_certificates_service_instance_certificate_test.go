package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Test SslCertificatesService InstanceCertificate. >>> Resource test cases.

// A certificate is issued by the certificate authority and cannot be created through an API, so
// this read-only resource needs one to already exist. Rather than depend on whatever the account
// happens to hold, the test issues its own: it buys an instance, submits the application, waits for
// issuance, and then reads back the certificate that came out.
func TestAccAliCloudSslCertificatesServiceInstanceCertificate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_instance_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceInstanceCertificateMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceInstanceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcascert%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceInstanceCertificateBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// Destroying this resource calls no API: the certificate was issued by the certificate
		// authority and goes on existing, so GetCertificateDetail keeps succeeding afterwards.
		// Asserting that it has disappeared would always fail, and asserting nothing about the
		// cloud is the honest check here.
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_id": "${alicloud_ssl_certificates_service_certificate_validation.default.certificate_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_id":     CHECKSET,
						"instance_id":        CHECKSET,
						"cert_identifier":    CHECKSET,
						"certificate_name":   CHECKSET,
						"certificate_status": CHECKSET,
						"certificate_source": CHECKSET,
						"common_name":        CHECKSET,
						"domain":             CHECKSET,
						"algorithm":          CHECKSET,
						"key_size":           CHECKSET,
						"not_before":         CHECKSET,
						"not_after":          CHECKSET,
						"issuer":             CHECKSET,
						"serial":             CHECKSET,
						"finger_print":       CHECKSET,
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

var AliCloudSslCertificatesServiceInstanceCertificateMap = map[string]string{
	"id":             CHECKSET,
	"certificate_id": CHECKSET,
}

// sslCertificatesServiceIssuedCertificateDependence builds the full issuance chain — instance,
// application, and the wait — so that a test can reference a certificate that has actually been
// issued. The subdomain is derived from the test name to keep concurrent runs from colliding.
//
// The domain is held in the test account and its authoritative DNS is Alibaba Cloud DNS, so the
// certificate authority validates it without any record having to be published and the certificate
// is issued a few minutes after the application is submitted. It is spelled out here rather than
// read from an environment variable because the remote acceptance-test runner injects only
// credentials — a getenv would read empty there and silently skip every test.
func sslCertificatesServiceIssuedCertificateDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "domain" {
  default = "%s.%s"
}

resource "alicloud_ssl_certificates_service_contact" "default" {
  name   = var.name
  mobile = "1390000${substr(var.name, -4, 4)}"
  email  = "${var.name}@example.com"
}

resource "alicloud_ssl_certificates_service_company" "default" {
  company_name    = var.name
  company_address = "杭州市"
  department      = "测试部门"
  city            = "杭州"
  province        = "浙江"
  country_code    = "111122"
  post_code       = "11112233"
  company_type    = "1"
  company_code    = "12312311"
  company_phone   = "15101081174"
  company_email   = "test@example.com"
  lang            = "zh"
}

resource "alicloud_ssl_certificates_service_instance" "default" {
  product_type        = "cas"
  period              = 12
  pricing_cycle       = 2
  instance_name       = var.name
  domain              = var.domain
  key_algorithm       = "RSA_2048"
  generate_csr_method = "online"
  validation_method   = "DNS"
  country_code        = "CN"
  contact_id_list     = [alicloud_ssl_certificates_service_contact.default.id]
  company_id          = alicloud_ssl_certificates_service_company.default.id

  parameter {
    code  = "fullSpec"
    value = "ws.dv.f"
  }
  parameter {
    code  = "fullDomainCount"
    value = "1"
  }
}

resource "alicloud_ssl_certificates_service_certificate_apply" "default" {
  instance_id         = alicloud_ssl_certificates_service_instance.default.id
  domain              = alicloud_ssl_certificates_service_instance.default.domain
  validation_method   = alicloud_ssl_certificates_service_instance.default.validation_method
  key_algorithm       = alicloud_ssl_certificates_service_instance.default.key_algorithm
  generate_csr_method = alicloud_ssl_certificates_service_instance.default.generate_csr_method
}

resource "alicloud_ssl_certificates_service_certificate_validation" "default" {
  instance_id = alicloud_ssl_certificates_service_instance.default.id
}
`, name, name, "aliterraform.com")
}

func AliCloudSslCertificatesServiceInstanceCertificateBasicDependence(name string) string {
	return sslCertificatesServiceIssuedCertificateDependence(name)
}

// Test SslCertificatesService InstanceCertificate. <<< Resource test cases.

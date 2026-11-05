package sslcertificateservice_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sslcertificateservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// sslCertificatesServiceUnissuableDomain is deliberately a domain the test account does not hold, so
// an application for it never passes validation and stays pending indefinitely. That is what the
// application resource needs to exercise its whole lifecycle: withdrawing only applies while an
// application is outstanding, and a domain the account does hold is validated by the certificate
// authority on its own within minutes.
//
// It is a constant rather than an environment variable because the remote acceptance-test runner
// injects only credentials — a getenv would read empty there and silently skip the test.
const sslCertificatesServiceUnissuableDomain = "tf-acc-cas-not-resolvable.com"

// Test SslCertificatesService CertificateApply. >>> Resource test cases.

// Applies for a certificate against a domain the account does not hold, so the application stays
// outstanding for the whole test. That is what makes the resource's lifecycle observable: the
// validation records this resource exists to report are only published while an application is
// pending, and withdrawing one on destroy is only possible for as long as it has not been issued.
// A domain the account does hold is validated by the certificate authority on its own within
// minutes, which would settle the application before either could be exercised.
func TestAccAliCloudSslCertificatesServiceCertificateApply_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate_apply.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateApplyMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sslcertificateservice.NewSslCertificatesServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSslCertificatesServiceCertificateApply")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcasapply%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSslCertificatesServiceCertificateApplyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers: acctest.SDKProviders,
		// Destroying the application withdraws it, but the instance it was submitted on survives,
		// so GetInstanceDetail keeps succeeding and the default "describe must report NotFound"
		// check would always fail. What marks an instance as no longer applying for a certificate
		// is the disappearance of its validation records, so that is what gets asserted.
		CheckDestroy: testAccCheckSslCertificatesServiceCertificateApplyDestroy,
		Steps: []resource.TestStep{
			{
				// The application settings are referenced from the instance rather than
				// hardcoded. ApplyCertificate submits whatever the instance holds, so Create
				// rejects values that disagree with it.
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_ssl_certificates_service_instance.default.id}",
					"domain":              "${alicloud_ssl_certificates_service_instance.default.domain}",
					"validation_method":   "${alicloud_ssl_certificates_service_instance.default.validation_method}",
					"key_algorithm":       "${alicloud_ssl_certificates_service_instance.default.key_algorithm}",
					"generate_csr_method": "${alicloud_ssl_certificates_service_instance.default.generate_csr_method}",
					// The instance generates the CSR itself (generate_csr_method = online), so this
					// resolves to an empty CSR. It is referenced rather than hardcoded because
					// ApplyCertificate submits whatever the instance holds, and Create rejects a
					// value that disagrees with it.
					"csr": "${alicloud_ssl_certificates_service_instance.default.csr}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         acctest.CHECKSET,
						"domain":              acctest.CHECKSET,
						"validation_method":   acctest.CHECKSET,
						"key_algorithm":       acctest.CHECKSET,
						"generate_csr_method": acctest.CHECKSET,
						// The validation records the certificate authority wants published. They are
						// what this resource exists to produce. certificate_status is deliberately
						// not asserted: it only gets a value once a certificate has been issued,
						// which cannot happen until these records are published and validated.
						"domain_validation_list.#":                  "1",
						"domain_validation_list.0.domain":           acctest.CHECKSET,
						"domain_validation_list.0.root_domain":      acctest.CHECKSET,
						"domain_validation_list.0.validation_type":  acctest.CHECKSET,
						"domain_validation_list.0.validation_key":   acctest.CHECKSET,
						"domain_validation_list.0.validation_value": acctest.CHECKSET,
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

// A configured application setting that disagrees with the instance is refused before anything is
// submitted. Without that refusal the value would be quietly ignored — ApplyCertificate reads the
// configuration off the instance — and then read back as the instance's own, leaving a plan that
// never converges. The domain stands in for all five settings; they share one code path.
func TestAccAliCloudSslCertificatesServiceCertificateApply_configurationMismatch(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcasapply%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: testAccCheckSslCertificatesServiceCertificateApplyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "alicloud_ssl_certificates_service_certificate_apply" "default" {
  instance_id = alicloud_ssl_certificates_service_instance.default.id
  domain      = "mismatched.%s"
}
`, AliCloudSslCertificatesServiceCertificateApplyBasicDependence(name), sslCertificatesServiceUnissuableDomain),
				ExpectError: regexp.MustCompile("ApplyCertificate submits the configuration recorded on the instance"),
			},
		},
	})
}

var AliCloudSslCertificatesServiceCertificateApplyMap = map[string]string{
	"id":          acctest.CHECKSET,
	"instance_id": acctest.CHECKSET,
}

// testAccCheckSslCertificatesServiceCertificateApplyDestroy asserts that the application really was
// withdrawn, rather than that the instance disappeared — it does not, the application is destroyed
// but the instance it belongs to lives on.
//
// The signal is DomainValidationList: an instance with an application in flight reports the records
// the certificate authority wants published, and an instance withdrawn from one reports none. A
// withdrawal that was accepted but never took effect leaves them in place, which is exactly the
// failure this check exists to catch.
func testAccCheckSslCertificatesServiceCertificateApplyDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := sslcertificateservice.NewSslCertificatesServiceServiceV2(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_certificates_service_certificate_apply" {
			continue
		}
		object, err := service.DescribeSslCertificatesServiceCertificateApply(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return errs.WrapError(err)
		}
		if records, ok := object["DomainValidationList"]; ok && records != nil {
			return errs.WrapError(errs.Error("certificate application %s is still outstanding after destroy: the instance is %v and still reports validation records %v",
				rs.Primary.ID, object["Status"], records))
		}
	}

	return nil
}

func AliCloudSslCertificatesServiceCertificateApplyBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
  domain              = "${var.name}.%s"
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
`, name, sslCertificatesServiceUnissuableDomain)
}

// Test SslCertificatesService CertificateApply. <<< Resource test cases.

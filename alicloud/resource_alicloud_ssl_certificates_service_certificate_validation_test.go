package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test SslCertificatesService CertificateValidation. >>> Resource test cases.

// The wait is exercised against an instance whose certificate has already been issued, adopted
// through the data source rather than issued here. Issuing a certificate inside the test would
// mean buying an instance that can never be cleaned up afterwards — once a certificate has been
// issued, the instance can be neither refunded nor deleted — so past issuance runs have left the
// test account holding such instances permanently, and this test adopts one of those instead of
// adding to them. Against an issued certificate the waiter returns on its first poll, which is
// the wait's success path; the failure path is covered by the unissuable case below.
func TestAccAliCloudSslCertificatesServiceCertificateValidation_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate_validation.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateValidationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificateValidation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcasvalidation%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateValidationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// Destroying a wait calls no API, and the adopted instance belongs to the account, not to
		// this test — there is nothing whose absence could be asserted.
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}",
					// The value is opaque to the resource — validation_record_ids maps to no API
					// field and exists to order the wait after whatever publishes the validation
					// records, so any string exercises exactly what the attribute does. The adopted
					// certificate was validated before this test ran; nothing here published records.
					"validation_record_ids": []string{"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"certificate_id":     CHECKSET,
						"cert_identifier":    CHECKSET,
						"certificate_status": "issued",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
				// Per-attribute verification is off because validation_record_ids exists only to
				// order the graph — it maps to no API field and is never read back, so an imported
				// wait cannot recover it.
				ImportStateVerify: false,
			},
		},
	})
}

var AliCloudSslCertificatesServiceCertificateValidationMap = map[string]string{
	"id":          CHECKSET,
	"instance_id": CHECKSET,
}

func AliCloudSslCertificatesServiceCertificateValidationBasicDependence(name string) string {
	return `
data "alicloud_ssl_certificates_service_instances" "default" {
  certificate_status = "issued"
}
`
}

// The failure path: an application for a domain the account does not hold never validates, so the
// certificate is never issued and the wait must surface that rather than report success. A short
// create timeout keeps the case fast — the assertion is that the waiter gives up loudly, and on
// what error.
//
// The unissuable domain is also what makes this case self-cleaning, unlike the issued case above:
// a certificate that was never issued leaves the application withdrawable and the instance
// refundable, so the destroy step tears the whole chain down and asserts that in passing.
func TestAccAliCloudSslCertificatesServiceCertificateValidation_unissuable(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcasvalidation%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		// The destroy step itself is the cleanup assertion: it fails the test if the application
		// cannot be withdrawn or the instance cannot be refunded and deleted.
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "alicloud_ssl_certificates_service_certificate_apply" "default" {
  instance_id         = alicloud_ssl_certificates_service_instance.default.id
  domain              = alicloud_ssl_certificates_service_instance.default.domain
  validation_method   = alicloud_ssl_certificates_service_instance.default.validation_method
  key_algorithm       = alicloud_ssl_certificates_service_instance.default.key_algorithm
  generate_csr_method = alicloud_ssl_certificates_service_instance.default.generate_csr_method
}

resource "alicloud_ssl_certificates_service_certificate_validation" "default" {
  instance_id = alicloud_ssl_certificates_service_certificate_apply.default.instance_id

  timeouts {
    create = "3m"
  }
}
`, AliCloudSslCertificatesServiceCertificateApplyBasicDependence(name)),
				ExpectError: regexp.MustCompile("timeout while waiting for state to become 'issued'"),
			},
		},
	})
}

// Test SslCertificatesService CertificateValidation. <<< Resource test cases.

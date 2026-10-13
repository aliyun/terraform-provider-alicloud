package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test SslCertificatesService InstanceCertificate. >>> Resource test cases.

// A certificate is issued by the certificate authority and cannot be created through an API, so
// this read-only resource adopts one that already exists. The test account durably holds issued
// certificates — an instance whose certificate has been issued can be neither refunded nor
// deleted, so past issuance runs have left several behind — and the test reads one of those
// through the data sources rather than buying an instance and issuing its own.
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
					"certificate_id": "${data.alicloud_ssl_certificates_service_instance_certificates.default.certificates.0.certificate_id}",
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

// The adoption chain: pick an instance whose certificate has been issued, then list that
// instance's certificates to obtain a certificate ID to adopt. Both lookups are server-side
// filters, so the fixture holds no resources of its own and leaves nothing behind.
func AliCloudSslCertificatesServiceInstanceCertificateBasicDependence(name string) string {
	return `
data "alicloud_ssl_certificates_service_instances" "default" {
  certificate_status = "issued"
}

data "alicloud_ssl_certificates_service_instance_certificates" "default" {
  instance_id = data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id
}
`
}

// Test SslCertificatesService InstanceCertificate. <<< Resource test cases.

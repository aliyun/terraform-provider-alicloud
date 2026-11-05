package sslcertificateservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sslcertificateservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSslCertificatesServiceInstanceCertificateMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sslcertificateservice.NewSslCertificatesServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSslCertificatesServiceInstanceCertificate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcascert%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSslCertificatesServiceInstanceCertificateBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers: acctest.SDKProviders,
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
						"certificate_id":     acctest.CHECKSET,
						"instance_id":        acctest.CHECKSET,
						"cert_identifier":    acctest.CHECKSET,
						"certificate_name":   acctest.CHECKSET,
						"certificate_status": acctest.CHECKSET,
						"certificate_source": acctest.CHECKSET,
						"common_name":        acctest.CHECKSET,
						"domain":             acctest.CHECKSET,
						"algorithm":          acctest.CHECKSET,
						"key_size":           acctest.CHECKSET,
						"not_before":         acctest.CHECKSET,
						"not_after":          acctest.CHECKSET,
						"issuer":             acctest.CHECKSET,
						"serial":             acctest.CHECKSET,
						"finger_print":       acctest.CHECKSET,
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
	"id":             acctest.CHECKSET,
	"certificate_id": acctest.CHECKSET,
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

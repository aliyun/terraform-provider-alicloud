package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// Certificates cannot be created through an API, so this data source needs an issued one to list.
// The test account durably holds issued certificates — an instance whose certificate has been
// issued can be neither refunded nor deleted — so the test queries one of those instances instead
// of buying an instance and issuing its own.
func TestAccAliCloudSslCertificatesServiceInstanceCertificateDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id": `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id": `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_status": `"issued"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_status": `"revoked"`,
		}),
	}

	sourceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_source": `"BUY"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_source": `"UPLOAD"`,
		}),
	}

	SslCertificatesServiceInstanceCertificateCheckInfo.dataSourceTestCheck(t, rand, instanceIdConf, statusConf, sourceConf)
}

var existSslCertificatesServiceInstanceCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":                    CHECKSET,
		"certificates.0.certificate_id":     CHECKSET,
		"certificates.0.instance_id":        CHECKSET,
		"certificates.0.cert_identifier":    CHECKSET,
		"certificates.0.certificate_name":   CHECKSET,
		"certificates.0.certificate_status": CHECKSET,
		"certificates.0.common_name":        CHECKSET,
		"certificates.0.domain":             CHECKSET,
		"certificates.0.algorithm":          CHECKSET,
		"certificates.0.not_before":         CHECKSET,
		"certificates.0.not_after":          CHECKSET,
		"certificates.0.issuer":             CHECKSET,
		"certificates.0.finger_print":       CHECKSET,
		"ids.#":                             CHECKSET,
	}
}

var fakeSslCertificatesServiceInstanceCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
		"ids.#":          "0",
	}
}

var SslCertificatesServiceInstanceCertificateCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_certificates_service_instance_certificates.default",
	existMapFunc: existSslCertificatesServiceInstanceCertificateMapFunc,
	fakeMapFunc:  fakeSslCertificatesServiceInstanceCertificateMapFunc,
}

func testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	// The instance under query is found through a server-side filter on certificate status, so the
	// query always targets an instance that actually holds an issued certificate.
	return fmt.Sprintf(`
data "alicloud_ssl_certificates_service_instances" "default" {
  certificate_status = "issued"
}

data "alicloud_ssl_certificates_service_instance_certificates" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
}

// Test SslCertificatesService InstanceCertificate. <<< Data source test cases.

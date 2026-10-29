package sslcertificateservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// Certificates cannot be created through an API, so this data source needs an issued one to list.
// The test account durably holds issued certificates — an instance whose certificate has been
// issued can be neither refunded nor deleted — so the test queries one of those instances instead
// of buying an instance and issuing its own.
func TestAccAliCloudSslCertificatesServiceInstanceCertificateDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(10000, 99999)

	instanceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id": `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id": `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_status": `"issued"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_status": `"revoked"`,
		}),
	}

	sourceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_source": `"BUY"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${data.alicloud_ssl_certificates_service_instances.default.instances.0.instance_id}"`,
			"certificate_source": `"UPLOAD"`,
		}),
	}

	SslCertificatesServiceInstanceCertificateCheckInfo.DataSourceTestCheck(t, rand, instanceIdConf, statusConf, sourceConf)
}

var existSslCertificatesServiceInstanceCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":                    acctest.CHECKSET,
		"certificates.0.certificate_id":     acctest.CHECKSET,
		"certificates.0.instance_id":        acctest.CHECKSET,
		"certificates.0.cert_identifier":    acctest.CHECKSET,
		"certificates.0.certificate_name":   acctest.CHECKSET,
		"certificates.0.certificate_status": acctest.CHECKSET,
		"certificates.0.common_name":        acctest.CHECKSET,
		"certificates.0.domain":             acctest.CHECKSET,
		"certificates.0.algorithm":          acctest.CHECKSET,
		"certificates.0.not_before":         acctest.CHECKSET,
		"certificates.0.not_after":          acctest.CHECKSET,
		"certificates.0.issuer":             acctest.CHECKSET,
		"certificates.0.finger_print":       acctest.CHECKSET,
		"ids.#":                             acctest.CHECKSET,
	}
}

var fakeSslCertificatesServiceInstanceCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
		"ids.#":          "0",
	}
}

var SslCertificatesServiceInstanceCertificateCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ssl_certificates_service_instance_certificates.default",
	ExistMapFunc: existSslCertificatesServiceInstanceCertificateMapFunc,
	FakeMapFunc:  fakeSslCertificatesServiceInstanceCertificateMapFunc,
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

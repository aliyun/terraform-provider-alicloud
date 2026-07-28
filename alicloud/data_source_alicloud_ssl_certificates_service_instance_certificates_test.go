package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// Certificates cannot be created through an API, so this data source needs an issued one to list.
// The test issues its own by running the full chain — instance, application, DNS validation record,
// wait — and then queries the certificates that instance produced.
func TestAccAliCloudSslCertificatesServiceInstanceCertificateDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_ssl_certificates_service_certificate_validation.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_ssl_certificates_service_certificate_validation.default.instance_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${alicloud_ssl_certificates_service_certificate_validation.default.instance_id}"`,
			"certificate_status": `"issued"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${alicloud_ssl_certificates_service_certificate_validation.default.instance_id}"`,
			"certificate_status": `"revoked"`,
		}),
	}

	sourceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${alicloud_ssl_certificates_service_certificate_validation.default.instance_id}"`,
			"certificate_source": `"BUY"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceCertificateSourceConfig(rand, map[string]string{
			"instance_id":        `"${alicloud_ssl_certificates_service_certificate_validation.default.instance_id}"`,
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
		"certificates.0.serial":             CHECKSET,
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

	// The filters reference the wait resource rather than the instance directly, so the query only
	// runs once the certificate has actually been issued and is therefore listable.
	return fmt.Sprintf(`
%s

data "alicloud_ssl_certificates_service_instance_certificates" "default" {
  %s
}
`, sslCertificatesServiceIssuedCertificateDependence(fmt.Sprintf("tfaccsslcasds%d", rand)), strings.Join(pairs, "\n  "))
}

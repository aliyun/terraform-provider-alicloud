package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSslCertificatesServicePcaCertificateDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)
	commonName := fmt.Sprintf("certqa-%d", rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_certificates_service_pca_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_certificates_service_pca_certificate.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf("%q", commonName),
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"nonexistent_pca_xyz"`,
		}),
	}

	certTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ssl_certificates_service_pca_certificate.default.id}"]`,
			"cert_type": `"subRoot"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ssl_certificates_service_pca_certificate.default.id}_fake"]`,
			"cert_type": `"externalCa"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ssl_certificates_service_pca_certificate.default.id}"]`,
			"name_regex":  fmt.Sprintf("%q", commonName),
			"cert_type":   `"subRoot"`,
			"output_file": `"/tmp/pca_certificates.txt"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_ssl_certificates_service_pca_certificate.default.id}_fake"]`,
			"cert_type": `"externalCa"`,
		}),
	}

	SslCertificatesServicePcaCertificatesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, certTypeConf, allConf)
}

var existSslCertificatesServicePcaCertificatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":                   "1",
		"certificates.0.id":                CHECKSET,
		"certificates.0.identifier":        CHECKSET,
		"certificates.0.serial_number":     CHECKSET,
		"certificates.0.x509_certificate":  CHECKSET,
		"certificates.0.certificate_type":  "SUB_ROOT",
		"certificates.0.algorithm":         CHECKSET,
		"certificates.0.sign_algorithm":    CHECKSET,
		"certificates.0.sha2":              CHECKSET,
		"certificates.0.md5":               CHECKSET,
		"certificates.0.locality":          "a",
		"certificates.0.organization":      "a",
		"certificates.0.organization_unit": "a",
		"certificates.0.common_name":       CHECKSET,
		"certificates.0.country_code":      CHECKSET,
		"certificates.0.state":             "a",
		"certificates.0.parent_identifier": CHECKSET,
		"certificates.0.status":            "ISSUE",
		"certificates.0.years":             "1",
		"certificates.0.before_date":       CHECKSET,
		"certificates.0.after_date":        CHECKSET,
		// resource_group_id is intentionally not asserted: DescribeCACertificateList returns ResourceGroupId with eventual consistency (absent on first Read after creation, populated on subsequent calls), so asserting it would be flaky.
	}
}

var fakeSslCertificatesServicePcaCertificatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
	}
}

var SslCertificatesServicePcaCertificatesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_certificates_service_pca_certificates.default",
	existMapFunc: existSslCertificatesServicePcaCertificatesMapFunc,
	fakeMapFunc:  fakeSslCertificatesServicePcaCertificatesMapFunc,
}

func testAccCheckAlicloudSslCertificatesServicePcaCertificatesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	commonName := fmt.Sprintf("certqa-%d", rand)
	rootCommonName := fmt.Sprintf("root-%d", rand)
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslCertificatesServicePcaCertificates%d"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "root" {
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = %q
  country_code      = "CN"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "default" {
  parent_identifier = "${alicloud_ssl_certificates_service_pca_certificate.root.id}"
  common_name       = %q
  locality          = "a"
  organization      = "a"
  organization_unit = "a"
  state             = "a"
  years             = "1"
  algorithm         = "RSA_2048"
  certificate_type  = "SUB_ROOT"
  country_code      = "CN"
}

data "alicloud_ssl_certificates_service_pca_certificates" "default" {
%s
}
`, rand, rootCommonName, commonName, strings.Join(pairs, "\n   "))
	return config
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDataWorksCertificateDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}_fake"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
	}

	ProjectIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}_fake"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}_fake"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
	}

	DataWorksCertificateCheckInfo.dataSourceTestCheck(t, rand, idsConf, ProjectIdConf, allConf)
}

var existDataWorksCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":                    "1",
		"certificates.0.description":        CHECKSET,
		"certificates.0.create_time":        CHECKSET,
		"certificates.0.create_user":        CHECKSET,
		"certificates.0.file_size_in_bytes": CHECKSET,
		"certificates.0.id":                 CHECKSET,
		"certificates.0.name":               CHECKSET,
	}
}

var fakeDataWorksCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
	}
}

var DataWorksCertificateCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_data_works_certificates.default",
	existMapFunc: existDataWorksCertificateMapFunc,
	fakeMapFunc:  fakeDataWorksCertificateMapFunc,
}

func testAccCheckAlicloudDataWorksCertificateSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tfaccdwcert%d"
}
resource "alicloud_data_works_project" "defaultxWQGsP" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}

resource "alicloud_oss_bucket" "cert" {
  bucket = "tfaccdwcert%d"
  acl    = "public-read"
}

resource "alicloud_oss_bucket_object" "cert" {
  bucket  = alicloud_oss_bucket.cert.id
  key     = "tf-acc-cert.pem"
  acl     = "public-read"
  content = <<-EOT
-----BEGIN CERTIFICATE-----
MIICqDCCAZACCQCbNs38rCeB6TANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDDAt0
Zi1hY2MtdGVzdDAeFw0yNjA5MDMwMDE4MDFaFw0zNjA4MzEwMDE4MDFaMBYxFDAS
BgNVBAMMC3RmLWFjYy10ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAwyKgyXSNG+NfmE6Ukz7oX4AgDbTaHAyJ3BWHfTK/DJNY0+0tgYE5Zle8bNui
rksdOpydcGAJeaoFNpCGwvd5FO9HklcYwLGdzASRXEHEN0XtG2oeHNnyH82yE2SH
hINbhYcM98WCEuIU5cIsohaQUXhdgBzQcUdc7q5mG/U87LQk4ZGq+JuqUhOAB2xG
8phoguTZ3ODk4+x9MEycpXfpTVkzz9y4JmL6jh3dEZvpyMmC94NYjrGRbioI9CDe
mpY1BsHRtSgxJ7Pvl5yYy7Tj2KC3jRv4FEP7C8vg9tqqV37yLtZsXw/UeVhn9nOm
tOQcqaS3o4BTp//Db7/ReoE2fwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBn+Hi1
Lq6pHiV0PPLN+4cQK8sLFSk1nHvJ1ONk9FsBmsnkrbky6I5fNQiBnzOm/rNl72Ao
JvsyROYVyyxX/M5LnZQwahUwDTaz4cNasUiPt33lNtueiLLhzC7UpB8NzH/v29MT
YnEmgc7msKrSnj8SfvZu6y9vbmLhJ0CbrpBTmcG6cXr/Em+l781p4hFdJmpOA7o8
oX+vfX/riPSP8bGgYlmsr2Iv4Uo/b8q6ZpF5p76M/MJirligYjc5fZG8UmsBu6eX
d5LC+4/r6W5uyIS/7fF7FLdzHfRhZknuu8Wzx+HXI+oKv0r/z8MGa1DsvJ6DHQIE
qWre020mtWvGRU1t
-----END CERTIFICATE-----
EOT
}

resource "alicloud_data_works_certificate" "default" {
  project_id       = alicloud_data_works_project.defaultxWQGsP.id
  name             = var.name
  certificate_file = "http://${alicloud_oss_bucket.cert.id}.${alicloud_oss_bucket.cert.extranet_endpoint}/${alicloud_oss_bucket_object.cert.key}"
}

data "alicloud_data_works_certificates" "default" {
%s
}
`, rand, rand, strings.Join(pairs, "\n   "))
	return config
}

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
			"ids": `["${alicloud_data_works_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_data_works_certificate.default.id}_fake"]`,
		}),
	}

	ProjectIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}_fake"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDataWorksCertificateSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_data_works_certificate.default.id}_fake"]`,
			"project_id": `"${alicloud_data_works_project.defaultxWQGsP.id}_fake"`,
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
	default = "tf-testAccDataWorksCertificate%d"
}
resource "alicloud_data_works_project" "defaultxWQGsP" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}



resource "alicloud_data_works_certificate" "default" {
  project_id       = alicloud_data_works_project.defaultxWQGsP.id
  certificate_file = "http://oxs-dataworks-openapi-console-cn-shanghai.oss-cn-shanghai.aliyuncs.com/d713_1107550004253538_8b2e08f49dd44414a8f27c4630b32c4b"
}

data "alicloud_data_works_certificates" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudSslCertificatesServiceInstanceDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"status":         `"inactive"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"status":         `"closed"`,
		}),
	}

	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"instance_type":  `"BUY"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"instance_type":  `"TEST"`,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"keyword":        `"${alicloud_ssl_certificates_service_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"keyword":        `"${alicloud_ssl_certificates_service_instance.default.instance_name}_fake"`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details":    `"true"`,
			"resource_group_id": `"${alicloud_ssl_certificates_service_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details":    `"true"`,
			"resource_group_id": `"rg-fake00000000000"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"status":            `"inactive"`,
			"instance_type":     `"BUY"`,
			"keyword":           `"${alicloud_ssl_certificates_service_instance.default.instance_name}"`,
			"resource_group_id": `"${alicloud_ssl_certificates_service_instance.default.resource_group_id}"`,
			"enable_details":    `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"status":            `"closed"`,
			"instance_type":     `"TEST"`,
			"keyword":           `"${alicloud_ssl_certificates_service_instance.default.instance_name}_fake"`,
			"resource_group_id": `"rg-fake00000000000"`,
			"enable_details":    `"true"`,
		}),
	}

	SslCertificatesServiceInstanceCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, instanceTypeConf, keywordConf, resourceGroupIdConf, allConf)
}

var existSslCertificatesServiceInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "1",
		"instances.#":               "1",
		"instances.0.id":            CHECKSET,
		"instances.0.instance_id":   CHECKSET,
		"instances.0.status":        CHECKSET,
		"instances.0.instance_type": CHECKSET,
		"instances.0.instance_name": CHECKSET,
		"instances.0.spec":          CHECKSET,
	}
}

var fakeSslCertificatesServiceInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"instances.#": "0",
	}
}

var SslCertificatesServiceInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_certificates_service_instances.default",
	existMapFunc: existSslCertificatesServiceInstanceMapFunc,
	fakeMapFunc:  fakeSslCertificatesServiceInstanceMapFunc,
}

func testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSslCasInstance%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ssl_certificates_service_instance" "default" {
  product_type      = "cas"
  period            = 12
  pricing_cycle     = 2
  instance_name     = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0

  parameter {
    code  = "fullSpec"
    value = "ws.dv.f"
  }
  parameter {
    code  = "fullDomainCount"
    value = "1"
  }
}

data "alicloud_ssl_certificates_service_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

package sslcertificateservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSslCertificatesServiceInstanceDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"status":         `"inactive"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"status":         `"closed"`,
		}),
	}

	instanceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"instance_type":  `"BUY"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"instance_type":  `"TEST"`,
		}),
	}

	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"keyword":        `"${alicloud_ssl_certificates_service_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details": `"true"`,
			"keyword":        `"${alicloud_ssl_certificates_service_instance.default.instance_name}_fake"`,
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details":    `"true"`,
			"resource_group_id": `"${alicloud_ssl_certificates_service_instance.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"enable_details":    `"true"`,
			"resource_group_id": `"rg-fake00000000000"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"status":            `"inactive"`,
			"instance_type":     `"BUY"`,
			"keyword":           `"${alicloud_ssl_certificates_service_instance.default.instance_name}"`,
			"resource_group_id": `"${alicloud_ssl_certificates_service_instance.default.resource_group_id}"`,
			"enable_details":    `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_certificates_service_instance.default.id}"]`,
			"status":            `"closed"`,
			"instance_type":     `"TEST"`,
			"keyword":           `"${alicloud_ssl_certificates_service_instance.default.instance_name}_fake"`,
			"resource_group_id": `"rg-fake00000000000"`,
			"enable_details":    `"true"`,
		}),
	}

	SslCertificatesServiceInstanceCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf, instanceTypeConf, keywordConf, resourceGroupIdConf, allConf)
}

var existSslCertificatesServiceInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "1",
		"instances.#":               "1",
		"instances.0.id":            acctest.CHECKSET,
		"instances.0.instance_id":   acctest.CHECKSET,
		"instances.0.status":        acctest.CHECKSET,
		"instances.0.instance_type": acctest.CHECKSET,
		"instances.0.instance_name": acctest.CHECKSET,
		"instances.0.spec":          acctest.CHECKSET,
	}
}

var fakeSslCertificatesServiceInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"instances.#": "0",
	}
}

var SslCertificatesServiceInstanceCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ssl_certificates_service_instances.default",
	ExistMapFunc: existSslCertificatesServiceInstanceMapFunc,
	FakeMapFunc:  fakeSslCertificatesServiceInstanceMapFunc,
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

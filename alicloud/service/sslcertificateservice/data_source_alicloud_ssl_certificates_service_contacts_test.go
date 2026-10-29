// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sslcertificateservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSslCertificatesServiceContactDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}"]`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	NameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"name":           `"${var.name}"`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"name":           `"${var.name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}"]`,
			"name":           `"${var.name}"`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}_fake"]`,
			"name":           `"${var.name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	SslCertificatesServiceContactCheckInfo.DataSourceTestCheck(t, rand, idsConf, NameConf, allConf)
}

var existSslCertificatesServiceContactMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#":            "1",
		"contacts.0.contact_id": acctest.CHECKSET,
		"contacts.0.id":         acctest.CHECKSET,
		"contacts.0.name":       acctest.CHECKSET,
		"contacts.0.email":      acctest.CHECKSET,
		"contacts.0.mobile":     acctest.CHECKSET,
	}
}

var fakeSslCertificatesServiceContactMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#": "0",
	}
}

var SslCertificatesServiceContactCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ssl_certificates_service_contacts.default",
	ExistMapFunc: existSslCertificatesServiceContactMapFunc,
	FakeMapFunc:  fakeSslCertificatesServiceContactMapFunc,
}

func testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslCertificatesServiceContact%d"
}


resource "alicloud_ssl_certificates_service_contact" "default" {
  name   = var.name
  mobile = "13312345678"
  email  = "test@example.com"
}

data "alicloud_ssl_certificates_service_contacts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

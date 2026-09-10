// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSslCertificatesServiceContactDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	NameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"name":           `"${var.name}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"name":           `"${var.name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}"]`,
			"name":           `"${var.name}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceContactSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ssl_certificates_service_contact.default.id}_fake"]`,
			"name":           `"${var.name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	SslCertificatesServiceContactCheckInfo.dataSourceTestCheck(t, rand, idsConf, NameConf, allConf)
}

var existSslCertificatesServiceContactMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#":            "1",
		"contacts.0.contact_id": CHECKSET,
		"contacts.0.id":         CHECKSET,
		"contacts.0.name":       CHECKSET,
		"contacts.0.email":      CHECKSET,
		"contacts.0.mobile":     CHECKSET,
	}
}

var fakeSslCertificatesServiceContactMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#": "0",
	}
}

var SslCertificatesServiceContactCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_certificates_service_contacts.default",
	existMapFunc: existSslCertificatesServiceContactMapFunc,
	fakeMapFunc:  fakeSslCertificatesServiceContactMapFunc,
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

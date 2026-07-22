// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudSslCertificatesServiceCompanyDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_certificates_service_company.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_certificates_service_company.default.id}_fake"]`,
		}),
	}

	companyIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"company_id": `"${alicloud_ssl_certificates_service_company.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"company_id": "999999999",
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"keyword": `"测试公司"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"keyword": `"nonexistent_keyword_xyz"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ssl_certificates_service_company.default.id}"]`,
			"company_id": `"${alicloud_ssl_certificates_service_company.default.id}"`,
			"keyword":    `"测试公司"`,
		}),
		fakeConfig: testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ssl_certificates_service_company.default.id}_fake"]`,
			"company_id": "999999999",
			"keyword":    `"nonexistent_keyword_xyz"`,
		}),
	}

	SslCertificatesServiceCompanyCheckInfo.dataSourceTestCheck(t, rand, idsConf, companyIdConf, keywordConf, allConf)
}

var existSslCertificatesServiceCompanyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"companies.#":                 "1",
		"companies.0.company_id":      CHECKSET,
		"companies.0.company_email":   CHECKSET,
		"companies.0.lang":            CHECKSET,
		"companies.0.city":            CHECKSET,
		"companies.0.company_type":    CHECKSET,
		"companies.0.province":        CHECKSET,
		"companies.0.company_address": CHECKSET,
		"companies.0.company_name":    CHECKSET,
		"companies.0.department":      CHECKSET,
		"companies.0.country_code":    CHECKSET,
		"companies.0.post_code":       CHECKSET,
		"companies.0.company_code":    CHECKSET,
		"companies.0.company_phone":   CHECKSET,
	}
}

var fakeSslCertificatesServiceCompanyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"companies.#": "0",
	}
}

var SslCertificatesServiceCompanyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_certificates_service_companies.default",
	existMapFunc: existSslCertificatesServiceCompanyMapFunc,
	fakeMapFunc:  fakeSslCertificatesServiceCompanyMapFunc,
}

func testAccCheckAlicloudSslCertificatesServiceCompanySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslCertificatesServiceCompany%d"
}


resource "alicloud_ssl_certificates_service_company" "default" {
        company_address = "西安市"
        company_name = "测试公司1"
        department = "测试部门1"
        city = "西安"
        company_type = "1"
        country_code = "111122"
        post_code = "11112233"
        company_code = "12312311"
        company_phone = "15101081174"
        province = "陕西"
        lang = "zh"
        company_email = "test@example.com"
}

data "alicloud_ssl_certificates_service_companies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

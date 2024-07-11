package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDirectMailDomainsDataSource_basic0(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids": `[alicloud_direct_mail_domain.default.id]`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids": `["${alicloud_direct_mail_domain.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}_fake"`,
		}),
	}

	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}.fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":   `"1"`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":   `"0"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_domain.default.id]`,
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"key_word":   `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":     `1`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":        `["${alicloud_direct_mail_domain.default.id}_fake"]`,
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}_fake"`,
			"key_word":   `"${alicloud_direct_mail_domain.default.domain_name}.fake"`,
			"status":     `0`,
		}),
	}

	var existAliCloudDirectMailDomainsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"domains.#":                      "1",
			"domains.0.id":                   CHECKSET,
			"domains.0.domain_id":            CHECKSET,
			"domains.0.domain_name":          CHECKSET,
			"domains.0.domain_record":        CHECKSET,
			"domains.0.domain_type":          "",
			"domains.0.cname_auth_status":    CHECKSET,
			"domains.0.cname_confirm_status": "",
			"domains.0.cname_record":         "",
			"domains.0.icp_status":           CHECKSET,
			"domains.0.mx_auth_status":       CHECKSET,
			"domains.0.mx_record":            "",
			"domains.0.spf_auth_status":      CHECKSET,
			"domains.0.spf_record":           "",
			"domains.0.default_domain":       "",
			"domains.0.host_record":          "",
			"domains.0.dns_mx":               "",
			"domains.0.dns_txt":              "",
			"domains.0.dns_spf":              "",
			"domains.0.dns_dmarc":            "",
			"domains.0.dkim_auth_status":     "",
			"domains.0.dkim_rr":              "",
			"domains.0.dkim_public_key":      "",
			"domains.0.dmarc_auth_status":    "",
			"domains.0.dmarc_record":         "",
			"domains.0.dmarc_host_record":    "",
			"domains.0.tl_domain_name":       "",
			"domains.0.tracef_record":        "",
			"domains.0.status":               "1",
			"domains.0.create_time":          CHECKSET,
		}
	}

	var fakeAliCloudDirectMailDomainsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"domains.#": "0",
		}
	}

	var aliCloudDirectMailDomainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_direct_mail_domains.default",
		existMapFunc: existAliCloudDirectMailDomainsDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudDirectMailDomainsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}

	aliCloudDirectMailDomainsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keyWordConf, statusConf, allConf)
}

func TestAccAliCloudDirectMailDomainsDataSource_basic1(t *testing.T) {
	rand := acctest.RandInt()

	idsDetailConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"enable_details": `false`,
		}),
	}

	nameRegexDetailConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `false`,
		}),
	}

	keyWordDetailConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `false`,
		}),
	}

	statusDetailConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `"1"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `"1"`,
			"enable_details": `false`,
		}),
	}

	allDetailConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `1`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `1`,
			"enable_details": `false`,
		}),
	}

	var existAliCloudDirectMailDomainsDetailDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"domains.#":                      "1",
			"domains.0.id":                   CHECKSET,
			"domains.0.domain_id":            CHECKSET,
			"domains.0.domain_name":          CHECKSET,
			"domains.0.domain_record":        CHECKSET,
			"domains.0.domain_type":          CHECKSET,
			"domains.0.cname_auth_status":    CHECKSET,
			"domains.0.cname_confirm_status": CHECKSET,
			"domains.0.cname_record":         CHECKSET,
			"domains.0.icp_status":           CHECKSET,
			"domains.0.mx_auth_status":       CHECKSET,
			"domains.0.mx_record":            CHECKSET,
			"domains.0.spf_auth_status":      CHECKSET,
			"domains.0.spf_record":           CHECKSET,
			"domains.0.default_domain":       CHECKSET,
			"domains.0.host_record":          CHECKSET,
			"domains.0.dns_mx":               "",
			"domains.0.dns_txt":              "",
			"domains.0.dns_spf":              "",
			"domains.0.dns_dmarc":            "",
			"domains.0.dkim_auth_status":     CHECKSET,
			"domains.0.dkim_rr":              CHECKSET,
			"domains.0.dkim_public_key":      CHECKSET,
			"domains.0.dmarc_auth_status":    CHECKSET,
			"domains.0.dmarc_record":         CHECKSET,
			"domains.0.dmarc_host_record":    CHECKSET,
			"domains.0.tl_domain_name":       CHECKSET,
			"domains.0.tracef_record":        CHECKSET,
			"domains.0.status":               "1",
			"domains.0.create_time":          CHECKSET,
		}
	}

	var existAliCloudDirectMailDomainsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"domains.#":                      "1",
			"domains.0.id":                   CHECKSET,
			"domains.0.domain_id":            CHECKSET,
			"domains.0.domain_name":          CHECKSET,
			"domains.0.domain_record":        CHECKSET,
			"domains.0.domain_type":          "",
			"domains.0.cname_auth_status":    CHECKSET,
			"domains.0.cname_confirm_status": "",
			"domains.0.cname_record":         "",
			"domains.0.icp_status":           CHECKSET,
			"domains.0.mx_auth_status":       CHECKSET,
			"domains.0.mx_record":            "",
			"domains.0.spf_auth_status":      CHECKSET,
			"domains.0.spf_record":           "",
			"domains.0.default_domain":       "",
			"domains.0.host_record":          "",
			"domains.0.dns_mx":               "",
			"domains.0.dns_txt":              "",
			"domains.0.dns_spf":              "",
			"domains.0.dns_dmarc":            "",
			"domains.0.dkim_auth_status":     "",
			"domains.0.dkim_rr":              "",
			"domains.0.dkim_public_key":      "",
			"domains.0.dmarc_auth_status":    "",
			"domains.0.dmarc_record":         "",
			"domains.0.dmarc_host_record":    "",
			"domains.0.tl_domain_name":       "",
			"domains.0.tracef_record":        "",
			"domains.0.status":               "1",
			"domains.0.create_time":          CHECKSET,
		}
	}

	var aliCloudDirectMailDomainsDetailCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_direct_mail_domains.default",
		existMapFunc: existAliCloudDirectMailDomainsDetailDataSourceMapFunc,
		fakeMapFunc:  existAliCloudDirectMailDomainsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}

	aliCloudDirectMailDomainsDetailCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsDetailConf, nameRegexDetailConf, keyWordDetailConf, statusDetailConf, allDetailConf)
}

func testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
 		default = "tf-testacc%d.pop.com"
	}

	resource "alicloud_direct_mail_domain" "default" {
 		domain_name = var.name
	}

	data "alicloud_direct_mail_domains" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}

package directmail_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDirectMailDomainsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids": `[alicloud_direct_mail_domain.default.id]`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids": `["${alicloud_direct_mail_domain.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}_fake"`,
		}),
	}

	keyWordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}.fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":   `"1"`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":   `"0"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_domain.default.id]`,
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"key_word":   `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":     `1`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
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
			"domains.0.id":                   acctest.CHECKSET,
			"domains.0.domain_id":            acctest.CHECKSET,
			"domains.0.domain_name":          acctest.CHECKSET,
			"domains.0.domain_record":        acctest.CHECKSET,
			"domains.0.domain_type":          "",
			"domains.0.cname_auth_status":    acctest.CHECKSET,
			"domains.0.cname_confirm_status": "",
			"domains.0.cname_record":         "",
			"domains.0.icp_status":           acctest.CHECKSET,
			"domains.0.mx_auth_status":       acctest.CHECKSET,
			"domains.0.mx_record":            "",
			"domains.0.spf_auth_status":      acctest.CHECKSET,
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
			"domains.0.create_time":          acctest.CHECKSET,
		}
	}

	var fakeAliCloudDirectMailDomainsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"domains.#": "0",
		}
	}

	var aliCloudDirectMailDomainsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_direct_mail_domains.default",
		ExistMapFunc: existAliCloudDirectMailDomainsDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudDirectMailDomainsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}

	aliCloudDirectMailDomainsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keyWordConf, statusConf, allConf)
}

func TestAccAliCloudDirectMailDomainsDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandInt()

	idsDetailConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"enable_details": `false`,
		}),
	}

	nameRegexDetailConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `false`,
		}),
	}

	keyWordDetailConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"enable_details": `false`,
		}),
	}

	statusDetailConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `"1"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `"1"`,
			"enable_details": `false`,
		}),
	}

	allDetailConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
			"ids":            `[alicloud_direct_mail_domain.default.id]`,
			"name_regex":     `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"key_word":       `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":         `1`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudDirectMailDomainsDefaultDataSource(rand, map[string]string{
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
			"domains.0.id":                   acctest.CHECKSET,
			"domains.0.domain_id":            acctest.CHECKSET,
			"domains.0.domain_name":          acctest.CHECKSET,
			"domains.0.domain_record":        acctest.CHECKSET,
			"domains.0.domain_type":          acctest.CHECKSET,
			"domains.0.cname_auth_status":    acctest.CHECKSET,
			"domains.0.cname_confirm_status": acctest.CHECKSET,
			"domains.0.cname_record":         acctest.CHECKSET,
			"domains.0.icp_status":           acctest.CHECKSET,
			"domains.0.mx_auth_status":       acctest.CHECKSET,
			"domains.0.mx_record":            acctest.CHECKSET,
			"domains.0.spf_auth_status":      acctest.CHECKSET,
			"domains.0.spf_record":           acctest.CHECKSET,
			"domains.0.default_domain":       acctest.CHECKSET,
			"domains.0.host_record":          acctest.CHECKSET,
			"domains.0.dns_mx":               "",
			"domains.0.dns_txt":              "",
			"domains.0.dns_spf":              "",
			"domains.0.dns_dmarc":            "",
			"domains.0.dkim_auth_status":     acctest.CHECKSET,
			"domains.0.dkim_rr":              acctest.CHECKSET,
			"domains.0.dkim_public_key":      acctest.CHECKSET,
			"domains.0.dmarc_auth_status":    acctest.CHECKSET,
			"domains.0.dmarc_record":         acctest.CHECKSET,
			"domains.0.dmarc_host_record":    acctest.CHECKSET,
			"domains.0.tl_domain_name":       acctest.CHECKSET,
			"domains.0.tracef_record":        acctest.CHECKSET,
			"domains.0.status":               "1",
			"domains.0.create_time":          acctest.CHECKSET,
		}
	}

	var existAliCloudDirectMailDomainsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"domains.#":                      "1",
			"domains.0.id":                   acctest.CHECKSET,
			"domains.0.domain_id":            acctest.CHECKSET,
			"domains.0.domain_name":          acctest.CHECKSET,
			"domains.0.domain_record":        acctest.CHECKSET,
			"domains.0.domain_type":          "",
			"domains.0.cname_auth_status":    acctest.CHECKSET,
			"domains.0.cname_confirm_status": "",
			"domains.0.cname_record":         "",
			"domains.0.icp_status":           acctest.CHECKSET,
			"domains.0.mx_auth_status":       acctest.CHECKSET,
			"domains.0.mx_record":            "",
			"domains.0.spf_auth_status":      acctest.CHECKSET,
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
			"domains.0.create_time":          acctest.CHECKSET,
		}
	}

	var aliCloudDirectMailDomainsDetailCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_direct_mail_domains.default",
		ExistMapFunc: existAliCloudDirectMailDomainsDetailDataSourceMapFunc,
		FakeMapFunc:  existAliCloudDirectMailDomainsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}

	aliCloudDirectMailDomainsDetailCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsDetailConf, nameRegexDetailConf, keyWordDetailConf, statusDetailConf, allDetailConf)
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

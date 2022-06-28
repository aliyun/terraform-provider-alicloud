package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDirectMailDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_domain.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_domain.default.id]`,
			"status": `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_domain.default.id]`,
			"status": `"0"`,
		}),
	}
	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_domain.default.domain_name}.fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_domain.default.id]`,
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"key_word":   `"${alicloud_direct_mail_domain.default.domain_name}"`,
			"status":     `1`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailDomainsDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_direct_mail_domain.default.domain_name}_fake"`,
			"key_word":   `"${alicloud_direct_mail_domain.default.domain_name}.fake"`,
			"status":     `0`,
		}),
	}
	var existAlicloudDirectMailDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"names.#":               "1",
			"domains.#":             "1",
			"domains.0.domain_id":   CHECKSET,
			"domains.0.domain_name": fmt.Sprintf("tf-testacc%d.pop.com", rand),
			"domains.0.status":      `1`,
		}
	}
	var fakeAlicloudDirectMailDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDirectMailDomainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_direct_mail_domains.default",
		existMapFunc: existAlicloudDirectMailDomainsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDirectMailDomainsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}
	alicloudDirectMailDomainsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, keyWordConf, allConf)
}
func testAccCheckAlicloudDirectMailDomainsDataSourceName(rand int, attrMap map[string]string) string {
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

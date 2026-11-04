package dcdn_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDcdnIpaDomainsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_ipa_domain.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_ipa_domain.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status": `"online"`,
		}),
		FakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status": `"offline"`,
		}),
	}
	domainNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}-fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status":      `"online"`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status":      `"offline"`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}-fake"`,
		}),
	}
	var existAlicloudDcdnIpaDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"domains.#":                    "1",
			"domains.0.domain_name":        fmt.Sprintf("tf-testacccn-%d.alicloud-provider.cn", rand),
			"domains.0.cert_name":          "",
			"domains.0.cname":              acctest.CHECKSET,
			"domains.0.create_time":        acctest.CHECKSET,
			"domains.0.description":        "",
			"domains.0.id":                 acctest.CHECKSET,
			"domains.0.resource_group_id":  acctest.CHECKSET,
			"domains.0.scope":              "overseas",
			"domains.0.sources.#":          "1",
			"domains.0.sources.0.content":  "www.alicloud-provider.cn",
			"domains.0.sources.0.port":     "8898",
			"domains.0.sources.0.priority": "20",
			"domains.0.sources.0.type":     "domain",
			"domains.0.sources.0.weight":   "10",
			"domains.0.ssl_protocol":       acctest.CHECKSET,
			"domains.0.ssl_pub":            "",
			"domains.0.status":             "online",
		}
	}
	var fakeAlicloudDcdnIpaDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDcdnIpaDomainsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dcdn_ipa_domains.default",
		ExistMapFunc: existAlicloudDcdnIpaDomainsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDcdnIpaDomainsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	alicloudDcdnIpaDomainsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, domainNameConf, allConf)
}
func testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "domain_name" {	
	default = "tf-testacccn-%d.alicloud-provider.cn"
}

resource "alicloud_dcdn_ipa_domain" "default" {
	domain_name = "${var.domain_name}"
	scope = "overseas"
    sources {
		content =  "www.alicloud-provider.cn"
		port =     8898
		priority = "20"
		type =     "domain"
		weight =   10
	}
}

data "alicloud_dcdn_ipa_domains" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

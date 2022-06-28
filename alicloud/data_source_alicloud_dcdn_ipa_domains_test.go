package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDCDNIpaDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_ipa_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_ipa_domain.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status": `"online"`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status": `"offline"`,
		}),
	}
	domainNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}-fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_dcdn_ipa_domain.default.id}"]`,
			"status":      `"online"`,
			"domain_name": `"${alicloud_dcdn_ipa_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand, map[string]string{
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
			"domains.0.domain_name":        fmt.Sprintf("tf-testacccn-%d.xiaozhu.com", rand),
			"domains.0.cert_name":          "",
			"domains.0.cname":              CHECKSET,
			"domains.0.create_time":        CHECKSET,
			"domains.0.description":        "",
			"domains.0.id":                 CHECKSET,
			"domains.0.resource_group_id":  CHECKSET,
			"domains.0.scope":              "domestic",
			"domains.0.sources.#":          "1",
			"domains.0.sources.0.content":  "www.xiaozhu.com",
			"domains.0.sources.0.port":     "8898",
			"domains.0.sources.0.priority": "20",
			"domains.0.sources.0.type":     "domain",
			"domains.0.sources.0.weight":   "10",
			"domains.0.ssl_protocol":       CHECKSET,
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
	var alicloudDcdnIpaDomainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dcdn_ipa_domains.default",
		existMapFunc: existAlicloudDcdnIpaDomainsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDcdnIpaDomainsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	alicloudDcdnIpaDomainsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, domainNameConf, allConf)
}
func testAccCheckAlicloudDcdnIpaDomainsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "domain_name" {	
	default = "tf-testacccn-%d.xiaozhu.com"
}

resource "alicloud_dcdn_ipa_domain" "default" {
	domain_name = "${var.domain_name}"
	scope = "domestic"
    sources {
		content =  "www.xiaozhu.com"
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

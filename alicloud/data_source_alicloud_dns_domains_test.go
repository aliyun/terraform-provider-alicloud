package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDnsDomainsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	aliDomainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"ali_domain":        `"false"`,
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"ali_domain":        `"true"`,
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
		}),
	}
	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"ali_domain":        `"false"`,
			"group_name_regex":  `"${alicloud_dns_group.default.name}"`,
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"ali_domain":        `"false"`,
			"group_name_regex":  `"${alicloud_dns_group.default.name}_fake"`,
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
			"instance_id":       `""`,
		}),
		fakeConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
			"instance_id":       `"fake"`,
		}),
	}
	versionCodeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
			"version_code":      `"mianfei"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
			"version_code":      `"bumianfei"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
			"version_code":      `"mianfei"`,
			"instance_id":       `""`,
			"ali_domain":        `"false"`,
			"group_name_regex":  `"${alicloud_dns_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsDomainsDataSourceConfig(rand, map[string]string{
			"domain_name_regex": `"${alicloud_dns.default.name}"`,
			"version_code":      `"mianfei"`,
			"instance_id":       `""`,
			"ali_domain":        `"true"`,
			"group_name_regex":  `"${alicloud_dns_group.default.name}"`,
		}),
	}
	dnsDomainsCheckInfo.dataSourceTestCheck(t, rand, aliDomainConf, groupNameConf, instanceIdConf, versionCodeConf, allConf)
}

func testAccCheckAlicloudDnsDomainsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "alicloud_dns_group" "default" {
	name = "tf-testaccdns%d"
}

resource "alicloud_dns" "default" {
	name = "tf-testaccdnsalidomain%d.abc"
	group_id = "${alicloud_dns_group.default.id}"
}

data "alicloud_dns_domains" "default" {
	%s
}`, rand, rand, strings.Join(pairs, "\n  "))
	return config
}

var existDnsDomainsMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"domains.#":               "1",
		"domains.0.domain_id":     CHECKSET,
		"domains.0.domain_name":   fmt.Sprintf("tf-testaccdnsalidomain%d.abc", rand),
		"domains.0.ali_domain":    "false",
		"domains.0.group_id":      CHECKSET,
		"domains.0.group_name":    fmt.Sprintf("tf-testaccdns%d", rand),
		"domains.0.instance_id":   "",
		"domains.0.version_code":  "mianfei",
		"domains.0.puny_code":     CHECKSET,
		"domains.0.dns_servers.#": CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"names.#":                 "1",
		"names.0":                 fmt.Sprintf("tf-testaccdnsalidomain%d.abc", rand),
	}
}

var fakeDnsDomainsMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"names.#":   "0",
		"ids.#":     "0",
		"domains.#": "0",
	}
}

var dnsDomainsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dns_domains.default",
	existMapFunc: existDnsDomainsMapCheck,
	fakeMapFunc:  fakeDnsDomainsMapCheck,
}

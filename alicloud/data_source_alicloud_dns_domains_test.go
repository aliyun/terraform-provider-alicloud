package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"strconv"
	"testing"
)

func TestAccAlicloudDnsDomainsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_dns_domains.default", strconv.FormatInt(int64(rand), 10), dataSourceDnsDomainsConfigDependence)
	aliDomainConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "false",
			"domain_name_regex": "${alicloud_dns.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "true",
			"domain_name_regex": "${alicloud_dns.default.name}",
		}),
	}
	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "false",
			"group_name_regex":  "${alicloud_dns_group.default.name}",
			"domain_name_regex": "${alicloud_dns.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "false",
			"group_name_regex":  "${alicloud_dns_group.default.name}_fake",
			"domain_name_regex": "${alicloud_dns.default.name}",
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns.default.name}",
			"instance_id":       "",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns.default.name}",
			"instance_id":       "fake",
		}),
	}
	versionCodeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns.default.name}",
			"version_code":      "mianfei",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns.default.name}",
			"version_code":      "bumianfei",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns.default.name}",
			"version_code":      "mianfei",
			"instance_id":       "",
			"ali_domain":        "false",
			"group_name_regex":  "${alicloud_dns_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns.default.name}",
			"version_code":      "mianfei",
			"instance_id":       "",
			"ali_domain":        "true",
			"group_name_regex":  "${alicloud_dns_group.default.name}",
		}),
	}
	dnsDomainsCheckInfo.dataSourceTestCheck(t, rand, aliDomainConf, groupNameConf, instanceIdConf, versionCodeConf, allConf)
}

func dataSourceDnsDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "dnsName"{
	default = "tf-testacc%sdnsalidomainbasic%s.abc"
}

variable "dnsGroupName"{
	default = "tf-testaccdns%s"
}

resource "alicloud_dns_group" "default" {
  name = "${var.dnsGroupName}"
}

resource "alicloud_dns" "default" {
	name = "${var.dnsName}"
	group_id = "${alicloud_dns_group.default.id}"
}
`, defaultRegionToTest, name, name)
}

var existDnsDomainsMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"domains.#":               "1",
		"domains.0.domain_id":     CHECKSET,
		"domains.0.domain_name":   fmt.Sprintf("tf-testacc%sdnsalidomainbasic%d.abc", defaultRegionToTest, rand),
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
		"names.0":                 fmt.Sprintf("tf-testacc%sdnsalidomainbasic%d.abc", defaultRegionToTest, rand),
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

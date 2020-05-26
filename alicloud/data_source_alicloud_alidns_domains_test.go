package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlidnsDomainsSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"ids": `["${alicloud_alidns_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"ids": `["${alicloud_alidns_domain.default.id}-fake"]`,
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}-fake"`,
		}),
	}

	groupIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
			"group_id":   `"${alicloud_alidns_domain.default.group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
			"group_id":   `"${alicloud_alidns_domain.default.group_id}-fake"`,
		}),
	}

	tagConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
			"tags": `{
							Created = "Terraform"
							For 	= "test"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
			"tags": `{
							Created = "Terraform-fake"
							For 	= "test-fake"
					  }`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"ids":        `["${alicloud_alidns_domain.default.id}"]`,
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
			"group_id":   `"${alicloud_alidns_domain.default.group_id}"`,
			"tags": `{
							Created = "Terraform"
							For 	= "test"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand, map[string]string{
			"ids":        `["${alicloud_alidns_domain.default.id}"]`,
			"name_regex": `"${alicloud_alidns_domain.default.domain_name}"`,
			"group_id":   `"${alicloud_alidns_domain.default.group_id}"`,
			"tags": `{
							Created = "Terraform-fake"
							For 	= "test-fake"
					  }`,
		}),
	}

	var existAlidnsDomainsMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"ids.0":                  CHECKSET,
			"names.#":                "1",
			"names.0":                fmt.Sprintf("tf-testadd%salidnsdomainbasic%d.abc", defaultRegionToTest, rand),
			"domains.#":              "1",
			"domains.0.domain_id":    CHECKSET,
			"domains.0.group_id":     CHECKSET,
			"domains.0.domain_name":  fmt.Sprintf("tf-testadd%salidnsdomainbasic%d.abc", defaultRegionToTest, rand),
			"domains.0.tags.%":       "2",
			"domains.0.tags.Created": "Terraform",
			"domains.0.tags.For":     "test",
		}
	}

	var fakeAlidnsDomainsMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"domains.#": "0",
		}
	}

	var alidnsDomainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_domains.default",
		fakeMapFunc:  fakeAlidnsDomainsMapCheck,
		existMapFunc: existAlidnsDomainsMapCheck,
	}

	alidnsDomainsCheckInfo.dataSourceTestCheck(t, rand, idConfig, nameRegexConfig, groupIdConfig, tagConfig, allConfig)

}

func testAccCheckAlicloudAlidnsDomainsDataSourceConfigBaisc(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "dnsName" {
	default = "tf-testadd%salidnsdomainbasic%d.abc"
}

variable "dnsGroupName" {
	default = "tf-testaccdns%d"
}

resource "alicloud_alidns_domain_group" "default" {
	group_name = "${var.dnsGroupName}"
}

resource "alicloud_alidns_domain" "default" {
	domain_name = "${var.dnsName}"
	group_id = "${alicloud_alidns_domain_group.default.id}"
	resource_group_id = "%s"
	tags = {
		Created = "Terraform"
		For 	= "test"
	  }
}

data "alicloud_alidns_domains" "default" {
	%s
}`, defaultRegionToTest, rand, rand, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

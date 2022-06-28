package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVodDomainDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vod_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vod_domain.default.domain_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vod_domain.default.id}"]`,
			"status": `"${alicloud_vod_domain.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"status": `"offline"`,
			"ids":    `["${alicloud_vod_domain.default.id}"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}"]`,
			"tags": `{ 
						"key1" = "value1"
						"key2" = "value2" 
						"Tftestacc123"= "Tftest123"
					}`,
		}),
		fakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}"]`,
			"tags": `{ 
						"key1" = "value1_fake"
						"key2" = "value2_fake" 
						"Tftestacc123"= "Tftest123"
					}`,
		}),
	}
	domainSearchTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_vod_domain.default.id}"]`,
			"domain_search_type": `"fuzzy_match"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vod_domain.default.id}"]`,
			"name_regex": `"${alicloud_vod_domain.default.domain_name}"`,
			"status":     `"${alicloud_vod_domain.default.status}"`,
			"tags": `{ 
						"key1" = "value1"
						"key2" = "value2"
						"Tftestacc123"= "Tftest123"
					}`,
		}),
		fakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vod_domain.default.id}_fake"]`,
			"name_regex": `"${alicloud_vod_domain.default.domain_name}_fake"`,
			"status":     `"offline"`,
			"tags": `{ 
						"key1" = "value1_fake"
						"key2" = "value2_fake"
						"Tftestacc123" = "Tftest123"
					}`,
		}),
	}
	var existAlicloudVodDomainDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"domains.#":                       "1",
			"domains.0.domain_name":           "kftwh.com",
			"domains.0.sources.0.source_type": "domain",
		}
	}
	var fakeAlicloudVodDomainDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVodDomainCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vod_domains.default",
		existMapFunc: existAlicloudVodDomainDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVodDomainDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
	}
	alicloudVodDomainCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, tagsConf, domainSearchTypeConf, allConf)

}
func testAccCheckAlicloudVodDomainDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccVodDomain-%d"
}

resource "alicloud_vod_domain" "default" {
  domain_name = "kftwh.com"
  scope       = "domestic"
  sources {
    source_type    = "domain"
    source_content = "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com"
    source_port    = "80"
  }
  tags = {
    key1 = "value1"
    key2 = "value2"
	Tftestacc123 = "Tftest123"
  }
}


data "alicloud_vod_domains" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

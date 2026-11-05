package vod_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudVodDomainDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vod_domain.default.domain_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vod_domain.default.domain_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vod_domain.default.id}"]`,
			"status": `"${alicloud_vod_domain.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"status": `"offline"`,
			"ids":    `["${alicloud_vod_domain.default.id}"]`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}"]`,
			"tags": `{ 
						"key1" = "value1"
						"key2" = "value2" 
						"Tftestacc123"= "Tftest123"
					}`,
		}),
		FakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vod_domain.default.id}"]`,
			"tags": `{ 
						"key1" = "value1_fake"
						"key2" = "value2_fake" 
						"Tftestacc123"= "Tftest123"
					}`,
		}),
	}
	domainSearchTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_vod_domain.default.id}"]`,
			"domain_search_type": `"fuzzy_match"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vod_domain.default.id}"]`,
			"name_regex": `"${alicloud_vod_domain.default.domain_name}"`,
			"status":     `"${alicloud_vod_domain.default.status}"`,
			"tags": `{ 
						"key1" = "value1"
						"key2" = "value2"
						"Tftestacc123"= "Tftest123"
					}`,
		}),
		FakeConfig: testAccCheckAlicloudVodDomainDataSourceName(rand, map[string]string{
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
	var alicloudVodDomainCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vod_domains.default",
		ExistMapFunc: existAlicloudVodDomainDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudVodDomainDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.VodSupportRegions)
	}
	// Inline resource.TestCase with ProtoV5ProviderFactories instead of the shared
	// acctest.DataSourceTestCheckWithPreCheck helper (which still pins the legacy
	// ProviderFactories field); the step/check construction below mirrors that
	// helper verbatim using exported acctest helpers. See migration-standard §5.
	var vodDomainDataSourceSteps []resource.TestStep
	for _, conf := range []acctest.DataSourceTestAccConfig{idsConf, nameRegexConf, statusConf, tagsConf, domainSearchTypeConf, allConf} {
		existCheck := acctest.ResourceAttrInit(alicloudVodDomainCheckInfo.ResourceId, alicloudVodDomainCheckInfo.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
		emptyCheck := acctest.ResourceAttrInit(alicloudVodDomainCheckInfo.ResourceId, alicloudVodDomainCheckInfo.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
		if conf.ExistConfig != "" {
			vodDomainDataSourceSteps = append(vodDomainDataSourceSteps, resource.TestStep{
				Config: conf.ExistConfig,
				Check: resource.ComposeTestCheckFunc(
					existCheck(conf.ExistChangMap),
				),
			})
		}
		if conf.FakeConfig != "" {
			vodDomainDataSourceSteps = append(vodDomainDataSourceSteps, resource.TestStep{
				Config: conf.FakeConfig,
				Check: resource.ComposeTestCheckFunc(
					emptyCheck(conf.FakeChangMap),
				),
			})
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 preCheck,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    vodDomainDataSourceSteps,
	})
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

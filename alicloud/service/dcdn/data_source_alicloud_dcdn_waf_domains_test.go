package dcdn_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDcdnWafDomainsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 999999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDcdnWafDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_domain.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDcdnWafDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_domain.default.id}_fake"]`,
		}),
	}
	var existAlicloudDcdnWafDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"domains.#":                  "1",
			"domains.0.domain_name":      acctest.CHECKSET,
			"domains.0.id":               acctest.CHECKSET,
			"domains.0.client_ip_tag":    "X-Forwarded-For",
			"domains.0.defense_scenes.#": "1",
			"domains.0.defense_scenes.0.defense_scene": acctest.CHECKSET,
			"domains.0.defense_scenes.0.policy_id":     acctest.CHECKSET,
		}
	}
	var fakeAlicloudDcdnWafDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudDcdnWafDomainsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dcdn_waf_domains.default",
		ExistMapFunc: existAlicloudDcdnWafDomainsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDcdnWafDomainsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDcdnWafDomainsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudDcdnWafDomainsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "domain_name" {
  default = "tf-testacc%s%d.alicloud-provider.cn"
}
resource "alicloud_dcdn_domain" "default" {
  domain_name = var.domain_name
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
  }
}
resource "alicloud_dcdn_waf_domain" "default" {
	domain_name = alicloud_dcdn_domain.default.domain_name
	client_ip_tag = "X-Forwarded-For"
}

data "alicloud_dcdn_waf_domains" "default" {	
	enable_details = true
	%s
}
`, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}

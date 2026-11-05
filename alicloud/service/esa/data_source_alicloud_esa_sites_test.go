package esa_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEsaSiteDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_esa_site.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_esa_site.default.id}_fake"]`,
		}),
	}

	ResourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	SiteNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_esa_site.default.id}"]`,
			"site_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_esa_site.default.id}_fake"]`,
			"site_name": `"${var.name}_fake"`,
		}),
	}
	CoverageConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_esa_site.default.id}"]`,
			"coverage": `"overseas"`,
		}),
		FakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_esa_site.default.id}_fake"]`,
			"coverage": `"overseas_fake"`,
		}),
	}
	AccessTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_esa_site.default.id}"]`,
			"access_type": `"NS"`,
		}),
		FakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_esa_site.default.id}_fake"]`,
			"access_type": `"NS_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,

			"site_name": `"${var.name}"`,

			"coverage": `"overseas"`,

			"access_type": `"NS"`,
		}),
		FakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,

			"site_name": `"${var.name}_fake"`,

			"coverage": `"overseas_fake"`,

			"access_type": `"NS_fake"`,
		}),
	}

	EsaSiteCheckInfo.DataSourceTestCheck(t, rand, idsConf, ResourceGroupIdConf, SiteNameConf, CoverageConf, AccessTypeConf, allConf)
}

var existEsaSiteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sites.#":                   "1",
		"sites.0.modify_time":       acctest.CHECKSET,
		"sites.0.site_id":           acctest.CHECKSET,
		"sites.0.name_server_list":  acctest.CHECKSET,
		"sites.0.site_name":         acctest.CHECKSET,
		"sites.0.resource_group_id": acctest.CHECKSET,
		"sites.0.access_type":       acctest.CHECKSET,
		"sites.0.status":            acctest.CHECKSET,
		"sites.0.instance_id":       acctest.CHECKSET,
		"sites.0.create_time":       acctest.CHECKSET,
		"sites.0.coverage":          acctest.CHECKSET,
	}
}

var fakeEsaSiteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sites.#": "0",
	}
}

var EsaSiteCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_esa_sites.default",
	ExistMapFunc: existEsaSiteMapFunc,
	FakeMapFunc:  fakeEsaSiteMapFunc,
}

func testAccCheckAlicloudEsaSiteSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "bcd%d.com"
}
resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

resource "alicloud_esa_site" "default" {
  site_name   = var.name
  coverage    = "overseas"
  access_type = "NS"
  instance_id = alicloud_esa_rate_plan_instance.defaultIEoDfU.id
  tags = {
    testkey1 = "testvalue1"
    testkey2 = "testvalue2"
    testkey3 = "testvalue3"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

data "alicloud_esa_sites" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

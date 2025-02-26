package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEsaSiteDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_esa_site.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_esa_site.default.id}_fake"]`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	SiteNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_esa_site.default.id}"]`,
			"site_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_esa_site.default.id}_fake"]`,
			"site_name": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,

			"site_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEsaSiteSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_esa_site.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,

			"site_name": `"${var.name}_fake"`,
		}),
	}

	EsaSiteCheckInfo.dataSourceTestCheck(t, rand, idsConf, ResourceGroupIdConf, SiteNameConf, allConf)
}

var existEsaSiteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sites.#":                   "1",
		"sites.0.status":            CHECKSET,
		"sites.0.modify_time":       CHECKSET,
		"sites.0.site_id":           CHECKSET,
		"sites.0.name_server_list":  CHECKSET,
		"sites.0.site_name":         CHECKSET,
		"sites.0.resource_group_id": CHECKSET,
		"sites.0.instance_id":       CHECKSET,
		"sites.0.create_time":       CHECKSET,
		"sites.0.coverage":          CHECKSET,
		"sites.0.access_type":       CHECKSET,
	}
}

var fakeEsaSiteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sites.#": "0",
	}
}

var EsaSiteCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_esa_sites.default",
	existMapFunc: existEsaSiteMapFunc,
	fakeMapFunc:  fakeEsaSiteMapFunc,
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

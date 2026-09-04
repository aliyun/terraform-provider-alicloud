package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApigPortalDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPortalSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_portal.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudApigPortalSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_portal.default.id}_fake"]`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPortalSourceConfig(rand, map[string]string{
			"name": `"${alicloud_apig_portal.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigPortalSourceConfig(rand, map[string]string{
			"name": `"${alicloud_apig_portal.default.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPortalSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_portal.default.id}"]`,
			"name":           `"${alicloud_apig_portal.default.name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudApigPortalSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_portal.default.id}_fake"]`,
			"name":           `"${alicloud_apig_portal.default.name}_fake"`,
			"enable_details": `true`,
		}),
	}

	ApigPortalCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, allConf)
}

var existApigPortalMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"portals.#":             "1",
		"portals.0.id":          CHECKSET,
		"portals.0.portal_id":   CHECKSET,
		"portals.0.name":        CHECKSET,
		"portals.0.description": CHECKSET,
	}
}

var fakeApigPortalMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"portals.#": "0",
		"ids.#":     "0",
	}
}

var ApigPortalCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_portals.default",
	existMapFunc: existApigPortalMapFunc,
	fakeMapFunc:  fakeApigPortalMapFunc,
}

func testAccCheckAliCloudApigPortalSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tfaccapigportal%d"
}

resource "alicloud_apig_portal" "default" {
  name        = var.name
  description = "tfacc portal datasource"
}

data "alicloud_apig_portals" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

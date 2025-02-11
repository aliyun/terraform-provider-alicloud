package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEfloVpdDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.EfloSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_vpd.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_vpd.default.id}_fake"]`,
		}),
	}
	VpdNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_eflo_vpd.default.id}"]`,
			"vpd_name": `"${alicloud_eflo_vpd.default.vpd_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"vpd_name": `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		}),
	}
	VpdIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}"]`,
			"vpd_id": `"${alicloud_eflo_vpd.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"vpd_id": `"${alicloud_eflo_vpd.default.id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_eflo_vpd.default.id}"]`,
			"name_regex": `"${alicloud_eflo_vpd.default.vpd_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"name_regex": `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}"]`,
			"status": `"${alicloud_eflo_vpd.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"status": `"Not Available"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_vpd.default.id}"]`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_vpd.default.id}"]`,
			"vpd_name":          `"${alicloud_eflo_vpd.default.vpd_name}"`,
			"name_regex":        `"${alicloud_eflo_vpd.default.vpd_name}"`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}"`,
			"status":            `"${alicloud_eflo_vpd.default.status}"`,
		}),
		// There is a api bug, and product will fix it in the next three weeks. Please open it after next submit
		//fakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
		//	"ids":               `["${alicloud_eflo_vpd.default.id}_fake"]`,
		//	"vpd_name":          `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		//	"name_regex":        `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		//	"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}_fake"`,
		//	"status":            `"Not Available"`,
		//}),
	}

	EfloVpdCheckInfo.dataSourceTestCheck(t, rand, idsConf, VpdNameConf, VpdIdConf, nameRegexConf, statusConf, resourceGroupIdConf, allConf)
}

var existEfloVpdMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                    "1",
		"names.#":                  "1",
		"vpds.#":                   "1",
		"vpds.0.id":                CHECKSET,
		"vpds.0.cidr":              "10.0.0.0/8",
		"vpds.0.create_time":       CHECKSET,
		"vpds.0.gmt_modified":      CHECKSET,
		"vpds.0.resource_group_id": CHECKSET,
		"vpds.0.status":            CHECKSET,
		"vpds.0.vpd_id":            CHECKSET,
		"vpds.0.vpd_name":          CHECKSET,
	}
}

var fakeEfloVpdMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vpds.#":  "0",
	}
}

var EfloVpdCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_eflo_vpds.default",
	existMapFunc: existEfloVpdMapFunc,
	fakeMapFunc:  fakeEfloVpdMapFunc,
}

func testAccCheckAlicloudEfloVpdSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEfloVpd%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_eflo_vpd" "default" {
  cidr      = "10.0.0.0/8"
  vpd_name  = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

data "alicloud_eflo_vpds" "default" {
  enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

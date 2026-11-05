package eflo_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEfloVpdDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EfloSupportRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_vpd.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_vpd.default.id}_fake"]`,
		}),
	}
	VpdNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_eflo_vpd.default.id}"]`,
			"vpd_name": `"${alicloud_eflo_vpd.default.vpd_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"vpd_name": `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		}),
	}
	VpdIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}"]`,
			"vpd_id": `"${alicloud_eflo_vpd.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"vpd_id": `"${alicloud_eflo_vpd.default.id}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_eflo_vpd.default.id}"]`,
			"name_regex": `"${alicloud_eflo_vpd.default.vpd_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"name_regex": `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}"]`,
			"status": `"${alicloud_eflo_vpd.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"status": `"Not Available"`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_vpd.default.id}"]`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_vpd.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_vpd.default.id}"]`,
			"vpd_name":          `"${alicloud_eflo_vpd.default.vpd_name}"`,
			"name_regex":        `"${alicloud_eflo_vpd.default.vpd_name}"`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}"`,
			"status":            `"${alicloud_eflo_vpd.default.status}"`,
		}),
		// There is a api bug, and product will fix it in the next three weeks. Please open it after next submit
		//FakeConfig: testAccCheckAlicloudEfloVpdSourceConfig(rand, map[string]string{
		//	"ids":               `["${alicloud_eflo_vpd.default.id}_fake"]`,
		//	"vpd_name":          `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		//	"name_regex":        `"${alicloud_eflo_vpd.default.vpd_name}_fake"`,
		//	"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}_fake"`,
		//	"status":            `"Not Available"`,
		//}),
	}

	EfloVpdCheckInfo.DataSourceTestCheck(t, rand, idsConf, VpdNameConf, VpdIdConf, nameRegexConf, statusConf, resourceGroupIdConf, allConf)
}

var existEfloVpdMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                    "1",
		"names.#":                  "1",
		"vpds.#":                   "1",
		"vpds.0.id":                acctest.CHECKSET,
		"vpds.0.cidr":              "10.0.0.0/8",
		"vpds.0.create_time":       acctest.CHECKSET,
		"vpds.0.gmt_modified":      acctest.CHECKSET,
		"vpds.0.resource_group_id": acctest.CHECKSET,
		"vpds.0.status":            acctest.CHECKSET,
		"vpds.0.vpd_id":            acctest.CHECKSET,
		"vpds.0.vpd_name":          acctest.CHECKSET,
	}
}

var fakeEfloVpdMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vpds.#":  "0",
	}
}

var EfloVpdCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_eflo_vpds.default",
	ExistMapFunc: existEfloVpdMapFunc,
	FakeMapFunc:  fakeEfloVpdMapFunc,
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

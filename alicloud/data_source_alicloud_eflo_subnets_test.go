package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEfloSubnetDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.EfloSupportRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_subnet.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_subnet.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_eflo_subnet.default.id}"]`,
			"name_regex": `"${alicloud_eflo_subnet.default.subnet_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_eflo_subnet.default.id}"]`,
			"name_regex": `"${alicloud_eflo_subnet.default.subnet_name}_fake"`,
		}),
	}
	VpdIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_subnet.default.id}"]`,
			"vpd_id": `"${alicloud_eflo_vpd.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_subnet.default.id}"]`,
			"vpd_id": `"${alicloud_eflo_vpd.default.id}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_subnet.default.id}"]`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_subnet.default.id}"]`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_subnet.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_eflo_subnet.default.id}"]`,
			"status": `"Executing"`,
		}),
	}
	SubnetIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"subnet_id": `"${alicloud_eflo_subnet.default.subnet_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"subnet_id": `"${alicloud_eflo_subnet.default.subnet_id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_eflo_subnet.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_eflo_subnet.default.id}"]`,
			"type": `"OOB"`,
		}),
	}
	SubnetNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_eflo_subnet.default.id}"]`,
			"subnet_name": `"${alicloud_eflo_subnet.default.subnet_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_eflo_subnet.default.id}"]`,
			"subnet_name": `"${alicloud_eflo_subnet.default.subnet_name}_fake"`,
		}),
	}
	ZoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_eflo_subnet.default.id}"]`,
			"zone_id": `"${alicloud_eflo_subnet.default.zone_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_eflo_subnet.default.id}"]`,
			"zone_id": `"${alicloud_eflo_subnet.default.zone_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_eflo_subnet.default.id}"]`,
			"vpd_id":            `"${alicloud_eflo_vpd.default.id}"`,
			"subnet_name":       `"${alicloud_eflo_subnet.default.subnet_name}"`,
			"zone_id":           `"${alicloud_eflo_subnet.default.zone_id}"`,
			"subnet_id":         `"${alicloud_eflo_subnet.default.subnet_id}"`,
			"status":            `"Available"`,
			"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}"`,
			"name_regex":        `"${alicloud_eflo_subnet.default.subnet_name}"`,
		}),
		// There is a api bug, and product will fix it in the next three weeks. Please open it after next submit
		//fakeConfig: testAccCheckAlicloudEfloSubnetSourceConfig(rand, map[string]string{
		//	"ids":               `["${alicloud_eflo_subnet.default.id}_fake"]`,
		//	"vpd_id":            `"${alicloud_eflo_vpd.default.id}_fake"`,
		//	"subnet_name":       `"${alicloud_eflo_subnet.default.subnet_name}_fake"`,
		//	"zone_id":           `"${alicloud_eflo_subnet.default.zone_id}_fake"`,
		//	"subnet_id":         `"${alicloud_eflo_subnet.default.subnet_id}_fake"`,
		//	"status":            `"Executing"`,
		//	"type":              `"OOB"`,
		//	"resource_group_id": `"${alicloud_eflo_vpd.default.resource_group_id}_fake"`,
		//	"name_regex":        `"${alicloud_eflo_subnet.default.subnet_name}_fake"`,
		//}),
	}

	EfloSubnetCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, VpdIdConf, resourceGroupIdConf, statusConf, SubnetIdConf, typeConf, SubnetNameConf, ZoneIdConf, allConf)
}

var existEfloSubnetMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"names.#":                     "1",
		"subnets.#":                   "1",
		"subnets.0.id":                CHECKSET,
		"subnets.0.cidr":              CHECKSET,
		"subnets.0.create_time":       CHECKSET,
		"subnets.0.gmt_modified":      CHECKSET,
		"subnets.0.message":           CHECKSET,
		"subnets.0.resource_group_id": CHECKSET,
		"subnets.0.status":            CHECKSET,
		"subnets.0.subnet_id":         CHECKSET,
		"subnets.0.subnet_name":       CHECKSET,
		"subnets.0.type":              "",
		"subnets.0.vpd_id":            CHECKSET,
		"subnets.0.zone_id":           CHECKSET,
	}
}

var fakeEfloSubnetMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"names.#":   "0",
		"subnets.#": "0",
	}
}

var EfloSubnetCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_eflo_subnets.default",
	existMapFunc: existEfloSubnetMapFunc,
	fakeMapFunc:  fakeEfloSubnetMapFunc,
}

func testAccCheckAlicloudEfloSubnetSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEfloSubnet%d"
}

data "alicloud_zones" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eflo_vpd" "default" {
  cidr      = "10.0.0.0/8"
  vpd_name  = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

resource "alicloud_eflo_subnet" "default" {
  subnet_name = var.name
  zone_id     = data.alicloud_zones.default.zones.0.id
  cidr        =  "10.0.0.0/16"
  vpd_id      = alicloud_eflo_vpd.default.id
}

data "alicloud_eflo_subnets" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

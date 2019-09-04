package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudVSwitchesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vswitch.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vswitch.default.id}_fake" ]`,
		}),
	}

	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"cidr_block": `"172.16.0.0/24"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"cidr_block": `"172.16.0.0/23"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"is_default": `"true"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"vpc_id":     `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"vpc_id":     `"${alicloud_vpc.default.id}_fake"`,
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"zone_id":    `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"zone_id":    `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test"
					  }`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"ids":        `[ "${alicloud_vswitch.default.id}" ]`,
			"cidr_block": `"172.16.0.0/24"`,
			"is_default": `"false"`,
			"vpc_id":     `"${alicloud_vpc.default.id}"`,
			"zone_id":    `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.name}"`,
			"ids":        `[ "${alicloud_vswitch.default.id}" ]`,
			"cidr_block": `"172.16.0.0/24"`,
			"is_default": `"false"`,
			"vpc_id":     `"${alicloud_vpc.default.id}"`,
			"zone_id":    `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}

	vswitchesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, cidrBlockConf, idDefaultConf, vpcIdConf, zoneIdConf, tagsConf, allConf)

}

func testAccCheckAlicloudVSwitchesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVSwitchDatasource%d"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.default.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}

data "alicloud_vswitches" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVSwitchesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                      "1",
		"names.#":                    "1",
		"vswitches.#":                "1",
		"vswitches.0.id":             CHECKSET,
		"vswitches.0.vpc_id":         CHECKSET,
		"vswitches.0.zone_id":        CHECKSET,
		"vswitches.0.name":           fmt.Sprintf("tf-testAccVSwitchDatasource%d", rand),
		"vswitches.0.instance_ids.#": "0",
		"vswitches.0.cidr_block":     "172.16.0.0/24",
		"vswitches.0.description":    "",
		"vswitches.0.is_default":     "false",
		"vswitches.0.creation_time":  CHECKSET,
	}
}

var fakeVSwitchesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"vswitches.#": "0",
	}
}

var vswitchesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vswitches.default",
	existMapFunc: existVSwitchesMapFunc,
	fakeMapFunc:  fakeVSwitchesMapFunc,
}

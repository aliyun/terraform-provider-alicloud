package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCVSwitchesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}_fake"`,
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
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"ids":    `[ "${alicloud_vswitch.default.id}" ]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"ids":    `[ "${alicloud_vswitch.default.id}_fake" ]`,
			"status": `"Pending"`,
		}),
	}

	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"cidr_block": `"172.16.0.0/24"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"cidr_block": `"172.16.0.0/23"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"is_default": `"true"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"vpc_id":     `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"vpc_id":     `"${alicloud_vpc.default.id}_fake"`,
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"zone_id":    `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"zone_id":    `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test"
					  }`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vswitch.default.vswitch_name}"`,
			// The resource route tables do not support resource_group_id, so it was set empty.
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_vswitch.default.vswitch_name}"`,
			"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_vswitch.default.vswitch_name}"`,
			"ids":               `[ "${alicloud_vswitch.default.id}" ]`,
			"cidr_block":        `"172.16.0.0/24"`,
			"is_default":        `"false"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"zone_id":           `"${data.alicloud_zones.default.zones.0.id}"`,
			"resource_group_id": `""`,
			"status":            `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_vswitch.default.vswitch_name}"`,
			"ids":               `[ "${alicloud_vswitch.default.id}" ]`,
			"cidr_block":        `"172.16.0.0/24"`,
			"is_default":        `"false"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"zone_id":           `"${data.alicloud_zones.default.zones.0.id}_fake"`,
			"resource_group_id": `""`,
			"status":            `"Pending"`,
		}),
	}

	vswitchesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, statusConf, cidrBlockConf, idDefaultConf, vpcIdConf, zoneIdConf, tagsConf, resourceGroupIdConf, allConf)

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
  vpc_name = "${var.name}"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}"
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
		"ids.#":                                  "1",
		"names.#":                                "1",
		"vswitches.#":                            "1",
		"vswitches.0.id":                         CHECKSET,
		"vswitches.0.vswitch_id":                 CHECKSET,
		"vswitches.0.vpc_id":                     CHECKSET,
		"vswitches.0.zone_id":                    CHECKSET,
		"vswitches.0.name":                       fmt.Sprintf("tf-testAccVSwitchDatasource%d", rand),
		"vswitches.0.vswitch_name":               fmt.Sprintf("tf-testAccVSwitchDatasource%d", rand),
		"vswitches.0.cidr_block":                 "172.16.0.0/24",
		"vswitches.0.description":                "",
		"vswitches.0.is_default":                 "false",
		"vswitches.0.creation_time":              CHECKSET,
		"vswitches.0.tags.%":                     "2",
		"vswitches.0.tags.Created":               "TF",
		"vswitches.0.tags.For":                   "acceptance test",
		"vswitches.0.status":                     "Available",
		"vswitches.0.route_table_id":             CHECKSET,
		"vswitches.0.resource_group_id":          CHECKSET,
		"vswitches.0.available_ip_address_count": CHECKSET,
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

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudVpcsDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	initVswitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d"`, rand),
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d_fake"`, rand),
		}),
	}
	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/12"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/0"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Pending"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"true"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${alicloud_vswitch.default.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/12"`,
			"status":     `"Available"`,
			"is_default": `"false"`,
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/16"`,
			"status":     `"Available"`,
			"is_default": `"false"`,
			"vswitch_id": `"${alicloud_vswitch.default.id}_fake"`,
		}),
	}

	vpcsCheckInfo.dataSourceTestCheck(t, rand, initVswitchConf, nameRegexConf, cidrBlockConf, statusConf, idDefaultConf, vswitchIdConf, allConf)
}

func testAccCheckAlicloudVpcsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVpcsdatasource%d"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

data "alicloud_zones" "default" {

}

resource "alicloud_vswitch" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
	vpc_id = "${alicloud_vpc.default.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_vpcs" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"names.#":               "1",
		"vpcs.#":                "1",
		"vpcs.0.id":             CHECKSET,
		"vpcs.0.region_id":      CHECKSET,
		"vpcs.0.status":         "Available",
		"vpcs.0.vpc_name":       fmt.Sprintf("tf-testAccVpcsdatasource%d", rand),
		"vpcs.0.vswitch_ids.#":  "1",
		"vpcs.0.cidr_block":     "172.16.0.0/12",
		"vpcs.0.vrouter_id":     CHECKSET,
		"vpcs.0.route_table_id": CHECKSET,
		"vpcs.0.description":    "",
		"vpcs.0.is_default":     "false",
		"vpcs.0.creation_time":  CHECKSET,
	}
}

var fakeVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vpcs.#":  "0",
	}
}

var vpcsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpcs.default",
	existMapFunc: existVpcsMapFunc,
	fakeMapFunc:  fakeVpcsMapFunc,
}

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudLindormInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	queryStrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_lindorm_instance.default.id}"]`,
			"query_str": `"${alicloud_lindorm_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_lindorm_instance.default.id}_fake"]`,
			"query_str": `"${alicloud_lindorm_instance.default.instance_name}_fake"`,
		}),
	}
	supportEngineStrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_lindorm_instance.default.id}"]`,
			"support_engine": `"4"`,
		}),
		fakeConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_lindorm_instance.default.id}_fake"]`,
			"support_engine": `"4"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_lindorm_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_lindorm_instance.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_lindorm_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_lindorm_instance.default.instance_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_lindorm_instance.default.id}"]`,
			"status": `"ACTIVATION"`,
		}),
		fakeConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_lindorm_instance.default.id}"]`,
			"status": `"CREATING"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_lindorm_instance.default.id}"]`,
			"name_regex":     `"${alicloud_lindorm_instance.default.instance_name}"`,
			"status":         `"ACTIVATION"`,
			"support_engine": `"4"`,
			"query_str":      `"${alicloud_lindorm_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudLindormInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_lindorm_instance.default.id}_fake"]`,
			"name_regex":     `"${alicloud_lindorm_instance.default.instance_name}_fake"`,
			"status":         `"CREATING"`,
			"support_engine": `"4"`,
			"query_str":      `"${alicloud_lindorm_instance.default.instance_name}_fake"`,
		}),
	}
	var existAlicloudLindormInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"instances.#":               "1",
			"instances.0.instance_name": fmt.Sprintf("tf-testAccLindormInstances-%d", rand),
		}
	}
	var fakeAlicloudLindormInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var alicloudLindormInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_lindorm_instances.default",
		existMapFunc: existAlicloudLindormInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudLindormInstancesDataSourceNameMapFunc,
	}
	alicloudLindormInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, nameRegexConf, supportEngineStrConf, queryStrConf, allConf)
}
func testAccCheckAlicloudLindormInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testAccLindormInstances-%d"
}
provider "alicloud" {
  region = "cn-shenzhen"
}
locals {
  zone_id = "cn-shenzhen-e"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id = local.zone_id
  vswitch_name              = var.name
}
resource "alicloud_lindorm_instance" "default" {
	disk_category = "cloud_efficiency"
	payment_type   =           "PayAsYouGo"
	zone_id =                   local.zone_id
	vswitch_id =                length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	instance_name =            var.name
	table_engine_specification = "lindorm.c.xlarge"
	table_engine_node_count =  "2"
	instance_storage   =       "480"
}
data "alicloud_lindorm_instances" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

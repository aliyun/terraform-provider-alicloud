package oceanbasepro_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudOceanBaseInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ocean_base_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ocean_base_instance.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ocean_base_instance.default.id}"]`,
			"status": `"${alicloud_ocean_base_instance.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ocean_base_instance.default.id}"]`,
			"status": `"PREPAID_EXPIRE_CLOSED"`,
		}),
	}
	instanceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_ocean_base_instance.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_ocean_base_instance.default.id}_fake"`,
		}),
	}
	instanceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ocean_base_instance.default.id}"]`,
			"instance_name": `"${alicloud_ocean_base_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ocean_base_instance.default.id}"]`,
			"instance_name": `"${alicloud_ocean_base_instance.default.instance_name}_fake"`,
		}),
	}
	searchKeyConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ocean_base_instance.default.id}"]`,
			"search_key": `"${alicloud_ocean_base_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ocean_base_instance.default.id}"]`,
			"search_key": `"${alicloud_ocean_base_instance.default.instance_name}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ocean_base_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ocean_base_instance.default.instance_name}_fake"`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ocean_base_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_ocean_base_instance.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ocean_base_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_ocean_base_instance.default.resource_group_id}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ocean_base_instance.default.id}"]`,
			"instance_id":       `"${alicloud_ocean_base_instance.default.id}"`,
			"instance_name":     `"${alicloud_ocean_base_instance.default.instance_name}"`,
			"name_regex":        `"${alicloud_ocean_base_instance.default.instance_name}"`,
			"resource_group_id": `"${alicloud_ocean_base_instance.default.resource_group_id}"`,
			"status":            `"${alicloud_ocean_base_instance.default.status}"`,
			"search_key":        `"${alicloud_ocean_base_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ocean_base_instance.default.id}_fake"]`,
			"instance_id":       `"${alicloud_ocean_base_instance.default.id}_fake"`,
			"instance_name":     `"${alicloud_ocean_base_instance.default.instance_name}_fake"`,
			"name_regex":        `"${alicloud_ocean_base_instance.default.instance_name}_fake"`,
			"resource_group_id": `"${alicloud_ocean_base_instance.default.resource_group_id}_fake"`,
			"status":            `"PREPAID_EXPIRE_CLOSED"`,
			"search_key":        `"${alicloud_ocean_base_instance.default.instance_name}_fake"`,
		}),
	}
	var existAlicloudOceanBaseInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.commodity_code":    acctest.CHECKSET,
			"instances.0.create_time":       acctest.CHECKSET,
			"instances.0.cpu":               acctest.CHECKSET,
			"instances.0.disk_size":         "100",
			"instances.0.instance_class":    "8C32GB",
			"instances.0.instance_id":       acctest.CHECKSET,
			"instances.0.instance_name":     acctest.CHECKSET,
			"instances.0.payment_type":      "PayAsYouGo",
			"instances.0.resource_group_id": acctest.CHECKSET,
			"instances.0.series":            "normal",
			"instances.0.status":            "ONLINE",
			"instances.0.zones.#":           "3",
			"instances.0.node_num":          "3",
		}
	}
	var fakeAlicloudOceanBaseInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}
	var alicloudOceanBaseInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ocean_base_instances.default",
		ExistMapFunc: existAlicloudOceanBaseInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudOceanBaseInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudOceanBaseInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, instanceIdConf, instanceNameConf, searchKeyConf, nameRegexConf, resourceGroupIdConf, allConf)
}
func testAccCheckAlicloudOceanBaseInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testacc-%d"
}

data "alicloud_zones" "default" {}

resource "alicloud_ocean_base_instance" "default" {
  instance_name  = var.name
  series         = "normal"
  disk_size      = 100
  instance_class = "8C32GB"
  zones = ["${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"]
  payment_type = "PayAsYouGo"
}

data "alicloud_ocean_base_instances" "default" {
enable_details = true
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

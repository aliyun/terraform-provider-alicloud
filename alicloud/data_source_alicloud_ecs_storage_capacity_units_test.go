package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSStorageCapacityUnitsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_storage_capacity_unit.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_storage_capacity_unit.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_storage_capacity_unit.default.storage_capacity_unit_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_storage_capacity_unit.default.storage_capacity_unit_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_storage_capacity_unit.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_storage_capacity_unit.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecs_storage_capacity_unit.default.id}"]`,
			"name_regex": `"${alicloud_ecs_storage_capacity_unit.default.storage_capacity_unit_name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecs_storage_capacity_unit.default.id}_fake"]`,
			"name_regex": `"${alicloud_ecs_storage_capacity_unit.default.storage_capacity_unit_name}_fake"`,
			"status":     `"Creating"`,
		}),
	}
	var existAlicloudEcsStorageCapacityUnitsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"units.#":                            "1",
			"units.0.allocation_status":          "",
			"units.0.capacity":                   "20",
			"units.0.create_time":                CHECKSET,
			"units.0.description":                fmt.Sprintf("tf-testAccStorageCapacityUnit-%d", rand),
			"units.0.expired_time":               CHECKSET,
			"units.0.start_time":                 CHECKSET,
			"units.0.status":                     CHECKSET,
			"units.0.id":                         CHECKSET,
			"units.0.storage_capacity_unit_id":   CHECKSET,
			"units.0.storage_capacity_unit_name": fmt.Sprintf("tf-testAccStorageCapacityUnit-%d", rand),
		}
	}
	var fakeAlicloudEcsStorageCapacityUnitsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcsStorageCapacityUnitsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_storage_capacity_units.default",
		existMapFunc: existAlicloudEcsStorageCapacityUnitsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsStorageCapacityUnitsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsStorageCapacityUnitsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEcsStorageCapacityUnitsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccStorageCapacityUnit-%d"
}

resource "alicloud_ecs_storage_capacity_unit" "default" {
	capacity                   = 20
	description                = var.name
	storage_capacity_unit_name = var.name
}

data "alicloud_ecs_storage_capacity_units" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

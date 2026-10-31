package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenRouteMapsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_route_map.default.id]`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `["fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_route_map.default.id]`,
			"status": `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_route_map.default.id]`,
			"status": `"Creating"`,
		}),
	}

	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"description_regex": `"datasource_test"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"description_regex": `"datasource_test_fake"`,
		}),
	}

	cenRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_route_map.default.id]`,
			"cen_region_id ": `"${var.child_region}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_route_map.default.id]`,
			"cen_region_id ": `"${var.child_region}_fake"`,
		}),
	}

	transmitDirectionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"transmit_direction": `"RegionIn"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"transmit_direction": `"RegionOut"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"cen_region_id ":     `"${var.child_region}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"status":             `"Active"`,
			"transmit_direction": `"RegionIn"`,
			"description_regex":  `"datasource_test"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"cen_region_id ":     `"${var.child_region}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"status":             `"Active"`,
			"transmit_direction": `"RegionIn"`,
			"description_regex":  `"datasource_test_fake"`,
		}),
	}

	var existCenRouteMapsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"maps.#":               "1",
			"ids.#":                "1",
			"maps.0.id":            acctest.CHECKSET,
			"maps.0.cen_region_id": acctest.DefaultRegionToTest,
			"maps.0.description":   "datasource_test",
			"maps.0.destination_instance_ids_reverse_match": acctest.CHECKSET,
			"maps.0.map_result":                             "Permit",
			"maps.0.priority":                               "3",
			"maps.0.source_instance_ids_reverse_match":      acctest.CHECKSET,
			"maps.0.status":                                 "Active",
			"maps.0.transmit_direction":                     "RegionIn",
			"maps.0.route_map_id":                           acctest.CHECKSET,
		}
	}

	var fakeCenRouteMapsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"maps.#": "0",
			"ids.#":  "0",
		}
	}

	var cenRouteMapsRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_route_maps.default",
		ExistMapFunc: existCenRouteMapsRecordsMapFunc,
		FakeMapFunc:  fakeCenRouteMapsRecordsMapFunc,
	}

	cenRouteMapsRecordsCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf, descriptionConf, cenRegionIdConf, transmitDirectionConf, allConf)

}

func testAccCheckAlicloudCenRouteMapsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCenRouteMaps-%d"
}

variable "child_region" {
  default = "%s"
}

resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_route_map" "default" {
  depends_on = [alicloud_cen_transit_router.default]
  cen_id = "${alicloud_cen_instance.default.id}"
  cen_region_id = "${var.child_region}"
  map_result = "Permit"
  priority = 3
  transmit_direction = "RegionIn"
  description = "datasource_test"
}

data "alicloud_cen_route_maps" "default" {
%s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}

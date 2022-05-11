package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenRouteMapsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_route_map.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `["fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_route_map.default.id]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_route_map.default.id]`,
			"status": `"Creating"`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"description_regex": `"datasource_test"`,
		}),
		fakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"description_regex": `"datasource_test_fake"`,
		}),
	}

	cenRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_route_map.default.id]`,
			"cen_region_id ": `"${var.child_region}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_route_map.default.id]`,
			"cen_region_id ": `"${var.child_region}_fake"`,
		}),
	}

	transmitDirectionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"transmit_direction": `"RegionIn"`,
		}),
		fakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"transmit_direction": `"RegionOut"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"cen_region_id ":     `"${var.child_region}"`,
			"ids":                `[alicloud_cen_route_map.default.id]`,
			"status":             `"Active"`,
			"transmit_direction": `"RegionIn"`,
			"description_regex":  `"datasource_test"`,
		}),
		fakeConfig: testAccCheckAlicloudCenRouteMapsSourceConfig(rand, map[string]string{
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
			"maps.0.id":            CHECKSET,
			"maps.0.cen_region_id": defaultRegionToTest,
			"maps.0.description":   "datasource_test",
			"maps.0.destination_instance_ids_reverse_match": CHECKSET,
			"maps.0.map_result":                             "Permit",
			"maps.0.priority":                               "3",
			"maps.0.source_instance_ids_reverse_match":      CHECKSET,
			"maps.0.status":                                 "Active",
			"maps.0.transmit_direction":                     "RegionIn",
			"maps.0.route_map_id":                           CHECKSET,
		}
	}

	var fakeCenRouteMapsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"maps.#": "0",
			"ids.#":  "0",
		}
	}

	var cenRouteMapsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_route_maps.default",
		existMapFunc: existCenRouteMapsRecordsMapFunc,
		fakeMapFunc:  fakeCenRouteMapsRecordsMapFunc,
	}

	cenRouteMapsRecordsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, descriptionConf, cenRegionIdConf, transmitDirectionConf, allConf)

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

resource "alicloud_cen_route_map" "default" {
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
`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}

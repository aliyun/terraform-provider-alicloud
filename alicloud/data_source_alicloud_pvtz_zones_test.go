package alicloud

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPvtzZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_pvtz_zones.default"
	keywordConf := dataSourceTestAccConfig{
		existConfig: dataSourcePvtzZonesConfigDependence(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.basic.zone_name}"`,
		}),
		fakeConfig: dataSourcePvtzZonesConfigDependence(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.basic.zone_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: dataSourcePvtzZonesConfigDependence(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone.basic.id}"]`,
		}),
		fakeConfig: dataSourcePvtzZonesConfigDependence(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone.basic.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: dataSourcePvtzZonesConfigDependence(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.basic.zone_name}"`,
			"ids":     `["${alicloud_pvtz_zone.basic.id}"]`,
		}),
		fakeConfig: dataSourcePvtzZonesConfigDependence(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.basic.zone_name}_fake"`,
			"ids":     `["${alicloud_pvtz_zone.basic.id}_fake"]`,
		}),
	}

	var existPvtzZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ids.0":                    CHECKSET,
			"names.#":                  "1",
			"names.0":                  fmt.Sprintf("tf-testacc%d.test.com", rand),
			"zones.#":                  "1",
			"zones.0.zone_name":        fmt.Sprintf("tf-testacc%d.test.com", rand),
			"zones.0.remark":           "",
			"zones.0.record_count":     "0",
			"zones.0.is_ptr":           "false",
			"zones.0.create_timestamp": CHECKSET,
			"zones.0.update_timestamp": CHECKSET,
			"zones.0.bind_vpcs.#":      "0",
		}
	}

	var fakePvtzZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"zones.#": "0",
		}
	}

	var pvtzZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPvtzZonesMapFunc,
		fakeMapFunc:  fakePvtzZonesMapFunc,
	}

	pvtzZonesCheckInfo.dataSourceTestCheck(t, rand, keywordConf, idsConf, allConf)
}

func dataSourcePvtzZonesConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "basic" {
		zone_name = "tf-testacc%d.test.com"
	}
	data "alicloud_pvtz_zones" "default"{
		%s
	} 
	`, rand, strings.Join(pairs, " \n "))
}

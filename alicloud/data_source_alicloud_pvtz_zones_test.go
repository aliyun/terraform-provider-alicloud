package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudPvtzZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_pvtz_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%d.test.com", rand),
		dataSourcePvtzZonesConfigDependence)

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword": alicloud_pvtz_zone.basic.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword": "${alicloud_pvtz_zone.basic.name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{alicloud_pvtz_zone.basic.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_pvtz_zone.basic.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword": alicloud_pvtz_zone.basic.name,
			"ids":     []string{alicloud_pvtz_zone.basic.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword": alicloud_pvtz_zone.basic.name,
			"ids":     []string{"${alicloud_pvtz_zone.basic.id}-fake"},
		}),
	}

	var existPvtzZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 CHECKSET,
			"names.#":               "1",
			"names.0":               fmt.Sprintf("tf-testacc%d.test.com", rand),
			"zones.#":               "1",
			"zones.0.name":          fmt.Sprintf("tf-testacc%d.test.com", rand),
			"zones.0.remark":        "",
			"zones.0.record_count":  "0",
			"zones.0.is_ptr":        "false",
			"zones.0.creation_time": CHECKSET,
			"zones.0.update_time":   CHECKSET,
			"zones.0.bind_vpcs.#":   "0",
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

func dataSourcePvtzZonesConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "basic" {
		name = "%s"
	}
	`, name)
}

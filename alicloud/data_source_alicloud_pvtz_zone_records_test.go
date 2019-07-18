package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudPvtzZoneRecordsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_pvtz_zone_records.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%d.test.com", rand),
		dataSourcePvtzZoneRecordsConfigDependence)

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
		}),
	}
	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
			"keyword": "${alicloud_pvtz_zone_record.foo.value}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
			"keyword": "${alicloud_pvtz_zone_record.foo.value}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
			"ids":     []string{"${alicloud_pvtz_zone_record.foo.record_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
			"ids":     []string{"${alicloud_pvtz_zone_record.foo.record_id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
			"ids":     []string{"${alicloud_pvtz_zone_record.foo.record_id}"},
			"keyword": "${alicloud_pvtz_zone_record.foo.value}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alicloud_pvtz_zone_record.foo.zone_id}",
			"ids":     []string{"${alicloud_pvtz_zone_record.foo.record_id}"},
			"keyword": "${alicloud_pvtz_zone_record.foo.value}-fake",
		}),
	}

	var existPvtzZoneRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"records.#":                 "1",
			"records.0.id":              CHECKSET,
			"records.0.resource_record": "www",
			"records.0.type":            "A",
			"records.0.ttl":             "60",
			"records.0.priority":        "0",
			"records.0.value":           "2.2.2.2",
			"records.0.status":          "ENABLE",
		}
	}

	var fakePvtzZoneRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"records.#": "0",
		}
	}

	var pvtzZoneRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPvtzZoneRecordsMapFunc,
		fakeMapFunc:  fakePvtzZoneRecordsMapFunc,
	}

	pvtzZoneRecordsCheckInfo.dataSourceTestCheck(t, rand, zoneIdConf, keyWordConf, idsConf, allConf)
}

func dataSourcePvtzZoneRecordsConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "basic" {
		name = "%s"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.basic.id}"
		resource_record = "www"
		type = "A"
		value = "2.2.2.2"
		ttl = "60"
	}
	`, name)
}

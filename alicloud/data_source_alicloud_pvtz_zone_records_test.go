package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudPvtzZoneRecordsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.TestPvtzRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone_record.default.record_id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone_record.default.record_id}_fake"]`,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone_record.default.rr}"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone_record.default.rr}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"status": `"ENABLE"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"status": `"DISABLE"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_pvtz_zone_record.default.record_id}"]`,
			"keyword":        `"${alicloud_pvtz_zone_record.default.rr}"`,
			"user_client_ip": `"127.0.0.1"`,
			"status":         `"ENABLE"`,
			"search_mode":    `"EXACT"`,
			"lang":           `"en"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_pvtz_zone_record.default.record_id}_fake"]`,
			"keyword":     `"${alicloud_pvtz_zone_record.default.rr}_fake"`,
			"tag":         `"ecs"`,
			"status":      `"DISABLE"`,
			"search_mode": `"LIKE"`,
			"lang":        `"zh"`,
		}),
	}

	var existAliCloudPvtzZoneRecordsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"records.#":                 "1",
			"records.0.id":              CHECKSET,
			"records.0.record_id":       CHECKSET,
			"records.0.priority":        CHECKSET,
			"records.0.remark":          CHECKSET,
			"records.0.rr":              CHECKSET,
			"records.0.resource_record": CHECKSET,
			"records.0.ttl":             CHECKSET,
			"records.0.type":            CHECKSET,
			"records.0.value":           CHECKSET,
			"records.0.status":          CHECKSET,
		}
	}

	var fakeAliCloudPvtzZoneRecordsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"records.#": "0",
		}
	}

	var aliCloudPvtzZoneRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pvtz_zone_records.default",
		existMapFunc: existAliCloudPvtzZoneRecordsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudPvtzZoneRecordsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudPvtzZoneRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf, statusConf, allConf)
}

func testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone" "default" {
  		zone_name = var.name
	}

	resource "alicloud_pvtz_zone_record" "default" {
  		zone_id  = alicloud_pvtz_zone.default.id
  		rr       = "www"
  		type     = "MX"
  		value    = var.name
  		ttl      = "60"
  		priority = 2
  		remark   = var.name
	}

	data "alicloud_pvtz_zone_records" "default" {
  		zone_id = alicloud_pvtz_zone_record.default.zone_id
		%s
	}
	`, rand, strings.Join(pairs, " \n "))
}

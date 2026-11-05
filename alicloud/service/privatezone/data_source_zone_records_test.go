package privatezone_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudPvtzZoneRecordsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone_record.default.record_id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone_record.default.record_id}_fake"]`,
		}),
	}

	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone_record.default.rr}"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone_record.default.rr}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"status": `"ENABLE"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"status": `"DISABLE"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_pvtz_zone_record.default.record_id}"]`,
			"keyword":        `"${alicloud_pvtz_zone_record.default.rr}"`,
			"user_client_ip": `"127.0.0.1"`,
			"status":         `"ENABLE"`,
			"search_mode":    `"EXACT"`,
			"lang":           `"en"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZoneRecordsDataSourceName(rand, map[string]string{
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
			"records.0.id":              acctest.CHECKSET,
			"records.0.record_id":       acctest.CHECKSET,
			"records.0.priority":        acctest.CHECKSET,
			"records.0.remark":          acctest.CHECKSET,
			"records.0.rr":              acctest.CHECKSET,
			"records.0.resource_record": acctest.CHECKSET,
			"records.0.ttl":             acctest.CHECKSET,
			"records.0.type":            acctest.CHECKSET,
			"records.0.value":           acctest.CHECKSET,
			"records.0.status":          acctest.CHECKSET,
		}
	}

	var fakeAliCloudPvtzZoneRecordsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"records.#": "0",
		}
	}

	var aliCloudPvtzZoneRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_pvtz_zone_records.default",
		ExistMapFunc: existAliCloudPvtzZoneRecordsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudPvtzZoneRecordsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudPvtzZoneRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf, statusConf, allConf)
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

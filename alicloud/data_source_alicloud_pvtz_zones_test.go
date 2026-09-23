package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudPvtzZonesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.TestPvtzRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_pvtz_zone.default.zone_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_pvtz_zone.default.zone_name}_fake"`,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.default.zone_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.default.zone_name}_fake"`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_pvtz_zone.default.id}"]`,
			"name_regex":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"keyword":           `"${alicloud_pvtz_zone.default.zone_name}"`,
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"search_mode":       `"EXACT"`,
			"lang":              `"en"`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_pvtz_zone.default.id}_fake"]`,
			"name_regex":        `"${alicloud_pvtz_zone.default.zone_name}_fake"`,
			"keyword":           `"${alicloud_pvtz_zone.default.zone_name}_fake"`,
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}_fake"`,
			"search_mode":       `"LIKE"`,
			"lang":              `"zh"`,
		}),
	}

	var existAliCloudPvtzZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"zones.#":                   "1",
			"zones.0.id":                CHECKSET,
			"zones.0.zone_id":           CHECKSET,
			"zones.0.zone_name":         CHECKSET,
			"zones.0.name":              CHECKSET,
			"zones.0.proxy_pattern":     CHECKSET,
			"zones.0.record_count":      CHECKSET,
			"zones.0.resource_group_id": CHECKSET,
			"zones.0.remark":            CHECKSET,
			"zones.0.is_ptr":            CHECKSET,
			"zones.0.create_timestamp":  CHECKSET,
			"zones.0.update_timestamp":  CHECKSET,
			"zones.0.bind_vpcs.#":       "0",
		}
	}

	var fakeAliCloudPvtzZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"zones.#": "0",
		}
	}

	var aliCloudPvtzZonesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pvtz_zones.default",
		existMapFunc: existAliCloudPvtzZonesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudPvtzZonesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudPvtzZonesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keywordConf, resourceGroupIdConf, allConf)
}

func TestAccAliCloudPvtzZonesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.TestPvtzRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_pvtz_zone_attachment.default.zone_id}"]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_pvtz_zone_attachment.default.zone_id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `false`,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `false`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"enable_details":    `false`,
		}),
	}

	queryVpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_vpc_id":   `"${alicloud_vpc.default.id}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_vpc_id":   `"${alicloud_vpc.default.id}"`,
			"enable_details": `false`,
		}),
	}

	queryRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_region_id": `"${data.alicloud_regions.default.regions.0.id}"`,
			"enable_details":  `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_region_id": `"${data.alicloud_regions.default.regions.0.id}"`,
			"enable_details":  `false`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_pvtz_zone.default.id}"]`,
			"name_regex":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"keyword":           `"${alicloud_pvtz_zone.default.zone_name}"`,
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"query_vpc_id":      `"${alicloud_vpc.default.id}"`,
			"query_region_id":   `"${data.alicloud_regions.default.regions.0.id}"`,
			"search_mode":       `"EXACT"`,
			"lang":              `"en"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_pvtz_zone.default.id}"]`,
			"name_regex":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"keyword":           `"${alicloud_pvtz_zone.default.zone_name}"`,
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"query_vpc_id":      `"${alicloud_vpc.default.id}"`,
			"query_region_id":   `"${data.alicloud_regions.default.regions.0.id}"`,
			"search_mode":       `"EXACT"`,
			"lang":              `"en"`,
			"enable_details":    `false`,
		}),
	}

	var existAliCloudPvtzZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"zones.#":                         "1",
			"zones.0.id":                      CHECKSET,
			"zones.0.zone_id":                 CHECKSET,
			"zones.0.zone_name":               CHECKSET,
			"zones.0.name":                    CHECKSET,
			"zones.0.proxy_pattern":           CHECKSET,
			"zones.0.record_count":            CHECKSET,
			"zones.0.resource_group_id":       CHECKSET,
			"zones.0.remark":                  CHECKSET,
			"zones.0.is_ptr":                  CHECKSET,
			"zones.0.create_timestamp":        CHECKSET,
			"zones.0.update_timestamp":        CHECKSET,
			"zones.0.slave_dns":               CHECKSET,
			"zones.0.bind_vpcs.#":             "1",
			"zones.0.bind_vpcs.0.vpc_id":      CHECKSET,
			"zones.0.bind_vpcs.0.vpc_name":    CHECKSET,
			"zones.0.bind_vpcs.0.region_id":   CHECKSET,
			"zones.0.bind_vpcs.0.region_name": CHECKSET,
		}
	}

	var fakeAliCloudPvtzZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"zones.#":                   "1",
			"zones.0.id":                CHECKSET,
			"zones.0.zone_id":           CHECKSET,
			"zones.0.zone_name":         CHECKSET,
			"zones.0.name":              CHECKSET,
			"zones.0.proxy_pattern":     CHECKSET,
			"zones.0.record_count":      CHECKSET,
			"zones.0.resource_group_id": CHECKSET,
			"zones.0.remark":            CHECKSET,
			"zones.0.is_ptr":            CHECKSET,
			"zones.0.create_timestamp":  CHECKSET,
			"zones.0.update_timestamp":  CHECKSET,
			"zones.0.bind_vpcs.#":       "0",
		}
	}

	var aliCloudPvtzZonesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pvtz_zones.default",
		existMapFunc: existAliCloudPvtzZonesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudPvtzZonesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudPvtzZonesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keywordConf, resourceGroupIdConf, queryVpcIdConf, queryRegionIdConf, allConf)
}

func testAccCheckAliCloudPvtzZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc%d.test.com"
	}

	data "alicloud_regions" "default" {
  		current = true
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_pvtz_zone" "default" {
  		zone_name         = var.name
  		remark            = var.name
  		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
	}

	resource "alicloud_pvtz_zone_attachment" "default" {
  		zone_id = alicloud_pvtz_zone.default.id
  		vpcs {
    		vpc_id = alicloud_vpc.default.id
  		}
	}

	data "alicloud_pvtz_zones" "default"{
		%s
	}
	`, rand, strings.Join(pairs, " \n "))
}

package privatezone_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudPvtzZonesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_pvtz_zone.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_pvtz_zone.default.zone_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_pvtz_zone.default.zone_name}_fake"`,
		}),
	}

	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.default.zone_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_pvtz_zone.default.zone_name}_fake"`,
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_pvtz_zone.default.id}"]`,
			"name_regex":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"keyword":           `"${alicloud_pvtz_zone.default.zone_name}"`,
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"search_mode":       `"EXACT"`,
			"lang":              `"en"`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
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
			"zones.0.id":                acctest.CHECKSET,
			"zones.0.zone_id":           acctest.CHECKSET,
			"zones.0.zone_name":         acctest.CHECKSET,
			"zones.0.name":              acctest.CHECKSET,
			"zones.0.proxy_pattern":     acctest.CHECKSET,
			"zones.0.record_count":      acctest.CHECKSET,
			"zones.0.resource_group_id": acctest.CHECKSET,
			"zones.0.remark":            acctest.CHECKSET,
			"zones.0.is_ptr":            acctest.CHECKSET,
			"zones.0.create_timestamp":  acctest.CHECKSET,
			"zones.0.update_timestamp":  acctest.CHECKSET,
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

	var aliCloudPvtzZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_pvtz_zones.default",
		ExistMapFunc: existAliCloudPvtzZonesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudPvtzZonesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudPvtzZonesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keywordConf, resourceGroupIdConf, allConf)
}

func TestAccAliCloudPvtzZonesDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_pvtz_zone_attachment.default.zone_id}"]`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_pvtz_zone_attachment.default.zone_id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `false`,
		}),
	}

	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"keyword":        `"${alicloud_pvtz_zone.default.zone_name}"`,
			"enable_details": `false`,
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"enable_details":    `true`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_pvtz_zone.default.resource_group_id}"`,
			"enable_details":    `false`,
		}),
	}

	queryVpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_vpc_id":   `"${alicloud_vpc.default.id}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_vpc_id":   `"${alicloud_vpc.default.id}"`,
			"enable_details": `false`,
		}),
	}

	queryRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_region_id": `"${data.alicloud_regions.default.regions.0.id}"`,
			"enable_details":  `true`,
		}),
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
			"query_region_id": `"${data.alicloud_regions.default.regions.0.id}"`,
			"enable_details":  `false`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
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
		FakeConfig: testAccCheckAliCloudPvtzZonesDataSourceName(rand, map[string]string{
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
			"zones.0.id":                      acctest.CHECKSET,
			"zones.0.zone_id":                 acctest.CHECKSET,
			"zones.0.zone_name":               acctest.CHECKSET,
			"zones.0.name":                    acctest.CHECKSET,
			"zones.0.proxy_pattern":           acctest.CHECKSET,
			"zones.0.record_count":            acctest.CHECKSET,
			"zones.0.resource_group_id":       acctest.CHECKSET,
			"zones.0.remark":                  acctest.CHECKSET,
			"zones.0.is_ptr":                  acctest.CHECKSET,
			"zones.0.create_timestamp":        acctest.CHECKSET,
			"zones.0.update_timestamp":        acctest.CHECKSET,
			"zones.0.slave_dns":               acctest.CHECKSET,
			"zones.0.bind_vpcs.#":             "1",
			"zones.0.bind_vpcs.0.vpc_id":      acctest.CHECKSET,
			"zones.0.bind_vpcs.0.vpc_name":    acctest.CHECKSET,
			"zones.0.bind_vpcs.0.region_id":   acctest.CHECKSET,
			"zones.0.bind_vpcs.0.region_name": acctest.CHECKSET,
		}
	}

	var fakeAliCloudPvtzZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"zones.#":                   "1",
			"zones.0.id":                acctest.CHECKSET,
			"zones.0.zone_id":           acctest.CHECKSET,
			"zones.0.zone_name":         acctest.CHECKSET,
			"zones.0.name":              acctest.CHECKSET,
			"zones.0.proxy_pattern":     acctest.CHECKSET,
			"zones.0.record_count":      acctest.CHECKSET,
			"zones.0.resource_group_id": acctest.CHECKSET,
			"zones.0.remark":            acctest.CHECKSET,
			"zones.0.is_ptr":            acctest.CHECKSET,
			"zones.0.create_timestamp":  acctest.CHECKSET,
			"zones.0.update_timestamp":  acctest.CHECKSET,
			"zones.0.bind_vpcs.#":       "0",
		}
	}

	var aliCloudPvtzZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_pvtz_zones.default",
		ExistMapFunc: existAliCloudPvtzZonesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudPvtzZonesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudPvtzZonesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keywordConf, resourceGroupIdConf, queryVpcIdConf, queryRegionIdConf, allConf)
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

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenRegionRouteEntriesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckCenRegionRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_route_entry.default.instance_id}"`,
			"region_id":   fmt.Sprintf(`"%s"`, defaultRegionToTest),
		}),
		fakeConfig: testAccCheckCenRegionRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_route_entry.default.instance_id}_fake"`,
			"region_id":   fmt.Sprintf(`"%s"`, defaultRegionToTest),
		}),
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	CenRegionRouteEntriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckCenRegionRouteEntriesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAcc%sRegionRouteEntriesDataSourceBasic-%d"
	}
	
	resource "alicloud_instance" "default" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
	
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
	
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}-region-route-entry"
	}
	
	resource "alicloud_cen_instance" "default" {
		name = "${var.name}-cen"
		description = "terraform01"
	}
	
	resource "alicloud_cen_instance_attachment" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_type = "VPC"
	    child_instance_region_id = "%s"
	    depends_on = [
	        "alicloud_vswitch.default"]
	}
	
	resource "alicloud_route_entry" "default" {
	    route_table_id = "${alicloud_vpc.default.route_table_id}"
	    destination_cidrblock = "11.0.0.0/16"
	    nexthop_type = "Instance"
	    nexthop_id = "${alicloud_instance.default.id}"
	}
	
	resource "alicloud_cen_route_entry" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    route_table_id = "${alicloud_vpc.default.route_table_id}"
	    cidr_block = "${alicloud_route_entry.default.destination_cidrblock}"
	    depends_on = [
			"alicloud_cen_instance_attachment.default"]
	}
	
	data "alicloud_cen_region_route_entries" "default" {
	%s
	}
	`, EcsInstanceCommonTestCase, defaultRegionToTest, rand, defaultRegionToTest, strings.Join(pairs, "\n  "))
	return config
}

var existCenRegionRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instance_id":                  CHECKSET,
		"region_id":                    CHECKSET,
		"entries.#":                    "3",
		"entries.0.cidr_block":         "100.64.0.0/10",
		"entries.0.type":               "System",
		"entries.0.next_hop_type":      "local_service",
		"entries.0.next_hop_id":        "",
		"entries.0.next_hop_region_id": CHECKSET,
	}
}

var fakeCenRegionRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var CenRegionRouteEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_region_route_entries.default",
	existMapFunc: existCenRegionRouteEntriesMapFunc,
	fakeMapFunc:  fakeCenRegionRouteEntriesMapFunc,
}

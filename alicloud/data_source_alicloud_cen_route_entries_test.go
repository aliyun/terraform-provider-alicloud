package alicloud

import (
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"fmt"
)

func TestAccAlicloudCenRouteEntriesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudCenRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_cen_route_entry.default.instance_id}"`,
			"route_table_id": `"${alicloud_cen_route_entry.default.route_table_id}"`,
			"cidr_block":     `"11.0.0.0/16"`,
		}),
		fakeConfig: testAccAlicloudCenRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_cen_route_entry.default.instance_id}"`,
			"route_table_id": `"${alicloud_cen_route_entry.default.route_table_id}"`,
			"cidr_block":     `"11.2.0.0/16"`,
		}),
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	CenRouteEntriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccAlicloudCenRouteEntriesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s
	variable "name" {
	    default = "tf-testAcc%sCenRouteEntries-%d"
	}

	resource "alicloud_instance" "default" {
	    vswitch_id = "${alicloud_vswitch.default.id}"
	    image_id = "${data.alicloud_images.default.images.0.id}"
	    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	    system_disk_category = "cloud_efficiency"
	    internet_charge_type = "PayByTraffic"
	    internet_max_bandwidth_out = 5
	    security_groups = ["${alicloud_security_group.default.id}"]
	    instance_name = "${var.name}"
	}

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
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

	data "alicloud_cen_route_entries" "default" {
		%s
	}
	`, EcsInstanceCommonTestCase, defaultRegionToTest, rand, defaultRegionToTest, strings.Join(pairs, "\n  "))
	return config
}

var existCenRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instance_id":                CHECKSET,
		"route_table_id":             CHECKSET,
		"cidr_block":                 "11.0.0.0/16",
		"entries.#":                  "1",
		"entries.0.cidr_block":       "11.0.0.0/16",
		"entries.0.next_hop_type":    "Instance",
		"entries.0.route_type":       "Custom",
		"entries.0.route_table_id":   CHECKSET,
		"entries.0.next_hop_id":      CHECKSET,
		"entries.0.operational_mode": "true",
		"entries.0.conflicts":        NOSET,
	}
}

var fakeCenRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var CenRouteEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_route_entries.default",
	existMapFunc: existCenRouteEntriesMapFunc,
	fakeMapFunc:  fakeCenRouteEntriesMapFunc,
}

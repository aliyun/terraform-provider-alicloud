package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCRouteEntriesDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
		}),
		fakeConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
			"type":           `"System"`,
		}),
	}

	cidrBlockConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
			"cidr_block":     `"${alicloud_route_entry.default.destination_cidrblock}"`,
		}),
		fakeConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
			"cidr_block":     `"${alicloud_route_entry.default.destination_cidrblock}_fake"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_instance.default.id}"`,
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
			"cidr_block":     `"${alicloud_route_entry.default.destination_cidrblock}"`,
		}),
		fakeConfig: testAccCheckAlicloudRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alicloud_instance.default.id}"`,
			"route_table_id": `"${alicloud_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
			"cidr_block":     `"${alicloud_route_entry.default.destination_cidrblock}_fake"`,
		}),
	}

	routeEntriesCheckInfo.dataSourceTestCheck(t, rand, instanceIdConf, typeConfig, cidrBlockConfig, allConfig)
}

func testAccCheckAlicloudRouteEntriesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
	name_regex = "^ubuntu"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "tf-testAcc-for-route-entries-datasource%d"
}
resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	description = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.default.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.default.id}"]

	vswitch_id = "${alicloud_vswitch.default.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

resource "alicloud_route_entry" "default" {
	route_table_id = "${alicloud_vpc.default.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.default.id}"
}

data "alicloud_route_entries" "default" {
  %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":                "1",
		"entries.0.route_table_id": CHECKSET,
		"entries.0.cidr_block":     CHECKSET,
		"entries.0.instance_id":    CHECKSET,
		"entries.0.status":         CHECKSET,
		"entries.0.type":           "Custom",
		"entries.0.next_hop_type":  "Instance",
	}
}

var fakeRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var routeEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_route_entries.default",
	existMapFunc: existRouteEntriesMapFunc,
	fakeMapFunc:  fakeRouteEntriesMapFunc,
}

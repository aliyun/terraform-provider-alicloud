package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenChildInstanceRouteEntryToAttachmentDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.CenChildInstanceRouteEntryToAttachmentSupportRegions)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenChildInstanceRouteEntryToAttachmentSourceConfig(rand, map[string]string{
			"ids":                           `["${alicloud_cen_child_instance_route_entry_to_attachment.default.id}"]`,
			"cen_id":                        `"${alicloud_cen_child_instance_route_entry_to_attachment.default.cen_id}"`,
			"child_instance_route_table_id": `"${alicloud_cen_child_instance_route_entry_to_attachment.default.child_instance_route_table_id}"`,
			"transit_router_attachment_id":  `"${alicloud_cen_child_instance_route_entry_to_attachment.default.transit_router_attachment_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenChildInstanceRouteEntryToAttachmentSourceConfig(rand, map[string]string{
			"ids":                           `["${alicloud_cen_child_instance_route_entry_to_attachment.default.id}_fake"]`,
			"child_instance_route_table_id": `"${alicloud_cen_child_instance_route_entry_to_attachment.default.child_instance_route_table_id}"`,
			"transit_router_attachment_id":  `"${alicloud_cen_child_instance_route_entry_to_attachment.default.transit_router_attachment_id}"`,
		}),
	}

	CenChildInstanceRouteEntryToAttachmentCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existCenChildInstanceRouteEntryToAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#":                               "1",
		"attachments.0.id":                            acctest.CHECKSET,
		"attachments.0.cen_id":                        acctest.CHECKSET,
		"attachments.0.child_instance_route_table_id": acctest.CHECKSET,
		"attachments.0.destination_cidr_block":        acctest.CHECKSET,
		"attachments.0.status":                        acctest.CHECKSET,
		"attachments.0.transit_router_attachment_id":  acctest.CHECKSET,
	}
}

var fakeCenChildInstanceRouteEntryToAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#": "0",
	}
}

var CenChildInstanceRouteEntryToAttachmentCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cen_child_instance_route_entry_to_attachments.default",
	ExistMapFunc: existCenChildInstanceRouteEntryToAttachmentMapFunc,
	FakeMapFunc:  fakeCenChildInstanceRouteEntryToAttachmentMapFunc,
}

func testAccCheckAlicloudCenChildInstanceRouteEntryToAttachmentSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCenChildInstanceRouteEntryToAttachment%d"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[2]
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mappings {
    zone_id    = alicloud_vswitch.default_master.zone_id
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.default_slave.zone_id
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.default.id
  route_table_name = var.name
  description      = var.name
}

resource "alicloud_cen_child_instance_route_entry_to_attachment" "default" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  cen_id                        = alicloud_cen_instance.default.id
  destination_cidr_block        = "10.0.0.0/24"
  child_instance_route_table_id = alicloud_route_table.foo.id
  depends_on = ["alicloud_route_table.foo"]
}

data "alicloud_cen_child_instance_route_entry_to_attachments" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

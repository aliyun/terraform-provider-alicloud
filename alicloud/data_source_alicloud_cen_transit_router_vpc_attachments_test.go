package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVpcAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.id}_fake"]`,
		}),
	}
	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.default.id}_fake"`,
		}),
	}
	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_vpc_attachment.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_vpc_attachment.default.status}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"ids":               `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"status":            `"${alicloud_cen_transit_router_vpc_attachment.default.status}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}_fake"`,
			"ids":               `["${alicloud_cen_transit_router_vpc_attachment.default.id}_fake"]`,
			"status":            `"${alicloud_cen_transit_router_vpc_attachment.default.status}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"attachments.#": "1",
			"attachments.0.transit_router_attachment_description": `descp`,
			"attachments.0.transit_router_attachment_name":        `name`,
			"attachments.0.vpc_id":                                CHECKSET,
			"attachments.0.vpc_owner_id":                          CHECKSET,
			"attachments.0.zone_mappings":                         `[{"ZoneId":"cn-hangzhou-h","VSwitchId":""},{"ZoneId":"cn-hangzhou-i","VSwitchId":"vsw-2"}]`,
		}
	}
	var fakeAlicloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterVpcAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_vpc_attachments.default",
		existMapFunc: existAlicloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterVpcAttachmentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, cenIdConf, transitRouterIdConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDataTransitRouterVpcAttachment-%d"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.name
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "192.168.1.0/24"
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.name
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "192.168.2.0/24"
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= "${alicloud_cen_instance.default.id}"
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id = "${alicloud_cen_instance.default.id}"
  transit_router_id = "${alicloud_cen_transit_router.default.transit_router_id}"
  vpc_id = "${alicloud_vpc.default.id}"
  zone_mappings {
    zone_id = "cn-hangzhou-h"
    vswitch_id = "${alicloud_vswitch.default_master.id}"
  }
  zone_mappings {
    zone_id = "cn-hangzhou-i"
    vswitch_id = "${alicloud_vswitch.default_slave.id}"
  }
  transit_router_attachment_description = "descp"
  transit_router_attachment_name = "name"
}

data "alicloud_cen_transit_router_vpc_attachments" "default" {	
	enable_details = true
	cen_id = "${alicloud_cen_instance.default.id}"
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

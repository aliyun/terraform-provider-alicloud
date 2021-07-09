package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVpcAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"]`,
			"status": `"Attached"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"]`,
			"status": `"Attaching"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"]`,
			"status": `"Attached"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}_fake"]`,
			"status": `"Attaching"`,
		}),
	}
	var existAlicloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"attachments.#": "1",
			"attachments.0.transit_router_attachment_description": `descp`,
			"attachments.0.transit_router_attachment_name":        fmt.Sprintf("tf-testAccDataTransitRouterVpcAttachment-%d", rand),
			"attachments.0.vpc_id":                                CHECKSET,
			"attachments.0.resource_type":                         "VPC",
			"attachments.0.vpc_owner_id":                          CHECKSET,
			"attachments.0.zone_mappings.0.vswitch_id":            CHECKSET,
			"attachments.0.zone_mappings.0.zone_id":               `cn-hangzhou-h`,
			"attachments.0.zone_mappings.1.vswitch_id":            CHECKSET,
			"attachments.0.zone_mappings.1.zone_id":               `cn-hangzhou-i`,
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
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	}
	alicloudCenTransitRouterVpcAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
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
  vpc_id = alicloud_vpc.default.id
  cidr_block = "192.168.1.0/24"
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.name
  vpc_id = alicloud_vpc.default.id
  cidr_block = "192.168.2.0/24"
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id = alicloud_vpc.default.id
  zone_mappings {
    zone_id = "cn-hangzhou-h"
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id = "cn-hangzhou-i"
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_description = "descp"
  transit_router_attachment_name = var.name
}

data "alicloud_cen_transit_router_vpc_attachments" "default" {	
	cen_id = alicloud_cen_instance.default.id
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

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
			"attachments.0.payment_type":                          "PayAsYouGo",
			"attachments.0.vpc_owner_id":                          CHECKSET,
			"attachments.0.zone_mappings.0.vswitch_id":            CHECKSET,
			"attachments.0.zone_mappings.0.zone_id":               CHECKSET,
			"attachments.0.zone_mappings.1.vswitch_id":            CHECKSET,
			"attachments.0.zone_mappings.1.zone_id":               CHECKSET,
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

data "alicloud_zones" "default" {
    available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.ids.0
}
data "alicloud_vswitches" "default_master" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.ids.1
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
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_mappings {
    zone_id = data.alicloud_vswitches.default_master.vswitches.0.zone_id
    vswitch_id = data.alicloud_vswitches.default_master.vswitches.0.id
  }
  zone_mappings {
    zone_id = data.alicloud_vswitches.default.vswitches.0.zone_id
    vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
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

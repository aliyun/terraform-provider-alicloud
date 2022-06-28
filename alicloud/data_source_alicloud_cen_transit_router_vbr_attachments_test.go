package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVbrAttachmentsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1, 2999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}"]`,
			"status": `"Attached"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}"]`,
			"status": `"Attaching"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}"]`,
			"status":            `"Attached"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}_fake"]`,
			"status":            `"Attaching"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
		}),
	}
	var existAlicloudCenTransitRouterVbrAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"attachments.#": "1",
			"attachments.0.auto_publish_route_enabled":            `true`,
			"attachments.0.transit_router_attachment_id":          CHECKSET,
			"attachments.0.transit_router_attachment_description": `desp`,
			"attachments.0.transit_router_attachment_name":        fmt.Sprintf("tf-testAccDataTransitRouterVbrAttachment-%d", rand),
			"attachments.0.vbr_id":                                CHECKSET,
			"attachments.0.resource_type":                         `VBR`,
		}
	}
	var fakeAlicloudCenTransitRouterVbrAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterVbrAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_vbr_attachments.default",
		existMapFunc: existAlicloudCenTransitRouterVbrAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterVbrAttachmentsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterVbrAttachmentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDataTransitRouterVbrAttachment-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
resource "alicloud_cen_transit_router_vbr_attachment" "default" {
	auto_publish_route_enabled = true
	cen_id = alicloud_cen_instance.default.id
	transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	vbr_id = alicloud_express_connect_virtual_border_router.default.id
	transit_router_attachment_description = "desp"
	transit_router_attachment_name = var.name
}

data "alicloud_cen_transit_router_vbr_attachments" "default" {	
	cen_id = alicloud_cen_transit_router_vbr_attachment.default.cen_id
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}

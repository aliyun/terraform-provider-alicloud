package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterPeerAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"]`,
			//"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"]`,
			//"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"`,
		}),
	}
	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"]`,
			//"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			//"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			//"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"`,
			"ids":               `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			"status":            `"Attached"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			//"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}_fake"`,
			"ids":               `["${alicloud_cen_transit_router_peer_attachment.default.id}_fake"]`,
			"status":            `"Attaching"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"attachments.#": "1",
			"attachments.0.auto_publish_route_enabled":            `true`,
			"attachments.0.bandwidth":                             `2`,
			"attachments.0.cen_bandwidth_package_id":              `cenbwp-buw65zk0606xh0ukvd`,
			"attachments.0.peer_transit_router_id":                CHECKSET,
			"attachments.0.peer_transit_router_region_id":         `us-east-1`,
			"attachments.0.transit_router_attachment_description": CHECKSET,
			"attachments.0.transit_router_attachment_name":        CHECKSET,
			"attachments.0.transit_router_id":                     CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterPeerAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_peer_attachments.default",
		existMapFunc: existAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterPeerAttachmentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, transitRouterIdConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTransitRouterPeerAttachment-%d"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  instance_id        = "cen-f6rslz7pzbnj8sshxc"
  bandwidth_package_id = "cenbwp-buw65zk0606xh0ukvd"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = "cen-f6rslz7pzbnj8sshxc"
  depends_on = [alicloud_cen_bandwidth_package_attachment.default]
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  cen_id = "cen-f6rslz7pzbnj8sshxc"
  transit_router_id = "tr-bp1p0oqyc5iv22yjpymgu"
  peer_transit_router_region_id = "us-east-1"
  peer_transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  cen_bandwidth_package_id = "cenbwp-buw65zk0606xh0ukvd"
  bandwidth = 2
  auto_publish_route_enabled = true
  transit_router_attachment_name = var.name
  transit_router_attachment_description = "desp"
}

data "alicloud_cen_transit_router_peer_attachments" "default" {	
  cen_id = "cen-f6rslz7pzbnj8sshxc"
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

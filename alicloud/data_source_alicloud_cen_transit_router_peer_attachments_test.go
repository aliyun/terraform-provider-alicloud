package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterPeerAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			//"ids":    `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.id}"`,
			"cen_id":                       `"${alicloud_cen_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			//"ids":    `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.id}_fake"`,
			"cen_id":                       `"${alicloud_cen_instance.default.id}_fake"`,
		}),
	}
	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			//"ids":               `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			"cen_id":                       `"${alicloud_cen_instance.default.id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.id}"`,
			"transit_router_id":            `"${alicloud_cen_transit_router.default_0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			//"ids":               `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			"cen_id":                       `"${alicloud_cen_instance.default.id}_fake"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.id}_fake"`,
			"transit_router_id":            `"${alicloud_cen_transit_router.default_0.id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"cen_id":                       `"${alicloud_cen_instance.default.id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.id}"`,
			//"ids":               `["${alicloud_cen_transit_router_peer_attachment.default.id}"]`,
			"status":            `"${alicloud_cen_transit_router_peer_attachment.default.status}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default_0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"cen_id":                       `"${alicloud_cen_instance.default.id}_fake"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_peer_attachment.default.id}_fake"`,
			//"ids":               `["${alicloud_cen_transit_router_peer_attachment.default.id}_fake"]`,
			"status":            `"${alicloud_cen_transit_router_peer_attachment.default.status}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default_0.id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			//"ids.#":         "1",
			"transit_router_attachments.#":                                       "1",
			"transit_router_attachments.0.auto_publish_route_enabled":            `true`,
			"transit_router_attachments.0.bandwidth":                             `2`,
			"transit_router_attachments.0.cen_bandwidth_package_id":              `cenbwp-buw65zk0606xh0ukvd`,
			"transit_router_attachments.0.peer_transit_router_id":                CHECKSET,
			"transit_router_attachments.0.peer_transit_router_region_id":         `us-east-1`,
			"transit_router_attachments.0.transit_router_attachment_description": CHECKSET,
			"transit_router_attachments.0.transit_router_attachment_name":        CHECKSET,
			"transit_router_attachments.0.transit_router_id":                     CHECKSET,
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
	alicloudCenTransitRouterPeerAttachmentsCheckInfo.dataSourceTestCheck(t, rand, cenIdConf, transitRouterIdConf, allConf)
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

provider "alicloud" {
  alias = "other_region_id"
  region = "us-east-1"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  protection_level = "REDUCED"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  instance_id        = "${alicloud_cen_instance.default.id}"
  bandwidth_package_id = "cenbwp-buw65zk0606xh0ukvd"
  depends_on = [
    alicloud_cen_instance.default]
}

resource "alicloud_cen_transit_router" "default_0" {
  cen_id = "${alicloud_cen_instance.default.id}"
  region_id = "cn-hangzhou"
  depends_on = [
    alicloud_cen_bandwidth_package_attachment.default]
}

resource "alicloud_cen_transit_router" "default_1" {
  cen_id = "${alicloud_cen_instance.default.id}"
  region_id = "us-east-1"
  depends_on = [
    alicloud_cen_transit_router.default_0]
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  cen_id = "${alicloud_cen_instance.default.id}"
  transit_router_id = "${alicloud_cen_transit_router.default_0.id}"
  peer_transit_router_region_id = "us-east-1"
  peer_transit_router_id = "${alicloud_cen_transit_router.default_1.id}"
  cen_bandwidth_package_id = "cenbwp-buw65zk0606xh0ukvd"
  bandwidth = 2
  auto_publish_route_enabled = true
  transit_router_attachment_name = "${var.name}"
  transit_router_attachment_description = "${var.name}"
  depends_on = [
    alicloud_cen_transit_router.default_0,
    alicloud_cen_transit_router.default_1]
}

data "alicloud_cen_transit_router_peer_attachments" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

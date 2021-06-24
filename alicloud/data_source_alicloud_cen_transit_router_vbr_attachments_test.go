package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVbrAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.id}_fake"]`,
			"cen_id": `"${alicloud_cen_instance.default.id}_fake"`,
		}),
	}
	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.default.id}_fake"`,
		}),
	}
	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_vbr_attachment.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_vbr_attachment.default.status}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"ids":               `["${alicloud_cen_transit_router_vbr_attachment.default.id}"]`,
			"status":            `"${alicloud_cen_transit_router_vbr_attachment.default.status}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}_fake"`,
			"ids":               `["${alicloud_cen_transit_router_vbr_attachment.default.id}_fake"]`,
			"status":            `"${alicloud_cen_transit_router_vbr_attachment.default.status}_fake"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterVbrAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"attachments.#": "1",
			"attachments.0.auto_publish_route_enabled":            `true`,
			"attachments.0.transit_router_attachment_description": `desp`,
			"attachments.0.transit_router_attachment_name":        `name`,
			"attachments.0.vbr_id":                                `vbr-j6cd9pm9y6d6e20atoi6w`,
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
	alicloudCenTransitRouterVbrAttachmentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, cenIdConf, transitRouterIdConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterVbrAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTransitRouterVbrAttachment-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= "${alicloud_cen_instance.default.id}"
region_id = "cn-hongkong"
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
auto_publish_route_enabled = true
cen_id = "${alicloud_cen_instance.default.id}"
transit_router_id = "${alicloud_cen_transit_router.default.id}"
vbr_id = "vbr-j6cd9pm9y6d6e20atoi6w"
transit_router_attachment_description = "desp"
transit_router_attachment_name = "name"
}

data "alicloud_cen_transit_router_vbr_attachments" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

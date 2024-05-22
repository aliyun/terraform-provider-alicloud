package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCenTransitRouterVpcAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"vpc_id": `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"vpc_id": `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}_fake"`,
		}),
	}
	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
			"status":            `"Attaching"`,
		}),
	}
	transitRouterAttachmentIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"status": `"Attached"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"status": `"Attaching"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"name_regex":                   `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}"`,
			"vpc_id":                       `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}"`,
			"transit_router_id":            `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"`,
			"status":                       `"Attached"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_vpc_attachment.default.id}_fake"]`,
			"name_regex":                   `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}_fake"`,
			"vpc_id":                       `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}_fake"`,
			"transit_router_id":            `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}_fake"`,
			"status":                       `"Attaching"`,
		}),
	}

	var existAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"names.#":              "1",
			"attachments.#":        "1",
			"attachments.0.id":     CHECKSET,
			"attachments.0.cen_id": CHECKSET,
			"attachments.0.transit_router_attachment_id":          CHECKSET,
			"attachments.0.vpc_id":                                CHECKSET,
			"attachments.0.transit_router_id":                     CHECKSET,
			"attachments.0.resource_type":                         "VPC",
			"attachments.0.payment_type":                          "PayAsYouGo",
			"attachments.0.vpc_owner_id":                          CHECKSET,
			"attachments.0.auto_publish_route_enabled":            CHECKSET,
			"attachments.0.transit_router_attachment_name":        CHECKSET,
			"attachments.0.transit_router_attachment_description": CHECKSET,
			"attachments.0.status":                                CHECKSET,
			"attachments.0.zone_mappings.#":                       "2",
			"attachments.0.zone_mappings.0.vswitch_id":            CHECKSET,
			"attachments.0.zone_mappings.0.zone_id":               CHECKSET,
			"attachments.0.zone_mappings.1.vswitch_id":            CHECKSET,
			"attachments.0.zone_mappings.1.zone_id":               CHECKSET,
		}
	}

	var fakeAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"attachments.#": "0",
		}
	}

	var alicloudCenTransitRouterVpcAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_vpc_attachments.default",
		existMapFunc: existAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudCenTransitRouterVpcAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, transitRouterIdConf, transitRouterAttachmentIdConf, statusConf, allConf)
}

func testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAcc-CenTransitRouterVpcAttachment-%d"
	}

	data "alicloud_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.ids.0
	}

	data "alicloud_vswitches" "default_master" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.ids.1
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
  		protection_level  = "REDUCED"
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id = alicloud_cen_instance.default.id
	}

	resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  		cen_id                                = alicloud_cen_instance.default.id
  		vpc_id                                = data.alicloud_vpcs.default.ids.0
  		transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  		transit_router_attachment_name        = var.name
  		transit_router_attachment_description = var.name
  		zone_mappings {
    		vswitch_id = data.alicloud_vswitches.default_master.vswitches.0.id
    		zone_id    = data.alicloud_vswitches.default_master.vswitches.0.zone_id
  		}
  		zone_mappings {
    		vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
    		zone_id    = data.alicloud_vswitches.default.vswitches.0.zone_id
  		}
	}

	data "alicloud_cen_transit_router_vpc_attachments" "default" {
  		cen_id = alicloud_cen_instance.default.id
 		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}

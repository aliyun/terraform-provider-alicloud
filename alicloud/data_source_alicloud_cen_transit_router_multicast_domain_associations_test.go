package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterMulticastDomainAssociationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_multicast_domain_association.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_multicast_domain_association.default.id}_fake"]`,
		}),
	}
	transitRouterAttachmentIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}_fake"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"vswitch_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"vswitch_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.vswitch_id}_fake"`,
		}),
	}
	resourceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"resource_id": `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"resource_id": `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}_fake"`,
		}),
	}
	resourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
			"resource_type":                `"VPC"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"status": `"Associated"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"status": `"Dissociating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_multicast_domain_association.default.id}"]`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
			"vswitch_id":                   `"${alicloud_cen_transit_router_multicast_domain_association.default.vswitch_id}"`,
			"resource_id":                  `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}"`,
			"resource_type":                `"VPC"`,
			"status":                       `"Associated"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_multicast_domain_association.default.id}_fake"]`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}_fake"`,
			"vswitch_id":                   `"${alicloud_cen_transit_router_multicast_domain_association.default.vswitch_id}_fake"`,
			"resource_id":                  `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}_fake"`,
			"status":                       `"Dissociating"`,
		}),
	}
	var existAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"associations.#":    "1",
			"associations.0.id": CHECKSET,
			"associations.0.transit_router_multicast_domain_id": CHECKSET,
			"associations.0.transit_router_attachment_id":       CHECKSET,
			"associations.0.vswitch_id":                         CHECKSET,
			"associations.0.resource_id":                        CHECKSET,
			"associations.0.resource_owner_id":                  CHECKSET,
			"associations.0.resource_type":                      "VPC",
			"associations.0.status":                             "Associated",
		}
	}
	var fakeAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}
	var alicloudCenTransitRouterMulticastDomainAssociationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_multicast_domain_associations.default",
		existMapFunc: existAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCenTransitRouterMulticastDomainAssociationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, transitRouterAttachmentIdConf, vswitchIdConf, resourceIdConf, resourceTypeConf, statusConf, allConf)
}

func testAccCheckAlicloudCenTransitRouterMulticastDomainAssociationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccCenTransitRouterMulticastDomainAssociation-%d"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		name_regex = "default-zone-i"
  		vpc_id     = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id            = alicloud_cen_instance.default.id
  		support_multicast = true
	}

	resource "alicloud_cen_transit_router_multicast_domain" "default" {
  		transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	}

	resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  		cen_id            = alicloud_cen_transit_router.default.cen_id
  		transit_router_id = alicloud_cen_transit_router_multicast_domain.default.transit_router_id
  		vpc_id            = data.alicloud_vpcs.default.ids.0
  		zone_mappings {
    		zone_id    = data.alicloud_vswitches.default.vswitches.0.zone_id
    		vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
  		}
	}

	resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
  		transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
  		transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  		vswitch_id                         = data.alicloud_vswitches.default.vswitches.0.id
	}

	data "alicloud_cen_transit_router_multicast_domain_associations" "default" {
  		transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}

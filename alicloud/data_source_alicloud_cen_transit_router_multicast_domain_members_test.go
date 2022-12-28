package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterMulticastDomainMemberDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.CENTransitRouterMulticastDomainMemberSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_member.default.id}"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_member.default.id}_fake"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"`,
		}),
	}

	TransitRouterMulticastDomainIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand, map[string]string{
			"network_interface_id":               `"${alicloud_cen_transit_router_multicast_domain_member.default.network_interface_id}"`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand, map[string]string{
			"network_interface_id":               `"${alicloud_cen_transit_router_multicast_domain_member.default.network_interface_id}_fake"`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_member.default.id}"]`,
			"network_interface_id":               `"${alicloud_cen_transit_router_multicast_domain_member.default.network_interface_id}"`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_member.default.id}_fake"]`,
			"network_interface_id":               `"${alicloud_cen_transit_router_multicast_domain_member.default.network_interface_id}"`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"`,
		}),
	}

	CenTransitRouterMulticastDomainMemberCheckInfo.dataSourceTestCheck(t, rand, idsConf, TransitRouterMulticastDomainIdConf, allConf)
}

var existCenTransitRouterMulticastDomainMemberMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"members.#":                                    "1",
		"members.0.id":                                 CHECKSET,
		"members.0.group_ip_address":                   CHECKSET,
		"members.0.network_interface_id":               CHECKSET,
		"members.0.status":                             CHECKSET,
		"members.0.transit_router_multicast_domain_id": CHECKSET,
		"members.0.vpc_id":                             CHECKSET,
	}
}

var fakeCenTransitRouterMulticastDomainMemberMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"members.#": "0",
	}
}

var CenTransitRouterMulticastDomainMemberCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_transit_router_multicast_domain_members.default",
	existMapFunc: existCenTransitRouterMulticastDomainMemberMapFunc,
	fakeMapFunc:  fakeCenTransitRouterMulticastDomainMemberMapFunc,
}

func testAccCheckAlicloudCenTransitRouterMulticastDomainMemberSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCenTransitRouterMulticastDomainMember%d"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = data.alicloud_vswitches.default.ids.0
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
data "alicloud_cen_instances" "default" {
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_routers" "default" {
  cen_id     = data.alicloud_cen_instances.default.instances.0.id
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_router_multicast_domains" "default" {
  transit_router_id = data.alicloud_cen_transit_routers.default.transit_routers.0.transit_router_id
  name_regex        = "no-deleting-cen"
}

resource "alicloud_cen_transit_router_multicast_domain_member" "default" {
  vpc_id                             = data.alicloud_vpcs.default.ids.0
  transit_router_multicast_domain_id = data.alicloud_cen_transit_router_multicast_domains.default.ids.0
  network_interface_id               = alicloud_ecs_network_interface.default.id
  group_ip_address                   = "224.0.0.8"
}

data "alicloud_cen_transit_router_multicast_domain_members" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

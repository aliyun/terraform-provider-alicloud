package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterMulticastDomainSourceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	domainIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainSourceSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_source.default.id}"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainSourceSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_source.default.id}_fake"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainSourceSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_source.default.id}"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainSourceSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_source.default.id}_fake"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}"`,
		}),
	}

	CenTransitRouterMulticastDomainSourceCheckInfo.dataSourceTestCheck(t, rand, domainIdConf, allConf)
}

var existCenTransitRouterMulticastDomainSourceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sources.#":    "1",
		"sources.0.id": CHECKSET,
	}
}

var fakeCenTransitRouterMulticastDomainSourceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sources.#": "0",
	}
}

var CenTransitRouterMulticastDomainSourceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_transit_router_multicast_domain_sources.default",
	existMapFunc: existCenTransitRouterMulticastDomainSourceMapFunc,
	fakeMapFunc:  fakeCenTransitRouterMulticastDomainSourceMapFunc,
}

func testAccCheckAlicloudCenTransitRouterMulticastDomainSourceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCenTransitRouterMulticastDomainSource%d"
}

variable "instance_name" {
  default = "tf-testacc-cen_instance"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.instance_name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.instance_name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.instance_name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = "cn-hangzhou-j"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.instance_name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mappings {
    zone_id    = alicloud_vswitch.default_master.zone_id
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.default_slave.zone_id
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.instance_name
  transit_router_attachment_description = var.instance_name
}


resource "alicloud_security_group" "default" {
    name = var.name
    vpc_id = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default"{
  status = "OK"
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  depends_on =           ["alicloud_cen_transit_router_vpc_attachment.default"]
  transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
}

resource "alicloud_ecs_network_interface" "default" {
    depends_on =           ["alicloud_cen_transit_router_multicast_domain.default"]
    network_interface_name = var.name
    vswitch_id = alicloud_vswitch.default_master.id
    security_group_ids = [alicloud_security_group.default.id]
  description = "Basic test"
  primary_ip_address = cidrhost(alicloud_vswitch.default_master.cidr_block, 100)
  tags = {
    Created = "TF",
    For =    "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
	transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
	transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
	vswitch_id                         = alicloud_vswitch.default_master.id
}

resource "alicloud_cen_transit_router_multicast_domain_source" "default" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id
  network_interface_id               = alicloud_ecs_network_interface.default.id
  group_ip_address                   = "230.1.1.1"
}

data "alicloud_cen_transit_router_multicast_domain_sources" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

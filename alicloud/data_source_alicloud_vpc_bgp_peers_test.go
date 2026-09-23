package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCBgpPeersDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1, 2999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_peer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_peer.default.id}_fake"]`,
		}),
	}
	bgpGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"bgp_group_id": `"${alicloud_vpc_bgp_peer.default.bgp_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"bgp_group_id": `"${alicloud_vpc_bgp_peer.default.bgp_group_id}_fake"`,
		}),
	}
	routerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_group.default.router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_group.default.router_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"status": `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"bgp_group_id": `"${alicloud_vpc_bgp_peer.default.bgp_group_id}"`,
			"ids":          `["${alicloud_vpc_bgp_peer.default.id}"]`,
			"router_id":    `"${alicloud_vpc_bgp_group.default.router_id}"`,
			"status":       `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpPeersDataSourceName(rand, map[string]string{
			"bgp_group_id": `"${alicloud_vpc_bgp_peer.default.bgp_group_id}_fake"`,
			"ids":          `["${alicloud_vpc_bgp_peer.default.id}_fake"]`,
			"router_id":    `"${alicloud_vpc_bgp_group.default.router_id}_fake"`,
			"status":       `"Deleting"`,
		}),
	}
	var existAlicloudVpcBgpPeersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"peers.#":                 "1",
			"peers.0.auth_key":        CHECKSET,
			"peers.0.bfd_multi_hop":   "10",
			"peers.0.bgp_group_id":    CHECKSET,
			"peers.0.id":              CHECKSET,
			"peers.0.bgp_peer_id":     CHECKSET,
			"peers.0.bgp_peer_name":   "",
			"peers.0.bgp_status":      "Connect",
			"peers.0.description":     "",
			"peers.0.enable_bfd":      "true",
			"peers.0.hold":            CHECKSET,
			"peers.0.ip_version":      "IPV4",
			"peers.0.is_fake":         CHECKSET,
			"peers.0.keepalive":       CHECKSET,
			"peers.0.local_asn":       CHECKSET,
			"peers.0.peer_asn":        CHECKSET,
			"peers.0.peer_ip_address": "1.1.1.1",
			"peers.0.route_limit":     CHECKSET,
			"peers.0.router_id":       CHECKSET,
			"peers.0.status":          "Available",
		}
	}
	var fakeAlicloudVpcBgpPeersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudVpcBgpPeersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_bgp_peers.default",
		existMapFunc: existAlicloudVpcBgpPeersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcBgpPeersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpcBgpPeersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, bgpGroupIdConf, routerIdConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcBgpPeersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccBgpPeer-%d"
}

data "alicloud_express_connect_physical_connections" "default" {
	name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_group" "default" {
  auth_key       = "YourPassword+12345678"
  bgp_group_name = var.name
  description    = var.name
  local_asn      = 64512
  peer_asn       = 1111
  router_id      = alicloud_express_connect_virtual_border_router.default.id
}

resource "alicloud_vpc_bgp_peer" "default" {
  bfd_multi_hop   = "10"
  bgp_group_id    = alicloud_vpc_bgp_group.default.id
  enable_bfd      = true
  ip_version      = "IPV4"
  peer_ip_address = "1.1.1.1"
}

data "alicloud_vpc_bgp_peers" "default" {	
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}

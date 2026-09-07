package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCBgpNetworksDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1, 2999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_network.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_network.default.id}_fake"]`,
		}),
	}
	routerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_network.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_network.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}"`,
			"status":    `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}_fake"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}_fake"`,
			"status":    `"Pending"`,
		}),
	}
	var existAlicloudVpcBgpNetworksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"networks.#":                "1",
			"networks.0.id":             CHECKSET,
			"networks.0.dst_cidr_block": "192.168.0.1",
			"networks.0.router_id":      CHECKSET,
			"networks.0.status":         "Available",
		}
	}
	var fakeAlicloudVpcBgpNetworksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudVpcBgpNetworksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_bgp_networks.default",
		existMapFunc: existAlicloudVpcBgpNetworksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcBgpNetworksDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpcBgpNetworksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, routerIdConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccBgpNetwork-%d"
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

resource "alicloud_vpc_bgp_network" "default" {
  dst_cidr_block = "192.168.0.1"
  router_id = alicloud_express_connect_virtual_border_router.default.id
}

data "alicloud_vpc_bgp_networks" "default" {	
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}

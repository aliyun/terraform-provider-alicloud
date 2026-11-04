package expressconnect_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudVPCBgpNetworksDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := tfacctest.RandIntRange(1, 2999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_network.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_network.default.id}_fake"]`,
		}),
	}
	routerIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_network.default.id}"]`,
			"status": `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_network.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}"`,
			"status":    `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcBgpNetworksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_network.default.id}_fake"]`,
			"router_id": `"${alicloud_vpc_bgp_network.default.router_id}_fake"`,
			"status":    `"Pending"`,
		}),
	}
	var existAlicloudVpcBgpNetworksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"networks.#":                "1",
			"networks.0.id":             acctest.CHECKSET,
			"networks.0.dst_cidr_block": "192.168.0.1",
			"networks.0.router_id":      acctest.CHECKSET,
			"networks.0.status":         "Available",
		}
	}
	var fakeAlicloudVpcBgpNetworksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudVpcBgpNetworksCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vpc_bgp_networks.default",
		ExistMapFunc: existAlicloudVpcBgpNetworksDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudVpcBgpNetworksDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudVpcBgpNetworksCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, routerIdConf, statusConf, allConf)
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

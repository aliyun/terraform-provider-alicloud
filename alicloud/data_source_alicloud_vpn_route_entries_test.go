package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPNPbrRouteEntriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnPbrRouteEntriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_pbr_route_entry.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpnPbrRouteEntriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_pbr_route_entry.default.id}_fake"]`,
		}),
	}
	var existAlicloudVpnPbrRouteEntriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"entries.#":                "1",
			"entries.0.create_time":    CHECKSET,
			"entries.0.vpn_gateway_id": CHECKSET,
			"entries.0.next_hop":       CHECKSET,
			"entries.0.id":             CHECKSET,
			"entries.0.route_dest":     "10.0.0.0/24",
			"entries.0.route_source":   "192.168.1.0/24",
			"entries.0.weight":         CHECKSET,
			"entries.0.status":         CHECKSET,
		}
	}
	var fakeAlicloudVpnPbrRouteEntriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudVpnPbrRouteEntriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpn_pbr_route_entries.default",
		existMapFunc: existAlicloudVpnPbrRouteEntriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpnPbrRouteEntriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpnPbrRouteEntriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func testAccCheckAlicloudVpnPbrRouteEntriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccIpsecServer-%d"
}

data "alicloud_vpn_gateways" "default" {
}

resource "alicloud_vpn_customer_gateway" "default" {
  name       = var.name
  ip_address = "192.168.1.1"
}

resource "alicloud_vpn_connection" "default" {
  name                = var.name
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  vpn_gateway_id      = data.alicloud_vpn_gateways.default.ids.0
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource alicloud_vpn_pbr_route_entry default {
  vpn_gateway_id = data.alicloud_vpn_gateways.default.ids.0
  route_source   = "192.168.1.0/24"
  route_dest     = "10.0.0.0/24"
  next_hop       = alicloud_vpn_connection.default.id
  weight         = 0
  publish_vpc    = false
}
data "alicloud_vpn_pbr_route_entries" "default" {
    vpn_gateway_id = data.alicloud_vpn_gateways.default.ids.0
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

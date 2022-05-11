package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCBgpGroupsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1, 2999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_bgp_group.default.id}_fake"]`,
		}),
	}
	routerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_group.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_group.default.router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_vpc_bgp_group.default.id}"]`,
			"router_id": `"${alicloud_vpc_bgp_group.default.router_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_bgp_group.default.bgp_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_bgp_group.default.bgp_group_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_group.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_bgp_group.default.id}"]`,
			"status": `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_bgp_group.default.id}"]`,
			"name_regex": `"${alicloud_vpc_bgp_group.default.bgp_group_name}"`,
			"router_id":  `"${alicloud_vpc_bgp_group.default.router_id}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_bgp_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_vpc_bgp_group.default.bgp_group_name}_fake"`,
			"router_id":  `"${alicloud_vpc_bgp_group.default.router_id}_fake"`,
			"status":     `"Deleting"`,
		}),
	}
	var existAlicloudVpcBgpGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"groups.#":                "1",
			"groups.0.auth_key":       "YourPassword+12345678",
			"groups.0.bgp_group_name": fmt.Sprintf("tf-testAccBgpGroup-%d", rand),
			"groups.0.description":    fmt.Sprintf("tf-testAccBgpGroup-%d", rand),
			"groups.0.local_asn":      `64512`,
			"groups.0.peer_asn":       `1111`,
			"groups.0.router_id":      CHECKSET,
			"groups.0.hold":           CHECKSET,
			"groups.0.ip_version":     CHECKSET,
			"groups.0.is_fake_asn":    CHECKSET,
			"groups.0.keepalive":      CHECKSET,
			"groups.0.route_limit":    CHECKSET,
			"groups.0.id":             CHECKSET,
			"groups.0.status":         "Available",
		}
	}
	var fakeAlicloudVpcBgpGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpcBgpGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_bgp_groups.default",
		existMapFunc: existAlicloudVpcBgpGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcBgpGroupsDataSourceNameMapFunc,
	}

	alicloudVpcBgpGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, routerIdConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcBgpGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
  default = "tf-testAccBgpGroup-%d"
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

data "alicloud_vpc_bgp_groups" "default" {	
  %s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}

package vpngateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudVpnGatewayVcoRoutesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vco_route.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vco_route.default.id}_fake"]`,
		}),
	}
	routeEntryTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpn_gateway_vco_route.default.id}"]`,
			"route_entry_type": `"custom"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpn_gateway_vco_route.default.id}"]`,
			"route_entry_type": `"bgp"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_vco_route.default.id}"]`,
			"status": `"published"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_vco_route.default.id}"]`,
			"status": `"normal"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpn_gateway_vco_route.default.id}"]`,
			"route_entry_type": `"custom"`,
			"status":           `"published"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpn_gateway_vco_route.default.id}_fake"]`,
			"route_entry_type": `"bgp"`,
			"status":           `"normal"`,
		}),
	}
	var existAlicloudVpnGatewayVcoRoutesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"routes.#":                   "1",
			"routes.0.weight":            "100",
			"routes.0.vpn_connection_id": acctest.CHECKSET,
			"routes.0.next_hop":          acctest.CHECKSET,
			"routes.0.route_dest":        "192.168.12.0/24",
			"routes.0.status":            "published",
			"routes.0.as_path":           "",
			"routes.0.create_time":       acctest.CHECKSET,
			"routes.0.source":            "",
			"routes.0.id":                acctest.CHECKSET,
		}
	}
	var fakeAlicloudVpnGatewayVcoRoutesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudVpnGatewayVcoRoutesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vpn_gateway_vco_routes.default",
		ExistMapFunc: existAlicloudVpnGatewayVcoRoutesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudVpnGatewayVcoRoutesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudVpnGatewayVcoRoutesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, routeEntryTypeConf, statusConf, allConf)
}
func testAccCheckAlicloudVpnGatewayVcoRoutesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccVcoRoute-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id                     = alicloud_cen_instance.default.id
  transit_router_description = var.name
  transit_router_name        = var.name
}

resource "alicloud_cen_transit_router_cidr" "default" {
  transit_router_id        = alicloud_cen_transit_router.default.transit_router_id
  cidr                     = "192.168.0.0/16"
  transit_router_cidr_name = var.name
  description              = var.name
  publish_cidr_route       = true
}

data "alicloud_cen_transit_router_available_resources" "default" {}

resource "alicloud_vpn_customer_gateway" "default" {
  customer_gateway_name = var.name
  ip_address            = "42.104.22.%d"
  asn                   = "45014"
  description           = var.name
}

resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  remote_subnet       = "0.0.0.0/0"
  effect_immediately  = false
  ike_config {
    ike_auth_alg = "md5"
    ike_enc_alg  = "des"
    ike_version  = "ikev2"
    ike_mode     = "main"
    ike_lifetime = 86400
    psk          = "tf-examplevpn2"
    ike_pfs      = "group1"
    remote_id    = "testbob2"
    local_id     = "testalice2"
  }
  ipsec_config {
    ipsec_pfs      = "group5"
    ipsec_enc_alg  = "des"
    ipsec_auth_alg = "md5"
    ipsec_lifetime = 86400
  }
  bgp_config {
    enable       = true
    local_asn    = 45014
    tunnel_cidr  = "169.254.11.0/30"
    local_bgp_ip = "169.254.11.1"
  }
  health_check_config {
    enable   = true
    sip      = "192.168.1.1"
    dip      = "10.0.0.1"
    interval = 10
    retry    = 10
    policy   = "revoke_route"
  }
  enable_dpd           = true
  enable_nat_traversal = true
  vpn_attachment_name  = var.name
}
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
  auto_publish_route_enabled            = false
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
  cen_id                                = alicloud_cen_transit_router.default.cen_id
  transit_router_id                     = alicloud_cen_transit_router_cidr.default.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.default.id
  zone {
    zone_id = data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0
  }
}


resource "alicloud_vpn_gateway_vco_route" "default" {
  next_hop          = alicloud_cen_transit_router_vpn_attachment.default.vpn_id
  vpn_connection_id = alicloud_cen_transit_router_vpn_attachment.default.vpn_id
  weight            = "100"
  route_dest        = "192.168.10.0/24"
}

data "alicloud_vpn_gateway_vco_routes" "default" {	
	vpn_connection_id = alicloud_cen_transit_router_vpn_attachment.default.vpn_id
	%s
}
`, rand, rand%200+1, strings.Join(pairs, " \n "))
	return config
}

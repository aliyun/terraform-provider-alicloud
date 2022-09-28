package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVpnAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}
	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
			"transit_router_id": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
		}),
		fakeConfig: "",
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_vpn_attachment.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
			"status": `"Detaching"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}"`,
			"ids":               `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
			"status":            `"${alicloud_cen_transit_router_vpn_attachment.default.status}"`,
			"transit_router_id": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}_fake"`,
			"ids":        `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
			"status":     `"Detaching"`,
		}),
	}
	var existAlicloudCenTransitRouterVpnAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"names.#":       "1",
			"attachments.#": "1",
			"attachments.0.auto_publish_route_enabled":            "false",
			"attachments.0.transit_router_attachment_description": fmt.Sprintf("tf-testAccTransitRouterVpnAttachment-%d", rand),
			"attachments.0.transit_router_attachment_name":        fmt.Sprintf("tf-testAccTransitRouterVpnAttachment-%d", rand),
			"attachments.0.transit_router_id":                     CHECKSET,
			"attachments.0.vpn_id":                                CHECKSET,
			"attachments.0.vpn_owner_id":                          CHECKSET,
			"attachments.0.zone.#":                                "1",
			"attachments.0.transit_router_attachment_id":          CHECKSET,
			"attachments.0.id":                                    CHECKSET,
			"attachments.0.status":                                CHECKSET,
			"attachments.0.create_time":                           CHECKSET,
			"attachments.0.resource_type":                         CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterVpnAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterVpnAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_vpn_attachments.default",
		existMapFunc: existAlicloudCenTransitRouterVpnAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterVpnAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCenTransitRouterVpnAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, transitRouterIdConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterVpnAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTransitRouterVpnAttachment-%d"
}
resource "alicloud_cen_instance" "default" {
	cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
	cen_id = alicloud_cen_instance.default.id
	transit_router_description = "desd"
	transit_router_name = var.name
}

data "alicloud_cen_transit_router_available_resources" "default" {}

resource "alicloud_vpn_customer_gateway" "default" {
  name        = "${var.name}"
  ip_address  = "42.104.22.210"
  asn         = "45014"
  description = "testAccVpnConnectionDesc"
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
    psk          = "tf-testvpn2"
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
	auto_publish_route_enabled = false
	transit_router_attachment_description = var.name
	transit_router_attachment_name = var.name
	cen_id = alicloud_cen_transit_router.default.cen_id
	transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	vpn_id = alicloud_vpn_gateway_vpn_attachment.default.id
	zone {
		zone_id = data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0
	}
}

data "alicloud_cen_transit_router_vpn_attachments" "default" {
	cen_id = alicloud_cen_transit_router.default.cen_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

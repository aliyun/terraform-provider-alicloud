package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVpnAttachmentDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}

	transitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"transit_router_id": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
			"transit_router_id": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
		}),
	}

	transitRouterAttachmentIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpn_attachment.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"status": `"${alicloud_cen_transit_router_vpn_attachment.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"status": `"Detaching"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"tags": `{
							Created = "TF"
							For 	= "VpnAttachment"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"tags": `{
							Created = "TF_Fake"
							For 	= "VpnAttachment_Fake"
					  }`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
			"name_regex":                   `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}"`,
			"transit_router_id":            `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpn_attachment.default.id}"`,
			"status":                       `"${alicloud_cen_transit_router_vpn_attachment.default.status}"`,
			"tags": `{
							Created = "TF"
							For 	= "VpnAttachment"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
			"name_regex":                   `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}_fake"`,
			"transit_router_id":            `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"`,
			"status":                       `"Detaching"`,
			"tags": `{
							Created = "TF_Fake"
							For 	= "VpnAttachment_Fake"
					  }`,
		}),
	}

	CenTransitRouterVpnAttachmentCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, transitRouterIdConf, transitRouterAttachmentIdConf, statusConf, tagsConf, allConf)
}

var existCenTransitRouterVpnAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                "1",
		"names.#":              "1",
		"attachments.#":        "1",
		"attachments.0.status": CHECKSET,
		"attachments.0.transit_router_attachment_id":          CHECKSET,
		"attachments.0.vpn_owner_id":                          CHECKSET,
		"attachments.0.create_time":                           CHECKSET,
		"attachments.0.transit_router_attachment_name":        CHECKSET,
		"attachments.0.auto_publish_route_enabled":            CHECKSET,
		"attachments.0.charge_type":                           CHECKSET,
		"attachments.0.cen_id":                                CHECKSET,
		"attachments.0.transit_router_attachment_description": CHECKSET,
		"attachments.0.tags.%":                                CHECKSET,
		"attachments.0.transit_router_id":                     CHECKSET,
		"attachments.0.vpn_id":                                CHECKSET,
		"attachments.0.id":                                    CHECKSET,
		"attachments.0.resource_type":                         CHECKSET,
	}
}

var fakeCenTransitRouterVpnAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":         "0",
		"names.#":       "0",
		"attachments.#": "0",
	}
}

var CenTransitRouterVpnAttachmentCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_transit_router_vpn_attachments.default",
	existMapFunc: existCenTransitRouterVpnAttachmentMapFunc,
	fakeMapFunc:  fakeCenTransitRouterVpnAttachmentMapFunc,
}

func testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCenTransitRouterVpnAttachment%d"
}
data "alicloud_account" "default" {
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "test-vpn-attachment"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_cidr" "default" {
  cidr              = "192.168.10.0/24"
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  publish_cidr_route = true
}

resource "alicloud_vpn_customer_gateway" "default" {
  ip_address            = "1.1.1.8"
  customer_gateway_name = "test-vpn-attachment"
  depends_on            = [alicloud_cen_transit_router_cidr.default]
}

resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  network_type = "public"
  local_subnet = "0.0.0.0/0"
  enable_tunnels_bgp = "false"
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id = alicloud_vpn_customer_gateway.default.id
    enable_dpd = "true"
    enable_nat_traversal = "true"
    tunnel_index = "1"
    tunnel_ike_config {
      remote_id = "2.2.2.2"
      ike_enc_alg = "aes"
      ike_mode = "main"
      ike_version = "ikev1"
      local_id = "1.1.1.1"
      ike_auth_alg = "md5"
      ike_lifetime = "86100"
      ike_pfs = "group2"
      psk = "12345678"
    }
    
      tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg = "aes"
      ipsec_lifetime = "86200"
      ipsec_pfs = "group5"
    }
    
  }
  tunnel_options_specification {
    enable_nat_traversal = "true"
    tunnel_index = "2"
      tunnel_ike_config {
      local_id = "4.4.4.4"
      remote_id = "5.5.5.5"
      ike_lifetime = "86400"
      ike_pfs = "group5"
      ike_mode = "main"
      ike_version = "ikev2"
      psk = "32333442"
      ike_auth_alg = "md5"
      ike_enc_alg = "aes"
    }
    
      tunnel_ipsec_config {
      ipsec_enc_alg = "aes"
      ipsec_lifetime = "86400"
      ipsec_pfs = "group5"
      ipsec_auth_alg = "sha256"
    }
    
    customer_gateway_id = alicloud_vpn_customer_gateway.default.id
    enable_dpd = "true"
  }
  
  remote_subnet = "0.0.0.0/0"
}
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
	auto_publish_route_enabled = false
	transit_router_attachment_description = var.name
	transit_router_attachment_name = var.name
	cen_id = alicloud_cen_transit_router.default.cen_id
	transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	vpn_id = alicloud_vpn_gateway_vpn_attachment.default.id
  	vpn_owner_id = data.alicloud_account.default.id
	charge_type = "POSTPAY"
	tags = {
		Created = "TF"
		For     = "VpnAttachment"
	}
}

data "alicloud_cen_transit_router_vpn_attachments" "default" {
	cen_id = alicloud_cen_transit_router.default.cen_id
	%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

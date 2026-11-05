package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenTransitRouterVpnAttachmentDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpn_attachment.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}

	transitRouterIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"transit_router_id": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"]`,
			"transit_router_id": `"${alicloud_cen_transit_router_vpn_attachment.default.transit_router_id}"`,
		}),
	}

	transitRouterAttachmentIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpn_attachment.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpn_attachment.default.id}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"status": `"${alicloud_cen_transit_router_vpn_attachment.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"status": `"Detaching"`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"tags": `{
							Created = "TF"
							For 	= "VpnAttachment"
					  }`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
			"tags": `{
							Created = "TF_Fake"
							For 	= "VpnAttachment_Fake"
					  }`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
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
		FakeConfig: testAccCheckAlicloudCenTransitRouterVpnAttachmentSourceConfig(rand, map[string]string{
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

	CenTransitRouterVpnAttachmentCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, transitRouterIdConf, transitRouterAttachmentIdConf, statusConf, tagsConf, allConf)
}

var existCenTransitRouterVpnAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                "1",
		"names.#":              "1",
		"attachments.#":        "1",
		"attachments.0.status": acctest.CHECKSET,
		"attachments.0.transit_router_attachment_id":          acctest.CHECKSET,
		"attachments.0.vpn_owner_id":                          acctest.CHECKSET,
		"attachments.0.create_time":                           acctest.CHECKSET,
		"attachments.0.transit_router_attachment_name":        acctest.CHECKSET,
		"attachments.0.auto_publish_route_enabled":            acctest.CHECKSET,
		"attachments.0.charge_type":                           acctest.CHECKSET,
		"attachments.0.cen_id":                                acctest.CHECKSET,
		"attachments.0.transit_router_attachment_description": acctest.CHECKSET,
		"attachments.0.tags.%":                                acctest.CHECKSET,
		"attachments.0.transit_router_id":                     acctest.CHECKSET,
		"attachments.0.vpn_id":                                acctest.CHECKSET,
		"attachments.0.id":                                    acctest.CHECKSET,
		"attachments.0.resource_type":                         acctest.CHECKSET,
	}
}

var fakeCenTransitRouterVpnAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":         "0",
		"names.#":       "0",
		"attachments.#": "0",
	}
}

var CenTransitRouterVpnAttachmentCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cen_transit_router_vpn_attachments.default",
	ExistMapFunc: existCenTransitRouterVpnAttachmentMapFunc,
	FakeMapFunc:  fakeCenTransitRouterVpnAttachmentMapFunc,
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

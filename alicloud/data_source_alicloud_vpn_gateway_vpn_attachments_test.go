package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpnGatewayVpnAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vpn_attachment.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
			"status": `"${alicloud_vpn_gateway_vpn_attachment.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
			"status": `"attaching"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
			"name_regex": `"${alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name}"`,
			"status":     `"${alicloud_vpn_gateway_vpn_attachment.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpn_gateway_vpn_attachment.default.id}_fake"]`,
			"name_regex": `"${alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name}_fake"`,
			"status":     `"attaching"`,
		}),
	}
	var existAlicloudVpnGatewayVpnAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"names.#":                                      "1",
			"attachments.#":                                "1",
			"attachments.0.effect_immediately":             "false",
			"attachments.0.local_subnet":                   "0.0.0.0/0",
			"attachments.0.network_type":                   "public",
			"attachments.0.remote_subnet":                  "0.0.0.0/0",
			"attachments.0.vpn_attachment_name":            fmt.Sprintf("tf-testAccVpnAttachment-%d", rand),
			"attachments.0.ike_config.#":                   "1",
			"attachments.0.ike_config.0.ike_auth_alg":      "md5",
			"attachments.0.ike_config.0.ike_enc_alg":       "des",
			"attachments.0.ike_config.0.ike_version":       "ikev2",
			"attachments.0.ike_config.0.ike_mode":          "main",
			"attachments.0.ike_config.0.ike_lifetime":      "86400",
			"attachments.0.ike_config.0.psk":               "tf-testvpn2",
			"attachments.0.ike_config.0.ike_pfs":           "group1",
			"attachments.0.ike_config.0.remote_id":         "testbob2",
			"attachments.0.ike_config.0.local_id":          "testalice2",
			"attachments.0.ipsec_config.#":                 "1",
			"attachments.0.ipsec_config.0.ipsec_pfs":       "group5",
			"attachments.0.ipsec_config.0.ipsec_enc_alg":   "des",
			"attachments.0.ipsec_config.0.ipsec_auth_alg":  "md5",
			"attachments.0.ipsec_config.0.ipsec_lifetime":  "86400",
			"attachments.0.bgp_config.#":                   "1",
			"attachments.0.bgp_config.0.status":            "",
			"attachments.0.bgp_config.0.local_asn":         "45014",
			"attachments.0.bgp_config.0.local_bgp_ip":      "169.254.11.1",
			"attachments.0.bgp_config.0.tunnel_cidr":       "169.254.11.0/30",
			"attachments.0.health_check_config.#":          "1",
			"attachments.0.health_check_config.0.enable":   "true",
			"attachments.0.health_check_config.0.dip":      "10.0.0.1",
			"attachments.0.health_check_config.0.retry":    "10",
			"attachments.0.health_check_config.0.sip":      "192.168.1.1",
			"attachments.0.health_check_config.0.interval": "10",
			"attachments.0.health_check_config.0.policy":   "revoke_route",
			"attachments.0.health_check_config.0.status":   "",
			"attachments.0.create_time":                    CHECKSET,
			"attachments.0.customer_gateway_id":            CHECKSET,
			"attachments.0.status":                         CHECKSET,
			"attachments.0.connection_status":              "",
			"attachments.0.id":                             CHECKSET,
			"attachments.0.vpn_connection_id":              CHECKSET,
		}
	}
	var fakeAlicloudVpnGatewayVpnAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpnGatewayVpnAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpn_gateway_vpn_attachments.default",
		existMapFunc: existAlicloudVpnGatewayVpnAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpnGatewayVpnAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpnGatewayVpnAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudVpnGatewayVpnAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccVpnAttachment-%d"
}
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

data "alicloud_vpn_gateway_vpn_attachments" "default" {	
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

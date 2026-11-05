package vpngateway_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudVpnGatewayVpnAttachmentDataSource(t *testing.T) {
	tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vpn_attachment.default.id}_fake"]`,
		}),
	}
	nameRegexConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name}_fake"`,
		}),
	}
	statusConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
			"status": `"${alicloud_vpn_gateway_vpn_attachment.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
			"status": `"attaching"`,
		}),
	}

	allConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vpn_attachment.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_vpn_attachment.default.id}_fake"]`,
		}),
	}

	VpnGatewayVpnAttachmentCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, statusConf, idsConf, allConf)
}

var existVpnGatewayVpnAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                                      "1",
		"attachments.#":                                "1",
		"attachments.0.attach_type":                    tfacctest.CHECKSET,
		"attachments.0.resource_group_id":              tfacctest.CHECKSET,
		"attachments.0.enable_tunnels_bgp":             tfacctest.CHECKSET,
		"attachments.0.effect_immediately":             tfacctest.CHECKSET,
		"attachments.0.bgp_config.#":                   tfacctest.CHECKSET,
		"attachments.0.remote_subnet":                  tfacctest.CHECKSET,
		"attachments.0.network_type":                   tfacctest.CHECKSET,
		"attachments.0.ipsec_config.#":                 tfacctest.CHECKSET,
		"attachments.0.enable_nat_traversal":           tfacctest.CHECKSET,
		"attachments.0.ike_config.#":                   tfacctest.CHECKSET,
		"attachments.0.tags.%":                         tfacctest.CHECKSET,
		"attachments.0.status":                         tfacctest.CHECKSET,
		"attachments.0.local_subnet":                   tfacctest.CHECKSET,
		"attachments.0.vpn_attachment_name":            tfacctest.CHECKSET,
		"attachments.0.create_time":                    tfacctest.CHECKSET,
		"attachments.0.tunnel_options_specification.#": "2",
		"attachments.0.vpn_connection_id":              tfacctest.CHECKSET,
		"attachments.0.health_check_config.#":          tfacctest.CHECKSET,
		"attachments.0.enable_dpd":                     tfacctest.CHECKSET,
		"attachments.0.connection_status":              tfacctest.CHECKSET,
	}
}

var fakeVpnGatewayVpnAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#": "0",
	}
}

var VpnGatewayVpnAttachmentCheckInfo = tfacctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpn_gateway_vpn_attachments.default",
	ExistMapFunc: existVpnGatewayVpnAttachmentMapFunc,
	FakeMapFunc:  fakeVpnGatewayVpnAttachmentMapFunc,
}

func testAccCheckAlicloudVpnGatewayVpnAttachmentSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpnGatewayVpnAttachment%d"
}
variable "region_id" {
  default = "cn-huhehaote"
}

variable "az2" {
  default = "cn-huhehaote-b"
}

variable "az1" {
  default = "cn-huhehaote-a"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpn_customer_gateway" "cgw1" {
  ip_address = "2.2.2.2"
  asn        = "1219001"
}

resource "alicloud_vpn_customer_gateway" "cgw2" {
  ip_address            = "43.43.3.22"
  asn                   = "44331"
  customer_gateway_name = "test_amp"
}



resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  local_subnet        = "0.0.0.0/0"
  enable_tunnels_bgp  = true
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.cgw1.id
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_index         = "1"
    tunnel_bgp_config {
      local_asn    = "1219001"
      local_bgp_ip = "169.254.10.1"
      tunnel_cidr  = "169.254.10.0/30"
    }
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
      ike_lifetime = "86100"
      ike_mode     = "main"
      ike_pfs      = "group2"
      ike_version  = "ikev1"
      local_id     = "1.1.1.1"
      psk          = "12345678"
      remote_id    = "2.2.2.2"
    }
    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes"
      ipsec_lifetime = "86200"
      ipsec_pfs      = "group5"
    }
  }
  tunnel_options_specification {
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_index         = "2"
    tunnel_bgp_config {
      local_asn    = "1219001"
      local_bgp_ip = "169.254.20.1"
      tunnel_cidr  = "169.254.20.0/30"
    }
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
      ike_lifetime = "86400"
      ike_mode     = "main"
      ike_pfs      = "group5"
      ike_version  = "ikev2"
      local_id     = "4.4.4.4"
      psk          = "32333442"
      remote_id    = "5.5.5.5"
    }
    tunnel_ipsec_config {
      ipsec_auth_alg = "sha256"
      ipsec_enc_alg  = "aes"
      ipsec_lifetime = "86400"
      ipsec_pfs      = "group5"
    }
    customer_gateway_id = alicloud_vpn_customer_gateway.cgw1.id
  }
  remote_subnet = "0.0.0.0/0"
  network_type  = "public"
  tags = {
    test  = "1"
    test2 = "2"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

data "alicloud_vpn_gateway_vpn_attachments" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPNConnectionsDataSourceBasic(t *testing.T) {

	resourceId := "data.alicloud_vpn_connections.default"
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccVpnConnDataResource%d", rand)
	preCheck := func() {
		testAccPreCheck(t)
	}
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		name, dataSourceVpnConnectionsConfigDependence)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpn_connection.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpn_connection.default.id}_fake"},
		}),
	}

	vpnGateWayIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpn_gateway_id": "${alicloud_vpn_connection.default.vpn_gateway_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpn_gateway_id": "${alicloud_vpn_connection.default.vpn_gateway_id}_fake",
		}),
	}

	//customerGatewayIdConfig := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
	//	}),
	//	fakeConfig: testAccConfig(map[string]interface{}{
	//		"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}_fake",
	//	}),
	//}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "_fake",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpn_connection.default.id}"},
			//"vpn_gateway_id":      "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
			"customer_gateway_id": "${alicloud_vpn_connection.default.customer_gateway_id}",
			"name_regex":          name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpn_connection.default.id}"},
			//"vpn_gateway_id":      "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
			"customer_gateway_id": "${alicloud_vpn_connection.default.customer_gateway_id}",
			"name_regex":          name + "_fake",
		}),
	}

	var existVpnConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"connections.#":                                "1",
			"ids.#":                                        "1",
			"names.#":                                      "1",
			"connections.0.name":                           fmt.Sprintf("tf-testAccVpnConnDataResource%d", rand),
			"connections.0.vpn_gateway_id":                 CHECKSET,
			"connections.0.customer_gateway_id":            "",
			"connections.0.status":                         "",
			"connections.0.local_subnet":                   "3.0.0.0/24",
			"connections.0.enable_dpd":                     "false",
			"connections.0.enable_nat_traversal":           "false",
			"connections.0.remote_subnet":                  "10.0.1.0/24,10.0.0.0/24",
			"connections.0.tunnel_options_specification.#": CHECKSET,
		}
	}

	var fakeVpnConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"connections.#": "0",
			"ids.#":         "0",
			"names.#":       "0",
		}
	}

	var vpnConnectionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpnConnectionsMapFunc,
		fakeMapFunc:  fakeVpnConnectionsMapFunc,
	}
	vpnConnectionsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConfig, vpnGateWayIdConfig,
		nameRegexConfig, allConfig)

}

func dataSourceVpnConnectionsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "spec" {
  default = "5"
}

data "alicloud_vpn_gateway_zones" "default" {
  spec = "5M"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 1)
  zone_id      = data.alicloud_vpn_gateway_zones.default.ids.0
  vswitch_name = var.name
}

data "alicloud_vswitches" "default2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.1
}

resource "alicloud_vswitch" "vswitch2" {
  count        = length(data.alicloud_vswitches.default2.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_vpn_gateway_zones.default.ids.1
  vswitch_name = var.name
}

locals {
  vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  vswitch_id2 = length(data.alicloud_vswitches.default2.ids) > 0 ? data.alicloud_vswitches.default2.ids[0] : concat(alicloud_vswitch.vswitch2.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "HA-VPN" {
  vpn_type                     = "Normal"
  disaster_recovery_vswitch_id = local.vswitch_id2
  vpn_gateway_name             = var.name

  vswitch_id   = local.vswitch_id
  auto_pay     = true
  vpc_id       = data.alicloud_vpcs.default.ids.0
  network_type = "public"
  payment_type = "Subscription"
  enable_ipsec = true
  bandwidth    = var.spec
}

resource "alicloud_vpn_customer_gateway" "defaultCustomerGateway" {
  description           = "defaultCustomerGateway"
  ip_address            = "2.2.2.5"
  asn                   = "2224"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_customer_gateway" "changeCustomerGateway" {
  description           = "changeCustomerGateway"
  ip_address            = "2.2.2.6"
  asn                   = "2225"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_connection" "default" {
  vpn_gateway_id      = alicloud_vpn_gateway.HA-VPN.id
  vpn_connection_name = var.name
  local_subnet = [
    "3.0.0.0/24"
  ]
  remote_subnet = [
    "10.0.0.0/24",
    "10.0.1.0/24"
  ]
  tags = {
    Created = "TF"
    For     = "example"
  }
  enable_tunnels_bgp = "true"
  tunnel_options_specification {
    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
      ipsec_lifetime = "16400"
      ipsec_pfs      = "group5"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "master"
    tunnel_bgp_config {
      local_asn    = "1219002"
      tunnel_cidr  = "169.254.30.0/30"
      local_bgp_ip = "169.254.30.1"
    }

    tunnel_ike_config {
      ike_mode     = "aggressive"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel2"
      psk          = "12345678"
      remote_id    = "remote2"
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes256"
      ike_lifetime = "3600"
      ike_pfs      = "group14"
    }

  }
  tunnel_options_specification {
    tunnel_ike_config {
      remote_id    = "remote24"
      ike_enc_alg  = "aes256"
      ike_lifetime = "27000"
      ike_mode     = "aggressive"
      ike_pfs      = "group5"
      ike_auth_alg = "md5"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel2"
      psk          = "12345678"
    }

    tunnel_ipsec_config {
      ipsec_lifetime = "2700"
      ipsec_pfs      = "group14"
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "slave"
    tunnel_bgp_config {
      local_asn    = "1219002"
      local_bgp_ip = "169.254.40.1"
      tunnel_cidr  = "169.254.40.0/30"
    }
  }
}

`, name)
}

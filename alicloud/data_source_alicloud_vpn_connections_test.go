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
		testAccPreCheckWithAccountSiteType(t, IntlSite)
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

	customerGatewayIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"customer_gateway_id": "${alicloud_vpn_connection.default.customer_gateway_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"customer_gateway_id": "${alicloud_vpn_connection.default.customer_gateway_id}_fake",
		}),
	}

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
			"ids":                 []string{"${alicloud_vpn_connection.default.id}"},
			"vpn_gateway_id":      "${alicloud_vpn_connection.default.vpn_gateway_id}",
			"customer_gateway_id": "${alicloud_vpn_connection.default.customer_gateway_id}",
			"name_regex":          name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_vpn_connection.default.id}"},
			"vpn_gateway_id":      "${alicloud_vpn_connection.default.vpn_gateway_id}",
			"customer_gateway_id": "${alicloud_vpn_connection.default.customer_gateway_id}",
			"name_regex":          name + "_fake",
		}),
	}

	var existVpnConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"connections.#":                               "1",
			"ids.#":                                       "1",
			"names.#":                                     "1",
			"connections.0.name":                          fmt.Sprintf("tf-testAccVpnConnDataResource%d", rand),
			"connections.0.vpn_gateway_id":                CHECKSET,
			"connections.0.customer_gateway_id":           CHECKSET,
			"connections.0.status":                        "ike_sa_not_established",
			"connections.0.local_subnet":                  "172.16.1.0/24",
			"connections.0.remote_subnet":                 "10.4.0.0/24",
			"connections.0.ike_config.0.ike_auth_alg":     "sha1",
			"connections.0.ike_config.0.ike_enc_alg":      "3des",
			"connections.0.ike_config.0.ike_version":      "ikev2",
			"connections.0.ike_config.0.ike_mode":         "aggressive",
			"connections.0.ike_config.0.ike_lifetime":     "8640",
			"connections.0.ike_config.0.psk":              "tf-testvpn1",
			"connections.0.ike_config.0.ike_pfs":          "group2",
			"connections.0.ike_config.0.ike_remote_id":    "testbob1",
			"connections.0.ike_config.0.ike_local_id":     "testalice1",
			"connections.0.ipsec_config.0.ipsec_pfs":      "group2",
			"connections.0.ipsec_config.0.ipsec_enc_alg":  "aes",
			"connections.0.ipsec_config.0.ipsec_auth_alg": "sha1",
			"connections.0.ipsec_config.0.ipsec_lifetime": "86400",
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
		customerGatewayIdConfig, nameRegexConfig, allConfig)

}

func dataSourceVpnConnectionsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}


resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name = "${var.name}"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "default" {
	name = "${var.name}"
	ip_address = "41.104.22.229"
	description = "${var.name}"
}

resource "alicloud_vpn_connection" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.default.id}"
	local_subnet = ["172.16.1.0/24"]
	remote_subnet = ["10.4.0.0/24"]
	effect_immediately = true
	ike_config = [{
        ike_auth_alg = "sha1"
        ike_enc_alg = "3des"
        ike_version = "ikev2"
        ike_mode = "aggressive"
        ike_lifetime = 8640
        psk = "tf-testvpn1"
        ike_pfs = "group2"
        ike_remote_id = "testbob1"
        ike_local_id = "testalice1"
        }
    ]
	ipsec_config = [{
        ipsec_pfs = "group2"
        ipsec_enc_alg = "aes"
        ipsec_auth_alg = "sha1"
        ipsec_lifetime = 86400
    }]
}

`, name)
}

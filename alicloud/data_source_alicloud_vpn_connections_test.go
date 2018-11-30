package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpnConnectionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnConnsDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_connections.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.foo", "connections.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.name", "tf-testAccVpnConnDataResource"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_connections.foo", "connections.0.vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_connections.foo", "connections.0.customer_gateway_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.status", "ike_sa_not_established"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.remote_subnet", "10.4.0.0/24"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_auth_alg", "sha1"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_enc_alg", "3des"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_version", "ikev2"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_mode", "aggressive"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_lifetime", "8640"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.psk", "tf-testvpn1"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_pfs", "group2"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_remote_id", "testbob1"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ike_config.0.ike_local_id", "testalice1"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ipsec_config.0.ipsec_pfs", "group2"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ipsec_config.0.ipsec_enc_alg", "aes"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ipsec_config.0.ipsec_auth_alg", "sha1"),
					resource.TestCheckResourceAttr(
						"data.alicloud_vpn_connections.foo", "connections.0.ipsec_config.0.ipsec_lifetime", "86400"),
				),
			},
		},
	})
}

func TestAccAlicloudVpnConnectionsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnConnsDataEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_connections.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.foo", "connections.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.vpn_gateway_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.customer_gateway_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.local_subnet"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.remote_subnet"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_connections.foo", "connections.0.ike_config.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnConnsDataCfg = `
variable "name" {
	default = "tf-testAccVpnConnDataResource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "${var.name}"
	ip_address = "41.104.22.229"
	description = "${var.name}"
}

resource "alicloud_vpn_connection" "foo" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
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

data "alicloud_vpn_connections" "foo" {
	ids = ["${alicloud_vpn_connection.foo.id}"]
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
}
`

const testAccCheckAlicloudVpnConnsDataEmpty = `
data "alicloud_vpn_connections" "foo" {
	name_regex = "^tf-testacc-fake-name"
}
`

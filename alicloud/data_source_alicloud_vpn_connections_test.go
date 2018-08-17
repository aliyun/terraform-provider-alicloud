package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const ipsec_config = "{\"IpsecAuthAlg\":\"sha1\",\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecLifetime\":7200}"

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
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_connections.vpn_conns"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.vpn_conns", "vpn_connections.0.name", "tf-vco_test1"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.vpn_conns", "vpn_connections.0.local_subnet", "172.16.0.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.vpn_conns", "vpn_connections.0.remote_subnet", "10.0.0.0/8"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_connections.vpn_conns", "vpn_connections.0.ike_config"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_connections.vpn_conns", "vpn_connections.0.ipsec_config", ipsec_config),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnConnsDataCfg = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "testAccVpcConfig"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchConfig"
}

resource "alicloud_vpn" "foo" {
        name = "testAccVpnConfig_create"
        vpc_id = "${alicloud_vpc.foo.id}"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "testAccVpnCgwName"
	ip_address = "41.104.22.229"
	description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "vpn_conns" {
	name = "tf-vco_test1"
	vpn_gateway_id = "${alicloud_vpn.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
	local_subnet = "172.16.0.0/16"
	remote_subnet = "10.0.0.0/8"
	effect_immediately = true
	ike_config = "{\"LocalId\":\"${alicloud_vpn.foo.internet_ip}\",\"IkeAuthAlg\":\"sha1\",\"IkePfs\":\"group2\",\"IkeMode\":\"main\",\"IkeEncAlg\":\"aes\",\"Psk\":\"tf-testvpn1\",\"RemoteId\":\"{$alicloud_vpn_customer_gateway.ip_address}\",\"IkeVersion\":\"ikev2\",\"IkeLifetime\":3600}"
	ipsec_config = "{\"IpsecAuthAlg\":\"sha1\",\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecLifetime\":7200}"
}

data "alicloud_vpn_connections" "vpn_conns" {
	vpn_connection_id = "${alicloud_vpn_connection.vpn_conns.id}"
	output_file = "/tmp/vpnconns"
}
`
